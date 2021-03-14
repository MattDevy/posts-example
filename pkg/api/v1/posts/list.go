package posts

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/MattDevy/posts-example/pkg/models"
	"github.com/gin-gonic/gin"
)

type ListPostsRequest struct {
	PageSize  int    `json:"page_size,omitempty" form:"page_size"`
	PageToken string `json:"page_token,omitempty" form:"page_token"`
}

const (
	MAX_PAGE_SIZE = 50
)

type ListPostsResponse struct {
	Posts         []models.Post `json:"posts,omitempty"`
	NextPageToken string        `json:"next_page_token,omitempty"`
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

	db := p.db

	offset := 0

	if req.PageToken != "" {
		q, err := url.QueryUnescape(req.PageToken)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		s := strings.Split(q, "#")
		if len(s) != 2 {
			c.Status(http.StatusBadRequest)
			return
		}
		offset, err = strconv.Atoi(s[1])
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		_, err = time.Parse(time.RFC3339Nano, s[0])
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		db = db.Offset(offset).Where("created_at < ?", s[0])
	}

	resultsArr := make([]models.Post, 0, req.PageSize)

	if err := db.Limit(req.PageSize).Order("created_at desc").Find(&resultsArr).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ListPostsResponse{Error: err})
		fmt.Println(err.Error())
		return
	}

	var nextPageToken string
	if len(resultsArr) == req.PageSize {
		nextPageToken = fmt.Sprintf("%v#%v", resultsArr[0].CreatedAt.Format(time.RFC3339Nano), offset+len(resultsArr))
	}

	c.JSON(http.StatusOK, ListPostsResponse{Posts: resultsArr, NextPageToken: url.QueryEscape(nextPageToken)})
}
