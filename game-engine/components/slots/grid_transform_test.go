package slots

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

func TestFlipHorizontal(t *testing.T) {
	testCases := []struct {
		name    string
		reels   int
		rows    int
		indexes utils.Indexes
		want    utils.Indexes
	}{
		{
			name:    "3x3",
			reels:   3,
			rows:    3,
			indexes: utils.Indexes{4, 5, 6, 7, 8, 9, 10, 11, 12},
			want:    utils.Indexes{6, 5, 4, 9, 8, 7, 12, 11, 10},
		},
		{
			name:    "3x6",
			reels:   3,
			rows:    6,
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18},
			want:    utils.Indexes{6, 5, 4, 3, 2, 1, 12, 11, 10, 9, 8, 7, 18, 17, 16, 15, 14, 13},
		},
		{
			name:    "5x3",
			reels:   5,
			rows:    3,
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
			want:    utils.Indexes{3, 2, 1, 6, 5, 4, 9, 8, 7, 12, 11, 10, 15, 14, 13},
		},
		{
			name:    "5x4",
			reels:   5,
			rows:    4,
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
			want:    utils.Indexes{4, 3, 2, 1, 8, 7, 6, 5, 12, 11, 10, 9, 16, 15, 14, 13, 20, 19, 18, 17},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := make(utils.Indexes, len(tc.indexes))
			copy(got, tc.indexes)
			FlipHorizontal(tc.reels, tc.rows, got)
			assert.EqualValues(t, tc.want, got)
		})
	}
}

func TestFlipVertical(t *testing.T) {
	testCases := []struct {
		name    string
		reels   int
		rows    int
		indexes utils.Indexes
		want    utils.Indexes
	}{
		{
			name:    "3x3",
			reels:   3,
			rows:    3,
			indexes: utils.Indexes{4, 5, 6, 7, 8, 9, 10, 11, 12},
			want:    utils.Indexes{10, 11, 12, 7, 8, 9, 4, 5, 6},
		},
		{
			name:    "4x4",
			reels:   4,
			rows:    4,
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
			want:    utils.Indexes{13, 14, 15, 16, 9, 10, 11, 12, 5, 6, 7, 8, 1, 2, 3, 4},
		},
		{
			name:    "5x3",
			reels:   5,
			rows:    3,
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
			want:    utils.Indexes{13, 14, 15, 10, 11, 12, 7, 8, 9, 4, 5, 6, 1, 2, 3},
		},
		{
			name:    "5x4",
			reels:   5,
			rows:    4,
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
			want:    utils.Indexes{17, 18, 19, 20, 13, 14, 15, 16, 9, 10, 11, 12, 5, 6, 7, 8, 1, 2, 3, 4},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := make(utils.Indexes, len(tc.indexes))
			copy(got, tc.indexes)
			FlipVertical(tc.reels, tc.rows, got)
			assert.EqualValues(t, tc.want, got)
		})
	}
}

func TestRotate(t *testing.T) {
	testCases := []struct {
		name    string
		reels   int
		rows    int
		indexes utils.Indexes
		want    utils.Indexes
	}{
		{
			name:    "3x3",
			reels:   3,
			rows:    3,
			indexes: utils.Indexes{4, 5, 6, 7, 8, 9, 10, 11, 12},
			want:    utils.Indexes{12, 11, 10, 9, 8, 7, 6, 5, 4},
		},
		{
			name:    "3x6",
			reels:   3,
			rows:    6,
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18},
			want:    utils.Indexes{18, 17, 16, 15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1},
		},
		{
			name:    "5x3",
			reels:   5,
			rows:    3,
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
			want:    utils.Indexes{15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1},
		},
		{
			name:    "5x4",
			reels:   5,
			rows:    4,
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
			want:    utils.Indexes{20, 19, 18, 17, 16, 15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := make(utils.Indexes, len(tc.indexes))
			copy(got, tc.indexes)
			Rotate(tc.reels, tc.rows, got)
			assert.EqualValues(t, tc.want, got)
		})
	}
}

