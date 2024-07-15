package storage

import "context"

type Storage interface {
	GetItem(ctx context.Context, tableName, id string, item interface{}) error
	PutItem(ctx context.Context, tableName string, item interface{}) error
	List(ctx context.Context, tableName string, objs interface{}) error
	DeleteItem(ctx context.Context, tableName, id string) error
}
