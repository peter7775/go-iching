package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Hexagram struct {
	Number    int               `json:"number"`
	Name      string            `json:"name"`
	Judgment  string            `json:"judgment"`
	Image     string            `json:"image"`
	LineTexts map[string]string `json:"line_texts"`
}

func main() {
	f, err := os.Open("data/hexagrams.sample.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	var items []Hexagram
	if err := json.NewDecoder(f).Decode(&items); err != nil {
		panic(err)
	}
	fmt.Printf("loaded %d hexagrams from sample dataset\n", len(items))
	fmt.Println("next step: insert into hexagrams + hexagram_lines tables")
}