func TestShift(t *testing.T) {
	testCases := []struct {
		name    string
		reels   int
		rows    int
		count   int
		indexes utils.Indexes
		want    utils.Indexes
	}{
		{
			name:    "3x3x0",
			reels:   3,
			rows:    3,
			indexes: utils.Indexes{4, 5, 6, 7, 8, 9, 10, 11, 12},
			want:    utils.Indexes{4, 5, 6, 7, 8, 9, 10, 11, 12},
		},
		{
			name:    "3x3x1",
			reels:   3,
			rows:    3,
			count:   1,
			indexes: utils.Indexes{4, 5, 6, 7, 8, 9, 10, 11, 12},
			want:    utils.Indexes{10, 11, 12, 4, 5, 6, 7, 8, 9},
		},
		{
			name:    "3x3x2",
			reels:   3,
			rows:    3,
			count:   2,
			indexes: utils.Indexes{4, 5, 6, 7, 8, 9, 10, 11, 12},
			want:    utils.Indexes{7, 8, 9, 10, 11, 12, 4, 5, 6},
		},
		{
			name:    "3x3x3",
			reels:   3,
			rows:    3,
			count:   3,
			indexes: utils.Indexes{4, 5, 6, 7, 8, 9, 10, 11, 12},
			want:    utils.Indexes{4, 5, 6, 7, 8, 9, 10, 11, 12},
		},
		{
			name:    "3x6x0",
			reels:   3,
			rows:    6,
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18},
			want:    utils.Indexes{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18},
		},
		{
			name:    "3x6x1",
			reels:   3,
			rows:    6,
			count:   1,
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18},
			want:    utils.Indexes{13, 14, 15, 16, 17, 18, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
		},
		{
			name:    "3x6x2",
			reels:   3,
			rows:    6,
			count:   2,
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18},
			want:    utils.Indexes{7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 1, 2, 3, 4, 5, 6},
		},
		{
			name:    "3x6x3",
			reels:   3,
			rows:    6,
			count:   3,
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18},
			want:    utils.Indexes{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18},
		},
		{
			name:    "5x3x0",
			reels:   5,
			rows:    3,
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
			want:    utils.Indexes{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
		},
		{
			name:    "5x3x1",
			reels:   5,
			rows:    3,
			count:   1,
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
			want:    utils.Indexes{13, 14, 15, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
		},
		{
			name:    "5x3x2",
			reels:   5,
			rows:    3,
			count:   2,
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
			want:    utils.Indexes{10, 11, 12, 13, 14, 15, 1, 2, 3, 4, 5, 6, 7, 8, 9},
		},
		{
			name:    "5x3x3",
			reels:   5,
			rows:    3,
			count:   3,
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
			want:    utils.Indexes{7, 8, 9, 10, 11, 12, 13, 14, 15, 1, 2, 3, 4, 5, 6},
		},
		{
			name:    "5x3x4",
			reels:   5,
			rows:    3,
			count:   4,
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
			want:    utils.Indexes{4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 1, 2, 3},
		},
		{
			name:    "5x3x5",
			reels:   5,
			rows:    3,
			count:   5,
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
			want:    utils.Indexes{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
		},
		{
			name:    "5x4x0",
			reels:   5,
			rows:    4,
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
			want:    utils.Indexes{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
		},
		{
			name:    "5x4x1",
			reels:   5,
			rows:    4,
			count:   1,
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
			want:    utils.Indexes{17, 18, 19, 20, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
		},
		{
			name:    "5x4x2",
			reels:   5,
			rows:    4,
			count:   2,
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
			want:    utils.Indexes{13, 14, 15, 16, 17, 18, 19, 20, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
		},
		{
			name:    "5x4x3",
			reels:   5,
			rows:    4,
			count:   3,
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
			want:    utils.Indexes{9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 1, 2, 3, 4, 5, 6, 7, 8},
		},
		{
			name:    "5x4x4",
			reels:   5,
			rows:    4,
			count:   4,
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
			want:    utils.Indexes{5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 1, 2, 3, 4},
		},
		{
			name:    "5x4x5",
			reels:   5,
			rows:    4,
			count:   5,
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
			want:    utils.Indexes{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := make(utils.Indexes, len(tc.indexes))
			copy(got, tc.indexes)
			Shift(tc.reels, tc.rows, tc.count, got)
			assert.EqualValues(t, tc.want, got)
		})
	}
}
