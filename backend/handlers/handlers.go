package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/itimofeev/task2trip/backend"
	"github.com/itimofeev/task2trip/backend/postgres"
	"github.com/itimofeev/task2trip/rest/restapi/operations/misc"
	"github.com/itimofeev/task2trip/rest/restapi/operations/users"
	"github.com/itimofeev/task2trip/util"
)

var store backend.Store

func Init() {
	store = postgres.NewStore("postgresql://postgres@db:5432/postgres?sslmode=disable")
}

var AuthFunc = func(token string) (interface{}, error) {
	claims := &util.Claims{}
	if err := util.ParseJWT(token, claims); err != nil {
		return nil, err
	}

	user, err := store.GetUserByID(token)
	if err != nil {
		e := util.ConvertHTTPErrorToResponse(err)
		t, _ := e.(error)
		return nil, t

	}
	return user, nil
}

var UserSignupHandlerFunc = users.UserSignupHandlerFunc(func(params users.UserSignupParams) middleware.Responder {
	return middleware.NotImplemented("operation users.UserLogin has not yet been implemented")
})

var AboutHandler = misc.AboutHandlerFunc(func(params misc.AboutParams) middleware.Responder {
	return misc.NewAboutOK().WithPayload("hello, there")
})
