package main

type Character struct {
	Faces []string `toml:"faces"`
}

func (c *Character) NextFace(currentIndex int) (string, int) {
	if len(c.Faces) == 0 {
		return "(o_o)", 0
	}
	nextIndex := (currentIndex + 1) % len(c.Faces)
	return c.Faces[nextIndex], nextIndex
}

func getArms(frame string) (string, string) {
	switch frame {
	case "<":
		return "<", "<"
	case ">":
		return ">", ">"
	default:
		return frame, frame
	}
}
