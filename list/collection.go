package list

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/stinkyfingers/shoppinglistapi/storage"
	"github.com/stinkyfingers/shoppinglistapi/user"
)

type Collection struct {
	ID      string      `json:"id" dynamodbav:"id"`
	Creator user.User   `json:"creator" dynamodbav:"creator"`
	Users   []user.User `json:"users" dynmodbav:"users"`
	Name    string      `json:"name" dynamodbav:"name"`
	Lists   []List      `json:"lists" dynamodbav:"lists"`
}

const (
	collectionsCollection = "sl_collections"
)

var (
	ErrCollectionNotFound = fmt.Errorf("collection not found")
)

func (c *Collection) Upsert(ctx context.Context, store storage.Storage) error {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	err := store.PutItem(ctx, collectionsCollection, c)
	return err
}

func (c *Collection) Get(ctx context.Context, store storage.Storage) error {
	return store.GetItem(ctx, collectionsCollection, c.ID, c)
}

func GetCollections(ctx context.Context, store storage.Storage) ([]Collection, error) {
	var collections []Collection
	err := store.List(ctx, collectionsCollection, &collections)
	if err != nil {
		return nil, err
	}
	return collections, nil
}

func (c *Collection) Delete(ctx context.Context, store storage.Storage) error {
	return store.DeleteItem(ctx, collectionsCollection, c.ID)
}
