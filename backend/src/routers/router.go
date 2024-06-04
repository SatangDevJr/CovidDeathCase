package router

import (
	"covid/src/api/middleware"
	"fmt"
	"net/http"
	"os"

	"covid/src/cmd/config"

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

	/* Handler */

	/* Router */
	middleware := middleware.NewMiddleware(routerConfig.Logs)
	router := mux.NewRouter()
	router.Use(middleware.Recover)
	router.HandleFunc("/version", versionHandler)
	router.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)

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
