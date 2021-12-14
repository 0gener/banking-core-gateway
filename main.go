package main

import (
	"log"
	"net/http"
	"os"

	"github.com/0gener/banking-core-gateway/client"
	"github.com/0gener/banking-core-gateway/middleware"
	"github.com/0gener/banking-core-gateway/router"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	// FIXME: gin not loading GIN_MODE env correctly
	log.Println(os.Getenv("GIN_MODE"))

	jwtMiddleware, err := middleware.NewJwtMiddleware(&middleware.JwtMiddlewareOptions{
		Domain:   os.Getenv("AUTH0_DOMAIN"),
		Audience: os.Getenv("AUTH0_AUDIENCE"),
	})
	if err != nil {
		log.Fatalf("There was an error creating auth0 jwt middleware: %v", err)
	}

	accountsClient, conn := client.NewAccountsClient(client.AccountsClientOptions{
		Url: os.Getenv("ACCOUNTS_SERVICE_URL"),
	})
	defer conn.Close()
	r := router.New(jwtMiddleware, accountsClient)

	log.Print("Server listening on http://localhost:8080")
	if err := http.ListenAndServe("0.0.0.0:8080", r); err != nil {
		log.Fatalf("There was an error with the http server: %v", err)
	}
}
