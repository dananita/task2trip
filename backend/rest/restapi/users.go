package restapi

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/itimofeev/task2trip/backend"
	"github.com/itimofeev/task2trip/backend/postgres"
	"github.com/itimofeev/task2trip/backend/rest/restapi/operations/misc"
	"github.com/itimofeev/task2trip/backend/rest/restapi/operations/users"
	"github.com/itimofeev/task2trip/util"
	"os"
	"strings"
)

var Store backend.Store

func Init() {
	dbAddr := os.Getenv("DB_ADDR")
	if len(dbAddr) == 0 {
		dbAddr = "postgresql://postgres@db:5432/postgres?sslmode=disable"
	}
	Store = postgres.NewStore(dbAddr)
}

var AuthFunc = func(token string) (interface{}, error) {
	token = strings.TrimPrefix(token, "Bearer ")
	claims := &util.Claims{}
	if err := util.ParseJWT(token, claims); err != nil {
		return nil, err
	}

	user, err := Store.GetUserByID(claims.UserID)
	if err != nil {
		e := util.ConvertHTTPErrorToResponse(err)
		t, _ := e.(error)
		return nil, t

	}
	return user, nil
}

var UserSignupHandlerFunc = users.UserSignupHandlerFunc(func(params users.UserSignupParams) middleware.Responder {
	user, err := Store.CreateUser(*params.User.Email, *params.User.Password)
	if err != nil {
		return util.ConvertHTTPErrorToResponse(err)
	}
	return users.NewUserSignupOK().WithPayload(convertUser(user))
})

var AboutHandler = misc.AboutHandlerFunc(func(params misc.AboutParams) middleware.Responder {
	return misc.NewAboutOK().WithPayload("hello, there")
})

var UsersCurrentUserHandler = users.CurrentUserHandlerFunc(func(params users.CurrentUserParams, principal interface{}) middleware.Responder {
	user := principal.(*backend.User)

	return users.NewCurrentUserOK().WithPayload(convertUser(user))
})

var UsersUserLoginHandler = users.UserLoginHandlerFunc(func(params users.UserLoginParams) middleware.Responder {
	user, err := Store.GetUserByEmailAndPassword(*params.Credentials.Email, *params.Credentials.Password)
	if err != nil {
		return util.ConvertHTTPErrorToResponse(err)
	}

	return users.NewUserLoginOK().WithPayload(&users.UserLoginOKBody{
		AuthToken: util.GenerateAuthToken(&util.Claims{
			UserName: user.Email,
			UserID:   user.ID,
		}),
	})
})
