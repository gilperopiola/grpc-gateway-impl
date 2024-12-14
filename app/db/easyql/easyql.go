package easyql

import (
	"context"
	"fmt"

	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/models"
)

// Type constraint including all models
type AllModels interface {
	models.User | models.Group
}

// Creates a new record on the DB
func Create[T AllModels](ctx context.Context, record T, db core.InnerDB) (*T, error) {
	if err := db.WithContext(ctx).Create(&record).Error(); err != nil {
		return nil, fmt.Errorf("error in easyql.Create: %w", err)
	}
	return &record, nil
}

// Gets a record on the DB by ID
func Get[T AllModels](ctx context.Context, id any, record *T, db core.InnerDB) (*T, error) {
	if err := db.WithContext(ctx).First(record, id).Error(); err != nil {
		return nil, fmt.Errorf("error in easyql.Get: %w", err)
	}
	return record, nil
}

// Updates a record on the DB by ID
func Update[T AllModels](ctx context.Context, record *T, db core.InnerDB) (*T, error) {
	if err := db.WithContext(ctx).Save(record).Error(); err != nil {
		return nil, fmt.Errorf("error in easyql.Update: %w", err)
	}
	return record, nil
}

// Deletes a record on the DB by ID
func Delete[T AllModels](ctx context.Context, id any, db core.InnerDB) error {
	if err := db.WithContext(ctx).Delete(new(T), id).Error(); err != nil {
		return fmt.Errorf("error in easyql.Delete: %w", err)
	}
	return nil
}
