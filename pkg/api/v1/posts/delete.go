package posts

import (
	"net/http"
	"strconv"

	"github.com/MattDevy/posts-example/pkg/models"
	"github.com/MattDevy/posts-example/pkg/scopes"
	"github.com/gin-gonic/gin"
)

// DeletePost removes a post with the corresponding ID
func (p *PostsAPI) DeletePost(c *gin.Context) {
	imageID, err := strconv.Atoi(c.Param("post_id"))
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	// Delete the post
	if err := p.db.Scopes(
		scopes.WithSpanFromContext(c.Request.Context()),
	).Delete(&models.Post{}, uint(imageID)).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.Status(http.StatusOK)
}
