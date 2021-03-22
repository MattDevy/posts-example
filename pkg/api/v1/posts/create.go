package posts

import (
	"net/http"

	"github.com/MattDevy/posts-example/pkg/models"
	"github.com/MattDevy/posts-example/pkg/otgorm"
	"github.com/gin-gonic/gin"
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
	if err := c.ShouldBindJSON(&post.PostData); err != nil {
		c.JSON(http.StatusBadRequest, CreatePostResponse{Error: err})
		return
	}

	if err := p.db.Scopes(
		otgorm.WithSpanFromContext(c.Request.Context()),
	).Create(post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, CreatePostResponse{Error: err})
		return
	}

	c.JSON(http.StatusCreated, CreatePostResponse{ID: post.ID})
}
