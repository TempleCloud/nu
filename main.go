package main

import (
	"log"
	"net/http"
	"time"

	"github.com/boltdb/bolt"
	"github.com/gin-gonic/gin"
	"github.com/templecloud/nu/boltdb"
	"github.com/templecloud/nu/handlers"
)

//--------------------------------------------------------------------------------------------------

// The System contains pointers to system wide objects
type System struct {
	db     *bolt.DB
	server *http.Server
}

//--------------------------------------------------------------------------------------------------

// func main() {
//
// 	// |*| : Initialise BoltDB
// 	//
//
// 	db := boltdb.NewDB()
// 	defer db.Close()
//
// 	// |*| : Initialise Httprouter
// 	//
// 	router := httprouter.New()
//
// 	// health check endpoints
// 	router.GET("/v1/health/ping", handlers.Ping)
// 	router.GET("/v1/health/echo", handlers.Echo)
//
// 	// nu endpoints - functions
// 	router.GET("/v1/nu/functions", handlers.ListFunctions(db))
// 	router.GET("/v1/nu/functions/:functionId", handlers.GetFunction(db))
// 	router.PUT("/v1/nu/functions", handlers.RegisterFunction(db))
// 	// router.POST("/v1/nu/functions/:functionId", nuFunctionHandler)
// 	// router.DELETE("/v1/nu/functions/:functionId", nuFunctionHandler)
//
// 	// docker proxy endpoints - testing only
// 	router.HEAD("/v1/docker/*command", handlers.DockerProxy)
// 	router.OPTIONS("/v1/docker/*command", handlers.DockerProxy)
// 	router.GET("/v1/docker/*command", handlers.DockerProxy)
// 	router.PUT("/v1/docker/*command", handlers.DockerProxy)
// 	router.POST("/v1/docker/*command", handlers.DockerProxy)
// 	router.DELETE("/v1/docker/*command", handlers.DockerProxy)
// 	router.PATCH("/v1/docker/*command", handlers.DockerProxy)
//
// 	// |*| : Initialise HTTP webserver
// 	//
// 	server := &http.Server{
// 		Addr:           ":9090",
// 		Handler:        router,
// 		ReadTimeout:    10 * time.Second,
// 		WriteTimeout:   10 * time.Second,
// 		MaxHeaderBytes: 1 << 20,
// 	}
//
// 	system := &System{db, server}
// 	log.Printf("System: %v", system)
//
// 	log.Fatal(server.ListenAndServe())
//
// }

func main() {

	// |*| : Initialise BoltDB
	//

	db := boltdb.NewDB()
	defer db.Close()

	// |*| : Initialise Httprouter
	//
	router := gin.Default()

	// health endpoints
	router.GET("/v1/health/ping", handlers.Ping)
	router.GET("/v1/health/echo", handlers.Echo)

	// nu endpoints - functions
	router.GET("/v1/nu/functions", handlers.ListFunctions(db))
	router.GET("/v1/nu/functions/:functionId", handlers.GetFunction(db))
	router.PUT("/v1/nu/functions", handlers.RegisterFunction(db))
	router.POST("/v1/nu/functions/:functionId", handlers.UpdateFunction(db))
	router.DELETE("/v1/nu/functions/:functionId", handlers.DeleteFunction(db))
	// nu endpoints - functions/code-archive
	router.POST("/v1/nu/functions/:functionId/code-archive", handlers.UpdateFunctionData(db))

	// docker proxy endpoints - testing only
	router.HEAD("/v1/docker/*command", handlers.DockerProxy)
	router.OPTIONS("/v1/docker/*command", handlers.DockerProxy)
	router.GET("/v1/docker/*command", handlers.DockerProxy)
	router.PUT("/v1/docker/*command", handlers.DockerProxy)
	router.POST("/v1/docker/*command", handlers.DockerProxy)
	router.DELETE("/v1/docker/*command", handlers.DockerProxy)
	router.PATCH("/v1/docker/*command", handlers.DockerProxy)

	// |*| : Initialise HTTP webserver
	//
	server := &http.Server{
		Addr:           ":9090",
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	system := &System{db, server}
	log.Printf("System: %v", system)

	log.Fatal(server.ListenAndServe())

}
