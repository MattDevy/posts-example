package posts

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// PostsAPI contains a *gorm.DB which is used by the http methods
type PostsAPI struct {
	db *gorm.DB
}

// NewPostsAPI returns a new PostsAPI with the *gorm.DB
func NewPostsAPI(db *gorm.DB) *PostsAPI {
	return &PostsAPI{
		db: db,
	}
}

// Registers the PostsAPI methods with the gin.RouterGroup
func (p *PostsAPI) Register(r *gin.RouterGroup) {
	r.POST("/posts", p.CreatePost)
	r.GET("/posts", p.ListPosts)
	r.GET("/posts/:post_id", p.GetPost)
	r.DELETE("/posts/:post_id", p.DeletePost)
	r.PUT("/posts/:post_id", p.UpdatePost)
}
