package graph

import (
	"time"

	"github.com/fogleman/gg"
	"golang.org/x/text/language"
	"golang.org/x/text/message"

	curData "git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/i18n/currency/data"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/series/math"
	series "git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/series/time"
)

// TimeSeriesCents is a time series graph generator where the value is expressed in currency cents.
type TimeSeriesCents struct {
	locale       language.Tag
	symbol       string
	tz           *time.Location
	formatter    *message.Printer
	fmtX         string
	fmtY         string
	appY         bool
	color        [3]int
	names        []string
	withLabels   bool
	withGrid     bool
	series       []series.IVs[int64]
	width        int
	height       int
	fontSize     float64
	leftMargin   float64
	rightMargin  float64
	topMargin    float64
	bottomMargin float64
	textPadding  float64
	lineWidth    float64
	dotSize      float64
	minTime      time.Time
	maxTime      time.Time
	stepsTime    int64
	scaleTime    []time.Time
	minVal       int64
	maxVal       int64
	scaleVal     []int64
	heightTime   float64
	widthVal     float64
	plotWidth    float64
	plotHeight   float64
	startX       float64
	startY       float64
	diffTime     time.Duration
	diffVal      int64
}

// NewTimeSeriesCents instantiates a new time series graph generator for currency cents.
func NewTimeSeriesCents(series ...series.IVs[int64]) *TimeSeriesCents {
	tsc := &TimeSeriesCents{
		locale:       language.English,
		symbol:       "€",
		tz:           time.UTC,
		formatter:    message.NewPrinter(language.English),
		fmtY:         "%s %.2f",
		appY:         false,
		color:        lines[0],
		withLabels:   true,
		withGrid:     true,
		series:       series,
		width:        graphWidth,
		height:       graphHeight,
		fontSize:     fontSize,
		leftMargin:   leftMargin,
		rightMargin:  rightMargin,
		topMargin:    topMargin,
		bottomMargin: bottomMargin,
		textPadding:  textPadding,
		lineWidth:    seriesLineWidth,
		dotSize:      seriesDotSize,
	}

	tsc.initX()
	tsc.initY()

	return tsc
}

// WithLocale sets the locale for the graph.
func (tsc *TimeSeriesCents) WithLocale(locale language.Tag) *TimeSeriesCents {
	tsc.locale = locale
	tsc.formatter = message.NewPrinter(locale)
	if fmt := curData.CurrencyFormats[locale]; fmt != "" {
		tsc.fmtY = fmt
		tsc.appY = curData.CurrencyAppend[locale]
	}
	return tsc
}

// WithCurrency sets the currency for the graph.
func (tsc *TimeSeriesCents) WithCurrency(currency string) *TimeSeriesCents {
	tsc.symbol = currency
	if c := curData.CurrencyFromCode(currency); c != nil {
		tsc.symbol = c.Symbol()
	}
	return tsc
}

// WithTimeZone sets the time zone for the graph.
func (tsc *TimeSeriesCents) WithTimeZone(tz *time.Location) *TimeSeriesCents {
	tsc.tz = tz
	return tsc
}

// WithColor sets the color for the first time series.
func (tsc *TimeSeriesCents) WithColor(color [3]int) *TimeSeriesCents {
	tsc.color = color
	return tsc
}

// WithNames sets the names for the time series.
func (tsc *TimeSeriesCents) WithNames(names ...string) *TimeSeriesCents {
	tsc.names = names
	return tsc
}

// WithSize sets the size for the output image.
func (tsc *TimeSeriesCents) WithSize(width, height int) *TimeSeriesCents {
	tsc.width = width
	tsc.height = height
	tsc.fontSize = fontSize * float64(width) / float64(graphWidth)
	tsc.leftMargin = leftMargin * float64(width) / float64(graphWidth)
	tsc.rightMargin = rightMargin * float64(width) / float64(graphWidth)
	tsc.topMargin = topMargin * float64(height) / float64(graphHeight)
	tsc.bottomMargin = bottomMargin * float64(height) / float64(graphHeight)
	tsc.textPadding = textPadding * float64(width) / float64(graphWidth)
	tsc.lineWidth = seriesLineWidth * float64(height) / float64(graphHeight)
	tsc.dotSize = seriesDotSize * float64(height) / float64(graphHeight)

	if tsc.lineWidth < 0.8 {
		tsc.lineWidth = 0.8
	}
	if tsc.dotSize < 1.2 {
		tsc.dotSize = 1.2
	}

	return tsc
}

