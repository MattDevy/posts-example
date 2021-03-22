package scopes

import (
	"context"
	"strconv"
	"time"

	"github.com/MattDevy/posts-example/pkg/otgorm"
	"gorm.io/gorm"
)

// Paginate gets the relevant information from the request and attaches it to the database query
func Paginate(pageSize int, pageToken string, Offset *int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		db = db.Limit(pageSize)

		if pageToken != "" {
			offset, err := strconv.Atoi(pageToken)
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

// StartTime filters query by start time if non-zero value
func StartTime(startTime time.Time) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if !startTime.IsZero() {
			return db.Where("created_at > ?", startTime)
		}
		return db
	}
}

// EndTime filters query by end time if non-zero value
func EndTime(endTime time.Time) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if !endTime.IsZero() {
			return db.Where("created_at <= ?", endTime)
		}
		return db
	}
}

// Order orders the query by "created_at", may support more in the future
func Order(orderBy string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		switch orderBy {
		case "", "created_at":
			return db.Order("created_at desc")
		default:
			return db
		}
	}
}

// WithSpanContext creates a new span, with the parent span from the context
// data is injected from using the query and the results into the span
func WithSpanFromContext(ctx context.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return otgorm.SetSpanToGorm(ctx, db)
	}
}
