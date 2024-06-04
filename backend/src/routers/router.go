package router

import (
	"covid/src/api/middleware"
	"fmt"
	"net/http"
	"os"

	"covid/src/cmd/config"

	deathCase "covid/src/pkg/deathcase"
	externalDdc "covid/src/pkg/external/ddc"

	deathCaseHandler "covid/src/api/deathcase/handler"

	"covid/src/pkg/utils/client"
	"covid/src/pkg/utils/logger"

	"github.com/gorilla/mux"

	httpSwagger "github.com/swaggo/http-swagger"
)

const (
	defaultPort   = "8000"
	defaultAppEnv = "LOCAL"
)

type RouterConfig struct {
	Logs   *logger.ELK
	Config config.Configuration
}

func InitRouter(routerConfig RouterConfig) http.Handler {
	fmt.Println("InitRouter :", routerConfig)

	/* Service */
	clientService := client.NewService(routerConfig.Logs)
	externalDdc.DDCURL = routerConfig.Config.DCCAPI
	externalDccService := externalDdc.NewService(clientService,routerConfig.Logs)
	deathCaseService := deathCase.NewService(externalDccService,routerConfig.Logs)

	/* Handler */
	deathCaseHandlerParam := deathCaseHandler.HandlerParam{
		Service: deathCaseService,
		Logs:    routerConfig.Logs,
	}
	deathCaseHandler := deathCaseHandler.MakeDeathCaseHandler(deathCaseHandlerParam)

	/* Router */
	middleware := middleware.NewMiddleware(routerConfig.Logs)
	router := mux.NewRouter()
	router.Use(middleware.Recover)
	router.HandleFunc("/version", versionHandler)
	router.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)

	deathCase := router.PathPrefix("/deathcases").Subrouter()
	deathCase.HandleFunc("/gettopthreeprovice", http.HandlerFunc(deathCaseHandler.GetTopThreeProvice)).Methods("GET")

	return router
}

func versionHandler(w http.ResponseWriter, r *http.Request) {
	appVersion := getEnvString("APP_VERSION", defaultPort)
	fmt.Fprintln(w, appVersion)
}

func getEnvString(env, fallback string) string {
	result := os.Getenv(env)
	if result == "" {
		return fallback
	}
	return result
}