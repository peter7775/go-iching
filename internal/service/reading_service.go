package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/example/iching-fiber-app/internal/domain"
	"github.com/example/iching-fiber-app/internal/iching"
)

type ReadingRepository interface {
	Save(context.Context, domain.Reading) (domain.Reading, error)
	List(context.Context) ([]domain.Reading, error)
	Get(context.Context, string) (domain.Reading, error)
	SaveReflection(context.Context, string, int, string, time.Time) error
}

type ReadingService struct {
	repo ReadingRepository
	en   map[int]iching.WilhelmHexagram
	cs   map[int]iching.CzechHexagram
}

func NewReadingService(repo ReadingRepository) *ReadingService {
	en, err := iching.LoadEnglish()
	if err != nil { panic(err) }
	cs, err := iching.LoadCzech()
	if err != nil { panic(err) }
	return &ReadingService{repo: repo, en: en, cs: cs}
}

type CreateReadingInput struct {
	Question string            `json:"question"`
	Method   domain.CastMethod `json:"method"`
	Language string            `json:"language"`
	Lines    []domain.Line     `json:"lines"`
	Random   bool              `json:"random"`
}

type ReflectionInput struct {
	Rating int    `json:"rating"`
	Note   string `json:"note"`
}

func (s *ReadingService) Create(ctx context.Context, in CreateReadingInput) (domain.Reading, error) {
	if in.Question == "" { return domain.Reading{}, errors.New("question is required") }
	lang := in.Language
	if lang == "" { lang = "cs" }
	lines := in.Lines
	if in.Random || in.Method == domain.MethodCoins {
		var err error
		lines, err = iching.RandomHexagramLines()
		if err != nil { return domain.Reading{}, err }
	}
	if len(lines) != 6 { return domain.Reading{}, errors.New("exactly 6 lines are required") }
	primary := iching.HexagramNumberFromLines(lines)
	relatingLines := iching.RelatingLines(lines)
	relating := iching.HexagramNumberFromLines(relatingLines)
	changing := changingLines(lines)
	interp := s.buildInterpretation(in.Question, lang, primary, relating, changing)
	now := time.Now().UTC()
	r := domain.Reading{ID: newID(), Question: in.Question, Method: in.Method, Language: lang, Lines: lines, PrimaryNumber: primary, RelatingNumber: relating, ChangingLines: changing, Interpretation: interp, CreatedAt: now, Reflection: reflectionState(now, nil, 0, "")}
	return s.repo.Save(ctx, r)
}

func (s *ReadingService) List(ctx context.Context) ([]domain.Reading, error) {
	items, err := s.repo.List(ctx)
	if err != nil { return nil, err }
	for i := range items { items[i].Reflection = normalizeReflection(items[i]) }
	return items, nil
}

func (s *ReadingService) Get(ctx context.Context, id string) (domain.Reading, error) {
	item, err := s.repo.Get(ctx, id)
	if err != nil { return domain.Reading{}, err }
	item.Reflection = normalizeReflection(item)
	return item, nil
}

func (s *ReadingService) SaveReflection(ctx context.Context, id string, in ReflectionInput) (domain.Reading, error) {
	if in.Rating < 0 || in.Rating > 5 { return domain.Reading{}, errors.New("rating must be between 0 and 5") }
	item, err := s.repo.Get(ctx, id)
	if err != nil { return domain.Reading{}, err }
	state := normalizeReflection(item)
	if !state.Eligible { return domain.Reading{}, fmt.Errorf("reflection available after %s", state.EligibleAt.Format(time.RFC3339)) }
	now := time.Now().UTC()
	if err := s.repo.SaveReflection(ctx, id, in.Rating, strings.TrimSpace(in.Note), now); err != nil { return domain.Reading{}, err }
	updated, err := s.repo.Get(ctx, id)
	if err != nil { return domain.Reading{}, err }
	updated.Reflection = normalizeReflection(updated)
	return updated, nil
}

