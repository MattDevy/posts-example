package posts

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/MattDevy/posts-example/pkg/models"
	"github.com/MattDevy/posts-example/pkg/otgorm"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const (
	MAX_PAGE_SIZE int32 = 50
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

// Paginate gets the relevant information from the request and attaches it to the database query
func Paginate(r *ListPostsRequest, Offset *int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		pageSize := MAX_PAGE_SIZE
		if r.PageSize > 0 && r.PageSize < MAX_PAGE_SIZE {
			pageSize = r.PageSize
		}
		db = db.Limit(int(pageSize))

		if r.PageToken != "" {
			offset, err := strconv.Atoi(r.PageToken)
			if err != nil {
				db.AddError(err)
				return db
			}
			db = db.Offset(offset)
			*Offset = offset
		}

		return db
	}
}

func StartTime(r *ListPostsRequest) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if !r.StartTime.IsZero() {
			return db.Where("created_at > ?", r.StartTime)
		}
		return db
	}
}

func EndTime(r *ListPostsRequest) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if !r.EndTime.IsZero() {
			return db.Where("created_at <= ?", r.EndTime)
		}
		return db
	}
}

func Order(r *ListPostsRequest) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		switch r.OrderBy {
		case "", "created_at":
			return db.Order("created_at desc")
		default:
			return db
		}
	}
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
		StartTime(req),
		EndTime(req),
		otgorm.WithSpanFromContext(c.Request.Context()),
	).Model(&models.Post{}).Count(&count).Error; err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	// Do the request
	var offset int
	resultsArr := make([]models.Post, 0, req.PageSize)
	if err := p.db.Scopes(
		Paginate(req, &offset),
		StartTime(req),
		EndTime(req),
		Order(req),
		otgorm.WithSpanFromContext(c.Request.Context()),
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
