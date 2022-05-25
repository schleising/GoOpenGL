package screen

type Screen struct {
	Width  int
	Height int
}

func (s *Screen) SetScreenSize(width, height int) {
	s.Width = width
	s.Height = height
}
