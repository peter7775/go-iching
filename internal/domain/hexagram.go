package domain

var SampleHexagrams = map[int]Hexagram{
	1: {
		Number: 1,
		Name: "The Creative",
		Judgment: "The Creative works sublime success, furthering through perseverance.",
		Image: "The movement of heaven is full of power.",
		LineTexts: map[int]string{
			1: "Hidden dragon. Do not act.",
			2: "Dragon appearing in the field.",
			3: "All day long the superior man is creatively active.",
			4: "Wavering flight over the depths.",
			5: "Flying dragon in the heavens.",
			6: "Arrogant dragon will have cause to repent.",
		},
	},
	2: {
		Number: 2,
		Name: "The Receptive",
		Judgment: "The Receptive brings about sublime success, furthering through the perseverance of a mare.",
		Image: "The earth's condition is receptive devotion.",
		LineTexts: map[int]string{
			1: "When there is hoarfrost underfoot, solid ice is not far off.",
			2: "Straight, square, great.",
			3: "Hidden lines. One is able to remain persevering.",
			4: "A tied-up sack. No blame, no praise.",
			5: "A yellow lower garment brings supreme good fortune.",
			6: "Dragons fight in the meadow.",
		},
	},
}
