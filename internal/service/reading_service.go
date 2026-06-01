package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/example/iching-app/internal/domain"
)

type ReadingRepository interface {
	Save(context.Context, domain.Reading) (domain.Reading, error)
	List(context.Context) ([]domain.Reading, error)
	Get(context.Context, string) (domain.Reading, error)
}

type ReadingService struct {
	repo ReadingRepository
}

func NewReadingService(repo ReadingRepository) *ReadingService {
	return &ReadingService{repo: repo}
}

type CreateReadingInput struct {
	Question string             `json:"question"`
	Method   domain.CastMethod  `json:"method"`
	Lines    []domain.Line      `json:"lines"`
}

func (s *ReadingService) Create(ctx context.Context, in CreateReadingInput) (domain.Reading, error) {
	if in.Question == "" {
		return domain.Reading{}, errors.New("question is required")
	}
	if len(in.Lines) != 6 {
		return domain.Reading{}, errors.New("exactly 6 lines are required")
	}

	primary := computePrimaryHexagram(in.Lines)
	relating := computeRelatingHexagram(in.Lines)
	changing := changingLines(in.Lines)

	interp := "MVP interpretace: doplň podle datasetu hexagramů a měnících se čar."
	if h, ok := domain.SampleHexagrams[primary]; ok {
		interp = h.Judgment
	}

	r := domain.Reading{
		ID:             newID(),
		Question:       in.Question,
		Method:         in.Method,
		Lines:          in.Lines,
		PrimaryNumber:  primary,
		RelatingNumber: relating,
		ChangingLines:  changing,
		Interpretation: interp,
		CreatedAt:      time.Now().UTC(),
	}

	return s.repo.Save(ctx, r)
}

func (s *ReadingService) List(ctx context.Context) ([]domain.Reading, error) {
	return s.repo.List(ctx)
}

func (s *ReadingService) Get(ctx context.Context, id string) (domain.Reading, error) {
	return s.repo.Get(ctx, id)
}

func changingLines(lines []domain.Line) []int {
	out := make([]int, 0)
	for _, l := range lines {
		if l.IsChanging() {
			out = append(out, l.Position)
		}
	}
	return out
}

func computePrimaryHexagram(lines []domain.Line) int {
	if allYang(lines) {
		return 1
	}
	if allYin(lines) {
		return 2
	}
	return 0
}

func computeRelatingHexagram(lines []domain.Line) int {
	flipped := make([]domain.Line, len(lines))
	copy(flipped, lines)
	for i := range flipped {
		if flipped[i].Value == domain.OldYin {
			flipped[i].Value = domain.YoungYang
		}
		if flipped[i].Value == domain.OldYang {
			flipped[i].Value = domain.YoungYin
		}
	}
	return computePrimaryHexagram(flipped)
}

func allYang(lines []domain.Line) bool {
	for _, l := range lines {
		if !l.IsYang() {
			return false
		}
	}
	return true
}

func allYin(lines []domain.Line) bool {
	for _, l := range lines {
		if l.IsYang() {
			return false
		}
	}
	return true
}

func newID() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
