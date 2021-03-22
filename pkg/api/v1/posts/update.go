package posts

import (
	"net/http"
	"strconv"

	"github.com/MattDevy/posts-example/pkg/models"
	"github.com/MattDevy/posts-example/pkg/otgorm"
	"github.com/gin-gonic/gin"
)

// UpdatePost will update any non-null fields
// This will only include editable fields of the models.Post
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

	// Update the post
	updatedPost := &models.Post{ID: uint(postID)}
	if err := p.db.Scopes(
		otgorm.WithSpanFromContext(c.Request.Context()),
	).Model(updatedPost).Updates(post).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	// Get the post
	if err := p.db.Scopes(
		otgorm.WithSpanFromContext(c.Request.Context()),
	).First(updatedPost, postID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	// Write back the complete post, with fields updated
	c.JSON(http.StatusOK, updatedPost)
}
