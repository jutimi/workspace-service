package postgres_repository

import (
	"fmt"
	"workspace-server/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func findById(id uuid.UUID, field string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("%s = ?", field), id)
	}
}

func paginate(limit, offset int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Limit(limit).Offset(offset)
	}
}

func findByText(text, field string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		searchText := "%" + text + "%"
		return db.Where(fmt.Sprintf("%s LIKE ?", field), searchText)
	}
}

func findBySlice[T []uuid.UUID | []string | []int](data T, field string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("%s IN ?", field), data)
	}
}

func findByName(name, field string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		searchText := "%" + utils.Slugify(name) + "%"
		return db.Where(fmt.Sprintf("%s LIKE ?", field), searchText)
	}
}

func orById(id uuid.UUID, field string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Or(fmt.Sprintf("%s = ?", field), id)
	}
}

func orBySlice[T []uuid.UUID | []string | []int](data T, field string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Or(fmt.Sprintf("%s IN ?", field), data)
	}
}

func orByText(text, field string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		searchText := "%" + text + "%"
		return db.Or(fmt.Sprintf("%s LIKE ?", field), searchText)
	}
}

func buildLockQuery(query *gorm.DB, lockOption string) *gorm.DB {
	switch lockOption {
	case clause.LockingOptionsNoWait:
		query = query.Clauses(clause.Locking{
			Strength: clause.LockingStrengthUpdate,
			Options:  clause.LockingOptionsNoWait,
		})
	case clause.LockingOptionsSkipLocked:
		query = query.Clauses(clause.Locking{
			Strength: clause.LockingStrengthUpdate,
			Options:  clause.LockingOptionsSkipLocked,
		})
	default:
		query = query.Clauses(clause.Locking{
			Strength: clause.LockingStrengthUpdate,
		})
	}

	return query
}
