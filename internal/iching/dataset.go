package iching

import (
	"encoding/json"
	"fmt"

	"github.com/example/iching-fiber-app/dataembed"
)

type WilhelmText struct {
	Text     string `json:"text"`
	Comments string `json:"comments"`
}

type WilhelmLine struct {
	Text     string `json:"text"`
	Comments string `json:"comments"`
}

type WilhelmTrigram struct {
	Chinese    string `json:"chinese"`
	Symbolic   string `json:"symbolic"`
	Alchemical string `json:"alchemical"`
}

type WilhelmHexagram struct {
	Number      int                    `json:"number"`
	Name        string                 `json:"name"`
	Title       string                 `json:"title"`
	Character   string                 `json:"character"`
	Traditional string                 `json:"traditional"`
	Pinyin      string                 `json:"pinyin"`
	Above       WilhelmTrigram         `json:"above"`
	Below       WilhelmTrigram         `json:"below"`
	Judgment    WilhelmText            `json:"judgment"`
	Image       WilhelmText            `json:"image"`
	Lines       map[string]WilhelmLine `json:"lines"`
}

type CzechHexagram struct {
	Number      int    `json:"number"`
	Character   string `json:"character"`
	Traditional string `json:"traditional"`
	Pinyin      string `json:"pinyin"`
	NameCS      string `json:"name_cs"`
	TitleCS     string `json:"title_cs"`
	Above       struct {
		Chinese      string `json:"chinese"`
		SymbolicCS   string `json:"symbolic_cs"`
		AlchemicalCS string `json:"alchemical_cs"`
	} `json:"above"`
	Below struct {
		Chinese      string `json:"chinese"`
		SymbolicCS   string `json:"symbolic_cs"`
		AlchemicalCS string `json:"alchemical_cs"`
	} `json:"below"`
	JudgmentCS WilhelmText            `json:"judgment_cs"`
	ImageCS    WilhelmText            `json:"image_cs"`
	LinesCS    map[string]WilhelmLine `json:"lines_cs"`
}

func LoadEnglish() (map[int]WilhelmHexagram, error) {
	var items []WilhelmHexagram
	if err := json.Unmarshal(dataembed.HexagramsEN, &items); err != nil {
		return nil, fmt.Errorf("unmarshal english dataset: %w", err)
	}
	out := make(map[int]WilhelmHexagram, len(items))
	for _, item := range items {
		out[item.Number] = item
	}
	return out, nil
}

func LoadCzech() (map[int]CzechHexagram, error) {
	var items []CzechHexagram
	if err := json.Unmarshal(dataembed.HexagramsCS, &items); err != nil {
		return nil, fmt.Errorf("unmarshal czech dataset: %w", err)
	}
	out := make(map[int]CzechHexagram, len(items))
	for _, item := range items {
		out[item.Number] = item
	}
	return out, nil
}
