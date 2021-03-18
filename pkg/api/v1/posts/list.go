package posts

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/MattDevy/posts-example/pkg/models"
	"github.com/gin-gonic/gin"
)

type ListPostsRequest struct {
	PageSize  int32  `json:"page_size,omitempty" form:"page_size"`
	PageToken string `json:"page_token,omitempty" form:"page_token"`
	OrderBy   string `json:"order_by,omitempty" form:"order_by"`
}

const (
	MAX_PAGE_SIZE = 50
)

type ListPostsResponse struct {
	Posts         []models.Post `json:"posts,omitempty"`
	NextPageToken string        `json:"next_page_token,omitempty"`
	TotalSize     int32         `json:"total_size,omitempty"`
	Error         error         `json:"error,omitempty"`
}

func (p *PostsAPI) ListPosts(c *gin.Context) {
	req := &ListPostsRequest{}
	if err := c.BindQuery(req); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	if req.PageSize == 0 || req.PageSize > MAX_PAGE_SIZE {
		req.PageSize = MAX_PAGE_SIZE
	}

	var count int64
	if err := p.db.Find(&models.Post{}).Count(&count).Error; err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	db := p.db

	offset := 0

	if req.PageToken != "" {
		offset, err := strconv.Atoi(req.PageToken)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		db = db.Offset(offset)
	}

	resultsArr := make([]models.Post, 0, req.PageSize)

	if err := db.Limit(int(req.PageSize)).Order("created_at desc").Find(&resultsArr).Error; err != nil {
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