// WithoutLabels hides the labels on the output image.
func (tsc *TimeSeriesCents) WithoutLabels() *TimeSeriesCents {
	tsc.withLabels = false
	return tsc
}

// WithoutGrid hides the grid on the output image.
func (tsc *TimeSeriesCents) WithoutGrid() *TimeSeriesCents {
	tsc.withGrid = false
	return tsc
}

// ThumbNail sets up the time-series to generate a thumbnail image (no grid/labels, small size)
func (tsc *TimeSeriesCents) ThumbNail() *TimeSeriesCents {
	tsc.WithSize(480, 270).WithoutGrid().WithoutLabels()
	tsc.lineWidth *= 2
	return tsc
}

// Graph generates the output image.
func (tsc *TimeSeriesCents) Graph() *gg.Context {
	dc := gg.NewContext(tsc.width, tsc.height)
	if err := dc.LoadFontFace(font, tsc.fontSize); err != nil {
		panic(err)
	}

	tsc.initTextHW(dc)

	// background
	dc.SetRGB255(background[0], background[1], background[2])
	dc.Clear()

	tsc.horizontalGrid(dc)
	tsc.verticalGrid(dc)

	for ix := range tsc.series {
		tsc.plotSeries(dc, ix, tsc.series[ix])
	}

	return dc
}

func (tsc *TimeSeriesCents) initX() {
	for ix := range tsc.series {
		if ix == 0 {
			tsc.minTime, tsc.maxTime = tsc.series[ix].MinMaxTime()
		} else {
			minTime, maxTime := tsc.series[ix].MinMaxTime()
			if minTime.Before(tsc.minTime) {
				tsc.minTime = minTime
			}
			if maxTime.After(tsc.maxTime) {
				tsc.maxTime = maxTime
			}
		}
	}

	if tsc.minTime == tsc.maxTime {
		tsc.minTime = tsc.minTime.Add(-time.Second)
		tsc.maxTime = tsc.maxTime.Add(time.Second)
	}
	diff := tsc.maxTime.Sub(tsc.minTime)

	switch {
	case diff/(24*time.Hour) > 0:
		tsc.fmtX = "01/02 15:04"
		tsc.stepsTime = 12
	case diff/time.Minute > 12:
		tsc.fmtX = "15:04"
		tsc.stepsTime = 24
	default:
		tsc.fmtX = "15:04:05"
		tsc.stepsTime = 18
	}

	tsc.scaleTime = math.AutoScaleValuesUTC(tsc.minTime.UTC(), tsc.maxTime.UTC(), tsc.stepsTime)
	tsc.minTime = tsc.scaleTime[0]
	tsc.maxTime = tsc.scaleTime[len(tsc.scaleTime)-1]
	tsc.diffTime = tsc.maxTime.Sub(tsc.minTime)
}

func (tsc *TimeSeriesCents) initY() {
	for ix := range tsc.series {
		if ix == 0 {
			tsc.minVal, tsc.maxVal = tsc.series[ix].MinMaxValue()
		} else {
			minVal, maxVal := tsc.series[ix].MinMaxValue()
			if minVal < tsc.minVal {
				tsc.minVal = minVal
			}
			if maxVal > tsc.maxVal {
				tsc.maxVal = maxVal
			}
		}
	}

	tsc.scaleVal = math.AutoScaleValuesInt64(tsc.minVal, tsc.maxVal, 20)
	tsc.minVal = tsc.scaleVal[0]
	tsc.maxVal = tsc.scaleVal[len(tsc.scaleVal)-1]

	if tsc.minVal == tsc.maxVal {
		tsc.minVal -= 100
		tsc.maxVal += 100
	}

	tsc.diffVal = tsc.maxVal - tsc.minVal
}

func (tsc *TimeSeriesCents) labelX(t time.Time) string {
	return t.In(tsc.tz).Format(tsc.fmtX)
}

func (tsc *TimeSeriesCents) labelY(a int64) string {
	amt := float64(a) / 100.0
	if tsc.appY {
		return tsc.formatter.Sprintf(tsc.fmtY, amt, tsc.symbol)
	}
	return tsc.formatter.Sprintf(tsc.fmtY, tsc.symbol, amt)
}

func (tsc *TimeSeriesCents) initTextHW(dc *gg.Context) {
	for ix := range tsc.scaleTime {
		_, h := dc.MeasureString(tsc.labelX(tsc.scaleTime[ix]))
		if h > tsc.heightTime {
			tsc.heightTime = h
		}
	}

	for ix := range tsc.scaleVal {
		w, _ := dc.MeasureString(tsc.labelY(tsc.scaleVal[ix]))
		if w > tsc.widthVal {
			tsc.widthVal = w
		}
	}
}

