package main

import (
	"github.com/go-openapi/loads"
	"log"
	"user-store/configs"
	"user-store/internal/restapi"
	"user-store/internal/restapi/operations"
)

func main()  {

	cfg, err := configs.NewConfig()
	if err != nil {
		log.Fatalln(err)
	}

	swaggerSpec,err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}

	api := operations.NewUserStoreAPI(swaggerSpec)
	server := restapi.NewServer(api)
	defer server.Shutdown()

	server.ConfigureAPI()
	server.Port = cfg.Port

	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}
}
