package graph

const (
	// pixel measurements
	graphWidth      = 1920
	graphHeight     = 1080
	leftMargin      = 15
	rightMargin     = 15
	topMargin       = 15
	bottomMargin    = 15
	textPadding     = 10
	seriesDotSize   = 3.9
	seriesLineWidth = 1.7

	// font
	fontSize = 18
	font     = "/usr/share/fonts/truetype/freefont/FreeSans.ttf"
)

var (
	// colors
	background = [3]int{16, 16, 16}
	gridline   = [3]int{80, 80, 80}
	scale      = [3]int{240, 240, 240}
	lines      = [][3]int{
		{0xe6, 0x00, 0x49},
		{0x0b, 0xb4, 0xff},
		{0x50, 0xe9, 0x91},
		{0xe6, 0xd8, 0x00},
		{0x9b, 0x19, 0xf5},
		{0xff, 0xa3, 0x00},
		{0xdc, 0x0a, 0xb4},
		{0xb3, 0xd4, 0xff},
		{0x00, 0xbf, 0xa0},
	}
	// Dutch Field: "#e60049", "#0bb4ff", "#50e991", "#e6d800", "#9b19f5", "#ffa300", "#dc0ab4", "#b3d4ff", "#00bfa0"
	// Retro Metro: "#ea5545", "#f46a9b", "#ef9b20", "#edbf33", "#ede15b", "#bdcf32", "#87bc45", "#27aeef", "#b33dc6"
	// River Nights: "#b30000", "#7c1158", "#4421af", "#1a53ff", "#0d88e6", "#00b7c7", "#5ad45a", "#8be04e", "#ebdc78"
	// Spring Pastels: "#fd7f6f", "#7eb0d5", "#b2e061", "#bd7ebe", "#ffb55a", "#ffee65", "#beb9db", "#fdcce5", "#8bd3c7"
)
