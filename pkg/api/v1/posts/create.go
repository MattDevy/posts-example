package posts

import (
	"net/http"

	"github.com/MattDevy/posts-example/pkg/models"
	"github.com/MattDevy/posts-example/pkg/scopes"
	"github.com/gin-gonic/gin"
)

// CreatePostRequest includes all fields in a CreatePost request
type CreatePostRequest struct {
	Post models.Post `json:"post" binding:"required"`
}

// CreatePostResponse includes all fields in a CreatePost response
type CreatePostResponse struct {
	// ID of the created post
	ID uint `json:"id,omitempty"`
	// Error encountered when creating new post
	Error error `json:"error,omitempty"`
}

// CreatePost creates a new post, returning either the ID of the created post or an error
func (p *PostsAPI) CreatePost(c *gin.Context) {
	post := &models.Post{}
	if err := c.ShouldBindJSON(&post.PostData); err != nil {
		c.JSON(http.StatusBadRequest, CreatePostResponse{Error: err})
		return
	}

	// Create new post in database
	if err := p.db.Scopes(
		scopes.WithSpanFromContext(c.Request.Context()),
	).Create(post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, CreatePostResponse{Error: err})
		return
	}

	c.JSON(http.StatusCreated, CreatePostResponse{ID: post.ID})
}
