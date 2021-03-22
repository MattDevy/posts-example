package posts

import (
	"net/http"
	"strconv"

	"github.com/MattDevy/posts-example/pkg/models"
	"github.com/MattDevy/posts-example/pkg/otgorm"
	"github.com/gin-gonic/gin"
)

func (p *PostsAPI) UpdatePost(c *gin.Context) {
	postID, err := strconv.Atoi(c.Param("post_id"))
	if err != nil {
		c.Status(http.StatusNotFound)
	}
	post := &models.Post{}
	if err := c.ShouldBindJSON(&post.PostData); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	updatedPost := &models.Post{ID: uint(postID)}
	// Update
	if err := p.db.Scopes(
		otgorm.WithSpanFromContext(c.Request.Context()),
	).Model(updatedPost).Updates(post).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	// Get data
	if err := p.db.Scopes(
		otgorm.WithSpanFromContext(c.Request.Context()),
	).First(updatedPost, postID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, updatedPost)
}