func (tsc *TimeSeriesCents) horizontalGrid(dc *gg.Context) {
	x1 := tsc.leftMargin
	x2 := float64(tsc.width) - tsc.rightMargin
	y1 := tsc.topMargin
	y2 := float64(tsc.height) - tsc.bottomMargin

	if tsc.withLabels {
		x1 += 2*tsc.textPadding + tsc.widthVal
		x2 -= 2*tsc.textPadding + tsc.widthVal
		y2 -= 2*tsc.textPadding + tsc.heightTime
	}

	stepSize := (y2 - y1) / float64(len(tsc.scaleVal))
	y3 := y1 + stepSize/2
	y4 := y2 - stepSize/2

	if tsc.withGrid {
		// horizontal grid (Y scale)
		dc.SetRGB255(gridline[0], gridline[1], gridline[2])
		dc.SetLineWidth(tsc.lineWidth)
		for y := y3; y < y2; y += stepSize {
			dc.DrawLine(x1, y, x2, y)
		}
		dc.Stroke()
	}

	if tsc.withLabels {
		// Y scale labels
		dc.SetRGB255(scale[0], scale[1], scale[2])
		for ix := range tsc.scaleVal {
			y := y4 - float64(ix)*stepSize
			s := tsc.labelY(tsc.scaleVal[ix])
			dc.DrawStringAnchored(s, x1-tsc.textPadding, y, 1, 0.5)
			dc.DrawStringAnchored(s, x2+tsc.textPadding, y, 0, 0.5)
		}
	}

	tsc.startY = y4
	tsc.plotHeight = y4 - y3
}

func (tsc *TimeSeriesCents) verticalGrid(dc *gg.Context) {
	x1 := tsc.leftMargin
	x2 := float64(tsc.width) - tsc.rightMargin
	y1 := tsc.topMargin
	y2 := float64(tsc.height) - tsc.bottomMargin

	if tsc.withLabels {
		x1 += 2*tsc.textPadding + tsc.widthVal
		x2 -= 2*tsc.textPadding + tsc.widthVal
		y2 -= 2*tsc.textPadding + tsc.heightTime
	}

	stepSize := (x2 - x1) / float64(len(tsc.scaleTime))
	x3 := x1 + stepSize/2
	x4 := x2 - stepSize/2

	if tsc.withGrid {
		// vertical grid (X scale)
		dc.SetRGB255(gridline[0], gridline[1], gridline[2])
		dc.SetLineWidth(tsc.lineWidth)
		for x := x3; x < x2; x += stepSize {
			dc.DrawLine(x, y1, x, y2)
		}
		dc.Stroke()
	}

	if tsc.withLabels {
		// X scale labels
		dc.SetRGB255(scale[0], scale[1], scale[2])
		for ix := range tsc.scaleTime {
			x := x3 + float64(ix)*stepSize
			dc.DrawStringAnchored(tsc.labelX(tsc.scaleTime[ix]), x, y2+tsc.textPadding, 0.5, 1)
		}
	}

	tsc.startX = x3
	tsc.plotWidth = x4 - x3
}

func (tsc *TimeSeriesCents) plotSeries(dc *gg.Context, seq int, series series.IVs[int64]) {
	var lineC [3]int
	if seq == 0 {
		lineC = tsc.color
	} else {
		lineC = lines[seq%len(lines)]
	}

	// plot time-series points
	dc.SetRGB255(lineC[0], lineC[1], lineC[2])
	for ix := range series {
		ts := series[ix]
		tt, tv := ts.Time(), ts.Value()
		x := tsc.startX + (float64(tt.Sub(tsc.minTime)) / float64(tsc.diffTime) * tsc.plotWidth)
		y := tsc.startY - (float64(tv-tsc.minVal) / float64(tsc.diffVal) * tsc.plotHeight)
		dc.DrawCircle(x, y, tsc.dotSize)
		dc.Fill()
	}

	// plot time-series line
	dc.SetLineWidth(tsc.lineWidth)
	for ix := range series {
		ts := series[ix]
		tt, tv := ts.Time(), ts.Value()
		x := tsc.startX + (float64(tt.Sub(tsc.minTime)) / float64(tsc.diffTime) * tsc.plotWidth)
		y := tsc.startY - (float64(tv-tsc.minVal) / float64(tsc.diffVal) * tsc.plotHeight)
		if ix == 0 {
			dc.MoveTo(x, y)
		} else {
			dc.LineTo(x, y)
		}
	}
	dc.Stroke()
}
