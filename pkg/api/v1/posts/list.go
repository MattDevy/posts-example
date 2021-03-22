package posts

import (
	"fmt"
	"net/http"
	"time"

	"github.com/MattDevy/posts-example/pkg/models"
	"github.com/MattDevy/posts-example/pkg/scopes"
	"github.com/gin-gonic/gin"
)

const (
	MAX_PAGE_SIZE int = 50
)

type ListPostsRequest struct {
	PageSize  int32     `json:"page_size,omitempty" form:"page_size"`
	PageToken string    `json:"page_token,omitempty" form:"page_token"`
	OrderBy   string    `json:"order_by,omitempty" form:"order_by"`
	StartTime time.Time `json:"start_time,omitempty" form:"start_time"`
	EndTime   time.Time `json:"end_time,omitempty" form:"end_time"`
}

type ListPostsResponse struct {
	Posts         []models.Post `json:"posts,omitempty"`
	NextPageToken string        `json:"next_page_token,omitempty"`
	TotalSize     int32         `json:"total_size,omitempty"`
	Error         error         `json:"error,omitempty"`
}

func (p *PostsAPI) ListPosts(c *gin.Context) {
	// Get the parameters
	req := &ListPostsRequest{}
	if err := c.BindQuery(req); err != nil {
		c.JSON(http.StatusBadRequest, ListPostsResponse{Error: err})
		return
	}

	// Find the total number of results
	var count int64
	if err := p.db.Scopes(
		scopes.StartTime(req.StartTime),
		scopes.EndTime(req.EndTime),
		scopes.WithSpanFromContext(c.Request.Context()),
	).Model(&models.Post{}).Count(&count).Error; err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	// Do the request
	var offset int
	resultsArr := make([]models.Post, 0, req.PageSize)
	if err := p.db.Scopes(
		scopes.Paginate(int(req.PageSize), MAX_PAGE_SIZE, req.PageToken, &offset),
		scopes.StartTime(req.StartTime),
		scopes.EndTime(req.EndTime),
		scopes.Order(req.OrderBy),
		scopes.WithSpanFromContext(c.Request.Context()),
	).Find(&resultsArr).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ListPostsResponse{Error: err})
		return
	}

	var nextPageToken string
	if len(resultsArr) == int(req.PageSize) {
		nextPageToken = fmt.Sprint(offset + len(resultsArr))
	}

	c.JSON(http.StatusOK,
		ListPostsResponse{
			Posts:         resultsArr,
			NextPageToken: nextPageToken,
			TotalSize:     int32(count),
		},
	)
}
