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

func userAuth(user *models.User) runtime.ClientAuthInfoWriter {
	return &UserAuth{user}
}

type UserAuth struct {
	user *models.User
}

func (a *UserAuth) AuthenticateRequest(r runtime.ClientRequest, _ strfmt.Registry) error {
	token := util.GenerateAuthToken(&util.Claims{
		UserName: *a.user.Name,
		UserID:   *a.user.ID,
	})

	return r.SetHeaderParam("Authorization", fmt.Sprintf("Bearer %s", token))
}

func createUser(t *testing.T) *models.User {
	email := util.RandEmail()
	pass := "hello, there"

	signUpOk, err := api.Users.UserSignup(users.NewUserSignupParams().WithUser(&models.UserCreateParams{Email: &email, Password: &pass}))
	require.NoError(t, err)
	require.Equal(t, email, *signUpOk.Payload.Name)

	_, err = api.Users.UserLogin(users.NewUserLoginParams().WithCredentials(users.UserLoginBody{Email: &email, Password: &pass}))
	require.NoError(t, err)

	currentUserOk, err := api.Users.CurrentUser(users.NewCurrentUserParams(), userAuth(signUpOk.Payload))
	require.NoError(t, err)
	require.Equal(t, email, *currentUserOk.Payload.Name)

	return currentUserOk.Payload
}
