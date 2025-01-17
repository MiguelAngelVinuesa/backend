package jwx

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	data1 = "eyJraWQiOiJ3c0M2a3QrNU9Da1NUNHF0REdhVVk0REJFS2Nod3FEUEI4ODdSVFZoQWZrPSIsImFsZyI6IlJTMjU2In0.eyJzdWIiOiJjOWNjNjcxZi01MzQ0LTRlYzgtYjk3My1lMjdjZDUxNzVmZmIiLCJkZXZpY2Vfa2V5IjoiZXUtY2VudHJhbC0xXzQ4OGEzYTY2LTdiMTEtNDQ1Ny1hZDNlLTMxNTAzZTA2NDQzZiIsImlzcyI6Imh0dHBzOlwvXC9jb2duaXRvLWlkcC5ldS1jZW50cmFsLTEuYW1hem9uYXdzLmNvbVwvZXUtY2VudHJhbC0xX1J6TzY5NjN4RCIsImNsaWVudF9pZCI6IjNhN2gwMGV0dmpsbDg1YjVuaGN0NXZ0M2psIiwib3JpZ2luX2p0aSI6IjhlYmI0NTdlLTRjNDAtNGVkNC05NmE1LTRiZmYxNzQ1NmZiMyIsImV2ZW50X2lkIjoiNTE2MjA4YjYtNjRmYS00YzE5LThhMTgtNjVjMDdkNWZmNDAyIiwidG9rZW5fdXNlIjoiYWNjZXNzIiwic2NvcGUiOiJhd3MuY29nbml0by5zaWduaW4udXNlci5hZG1pbiIsImF1dGhfdGltZSI6MTY3NTg2NTgyNiwiZXhwIjoxNjc1ODY5NDI2LCJpYXQiOjE2NzU4NjU4MjYsImp0aSI6ImE5MDI4ZjE5LTAyZDEtNDQ2MS05NTIzLTVjNDg3ODgyZDU2OSIsInVzZXJuYW1lIjoiYzljYzY3MWYtNTM0NC00ZWM4LWI5NzMtZTI3Y2Q1MTc1ZmZiIn0.EFo-WOies1Yd0QJZksGP2TngzUaoOHeLmpzJil7S-ncBJBujvsHQ0LTv675NrhtdTv5RqTBmbjB_HtJfIxpLtS70TUol0ano8Q0A_Oksp_-ArQntKz0d5Yw5PyZhAA7_Zh_L0TYDjddQcOqSdKp--ivp5PZZ0KnLzPwYCq2IfAT6o8Ralig7Ds_zotDCFrBrUxQGS5MDY2Y9BnMfSelBm28NDgFS7b-MH2TK2yi3yg6psekYI4O96myK-laaTPfenbbPqkLa6uFdpas4bzYkkKIn8k4V2vYQ_pu9-_uDYpYbadT0RUpfoDPTkfYdon95DsXP38NfMRcK3Dbsk9BKjA"
	data2 = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
	data3 = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.NHVaYe26MbtOYhSKkoKYdFVomg4i8ZJd8_-RU8VNbftc4TSMb4bXP3l3YlNWACwyXPGffz5aXHc6lty1Y2t4SWRqGteragsVdZufDn5BlnJl9pdR_kdVFUsra2rWKEofkZeIC4yWytE58sMIihvo9H1ScmmVwBcQP6XETqYd0aSHp1gOa9RdUPDvoXQ5oqygTqVtxaDr6wUFKrKItgBMzWIdNZ6y7O9E0DhEPTbE9rfBo6KTFsHAZnMg4k68CDp2woYIaXbmYTWcvbzIuHO7_37GT79XdIwkm95QJ7hYC9RiwrV7mesbY4PAahERJawntho0my942XheVLmGwLMBkQ"
	data4 = "eyJraWQiOiJ3c0M2a3QrNU9Da1NUNHF0REdhVVk0REJFS2Nod3FEUEI4ODdSVFZoQWZrPSIsImFsZyI6IlJTMjU2In0.eyJzdWIiOiJjOWNjNjcxZi01MzQ0LTRlYzgtYjk3My1lMjdjZDUxNzVmZmIiLCJkZXZpY2Vfa2V5IjoiZXUtY2VudHJhbC0xXzVmNTQxMmFiLWYwNzktNDUzNi04MzEwLTc5NWZmMDY3ZTVlNiIsImlzcyI6Imh0dHBzOlwvXC9jb2duaXRvLWlkcC5ldS1jZW50cmFsLTEuYW1hem9uYXdzLmNvbVwvZXUtY2VudHJhbC0xX1J6TzY5NjN4RCIsImNsaWVudF9pZCI6IjNhN2gwMGV0dmpsbDg1YjVuaGN0NXZ0M2psIiwib3JpZ2luX2p0aSI6ImM2MzJiM2I0LWMwMzEtNGY2ZS1hY2E2LTRjMGIyMjhjMjY5OCIsImV2ZW50X2lkIjoiNDc4YjA3YWItZTkyNC00YTI5LThmZTctNGRmM2JjMzRhZGNmIiwidG9rZW5fdXNlIjoiYWNjZXNzIiwic2NvcGUiOiJhd3MuY29nbml0by5zaWduaW4udXNlci5hZG1pbiIsImF1dGhfdGltZSI6MTY3NTk1NTEwNSwiZXhwIjoxNjc1OTU4NzA1LCJpYXQiOjE2NzU5NTUxMDUsImp0aSI6IjJiOWRmZGUxLTQ1M2MtNGIyZS1iMTU5LTRjMWFmMGYwOTZhOCIsInVzZXJuYW1lIjoiYzljYzY3MWYtNTM0NC00ZWM4LWI5NzMtZTI3Y2Q1MTc1ZmZiIn0.Rjszu4KK4jXmtgBYmmQSltgX03L-9EJUPCV0dU6eA5zIrB-z8TLC7S4Fl6iqxNwL4cbCJPSGb9oKXe1wBNwUvP4Q7eqFVlXFWRm5pCbRY_z4XEiJgEKg1ZB-bHrsgacKxsVfrTRzrhPfigwz-DK70_ooNLgGbl7dqgWiE5bhdV3bqRFGn1l8NA5vxhozKxen7kTq8UL2SozqDZOMgA_wkwSi88_GtmO74USeg3RIkYlGo1Q9lTBRIJyWD7Jc8dkunKCWIrHT3KgLnHWqlOHWpj09PIgkNlIb7DOcQQlqZS_xpLG81UM-kTxALJdt-m1Z6d7NKALlXZdrfIVkaSxnhg"

	issuer = "https://cognito-idp.eu-central-1.amazonaws.com/eu-central-1_RzO6963xD"
)

