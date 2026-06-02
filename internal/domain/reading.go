package domain

import "time"

type CastMethod string

const (
	MethodManual CastMethod = "manual"
	MethodCoins  CastMethod = "coins"
)

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

func (l Line) IsChanging() bool {
	return l.Value == OldYin || l.Value == OldYang
}

func (l Line) IsYang() bool {
	return l.Value == YoungYang || l.Value == OldYang
}

type InterpretationSection struct {
	Text     string `json:"text"`
	Comments string `json:"comments,omitempty"`
}

type InterpretationTrigram struct {
	Chinese    string `json:"chinese,omitempty"`
	Symbolic   string `json:"symbolic,omitempty"`
	Alchemical string `json:"alchemical,omitempty"`
}

type InterpretationHexagram struct {
	Number      int                   `json:"number"`
	Name        string                `json:"name,omitempty"`
	Title       string                `json:"title,omitempty"`
	Character   string                `json:"character,omitempty"`
	Traditional string                `json:"traditional,omitempty"`
	Pinyin      string                `json:"pinyin,omitempty"`
	Above       InterpretationTrigram `json:"above,omitempty"`
	Below       InterpretationTrigram `json:"below,omitempty"`
	Judgment    InterpretationSection `json:"judgment,omitempty"`
	Image       InterpretationSection `json:"image,omitempty"`
}

type InterpretationLine struct {
	Line     int    `json:"line"`
	Text     string `json:"text,omitempty"`
	Comments string `json:"comments,omitempty"`
}

type ReadingInterpretation struct {
	Language      string                 `json:"language"`
	Primary       InterpretationHexagram `json:"primary"`
	ChangingLines []InterpretationLine   `json:"changing_lines,omitempty"`
	Relating      InterpretationHexagram `json:"relating"`
	Summary       string                 `json:"summary,omitempty"`
	Markdown      string                 `json:"markdown,omitempty"`
}

type Reflection struct {
	Rating     int        `json:"rating"`
	Note       string     `json:"note,omitempty"`
	CreatedAt  *time.Time `json:"created_at,omitempty"`
	Eligible   bool       `json:"eligible"`
	EligibleAt time.Time  `json:"eligible_at"`
}

type Reading struct {
	ID             string                `json:"id"`
	Question       string                `json:"question"`
	Method         CastMethod            `json:"method"`
	Language       string                `json:"language,omitempty"`
	Lines          []Line                `json:"lines"`
	PrimaryNumber  int                   `json:"primary_number"`
	RelatingNumber int                   `json:"relating_number"`
	ChangingLines  []int                 `json:"changing_lines,omitempty"`
	Interpretation ReadingInterpretation `json:"interpretation"`
	CreatedAt      time.Time             `json:"created_at"`
	Reflection     Reflection            `json:"reflection"`
}

type RelevanceStats struct {
	RatedCount    int         `json:"rated_count"`
	TotalCount    int         `json:"total_count"`
	AverageRating float64     `json:"average_rating"`
	RelevancePct  float64     `json:"relevance_pct"`
	Distribution  map[int]int `json:"distribution"`
}
