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

func (l Line) IsYang() bool {
	return l.Value == YoungYang || l.Value == OldYang
}

func (l Line) IsChanging() bool {
	return l.Value == OldYin || l.Value == OldYang
}

type CastMethod string

const (
	MethodCoins  CastMethod = "coins"
	MethodManual CastMethod = "manual"
)

type Reading struct {
	ID             string     `json:"id"`
	Question       string     `json:"question"`
	Method         CastMethod `json:"method"`
	Lines          []Line     `json:"lines"`
	PrimaryNumber  int        `json:"primary_number"`
	RelatingNumber int        `json:"relating_number"`
	ChangingLines  []int      `json:"changing_lines"`
	Interpretation string     `json:"interpretation,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
}

type Hexagram struct {
	Number    int            `json:"number"`
	Name      string         `json:"name"`
	Judgment  string         `json:"judgment"`
	Image     string         `json:"image"`
	LineTexts map[int]string `json:"line_texts"`
}
