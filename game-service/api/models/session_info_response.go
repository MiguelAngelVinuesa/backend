// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

// SessionInfoResponse Session info response.
//
// Response for successfully retrievinf session information for replay.
//
// swagger:model SessionInfoResponse
type SessionInfoResponse struct {

	// Unique ID of the game played.
	// Example: mgd
	// Required: true
	GameID string `json:"gameId"`

	// rounds
	// Required: true
	Rounds []*SessionInfoResponseRoundsItems0 `json:"rounds"`

	// RTP of the game played.
	// Example: 96
	// Required: true
	Rtp int64 `json:"rtp"`

	// timestamp the session started (ISO-8601: yyyy-mm-ddThh:nn:ssZ)
	// Example: 2023-11-22T13:45:12Z
	// Required: true
	Started string `json:"started"`

	// Indicates if the request was successful.
	// Example: true
	// Required: true
	Success bool `json:"success"`

	// approximate total number of rounds in the session (may lag by a few seconds)
	// Example: 4810
	// Required: true
	Total int64 `json:"total"`

	// timestamp the session was last updated (ISO-8601: yyyy-mm-ddThh:nn:ssZ)
	// Example: 2023-11-22T13:48:33Z
	// Required: true
	Updated string `json:"updated"`
}
type SessionInfoResponseRoundsItems0 struct {

	// Unique ID of the round.
	// Example: 5418324dc7884ad7b7d6e0fff31e4d1a
	// Required: true
	RoundID string `json:"roundId"`

	// timestamp the round started (ISO-8601: yyyy-mm-ddThh:nn:ssZ)
	// Example: 2023-11-22T13:48:33Z
	// Required: true
	Started string `json:"started"`
}
