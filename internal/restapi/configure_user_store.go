// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"database/sql"
	"log"
	"net/http"
	"user-store/internal/models"
	"user-store/repository"
	"user-store/storage/sqlite"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	"user-store/internal/restapi/operations"
	"user-store/internal/restapi/operations/users"
)

//go:generate swagger generate server --target ../../internal --name UserStore --spec ../../swagger/user-store.yaml --principal interface{} --exclude-main

func configureFlags(api *operations.UserStoreAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.UserStoreAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	db := sqlite.New()
	userRepo := repository.CreateUserHandler()

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.UseSwaggerUI()
	// To continue using redoc as your UI, uncomment the following line
	// api.UseRedoc()

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	//Create User in db
	api.UsersAddOneHandler = users.AddOneHandlerFunc(func(params users.AddOneParams) middleware.Responder {
		usr, err := userRepo.CreateUser(db, params.Body)
		if err != nil {
			log.Printf("Internal server error. Cannot create user %v\n", err)
			return users.NewAddOneDefault(500).WithPayload(&models.Error{Code: 500, Message: "Internal server error"})
		}
		log.Printf("User with id %v created", usr.ID)
		return  users.NewAddOneCreated().WithPayload(usr)
	})

	//Get all users from db
	api.UsersGetUsersHandler = users.GetUsersHandlerFunc(func(params users.GetUsersParams) middleware.Responder {
		usrs,err := userRepo.GetAllUsers(db)
		if err != nil {
			log.Printf("Internal Server Error %v. Cannot get all users\n", err)
			return users.NewGetUsersDefault(500).WithPayload(&models.Error{Code: 500, Message: "Internal server error"})
		}
		log.Println("Get all users")
		return  users.NewGetUsersOK().WithPayload(usrs)
	})

	//Get user by id
	api.UsersGetOneHandler = users.GetOneHandlerFunc(func(params users.GetOneParams) middleware.Responder {
		user, err := userRepo.GetUser(db, params.ID)
		if err != nil {
			switch err {
			case sql.ErrNoRows:
				log.Printf("User with id %v not found\n", params.ID)
				return users.NewGetOneNotFound()
			default:
				log.Printf("Internal Server Error %v. Cannot get users id %v\n", err, params.ID)
				return users.NewGetOneDefault(500).WithPayload(&models.Error{Code: 500, Message: "Internal server error"})
			}
		}
		log.Printf("Get user with id %v ", params.ID)
		return users.NewGetOneOK().WithPayload(user)
	})

	//Delete user from db
	api.UsersDeleteOneHandler = users.DeleteOneHandlerFunc(func(params users.DeleteOneParams) middleware.Responder {
		_, err := userRepo.GetUser(db, params.ID)
		if err != nil {
			switch err {
			case sql.ErrNoRows:
				log.Printf("User with id %v not found\n", params.ID)
				return users.NewDeleteOneNotFound()
			default:
				log.Printf("Internal Server Error %v. Cannot delete users id %v\n", err, params.ID)
				return users.NewDeleteOneDefault(500).WithPayload(&models.Error{Code: 500, Message: "Internal server error"})
			}
		}

		if err := userRepo.DeleteUser(db, params.ID); err != nil {
			log.Printf("Internal Server Error %v. Cannot delete users id %v\n", err, params.ID)
			return users.NewDeleteOneDefault(500).WithPayload(&models.Error{Code: 500, Message: "Internal server error"})
		}
		log.Printf("User with id %v deleted", params.ID)
		return users.NewDeleteOneNoContent()
	})

	//Update user from db
	api.UsersUpdateOneHandler = users.UpdateOneHandlerFunc(func(params users.UpdateOneParams) middleware.Responder {
		_, err := userRepo.GetUser(db, params.ID)
		if err != nil {
			switch err {
			case sql.ErrNoRows:
				log.Printf("User with id %v not found\n", params.ID)
				return users.NewUpdateOneNotFound()
			default:
				log.Printf("Internal Server Error %v. Cannot update users id %v\n", err, params.ID)
				return users.NewUpdateOneDefault(500).WithPayload(&models.Error{Code: 500, Message: "Internal server error"})
			}
		}
		usr, err := userRepo.UpdateUser(db, params.ID, params.Body)
		if err != nil {
			log.Printf("Internal Server Error %v. Cannot update users id %v\n", err, params.ID)
			return users.NewUpdateOneDefault(500).WithPayload(&models.Error{Code: 500, Message: "Internal server error"})
		}
		log.Printf("User with id %v updated", params.ID)
		return users.NewUpdateOneOK().WithPayload(usr)
	})


	api.PreServerShutdown = func() {
		if err := db.Sql.Close(); err!=nil {
			log.Fatal(err)
		}
		log.Println("Connect to DB close")
	}

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
// scheme value will be set accordingly: "http", "https" or "unix".
func configureServer(s *http.Server, scheme, addr string) {

}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation.
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics.
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
