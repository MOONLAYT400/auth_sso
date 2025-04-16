package tests

import (
	"auth-service/tests/suite"
	"strconv"
	"testing"

	ssov1 "github.com/MOONLAYT400/Proto_sso/gen/go/sso"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	emptyAppID = 0
	appID      = 1
	appSecret  = "test-secret"

	passDefaultLen = 10
)

func TestAuth_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)

	email:= gofakeit.Email()
	pass := randomFakePassword()

	responseReg, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: pass,
	})

	require.NoError(t, err)
	assert.NotEmpty(t, responseReg.GetUserId()) 

	responseLogin,err := st.AuthClient.Login(ctx,&ssov1.LoginRequest{Email: email,Password: randomFakePassword(), AppId: appID})

	require.NoError(t, err)
	
	token:=responseLogin.GetToken()
	require.NotEmpty(t, token)

	parsedToken,err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
			return []byte(appSecret), nil
	})
		require.NoError(t, err)

	claims,ok := parsedToken.Claims.(jwt.MapClaims)
	assert.True(t, ok)

	assert.Equal(t, responseReg.GetUserId(),floatToString( claims["user_id"].(float64)))
	assert.Equal(t, email, claims["email"].(string))
	assert.Equal(t, appID, int(claims["app_id"].(float64)))


}

func randomFakePassword () string {
		password := gofakeit.Password(true, true, true, true, true, passDefaultLen)
		return password
}

func floatToString(input_num float64) string {
    // to convert a float number to a string
    return strconv.FormatFloat(input_num, 'f', 0, 64)
}