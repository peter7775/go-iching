package dataembed

import _ "embed"

//go:embed hexagrams.en.json
var HexagramsEN []byte

//go:embed hexagrams.cs.json
var HexagramsCS []byte
