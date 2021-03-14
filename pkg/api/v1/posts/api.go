package posts

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PostsAPI struct {
	db *gorm.DB
}

func NewPostsAPI(db *gorm.DB) *PostsAPI {
	return &PostsAPI{
		db: db,
	}
}

func (p *PostsAPI) Register(r *gin.RouterGroup) {
	r.POST("/posts", p.CreatePost)
	r.GET("/posts", p.ListPosts)
	r.GET("/posts/:post_id", p.GetPost)
	r.DELETE("/posts/:post_id", p.DeletePost)
}
