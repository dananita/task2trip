package tests

import (
	"github.com/go-openapi/loads"
	"github.com/itimofeev/task2trip/backend/handlers"
	client2 "github.com/itimofeev/task2trip/rest/client"
	"github.com/itimofeev/task2trip/rest/client/tasks"
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

var api = Init()

func Test_User_SignUP(t *testing.T) {
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
}

func withRandomUser(t *testing.T, f func(authToken string)) {
	email := util.RandEmail()
	pass := "hello, there"

	signUpOk, err := api.Users.UserSignup(users.NewUserSignupParams().WithUser(&models.UserCreateParams{Email: &email, Password: &pass}))
	require.NoError(t, err)
	require.Equal(t, email, *signUpOk.Payload.Name)

	loginOk, err := api.Users.UserLogin(users.NewUserLoginParams().WithCredentials(users.UserLoginBody{Email: &email, Password: &pass}))
	require.NoError(t, err)
	f(loginOk.Payload.AuthToken)
}

func Test_User_CreateTask(t *testing.T) {
	withRandomUser(t, func(authToken string) {
		cats, err := handlers.Store.ListCategories()
		require.NoError(t, err)

		taskCreatedOk, err := api.Tasks.CreateTask(tasks.NewCreateTaskParams().WithTask(&models.TaskCreateParams{
			Name:           util.PtrFromString("my super Task"),
			BudgetEstimate: util.PtrFromInt64(100),
			CategoryID:     util.PtrFromString(cats[0].ID),
			Description:    util.PtrFromString("my super Description"),
		}), &TokenAuth{AuthToken: authToken})

		require.NoError(t, err)
		require.Equal(t, taskCreatedOk.Payload.Name, util.PtrFromString("my super Task"))
		require.Equal(t, taskCreatedOk.Payload.BudgetEstimate, util.PtrFromInt64(100))
		require.Equal(t, taskCreatedOk.Payload.Category.ID, util.PtrFromString(cats[0].ID))
		require.Equal(t, taskCreatedOk.Payload.Description, util.PtrFromString("my super Description"))
	})
}
