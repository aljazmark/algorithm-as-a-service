package api

import (
	"algoAPI/models"
	"algoAPI/routes"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/didip/tollbooth_echo"
	openapi "github.com/go-openapi/runtime/middleware"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//Gets JWT secret from enviroment variable
//var secret = os.Getenv("JWTSecret")
var secret = "ExampleSecret42"

//InitAlgoAPI initiates the server
func InitAlgoAPI() {
	//Logger
	log := log.New(os.Stdout, "[ALGO]: ", log.LstdFlags)
	//Database connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	//client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb+srv://algoAPI:"+os.Getenv("DBpw")+"@cluster0.wc2t4.azure.mongodb.net/Algo?retryWrites=true&w=majority"))
	if err != nil {
		panic(err)
	}
	log.Println("| DB Connected")
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	//Controllers
	collection := client.Database("Algo")
	requestsHandlers := routes.NewRequests(log, collection)
	usersHandlers := routes.NewUsers(log, collection)
	datasHandlers := routes.NewDatas(log, collection)
	helpsHandlers := routes.NewHelps(log, models.Help{})
	//JWT Authentication Middleware
	var loggedIn = middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(secret),
	})
	//Rate limit middlware

	//Router
	e := echo.New()
	e.GET("/", checkHandler)
	//Users
	e.GET("/user/:id", loggedIn(usersHandlers.GetUser))
	e.POST("/user", usersHandlers.NewUser)
	e.PUT("/user/:id", loggedIn(usersHandlers.UpdateUser))
	e.POST("/user/login", usersHandlers.LoginUser)
	e.DELETE("/user/:id", loggedIn(usersHandlers.DeleteUser))
	//Requests
	e.GET("/request/:id", switchMiddleware(requestsHandlers.GetRequest))
	e.GET("/request/user/:id", loggedIn(requestsHandlers.GetRequestsByUser))
	e.POST("/request/:id", switchMiddleware(requestsHandlers.NewRequest))
	e.POST("/request/:id/:data", loggedIn(requestsHandlers.NewRequestWithData))
	e.DELETE("/request/:id", loggedIn(requestsHandlers.DeleteRequest))
	//Datas
	e.GET("/data/:id", loggedIn(datasHandlers.GetData))
	e.GET("/data/user/:id", loggedIn(datasHandlers.GetDatasByUser))
	e.POST("/data", loggedIn(datasHandlers.NewData))
	e.PUT("/data/:id", loggedIn(datasHandlers.UpdateData))
	e.DELETE("/data/:id", loggedIn(datasHandlers.DeleteData))
	//Help
	e.GET("/help/algorithms", helpsHandlers.GetAlgorithms)
	e.GET("/help/:algorithm", helpsHandlers.GetAlgoHelp)
	if x := helpsHandlers.PrepHelp(); x {
		log.Println("| help.json parsed successfully")
	} else {
		log.Println("| help.json parsed unsuccessfully")
	}
	//Documentation
	//docOptions := openapi.RedocOpts{SpecURL: "/swagger.yaml"}
	docOptionsRedoc := openapi.RedocOpts{Path: "/docsRedoc"}
	docsHandlerRedoc := openapi.Redoc(docOptionsRedoc, nil)
	docOptionsSwaggerUI := openapi.SwaggerUIOpts{}
	docsHandlerSwaggerUI := openapi.SwaggerUI(docOptionsSwaggerUI, nil)
	e.GET("/docsRedoc", echo.WrapHandler(docsHandlerRedoc))
	e.GET("/docs", echo.WrapHandler(docsHandlerSwaggerUI))
	e.File("/swagger.yaml", "./swagger.yaml")
	e.File("/swagger.json", "./swagger.json")
	//Server address set to localhost with port 8090 if not set in env. variable
	setPort := ":8090"
	if os.Getenv("PORT") != "" {
		setPort = ":" + os.Getenv("PORT")
	}
	//Timeouts at 3 minutes
	server := &http.Server{
		Addr:         setPort,
		Handler:      e,
		IdleTimeout:  180 * time.Second,
		ReadTimeout:  180 * time.Second,
		WriteTimeout: 180 * time.Second,
	}
	//Logger Middleware
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "time=${time_rfc3339_nano}, method=${method}, uri=${uri}, status=${status}\n",
	}))
	//CORS Middleware
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
	}))
	//Ratelimit Middleware(1 request/second per IP)
	rateLimiter := tollbooth.NewLimiter(1, &limiter.ExpirableOptions{DefaultExpirationTTL: time.Hour})
	rateLimiter.SetBurst(10)
	rateLimiter.SetMessage("Rate limit reached. Limit: 1 request/second Burst: 10")
	e.Use(tollbooth_echo.LimitHandler(rateLimiter))
	//Recover Middleware to handle any panic
	e.Use(middleware.Recover())
	//XSS protection middleware
	e.Use(middleware.Secure())
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
	log.Println("| Shutdown by signal:", shutdownSignal)
	shutdownContext, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	server.Shutdown(shutdownContext)

}

//switchMiddleware is middleware use for letting request through even if user is not logged in, otherwise it passes JWT
func switchMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	handler := func(c echo.Context) error {
		return nil
	}
	return func(c echo.Context) error {
		handler = middleware.JWTWithConfig(middleware.JWTConfig{
			SigningKey: []byte(secret),
		})(handler)
		handler(c)
		return next(c)
	}
}

func checkHandler(c echo.Context) error {
	return c.String(http.StatusOK, "Welcome to algo API, please visit /docs for more information.")
}
