// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"github.com/itimofeev/task2trip/backend/handlers"
	"github.com/itimofeev/task2trip/util"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	interpose "github.com/carbocation/interpose/middleware"
	"github.com/dre1080/recover"

	"github.com/itimofeev/task2trip/rest/restapi/operations"
	"github.com/itimofeev/task2trip/rest/restapi/operations/offers"
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
	api.AuthTokenAuth = handlers.AuthFunc

	handlers.Init()

	// Set your custom authorizer if needed. Default one is security.Authorized()
	// Expected interface runtime.Authorizer
	//
	// Example:
	// api.APIAuthorizer = security.Authorized()

	api.MiscAboutHandler = handlers.AboutHandler

	api.TasksCreateTaskHandler = handlers.TasksCreateTaskHandler

	api.CategoriesListCategoriesHandler = handlers.CategoriesListCategoriesHandler
	api.OffersListOffersHandler = offers.ListOffersHandlerFunc(func(params offers.ListOffersParams, principal interface{}) middleware.Responder {
		return middleware.NotImplemented("operation offers.ListOffers has not yet been implemented")
	})

	api.TasksSearchTasksHandler = handlers.TasksSearchTasksHandler
	api.UsersUserLoginHandler = handlers.UsersUserLoginHandler
	api.UsersUserSignupHandler = handlers.UserSignupHandlerFunc
	api.UsersCurrentUserHandler = handlers.UsersCurrentUserHandler

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
	handlePanic := recover.New(&recover.Options{
		Log: util.Log.Debug,
	})

	logViaLogrus := interpose.NegroniLogrus()

	return handlePanic(
		logViaLogrus(
			handler,
		),
	)
}
