package domain

func CloneLines(in []Line) []Line {
	out := make([]Line, len(in))
	copy(out, in)
	return out
}
