package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/itimofeev/task2trip/backend"
	"github.com/itimofeev/task2trip/rest/restapi/operations/categories"
	"github.com/itimofeev/task2trip/rest/restapi/operations/tasks"
	"github.com/itimofeev/task2trip/util"
)

var TasksCreateTaskHandler = tasks.CreateTaskHandlerFunc(func(params tasks.CreateTaskParams, principal interface{}) middleware.Responder {
	user := principal.(*backend.User)
	task, err := Store.CreateTask(user, params.Task)
	if err != nil {
		return util.ConvertHTTPErrorToResponse(err)
	}

	return tasks.NewCreateTaskCreated().WithPayload(convertTask(task))
})

var TasksSearchTasksHandler = tasks.SearchTasksHandlerFunc(func(params tasks.SearchTasksParams) middleware.Responder {
	var user *backend.User
	if params.Authorization != nil {
		principal, err := AuthFunc(*params.Authorization)
		if err == nil {
			user = principal.(*backend.User)
		}
	}

	tasks_, total, err := Store.SearchTasks(user, params)
	if err != nil {
		return util.ConvertHTTPErrorToResponse(err)
	}
	return tasks.NewSearchTasksOK().WithPayload(convertTasksPage(tasks_, total))
})

var CategoriesListCategoriesHandler = categories.ListCategoriesHandlerFunc(func(params categories.ListCategoriesParams) middleware.Responder {
	categories_, err := Store.ListCategories()
	if err != nil {
		return util.ConvertHTTPErrorToResponse(err)
	}

	return categories.NewListCategoriesOK().WithPayload(convertCategories(categories_))
})
