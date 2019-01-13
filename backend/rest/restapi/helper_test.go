package restapi

import (
	"fmt"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/itimofeev/task2trip/backend/rest/client/tasks"
	"github.com/itimofeev/task2trip/backend/rest/client/users"
	"github.com/itimofeev/task2trip/backend/rest/models"
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

func generateAuthToken(user *models.User) string {
	return util.GenerateAuthToken(&util.Claims{
		UserName: *user.Name,
		UserID:   *user.ID,
	})
}

func (a *UserAuth) AuthenticateRequest(r runtime.ClientRequest, _ strfmt.Registry) error {
	token := generateAuthToken(a.user)

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

func createTask(t *testing.T, user *models.User) *models.Task {
	cats, err := Store.ListCategories()
	require.NoError(t, err)

	taskCreatedOk, err := api.Tasks.CreateTask(tasks.NewCreateTaskParams().WithTask(&models.TaskCreateParams{
		Name:           util.PtrFromString("my super Task"),
		BudgetEstimate: util.PtrFromInt64(100),
		CategoryID:     util.PtrFromString(cats[0].ID),
		Description:    util.PtrFromString("my super Description"),
	}), userAuth(user))

	require.NoError(t, err)
	require.Equal(t, taskCreatedOk.Payload.Name, util.PtrFromString("my super Task"))
	require.Equal(t, taskCreatedOk.Payload.BudgetEstimate, util.PtrFromInt64(100))
	require.Equal(t, taskCreatedOk.Payload.Category.ID, util.PtrFromString(cats[0].ID))
	require.Equal(t, taskCreatedOk.Payload.Description, util.PtrFromString("my super Description"))

	return taskCreatedOk.Payload
}
