package domain

import "time"

type LineValue int

const (
	OldYin    LineValue = 6
	YoungYang LineValue = 7
	YoungYin  LineValue = 8
	OldYang   LineValue = 9
)

type Line struct {
	Position int       `json:"position"`
	Value    LineValue `json:"value"`
}

func (l Line) IsYang() bool { return l.Value == YoungYang || l.Value == OldYang }
func (l Line) IsChanging() bool { return l.Value == OldYin || l.Value == OldYang }

type CastMethod string

const (
	MethodCoins  CastMethod = "coins"
	MethodManual CastMethod = "manual"
)

type InterpretationSection struct {
	Text     string `json:"text"`
	Comments string `json:"comments,omitempty"`
}

type InterpretationTrigram struct {
	Chinese    string `json:"chinese"`
	Symbolic   string `json:"symbolic"`
	Alchemical string `json:"alchemical"`
}

type InterpretationHexagram struct {
	Number      int                   `json:"number"`
	Name        string                `json:"name"`
	Title       string                `json:"title"`
	Character   string                `json:"character,omitempty"`
	Traditional string                `json:"traditional,omitempty"`
	Pinyin      string                `json:"pinyin,omitempty"`
	Above       InterpretationTrigram `json:"above"`
	Below       InterpretationTrigram `json:"below"`
	Judgment    InterpretationSection `json:"judgment"`
	Image       InterpretationSection `json:"image"`
}

type InterpretationLine struct {
	Line     int    `json:"line"`
	Text     string `json:"text"`
	Comments string `json:"comments,omitempty"`
}

type ReadingInterpretation struct {
	Language      string                 `json:"language"`
	Primary       InterpretationHexagram `json:"primary"`
	ChangingLines []InterpretationLine   `json:"changing_lines"`
	Relating      InterpretationHexagram `json:"relating"`
	Summary       string                 `json:"summary"`
	Markdown      string                 `json:"markdown,omitempty"`
}

type Reading struct {
	ID             string                `json:"id"`
	Question       string                `json:"question"`
	Method         CastMethod            `json:"method"`
	Language       string                `json:"language"`
	Lines          []Line                `json:"lines"`
	PrimaryNumber  int                   `json:"primary_number"`
	RelatingNumber int                   `json:"relating_number"`
	ChangingLines  []int                 `json:"changing_lines"`
	Interpretation ReadingInterpretation `json:"interpretation"`
	CreatedAt      time.Time             `json:"created_at"`
}
