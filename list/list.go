package list

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/stinkyfingers/shoppinglistapi/storage"
)

type List struct {
	ID        string     `json:"id" dynamodbav:"id"`
	Name      string     `json:"name" dynamodbav:"name"`
	ListItems []ListItem `json:"listItems" dynamodbav:"listItems"`
	Store     string     `json:"store" dynamodbav:"store"`
}

const (
	listsCollection = "sl_lists"
)

var (
	ErrListNotFound = fmt.Errorf("list not found")
)

func (l *List) Insert(ctx context.Context, store storage.Storage) error {
	l.ID = uuid.New().String()
	err := store.PutItem(ctx, listsCollection, l)
	return err
}

func (l *List) Get(ctx context.Context, store storage.Storage) error {
	return store.GetItem(ctx, listsCollection, l.ID, l)
}

func GetLists(ctx context.Context, store storage.Storage) ([]List, error) {
	var lists []List
	err := store.List(ctx, listsCollection, &lists)
	if err != nil {
		return nil, err
	}
	return lists, nil
}

func (l *List) Delete(ctx context.Context, store storage.Storage) error {
	return store.DeleteItem(ctx, listsCollection, l.ID)
}
