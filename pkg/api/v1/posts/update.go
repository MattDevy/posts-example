package posts

import (
	"net/http"
	"strconv"
	"time"

	"github.com/MattDevy/posts-example/pkg/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (p *PostsAPI) UpdatePost(c *gin.Context) {
	postID, err := strconv.Atoi(c.Param("post_id"))
	if err != nil {
		c.Status(http.StatusNotFound)
	}
	post := &models.Post{}
	if err := c.ShouldBindJSON(post); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	// Zero out fields we don't want user to set
	post.ID = 0
	post.CreatedAt = time.Time{}
	post.UpdatedAt = time.Time{}
	post.DeletedAt = gorm.DeletedAt{}

	updatedPost := &models.Post{ID: uint(postID)}
	// Update
	if err := p.db.Model(updatedPost).Updates(post).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	// Get data
	if err := p.db.First(updatedPost, postID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, updatedPost)
}