func (s *ReadingService) RelevanceStats(ctx context.Context) (domain.RelevanceStats, error) {
	items, err := s.repo.List(ctx)
	if err != nil { return domain.RelevanceStats{}, err }
	dist := map[int]int{0:0,1:0,2:0,3:0,4:0,5:0}
	var sum int
	var rated int
	for _, item := range items {
		state := normalizeReflection(item)
		if state.CreatedAt != nil {
			dist[state.Rating]++
			sum += state.Rating
			rated++
		}
	}
	stats := domain.RelevanceStats{RatedCount: rated, TotalCount: len(items), Distribution: dist}
	if rated > 0 {
		stats.AverageRating = float64(sum) / float64(rated)
		stats.RelevancePct = (float64(sum) / float64(rated*5)) * 100
	}
	return stats, nil
}

func normalizeReflection(r domain.Reading) domain.Reflection {
	var created *time.Time
	if r.Reflection.CreatedAt != nil { created = r.Reflection.CreatedAt }
	return reflectionState(r.CreatedAt, created, r.Reflection.Rating, r.Reflection.Note)
}

func reflectionState(createdAt time.Time, reflectionAt *time.Time, rating int, note string) domain.Reflection {
	eligibleAt := createdAt.Add(7 * 24 * time.Hour)
	eligible := !time.Now().UTC().Before(eligibleAt)
	return domain.Reflection{Rating: rating, Note: note, CreatedAt: reflectionAt, Eligible: eligible, EligibleAt: eligibleAt}
}

func (s *ReadingService) buildInterpretation(question, lang string, primary, relating int, changing []int) domain.ReadingInterpretation {
	result := domain.ReadingInterpretation{Language: lang, Summary: summaryFor(lang)}
	if lang == "cs" {
		if p, ok := s.cs[primary]; ok {
			result.Primary = mapCzech(primary, p)
			for _, lineNo := range changing {
				if line, ok := p.LinesCS[fmt.Sprint(lineNo)]; ok {
					result.ChangingLines = append(result.ChangingLines, domain.InterpretationLine{Line: lineNo, Text: line.Text, Comments: line.Comments})
				}
			}
		}
		if r, ok := s.cs[relating]; ok { result.Relating = mapCzech(relating, r) }
	} else {
		if p, ok := s.en[primary]; ok {
			result.Primary = mapEnglish(primary, p)
			for _, lineNo := range changing {
				if line, ok := p.Lines[fmt.Sprint(lineNo)]; ok {
					result.ChangingLines = append(result.ChangingLines, domain.InterpretationLine{Line: lineNo, Text: line.Text, Comments: line.Comments})
				}
			}
		}
		if r, ok := s.en[relating]; ok { result.Relating = mapEnglish(relating, r) }
	}
	result.Markdown = renderMarkdown(question, result, changing)
	return result
}

