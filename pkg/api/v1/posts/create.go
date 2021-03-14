package posts

import (
	"net/http"
	"time"

	"github.com/MattDevy/posts-example/pkg/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CreatePostRequest struct {
	Post models.Post `json:"post" binding:"required"`
}

type CreatePostResponse struct {
	ID    uint  `json:"id,omitempty"`
	Error error `json:"error,omitempty"`
}

func (p *PostsAPI) CreatePost(c *gin.Context) {
	post := &models.Post{}
	if err := c.ShouldBindJSON(post); err != nil {
		c.JSON(http.StatusBadRequest, CreatePostResponse{Error: err})
		return
	}

	// Reset these incase set by user
	post.ID = 0
	post.CreatedAt = time.Time{}
	post.UpdatedAt = time.Time{}
	post.DeletedAt = gorm.DeletedAt{}

	if err := p.db.Create(post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, CreatePostResponse{Error: err})
		return
	}

	c.JSON(http.StatusCreated, CreatePostResponse{ID: post.ID})
}
