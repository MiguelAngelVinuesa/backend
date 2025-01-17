package metrics

import (
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/metrics"
)

const (
	ApiPing metrics.DurationType = iota
	ApiBinHashes
	ApiGameHash
	ApiStrings
	ApiPlural
	ApiPlurals
	ApiGameInfo
	ApiPreferences
	ApiCcbFlags
	ApiRound
	ApiRoundPaid
	ApiRoundSecond
	ApiRoundResume
	ApiRoundDebug
	ApiRoundDebugSecond
	ApiRoundDebugResume
	ApiRoundNext
	ApiRoundFinish
	ApiSessionInfo
	ApiRngConditions
	ApiRngMagicTest
	ApiRngMagic
	GeNewGame
	GeRound
	GeRoundResume
	GeRoundDebug
	IsStrings
	IsString
	IsPlural
	IsPlurals
	BoGamePrefs
	BoSession
	BoSessionRounds
	BoSimulatorRTP
	BoCasinoGameRTP
	DsRound
	DsRoundInit
	DsRoundComplete
	DsRoundNext
	DsRoundState
	DsSessionPut
	DsSessionGet
	DsGamePrefsPut
	DsGamePrefsGet
	DsPlayerPrefsPut
	DsPlayerPrefsGet
	MaxDuration = DsPlayerPrefsGet
)

var durationNames = []string{
	"API ping",
	"API bin-hashes",
	"API game hash",
	"API strings",
	"API plural",
	"API plurals",
	"API game-info",
	"API preferences",
	"API ccb-flags",
	"API round",
	"API round paid",
	"API round second",
	"API round resume",
	"API round debug",
	"API round debug-second",
	"API round debug-resume",
	"API round next",
	"API round finish",
	"API session",
	"API rng-conditions-lu",
	"API rng-magic test",
	"API rng-magic",
	"GE new game",
	"GE round",
	"GE round resume",
	"GE round debug",
	"IS strings",
	"IS string",
	"IS plural",
	"IS plurals",
	"BO game preferences",
	"BO session",
	"BO session rounds",
	"BO simulator rtp",
	"BO casino-game rtp",
	"DS round",
	"DS round init",
	"DS round complete",
	"DS round next",
	"DS round state",
	"DS put session",
	"DS get session",
	"DS put game_prefs",
	"DS get game-prefs",
	"DS put player-prefs",
	"DS get player-prefs",
}
