package srt

type Srt struct {
	Sequence int
	Start    string
	End      string
	Subtitle string
}

type SrtSort []*Srt

func (s SrtSort) Less(i, j int) bool {
	if s[i].Sequence < s[j].Sequence {
		return true
	}
	return false
}

func (s SrtSort) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s SrtSort) Len() int {
	return len(s)
}
