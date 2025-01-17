package slots

import (
	"fmt"
	"strings"
	"time"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/utils/conv"
)

type Bracket struct {
	Keys   []string   `json:"keys"`   // list of keys
	Values [][]string `json:"values"` // list of data rows.
}

type SimKV struct {
	Key    string   `json:"key"`    // key of a data row.
	Values []string `json:"values"` // one or more values of a data row.
}

type SimTable struct {
	Keys   []string   `json:"keys"`   // list of keys
	Values [][]string `json:"values"` // list of data rows.
}

type SimRound struct {
	TotalFactor float64 `json:"totalFactor,omitempty"` // total win factor.
	Results     any     `json:"results,omitempty"`     // slice of results.
}

type SimBlock struct {
	GameCode    string       `json:"game"`                  // unique code of the game.
	RTP         uint8        `json:"rtp"`                   // target RTP of the simulation.
	Choices     string       `json:"choices,omitempty"`     // players choices (e.g. BB, NORTH/SOUTH, etc)
	DataType    SimDataType  `json:"dataType"`              // type of sim data.
	DataSubType string       `json:"dataSubType,omitempty"` // sub-type of sim data.
	BlockType   SimBlockType `json:"blockType"`             // type of sim block.
	ColCount    uint8        `json:"colCount"`              // number of columns in sim block.
	BlockKV     SimKVs       `json:"blockKV,omitempty"`     // sim key-values block.
	BlockTable  SimTable     `json:"blockTable,omitempty"`  // sim table block.
	BestFree    SimRounds    `json:"bestFree,omitempty"`    // sim block best rounds with free spins.
	BestNoFree  SimRounds    `json:"bestNoFree,omitempty"`  // sim block best rounds without free spins.
	Created     *time.Time   `json:"created,omitempty"`     // timestamp the simulation was created.
	CreatedBy   string       `json:"createdBy,omitempty"`   // link to user who created the simulation.
}

type SimBlockType uint8

// The sequence must not be changed as DynamoDB data depends on it!
// Always add new types at the end.
const (
	BlockKV SimBlockType = iota + 1
	BlockTable
	BlockBestFree
	BlockBestNoFree
)

type SimDataType uint8

// The sequence must not be changed as DynamoDB objects depend on it!
// Always add new types at the end.
const (
	HeaderData SimDataType = iota + 1
	StatsData
	BetTotalsData
	BetSpreadData
	WinSpreadData
	RoundSpinsData
	MultiplierMarksData
	MultipliersData
	RoundSymbolsData
	SymbolsData
	PayoutsData
	PaylinesData
	ScatterPayoutsData
	BonusWheelData
	InstantBonusData
	PlayerChoiceData
	ActionsData
	FlagsData
	ScriptsData
	PlayersData
	RoundsFreeData
	RoundsNoFreeData

	SimDataTypeMin = HeaderData
	SimDataTypeMax = RoundsNoFreeData
)

func (d SimDataType) String() string {
	switch d {
	case HeaderData:
		return "header"
	case StatsData:
		return "stats"
	case BetTotalsData:
		return "bet-totals"
	case BetSpreadData:
		return "bet-spread"
	case WinSpreadData:
		return "win-spread"
	case RoundSpinsData:
		return "round-spins"
	case MultiplierMarksData:
		return "multiplier-marks"
	case MultipliersData:
		return "multipliers"
	case RoundSymbolsData:
		return "round-symbols"
	case SymbolsData:
		return "symbols"
	case PayoutsData:
		return "payouts"
	case PaylinesData:
		return "paylines"
	case ScatterPayoutsData:
		return "scatter-payouts"
	case BonusWheelData:
		return "bonus-wheel"
	case InstantBonusData:
		return "instant-bonus"
	case PlayerChoiceData:
		return "player-choice"
	case ActionsData:
		return "actions"
	case FlagsData:
		return "flags"
	case ScriptsData:
		return "scripts"
	case PlayersData:
		return "players"
	case RoundsFreeData:
		return "rounds-free"
	case RoundsNoFreeData:
		return "rounds-no-free"
	default:
		return "???"
	}
}

func SimDataTypeFromString(s string) SimDataType {
	return simDataTypes[s]
}

var simDataTypes = make(map[string]SimDataType, SimDataTypeMax)

func init() {
	for t := SimDataTypeMin; t <= SimDataTypeMax; t++ {
		simDataTypes[t.String()] = t
	}
}

func (s *SimBlock) DataTypeStr() string {
	if s.DataSubType == "" {
		return s.DataType.String()
	}
	return fmt.Sprintf("%s:%s", s.DataType.String(), s.DataSubType)
}

type SimRounds []SimRound

type SimKVs []SimKV

type Reporter interface {
	Printf(string, ...any)
	Println(...any)
}

type MyLogger struct {
	report strings.Builder
}

func (l *MyLogger) Printf(f string, params ...any) {
	l.report.WriteString(fmt.Sprintf(f, params...))
	l.report.WriteRune('\n')
}

func (l *MyLogger) Println(params ...any) {
	for ix := range params {
		l.report.WriteString(conv.StringFromAny(params[ix]))
		l.report.WriteRune('\n')
	}
}

func (l *MyLogger) Report() string {
	return l.report.String()
}
