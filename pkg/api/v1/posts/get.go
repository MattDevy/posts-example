package posts

import (
	"net/http"
	"strconv"

	"github.com/MattDevy/posts-example/pkg/models"
	"github.com/MattDevy/posts-example/pkg/otgorm"
	"github.com/gin-gonic/gin"
)

func (p *PostsAPI) GetPost(c *gin.Context) {
	postID, err := strconv.Atoi(c.Param("post_id"))
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}
	post := &models.Post{}

	if err := p.db.Scopes(
		otgorm.WithSpanFromContext(c.Request.Context()),
	).First(post, uint(postID)).Error; err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, post)
}
