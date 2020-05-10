package bcc

type Subtitle struct {
	From     float32 `json:"from"`
	To       float32 `json:"to"`
	Location int     `json:"location"`
	Content  string  `json:"content"`
}

type BCC struct {
	FontSize        float32     `json:"font_size"`
	FontColor       string      `json:"font_color"`
	BackgroundAlpha float32     `json:"background_alpha"`
	BackgroundColor string      `json:"background_color"`
	Stroke          string      `json:"stroke"`
	Body            []*Subtitle `json:"body"`
}

type SortSubtitle []*Subtitle

func (s SortSubtitle) Less(i, j int) bool {
	if s[i].From < s[j].From {
		return true
	}
	return false
}

func (s SortSubtitle) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s SortSubtitle) Len() int {
	return len(s)
}
