package restapi

import (
	"github.com/go-openapi/loads"
	client2 "github.com/itimofeev/task2trip/rest/client"
	"github.com/itimofeev/task2trip/rest/client/tasks"
	"github.com/itimofeev/task2trip/rest/models"
	"github.com/itimofeev/task2trip/rest/restapi/operations"
	"github.com/itimofeev/task2trip/util"
	"github.com/itimofeev/task2trip/util/client"
	"github.com/stretchr/testify/require"
	"testing"
)

func InitTestAPI() *client2.Task2Trip {
	swaggerSpec, err := loads.Spec("../../tools/swagger.yml")
	if err != nil {
		util.Log.Fatalln(err)
	}
	api := operations.NewTask2TripAPI(swaggerSpec)
	server := NewServer(api)
	server.ConfigureAPI()

	handler := server.GetHandler()

	c := client.New(client2.DefaultHost, "/api/v1", client2.DefaultSchemes)
	do := client.NewDO(handler)
	c.WithDO(do.Do)

	return client2.New(c, nil)
}

var api = InitTestAPI()

func Test_User_CreateTask(t *testing.T) {
	withRandomUser(t, func(authToken string) {
		cats, err := Store.ListCategories()
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