func TestParseJWT(t *testing.T) {
	testCases := []struct {
		name string
		data string
		fail bool
		msg  string
	}{
		{
			name: "empty",
			fail: true,
			msg:  `invalid byte sequence`,
		},
		{
			name: "bad",
			data: "123456789",
			fail: true,
			msg:  `invalid JWT`,
		},
		{
			name: "data1 - expired",
			data: data1,
			fail: true,
			msg:  `"exp" not satisfied`,
		},
		{
			name: "data2 - invalid",
			data: data2,
			fail: true,
			msg:  `failed to find matching key`,
		},
		{
			name: "data3 - invalid",
			data: data3,
			fail: true,
			msg:  `failed to find matching key`,
		},
		{
			name: "data4 - invalid",
			data: data4,
			fail: true,
			msg:  `"exp" not satisfied`,
		},
	}

	err := RegisterJWK(cognito)
	require.Nil(t, err)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			token, err2 := ParseJWT([]byte(tc.data), cognito)
			if tc.fail {
				require.NotNil(t, err2)
				require.Nil(t, token)

				assert.True(t, strings.Contains(err2.Error(), tc.msg))
			} else {
				require.Nil(t, err2)
				require.NotNil(t, token)

				assert.Equal(t, issuer, token.Issuer())
				assert.Less(t, token.IssuedAt().UnixMicro(), time.Now().UTC().UnixMicro())
				assert.Greater(t, token.Expiration().UnixMicro(), time.Now().UTC().Add(5*time.Minute).UnixMicro())
				assert.Equal(t, "access", token.PrivateClaims()["token_use"])
			}
		})
	}
}
