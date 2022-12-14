package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/pallat/todoapi/router"
	"github.com/pallat/todoapi/store"
	"github.com/pallat/todoapi/todo"
)

// var (
// 	buildcommit = "dev"
// 	buildtime   = time.Now().String()
// )

func main() {
	err := godotenv.Load("local.env")
	if err != nil {
		log.Printf("please consider environment variables: %s\n", err)
	}

	db, err := gorm.Open(sqlite.Open(os.Getenv("DB_CONN")), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	if err := db.AutoMigrate(&todo.Todo{}); err != nil {
		log.Println("auto migrate db", err)
	}

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://mongoadmin:secret@localhost:27017"))
	if err != nil {
		panic("failed to connect database")
	}
	collection := client.Database("myapp").Collection("todos")

	r := router.NewFiberRouter()

	// r := router.NewMyRouter()

	// r.GET("/healthz", func(c *gin.Context) {
	// 	c.Status(200)
	// })
	// r.GET("/limitz", limitedHandler)
	// r.GET("/x", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{
	// 		"buildcommit": buildcommit,
	// 		"buildtime":   buildtime,
	// 	})
	// })

	// gormStore := store.NewGormStore(db)
	mongoStore := store.NewMongoDBStore(collection)

	handler := todo.NewTodoHandler(mongoStore)
	r.POST("/todos", handler.NewTask)
	// r.GET("/todos", handler.List)
	// r.DELETE("/todos/:id", handler.Remove)

	// ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	// defer stop()

	// s := &http.Server{
	// 	Addr:           ":" + os.Getenv("PORT"),
	// 	Handler:        r,
	// 	ReadTimeout:    10 * time.Second,
	// 	WriteTimeout:   10 * time.Second,
	// 	MaxHeaderBytes: 1 << 20,
	// }

	// go func() {
	if err := r.Listen(os.Getenv("PORT")); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}
	// }()

	// <-ctx.Done()
	// stop()
	// fmt.Println("shutting down gracefully, press Ctrl+C again to force")

	// timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()

	// if err := s.Shutdown(timeoutCtx); err != nil {
	// 	fmt.Println(err)
	// }
}

// var limiter = rate.NewLimiter(5, 5)

// func limitedHandler(c *gin.Context) {
// 	if !limiter.Allow() {
// 		c.AbortWithStatus(http.StatusTooManyRequests)
// 		return
// 	}
// 	c.JSON(200, gin.H{
// 		"message": "pong",
// 	})
// }
