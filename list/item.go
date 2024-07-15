package list

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/stinkyfingers/shoppinglistapi/storage"
)

type ListItem struct {
	Item     Item   `json:"item" dynamodbav:"item"`
	Status   string `json:"status" dynamodbav:"status"`
	Quantity int    `json:"quantity" dynamodbav:"quantity"`
	Location string `json:"location" dynamodbav:"location"`
}

type Item struct {
	ID   string `json:"id,omitempty" dynamodbav:"id"`
	Name string `json:"name" dynamodbav:"name"`
}

const (
	itemsCollection = "sl_items"
)

var (
	ErrItemNotFound = fmt.Errorf("item not found")
)

func (i *Item) Insert(ctx context.Context, store storage.Storage) error {
	i.ID = uuid.New().String()
	err := store.PutItem(ctx, itemsCollection, i)
	return err
}

func (i *Item) Get(ctx context.Context, store storage.Storage) error {
	return store.GetItem(ctx, itemsCollection, i.ID, i)
}

func GetItems(ctx context.Context, store storage.Storage) ([]Item, error) {
	var items []Item
	err := store.List(ctx, itemsCollection, &items)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (i *Item) Delete(ctx context.Context, store storage.Storage) error {
	return store.DeleteItem(ctx, itemsCollection, i.ID)
}
