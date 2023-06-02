package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/totegamma/concurrent/x/auth"
	"github.com/totegamma/concurrent/x/core"
	"github.com/totegamma/concurrent/x/util"
)

func main() {

    fmt.Print(concurrentBanner)

    e := echo.New()

    config := util.Config{}
    err := config.Load("/etc/concurrent/config.yaml")
    if err != nil {
        e.Logger.Fatal(err)
    }

    log.Print("Config loaded! I am: ", config.CCAddr)

    db, err := gorm.Open(postgres.Open(config.Dsn), &gorm.Config{})
    if err != nil {
        log.Println("failed to connect database");
        panic("failed to connect database")
    }

    // Migrate the schema
    log.Println("start migrate")
    db.AutoMigrate(
        &core.Message{},
        &core.Character{},
        &core.Association{},
        &core.Stream{},
        &core.Host{},
        &core.Entity{},
    )

    rdb := redis.NewClient(&redis.Options{
        Addr:     config.RedisAddr,
        Password: "", // no password set
        DB:       0,  // use default DB
    })

    agent := SetupAgent(db, rdb, config)

    socketHandler := SetupSocketHandler(rdb, config)
    messageHandler := SetupMessageHandler(db, rdb, config)
    characterHandler := SetupCharacterHandler(db, config)
    associationHandler := SetupAssociationHandler(db, rdb, config)
    streamHandler := SetupStreamHandler(db, rdb, config)
    hostHandler := SetupHostHandler(db, config)
    entityHandler := SetupEntityHandler(db, config)
    authHandler := SetupAuthHandler(db, config)

    e.HideBanner = true
    e.Use(middleware.CORS())
    e.Use(middleware.Logger())
    e.Use(middleware.Recover())

    apiV1 := e.Group("/api/v1")
    apiV1.GET("/messages/:id", messageHandler.Get)
    apiV1.GET("/characters", characterHandler.Get)
    apiV1.GET("/associations/:id", associationHandler.Get)
    apiV1.GET("/stream", streamHandler.Get)
    apiV1.GET("/stream/recent", streamHandler.Recent)
    apiV1.GET("/stream/list", streamHandler.List)
    apiV1.GET("/stream/range", streamHandler.Range)
    apiV1.GET("/socket", socketHandler.Connect)
    apiV1.GET("/host/:id", hostHandler.Get) //TODO deprecated. remove later
    apiV1.GET("/host", hostHandler.Profile)
    apiV1.GET("/host/list", hostHandler.List)
    apiV1.GET("/entity/:id", entityHandler.Get)
    apiV1.GET("/entity/list", entityHandler.List)
    apiV1.GET("/auth/claim", authHandler.Claim)

    apiV1R := apiV1.Group("", auth.JWT)
    apiV1R.POST("/messages", messageHandler.Post)
    apiV1R.DELETE("/messages", messageHandler.Delete)
    apiV1R.PUT("/characters", characterHandler.Put)
    apiV1R.POST("/associations", associationHandler.Post)
    apiV1R.DELETE("/associations", associationHandler.Delete)
    apiV1R.PUT("/stream", streamHandler.Put)
    apiV1R.POST("/stream/checkpoint", streamHandler.Checkpoint)
    apiV1R.PUT("/host", hostHandler.Upsert)
    apiV1R.POST("/host/hello", hostHandler.Hello)
    apiV1R.POST("/entity", entityHandler.Post)
    apiV1R.GET("/admin/sayhello/:fqdn", hostHandler.SayHello)

    e.GET("/*", spa)
    e.GET("/health", func(c echo.Context) (err error) {
        return c.String(http.StatusOK, "ok")
    })

    agent.Boot()

    e.Logger.Fatal(e.Start(":8000"))
}

func spa(c echo.Context) error {
    path := c.Request().URL.Path
    fpath := filepath.Join("/etc/www/concurrent", path)
    if _, err := os.Stat(fpath); os.IsNotExist(err) {
        return c.File("/etc/www/concurrent/index.html")
    }
    return c.File(fpath)
}
