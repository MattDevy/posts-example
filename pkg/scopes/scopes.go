package scopes

import (
	"context"
	"strconv"
	"time"

	"github.com/MattDevy/posts-example/pkg/otgorm"
	"gorm.io/gorm"
)

// Paginate gets the relevant information from the request and attaches it to the database query
func Paginate(pageSize, maxPageSize int, pageToken string, Offset *int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		ps := maxPageSize
		if pageSize > 0 && pageSize < maxPageSize {
			ps = pageSize
		}
		db = db.Limit(int(ps))

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

func StartTime(startTime time.Time) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if !startTime.IsZero() {
			return db.Where("created_at > ?", startTime)
		}
		return db
	}
}

func EndTime(endTime time.Time) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if !endTime.IsZero() {
			return db.Where("created_at <= ?", endTime)
		}
		return db
	}
}

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

// WithSpanCtx scoping middleware
func WithSpanFromContext(ctx context.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return otgorm.SetSpanToGorm(ctx, db)
	}
}
