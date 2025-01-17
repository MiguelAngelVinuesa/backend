package cognito

import (
	"time"
)

type User struct {
	CognitoID     string
	Email         string
	ZoneInfo      string
	Locale        string
	EmailVerified bool
	Enabled       bool
	Created       *time.Time
	Updated       *time.Time
}
