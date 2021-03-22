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

// ListPostsRequest contains the fields for the ListPosts endpoint
type ListPostsRequest struct {
	// PageSize will be limited to a maximum of MAX_PAGE_SIZE, 50.
	// This is the maximum amount of results to be returned per page
	PageSize int32 `json:"page_size,omitempty" form:"page_size"`
	// PageToken returned from the previous ListPostsResponse indicates which page to
	// return next
	PageToken string `json:"page_token,omitempty" form:"page_token"`
	// OrderBy specifies what field the results are ordered by
	OrderBy string `json:"order_by,omitempty" form:"order_by"`
	// StartTime specifies the query start time, this is exclusive
	StartTime time.Time `json:"start_time,omitempty" form:"start_time"`
	// EndTime specifies the query end time, this is inclusive
	EndTime time.Time `json:"end_time,omitempty" form:"end_time"`
}

// ListPostsResponse is the  response returned from the ListPosts endpoint
type ListPostsResponse struct {
	// Posts, the page of posts returned, maximum length 50
	Posts []models.Post `json:"posts,omitempty"`
	// NextPageToken if empty, no more results to be returned
	// else to be used in next request for next page
	NextPageToken string `json:"next_page_token,omitempty"`
	// TotalSize is the total number of results, not just in this page
	TotalSize int32 `json:"total_size,omitempty"`
	// Error is non-empty if an error has occurred
	Error error `json:"error,omitempty"`
}

// ListPosts returns a list of the posts that match the query
func (p *PostsAPI) ListPosts(c *gin.Context) {
	// Get the parameters
	req := &ListPostsRequest{}
	if err := c.BindQuery(req); err != nil {
		c.JSON(http.StatusBadRequest, ListPostsResponse{Error: err})
		return
	}

	// Limit page size
	if req.PageSize > int32(MAX_PAGE_SIZE) {
		req.PageSize = int32(MAX_PAGE_SIZE)
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
		scopes.Paginate(int(req.PageSize), req.PageToken, &offset),
		scopes.StartTime(req.StartTime),
		scopes.EndTime(req.EndTime),
		scopes.Order(req.OrderBy),
		scopes.WithSpanFromContext(c.Request.Context()),
	).Find(&resultsArr).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ListPostsResponse{Error: err})
		return
	}

	// Create the next page token
	var nextPageToken string
	if len(resultsArr) == int(req.PageSize) {
		nextPageToken = fmt.Sprint(offset + len(resultsArr))
	}

	// Write back response
	c.JSON(http.StatusOK,
		ListPostsResponse{
			Posts:         resultsArr,
			NextPageToken: nextPageToken,
			TotalSize:     int32(count),
		},
	)
}