func mapEnglish(n int, h iching.WilhelmHexagram) domain.InterpretationHexagram {
	return domain.InterpretationHexagram{Number:n, Name:h.Name, Title:h.Title, Character:h.Character, Traditional:h.Traditional, Pinyin:h.Pinyin, Above:domain.InterpretationTrigram{Chinese:h.Above.Chinese, Symbolic:h.Above.Symbolic, Alchemical:h.Above.Alchemical}, Below:domain.InterpretationTrigram{Chinese:h.Below.Chinese, Symbolic:h.Below.Symbolic, Alchemical:h.Below.Alchemical}, Judgment:domain.InterpretationSection{Text:h.Judgment.Text, Comments:h.Judgment.Comments}, Image:domain.InterpretationSection{Text:h.Image.Text, Comments:h.Image.Comments}}
}
func mapCzech(n int, h iching.CzechHexagram) domain.InterpretationHexagram {
	return domain.InterpretationHexagram{Number:n, Name:h.NameCS, Title:h.TitleCS, Character:h.Character, Traditional:h.Traditional, Pinyin:h.Pinyin, Above:domain.InterpretationTrigram{Chinese:h.Above.Chinese, Symbolic:h.Above.SymbolicCS, Alchemical:h.Above.AlchemicalCS}, Below:domain.InterpretationTrigram{Chinese:h.Below.Chinese, Symbolic:h.Below.SymbolicCS, Alchemical:h.Below.AlchemicalCS}, Judgment:domain.InterpretationSection{Text:h.JudgmentCS.Text, Comments:h.JudgmentCS.Comments}, Image:domain.InterpretationSection{Text:h.ImageCS.Text, Comments:h.ImageCS.Comments}}
}
func summaryFor(lang string) string {
	if lang == "en" { return "Interpretation combines the primary hexagram, changing lines, and the relating hexagram." }
	return "Výklad kombinuje primární hexagram, proměnné čáry a výsledný hexagram."
}
func renderMarkdown(question string, i domain.ReadingInterpretation, changing []int) string {
	var b strings.Builder
	if i.Language == "en" { b.WriteString("# I Ching reading\n\n## Question\n\n"+question+"\n\n") } else { b.WriteString("# I-ťing věštba\n\n## Otázka\n\n"+question+"\n\n") }
	if i.Language == "en" { b.WriteString("## Primary hexagram\n\n") } else { b.WriteString("## Primární hexagram\n\n") }
	fmt.Fprintf(&b, "**%d — %s**  \n", i.Primary.Number, i.Primary.Title)
	fmt.Fprintf(&b, "%s %s %s\n\n", i.Primary.Character, i.Primary.Traditional, i.Primary.Pinyin)
	if i.Language == "en" {
		fmt.Fprintf(&b, "Upper trigram: %s / %s / %s  \n", i.Primary.Above.Chinese, i.Primary.Above.Symbolic, i.Primary.Above.Alchemical)
		fmt.Fprintf(&b, "Lower trigram: %s / %s / %s\n\n", i.Primary.Below.Chinese, i.Primary.Below.Symbolic, i.Primary.Below.Alchemical)
		b.WriteString("### Judgment\n\n" + i.Primary.Judgment.Text + "\n\n")
		b.WriteString("### Image\n\n" + i.Primary.Image.Text + "\n\n")
	} else {
		fmt.Fprintf(&b, "Horní trigram: %s / %s / %s  \n", i.Primary.Above.Chinese, i.Primary.Above.Symbolic, i.Primary.Above.Alchemical)
		fmt.Fprintf(&b, "Dolní trigram: %s / %s / %s\n\n", i.Primary.Below.Chinese, i.Primary.Below.Symbolic, i.Primary.Below.Alchemical)
		b.WriteString("### Rozsudek\n\n" + i.Primary.Judgment.Text + "\n\n")
		b.WriteString("### Obraz\n\n" + i.Primary.Image.Text + "\n\n")
	}
	if len(i.ChangingLines) > 0 {
		if i.Language == "en" { b.WriteString("## Changing lines\n\n") } else { b.WriteString("## Proměnné čáry\n\n") }
		for _, line := range i.ChangingLines {
			if i.Language == "en" { fmt.Fprintf(&b, "### Line %d\n\n%s\n\n", line.Line, line.Text) } else { fmt.Fprintf(&b, "### Čára %d\n\n%s\n\n", line.Line, line.Text) }
		}
	}
	if i.Language == "en" { b.WriteString("## Relating hexagram\n\n") } else { b.WriteString("## Výsledný hexagram\n\n") }
	fmt.Fprintf(&b, "**%d — %s**  \n", i.Relating.Number, i.Relating.Title)
	fmt.Fprintf(&b, "%s %s %s\n\n", i.Relating.Character, i.Relating.Traditional, i.Relating.Pinyin)
	return b.String()
}
func changingLines(lines []domain.Line) []int { out := []int{}; for _, l := range lines { if l.IsChanging() { out = append(out, l.Position) } }; return out }
func newID() string { b := make([]byte, 8); _, _ = rand.Read(b); return hex.EncodeToString(b) }
