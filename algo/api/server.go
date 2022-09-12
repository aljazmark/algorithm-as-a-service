package api

import (
	"algo/routes"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func InitAlgo() {
	//Router and server setup
	e := echo.New()
	e.GET("/", checkHandler)
	e.GET("/algorithm", routes.GetAlgorithms)
	e.POST("/algorithm/:algo", routes.RunAlgorithm)
	setPort := ":8080"
	if os.Getenv("PORT") != "" {
		setPort = ":" + os.Getenv("PORT")
	}
	server := &http.Server{
		Addr:         setPort,
		Handler:      e,
		IdleTimeout:  180 * time.Second,
		ReadTimeout:  180 * time.Second,
		WriteTimeout: 180 * time.Second,
	}
	//Middleware
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "time=${time_rfc3339_nano}, method=${method}, uri=${uri}, status=${status}\n",
	}))
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		//AllowOrigins: []string{"http://localhost:4200", "http://localhost:8090"},
		AllowOrigins: []string{"*"},
	}))
	//Recover Middleware to handle any panic
	e.Use(middleware.Recover())
	//Server go routine
	go func() {
		e.StartServer(server)
	}()
	e.HideBanner = true
	//Graceful shutdown
	shutdownChannel := make(chan os.Signal)
	signal.Notify(shutdownChannel, os.Interrupt)
	signal.Notify(shutdownChannel, os.Kill)
	shutdownSignal := <-shutdownChannel
	log.Println("[ALGO]: Shutdown by signal:", shutdownSignal)
	shutdownContext, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	server.Shutdown(shutdownContext)
}
func checkHandler(c echo.Context) error {
	return c.String(http.StatusOK, "Welcome to algo")
}
