package posts

import (
	"net/http"
	"strconv"

	"github.com/MattDevy/posts-example/pkg/models"
	"github.com/MattDevy/posts-example/pkg/scopes"
	"github.com/gin-gonic/gin"
)

func (p *PostsAPI) DeletePost(c *gin.Context) {
	imageID, err := strconv.Atoi(c.Param("post_id"))
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	if err := p.db.Scopes(
		scopes.WithSpanFromContext(c.Request.Context()),
	).Delete(&models.Post{}, uint(imageID)).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.Status(http.StatusOK)
}
