package tests

import (
	"fmt"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	client2 "github.com/itimofeev/task2trip/rest/client"
	"github.com/itimofeev/task2trip/rest/client/users"
	"github.com/itimofeev/task2trip/rest/models"
	"github.com/itimofeev/task2trip/rest/restapi"
	"github.com/itimofeev/task2trip/rest/restapi/operations"
	"github.com/itimofeev/task2trip/util"
	"github.com/itimofeev/task2trip/util/client"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
)

func Init() *client2.Task2Trip {
	swaggerSpec, err := loads.Embedded(restapi.SwaggerJSON, restapi.FlatSwaggerJSON)
	if err != nil {
		log.Fatalln(err)
	}
	api := operations.NewTask2TripAPI(swaggerSpec)
	server := restapi.NewServer(api)
	server.ConfigureAPI()

	handler := server.GetHandler()

	c := client.New(client2.DefaultHost, "/api/v1", client2.DefaultSchemes)
	do := client.NewDO(handler)
	c.WithDO(do.Do)

	return client2.New(c, nil)
}

type TokenAuth struct {
	AuthToken string
}

func (a *TokenAuth) AuthenticateRequest(r runtime.ClientRequest, _ strfmt.Registry) error {
	return r.SetHeaderParam("Authorization", fmt.Sprintf("Bearer %s", a.AuthToken))
}

func Test_User_SignUP(t *testing.T) {
	task2trip := Init()

	email := util.RandEmail()
	pass := "hello, there"

	signUpOk, err := task2trip.Users.UserSignup(users.NewUserSignupParams().WithUser(&models.UserCreateParams{Email: &email, Password: &pass}))
	require.NoError(t, err)
	require.Equal(t, email, *signUpOk.Payload.Name)

	loginOk, err := task2trip.Users.UserLogin(users.NewUserLoginParams().WithCredentials(users.UserLoginBody{Email: &email, Password: &pass}))
	require.NoError(t, err)

	currentUserOk, err := task2trip.Users.CurrentUser(users.NewCurrentUserParams(), &TokenAuth{AuthToken: loginOk.Payload.AuthToken})
	require.NoError(t, err)
	require.Equal(t, email, *currentUserOk.Payload.Name)
}
