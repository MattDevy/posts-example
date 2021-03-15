package main

import (
	"log"
	"os"
	"time"

	"github.com/MattDevy/posts-example/pkg/api/v1/posts"
	"github.com/MattDevy/posts-example/pkg/models"
	"github.com/gin-gonic/gin"
	cors "github.com/itsjamie/gin-cors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	BaseURL string = ":3000"
)

func main() {

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logger.Info, // Log level
			Colorful:      false,       // Disable color
		},
	)

	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Fatalf("Failed to open db, err: %v", err)
	}

	err = db.AutoMigrate(models.Post{})
	if err != nil {
		log.Fatal(err)
	}

	PostsAPI := posts.NewPostsAPI(db)

	router := gin.Default()
	cfg := cors.Config{
		Origins:         "*",
		Methods:         "GET, PUT, POST, DELETE",
		RequestHeaders:  "Origin, Authorization, Content-Type",
		ExposedHeaders:  "",
		MaxAge:          1 * time.Minute,
		Credentials:     true,
		ValidateHeaders: false,
	}
	router.Use(cors.Middleware(cfg))
	v1 := router.Group("/api/v1")
	PostsAPI.Register(v1)

	router.Run(BaseURL)
}
