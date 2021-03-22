package posts

import (
	"net/http"
	"strconv"

	"github.com/MattDevy/posts-example/pkg/models"
	"github.com/MattDevy/posts-example/pkg/scopes"
	"github.com/gin-gonic/gin"
)

// GetPost returns the post with the requested ID
func (p *PostsAPI) GetPost(c *gin.Context) {
	postID, err := strconv.Atoi(c.Param("post_id"))
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	// Get post from database
	post := &models.Post{}
	if err := p.db.Scopes(
		scopes.WithSpanFromContext(c.Request.Context()),
	).First(post, uint(postID)).Error; err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, post)
}
