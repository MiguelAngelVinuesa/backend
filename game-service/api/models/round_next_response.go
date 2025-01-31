// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

// RoundNextResponse Next result response.
//
// Response for successfully retrieving a specific game result from a stored slot-machine round.
//
// swagger:model RoundNextResponse
type RoundNextResponse struct {

	// Details of the requested result (if it is not a spin).
	Data interface{} `json:"data,omitempty"`

	// round data
	// Required: true
	RoundData *RoundNextResponseRoundData `json:"roundData"`

	// Details of the requested result (if it is a spin).
	SpinData interface{} `json:"spinData,omitempty"`

	// Indicates if the request was successful.
	// Example: true
	// Required: true
	Success bool `json:"success"`
}
type RoundNextResponseRoundData struct {

	// Balance in cents after the spin.
	// Example: 1000
	// Required: true
	BalanceAfter int64 `json:"balanceAfter"`

	// Balance in cents before the spin.
	// Example: 1000
	// Required: true
	BalanceBefore int64 `json:"balanceBefore"`

	// Original bet amount in cents.
	// Example: 1000
	Bet int64 `json:"bet,omitempty"`

	// Total win amount in cents from start of free spins.
	// Example: 500
	BonusWin int64 `json:"bonusWin,omitempty"`

	// Round result data kind.
	// Example: 1
	DataKind int32 `json:"dataKind,omitempty"`

	// Indicates that this spin result hit the max payout limit.
	// Example: 1
	MaxPayout int8 `json:"maxPayout,omitempty"`

	// Total win amount in cents from start of round.
	// Example: 11500
	// Required: true
	TotalWin int64 `json:"totalWin"`

	// Win amount in cents for this spin.
	// Example: 1000
	// Required: true
	Win int64 `json:"win"`
}
