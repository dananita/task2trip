package restapi

import (
	"fmt"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/itimofeev/task2trip/rest/client/users"
	"github.com/itimofeev/task2trip/rest/models"
	"github.com/itimofeev/task2trip/util"
	"github.com/stretchr/testify/require"
	"testing"
)

type TokenAuth struct {
	AuthToken string
}

func (a *TokenAuth) AuthenticateRequest(r runtime.ClientRequest, _ strfmt.Registry) error {
	return r.SetHeaderParam("Authorization", fmt.Sprintf("Bearer %s", a.AuthToken))
}

func withRandomUser(t *testing.T, f func(authToken string)) {
	email := util.RandEmail()
	pass := "hello, there"

	signUpOk, err := api.Users.UserSignup(users.NewUserSignupParams().WithUser(&models.UserCreateParams{Email: &email, Password: &pass}))
	require.NoError(t, err)
	require.Equal(t, email, *signUpOk.Payload.Name)

	loginOk, err := api.Users.UserLogin(users.NewUserLoginParams().WithCredentials(users.UserLoginBody{Email: &email, Password: &pass}))
	require.NoError(t, err)

	currentUserOk, err := api.Users.CurrentUser(users.NewCurrentUserParams(), &TokenAuth{AuthToken: loginOk.Payload.AuthToken})
	require.NoError(t, err)
	require.Equal(t, email, *currentUserOk.Payload.Name)

	f(loginOk.Payload.AuthToken)
}
