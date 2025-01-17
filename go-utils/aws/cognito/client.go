package cognito

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/lestrrat-go/jwx/v2/jwt"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/crypto/jwx"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/tg"
)

// Cognito implements a simple interface for the Cognito SDK.
type Cognito interface {
	VerifyToken(ctx context.Context, tokenData string) (token jwt.Token, err error)
	RefreshTokens(ctx context.Context, refreshToken, fp, ip string) (idToken string, accessToken string, err error)

	ListUsers(ctx context.Context, pool string, attrs ...string) ([]*User, error)
	CreateUser(ctx context.Context, pool string, data *User) (*User, error)
	EnableUser(ctx context.Context, pool, email string) error
	DisableUser(ctx context.Context, pool, email string) error
	DeleteUser(ctx context.Context, pool, email string) error

	PasswordReset(ctx context.Context, pool, email string) error
	ForcePassword(ctx context.Context, pool, email, pswd string) error
	ResendWelcomeEmail(ctx context.Context, email string) error
}

// NewClient instantiates a new wrapper interface for communicationg with AWS Cognito.
// It implements generic functions for all TopGaming use cases.
func NewClient(cfg *aws.Config, jwksURL, clientID, analyticsEndpoint string) (Cognito, error) {
	if err := jwx.RegisterJWK(jwksURL); err != nil {
		return nil, err
	}

	return &client{
		cfg:               cfg,
		provider:          cognitoidentityprovider.NewFromConfig(*cfg),
		jwks:              jwksURL,
		clientID:          clientID,
		analyticsEndpoint: analyticsEndpoint,
	}, nil
}

type client struct {
	cfg               *aws.Config
	provider          *cognitoidentityprovider.Client
	jwks              string
	clientID          string
	analyticsEndpoint string
}

// VerifyToken implements the Cognito interface.
func (c *client) VerifyToken(_ context.Context, tokenData string) (jwt.Token, error) {
	return jwx.ParseJWT([]byte(tokenData), c.jwks)
}

// RefreshTokens implements the Cognito interface.
func (c *client) RefreshTokens(ctx context.Context, refreshToken, fp, ip string) (string, string, error) {
	req := &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow:       "REFRESH_TOKEN_AUTH",
		AuthParameters: map[string]string{"REFRESH_TOKEN": refreshToken},
		ClientId:       aws.String(c.clientID),
	}

	if fp != "" || ip != "" {
		req.UserContextData = &types.UserContextDataType{EncodedData: aws.String(fp), IpAddress: aws.String(ip)}
	}

	if c.analyticsEndpoint != "" {
		req.AnalyticsMetadata = &types.AnalyticsMetadataType{AnalyticsEndpointId: aws.String(c.analyticsEndpoint)}
	}

	resp, err := c.provider.InitiateAuth(ctx, req)
	if err != nil {
		return "", "", err
	}
	if resp == nil || resp.AuthenticationResult == nil {
		return "", "", fmt.Errorf("RefreshTokens failed: empty/invalid response")
	}

	id := aws.ToString(resp.AuthenticationResult.IdToken)
	a := aws.ToString(resp.AuthenticationResult.AccessToken)

	_, err = c.VerifyToken(ctx, id)
	if err != nil {
		return "", "", err
	}

	_, err = c.VerifyToken(ctx, a)
	if err != nil {
		return "", "", err
	}

	return id, a, nil
}

// ListUsers returns a list or user details from the given pool.
func (c *client) ListUsers(ctx context.Context, pool string, attrs ...string) ([]*User, error) {
	req := &cognitoidentityprovider.ListUsersInput{
		UserPoolId:      aws.String(pool),
		AttributesToGet: attrs,
	}

	resp, err := c.provider.ListUsers(ctx, req)
	if err != nil {
		return nil, err
	}

	out := make([]*User, 0)
	for ix := range resp.Users {
		out = append(out, userFromCognito(&resp.Users[ix]))
	}
	return out, nil
}

// CreateUser creates a new user in the given pool.
func (c *client) CreateUser(ctx context.Context, pool string, data *User) (*User, error) {
	pwd := tg.RandomPassword()
	req := &cognitoidentityprovider.AdminCreateUserInput{
		UserPoolId:             &pool,
		Username:               &data.Email,
		TemporaryPassword:      &pwd,
		DesiredDeliveryMediums: []types.DeliveryMediumType{"EMAIL"},
	}

	resp, err := c.provider.AdminCreateUser(ctx, req)
	if err != nil {
		return nil, err
	}

	if !data.Enabled {
		err = c.DisableUser(ctx, pool, data.CognitoID)
		if err != nil {
			return nil, err
		}
	}

	return userFromCognito(resp.User), nil
}

// EnableUser enables a user in the given pool.
func (c *client) EnableUser(ctx context.Context, pool, email string) error {
	req := &cognitoidentityprovider.AdminEnableUserInput{
		UserPoolId: aws.String(pool),
		Username:   aws.String(email),
	}
	_, err := c.provider.AdminEnableUser(ctx, req)
	return err
}

// DisableUser enables a user in the given pool.
func (c *client) DisableUser(ctx context.Context, pool, email string) error {
	req := &cognitoidentityprovider.AdminDisableUserInput{
		UserPoolId: aws.String(pool),
		Username:   aws.String(email),
	}
	_, err := c.provider.AdminDisableUser(ctx, req)
	return err
}

// DeleteUser removes a user in the given pool.
func (c *client) DeleteUser(ctx context.Context, pool, email string) error {
	req := &cognitoidentityprovider.AdminDeleteUserInput{
		UserPoolId: aws.String(pool),
		Username:   aws.String(email),
	}
	_, err := c.provider.AdminDeleteUser(ctx, req)
	return err
}

// PasswordReset forces a password reset for the given email in the given pool.
func (c *client) PasswordReset(ctx context.Context, pool, email string) error {
	req := &cognitoidentityprovider.AdminResetUserPasswordInput{
		UserPoolId: &pool,
		Username:   &email,
	}
	_, err := c.provider.AdminResetUserPassword(ctx, req)
	return err
}

// ForcePassword forces a new temporary password for the given email in the given pool.
func (c *client) ForcePassword(ctx context.Context, pool, email, pswd string) error {
	req := &cognitoidentityprovider.AdminSetUserPasswordInput{
		UserPoolId: &pool,
		Username:   &email,
		Password:   &pswd,
		Permanent:  false,
	}
	_, err := c.provider.AdminSetUserPassword(ctx, req)
	return err
}

// ResendWelcomeEmail resends the welcome message for the given email in the given pool.
func (c *client) ResendWelcomeEmail(ctx context.Context, email string) error {
	req := &cognitoidentityprovider.ResendConfirmationCodeInput{
		ClientId: &c.clientID,
		Username: &email,
	}
	_, err := c.provider.ResendConfirmationCode(ctx, req)
	return err
}

func userFromCognito(ut *types.UserType) *User {
	u := &User{
		CognitoID: aws.ToString(ut.Username),
		Enabled:   ut.Enabled,
		Created:   ut.UserCreateDate,
		Updated:   ut.UserLastModifiedDate,
	}

	for iy := range ut.Attributes {
		a := ut.Attributes[iy]
		switch aws.ToString(a.Name) {
		case "email":
			u.Email = aws.ToString(a.Value)
		case "zoneinfo":
			u.ZoneInfo = aws.ToString(a.Value)
		case "locale":
			u.Locale = aws.ToString(a.Value)
		case "email_verified":
			u.EmailVerified = aws.ToString(a.Value) == "true"
		}
	}

	return u
}
