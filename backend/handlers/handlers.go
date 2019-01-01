package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/itimofeev/task2trip/rest/restapi/operations/users"
)

var UserSignupHandlerFunc = users.UserSignupHandlerFunc(func(params users.UserSignupParams) middleware.Responder {
	return middleware.NotImplemented("operation users.UserLogin has not yet been implemented")
})
