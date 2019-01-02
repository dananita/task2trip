package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/itimofeev/task2trip/backend"
	"github.com/itimofeev/task2trip/rest/models"
	"github.com/itimofeev/task2trip/rest/restapi/operations/categories"
	"github.com/itimofeev/task2trip/rest/restapi/operations/tasks"
	"github.com/itimofeev/task2trip/util"
)

var TasksCreateTaskHandler = tasks.CreateTaskHandlerFunc(func(params tasks.CreateTaskParams, principal interface{}) middleware.Responder {
	user := principal.(*backend.User)
	task, err := store.CreateTask(user, params.Task)
	if err != nil {
		return util.ConvertHTTPErrorToResponse(err)
	}

	return tasks.NewCreateTaskCreated().WithPayload(convertTask(task))
})

func convertTask(task *backend.Task) *models.Task {
	return &models.Task{
		ID:         &task.ID,
		Name:       &task.Name,
		Category:   convertCategory(task.Category),
		CreateTime: strfmt.DateTime(task.CreateTime),
	}
}

func convertCategory(category *backend.Category) *models.Category {
	return &models.Category{
		ID:           &category.ID,
		Key:          &category.Key,
		DefaultValue: category.DefaultValue,
	}
}

var CategoriesListCategoriesHandler = categories.ListCategoriesHandlerFunc(func(params categories.ListCategoriesParams) middleware.Responder {
	categs, err := store.ListCategories()
	if err != nil {
		return util.ConvertHTTPErrorToResponse(err)
	}

	return categories.NewListCategoriesOK().WithPayload(convertCategories(categs))
})

func convertCategories(categs []*backend.Category) (res []*models.Category) {
	for _, category := range categs {
		res = append(res, &models.Category{
			ID:           &category.ID,
			Key:          &category.Key,
			DefaultValue: category.DefaultValue,
		})
	}
	return res
}
