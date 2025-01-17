package graph

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"

	series "git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/series/time"
)

func TestGraph1(t *testing.T) {
	now := time.Now().Round(time.Second).Add(-60 * time.Minute)

	ts1 := make(series.IVs[int64], 16)
	for ix := range ts1 {
		ts1[ix] = series.NewIV[int64](
			now.Add(time.Duration(ix)*125*time.Second),
			int64(100000+rand.Intn(200000)),
		)
	}

	ts2 := make(series.IVs[int64], 32)
	for ix := range ts2 {
		ts2[ix] = series.NewIV[int64](
			now.Add(time.Duration(ix)*125*time.Second/2),
			int64(100000+rand.Intn(400000)),
		)
	}

	ts3 := make(series.IVs[int64], 48)
	for ix := range ts3 {
		ts3[ix] = series.NewIV[int64](
			now.Add(time.Duration(ix)*125*time.Second/3),
			int64(200000+rand.Intn(500000)),
		)
	}

	ts4 := make(series.IVs[int64], 96)
	for ix := range ts4 {
		ts4[ix] = series.NewIV[int64](
			now.Add(time.Duration(ix)*125*time.Second/6),
			int64(100000+rand.Intn(400000)),
		)
	}

	tsc := NewTimeSeriesCents(ts1, ts2, ts3, ts4)
	dc := tsc.Graph()
	dc.SavePNG("/tmp/graph.png")
}

func TestGraph2(t *testing.T) {
	dates := []string{"20221128", "20221129", "20221130", "20221201", "20221202"}
	for ix := range dates {
		date := dates[ix]
		ts := make(series.IVs[int64], 0, 512)

		f, err := os.Open(fmt.Sprintf("../../testdata/rounds-%s.csv", date))
		require.NoError(t, err)

		s := bufio.NewScanner(f)
		s.Split(bufio.ScanLines)

		for s.Scan() {
			fields := strings.Split(s.Text(), ",")
			require.Equal(t, 2, len(fields))

			f1 := strings.TrimSuffix(strings.TrimPrefix(fields[0], `"`), `"`)
			f2 := strings.TrimSuffix(strings.TrimPrefix(fields[1], `"`), `"`)

			c, err2 := time.Parse("2006-01-02 15:04:05.999", f1)
			v, err3 := strconv.ParseInt(f2, 10, 64)
			if err2 == nil && err3 == nil {
				ts = append(ts, series.NewIV[int64](c.UTC(), v))
			}
		}

		f.Close()
		require.NotEmpty(t, ts)

		tsc := NewTimeSeriesCents(ts).WithLocale(language.English).WithCurrency("gbp").WithTimeZone(time.Local)
		dc := tsc.Graph()
		dc.SavePNG(fmt.Sprintf("../../testdata/rounds-%s.png", date))
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
