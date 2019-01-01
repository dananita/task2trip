// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"net/http"

	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"

	"github.com/itimofeev/task2trip/rest/restapi/operations"
	"github.com/itimofeev/task2trip/rest/restapi/operations/categories"
	"github.com/itimofeev/task2trip/rest/restapi/operations/offers"
	"github.com/itimofeev/task2trip/rest/restapi/operations/tasks"
	"github.com/itimofeev/task2trip/rest/restapi/operations/users"
)

//go:generate swagger generate server --target ../rest --name Task2Trip --spec ../tools/swagger.yml

func configureFlags(api *operations.Task2TripAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.Task2TripAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	// Applies when the "X-Auth-Token" header is set
	api.AuthTokenAuth = func(token string) (interface{}, error) {
		return nil, errors.NotImplemented("api key auth (AuthToken) X-Auth-Token from header param [X-Auth-Token] has not yet been implemented")
	}

	// Set your custom authorizer if needed. Default one is security.Authorized()
	// Expected interface runtime.Authorizer
	//
	// Example:
	// api.APIAuthorizer = security.Authorized()
	api.TasksCreateTaskHandler = tasks.CreateTaskHandlerFunc(func(params tasks.CreateTaskParams, principal interface{}) middleware.Responder {
		return middleware.NotImplemented("operation tasks.CreateTask has not yet been implemented")
	})
	api.CategoriesListCategoriesHandler = categories.ListCategoriesHandlerFunc(func(params categories.ListCategoriesParams) middleware.Responder {
		return middleware.NotImplemented("operation categories.ListCategories has not yet been implemented")
	})
	api.OffersListOffersHandler = offers.ListOffersHandlerFunc(func(params offers.ListOffersParams, principal interface{}) middleware.Responder {
		return middleware.NotImplemented("operation offers.ListOffers has not yet been implemented")
	})
	api.TasksSearchTasksHandler = tasks.SearchTasksHandlerFunc(func(params tasks.SearchTasksParams, principal interface{}) middleware.Responder {
		return middleware.NotImplemented("operation tasks.SearchTasks has not yet been implemented")
	})
	api.UsersUserLoginHandler = users.UserLoginHandlerFunc(func(params users.UserLoginParams) middleware.Responder {
		return middleware.NotImplemented("operation users.UserLogin has not yet been implemented")
	})
	api.UsersUserSignupHandler = users.UserSignupHandlerFunc(func(params users.UserSignupParams) middleware.Responder {
		return middleware.NotImplemented("operation users.UserSignup has not yet been implemented")
	})

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix"
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
