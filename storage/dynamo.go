package storage

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Dynamo struct {
	Client *dynamodb.Client
}

var (
	ErrItemNotFound = fmt.Errorf("item not found")
)

func NewDynamo(profile string) (*Dynamo, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithSharedConfigProfile(profile),
		config.WithRegion(region),
	)
	if err != nil {
		return nil, err
	}
	client := dynamodb.NewFromConfig(cfg)
	return &Dynamo{
		Client: client,
	}, nil
}

func (d *Dynamo) GetItem(ctx context.Context, tableName, id string, item interface{}) error {
	idAttr, err := attributevalue.Marshal(id)
	out, err := d.Client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: &tableName,
		Key: map[string]types.AttributeValue{
			"id": idAttr,
		},
	})
	if err != nil {
		return err
	}
	if len(out.Item) == 0 {
		return ErrItemNotFound
	}
	fmt.Println(out.Item, err)
	return attributevalue.UnmarshalMap(out.Item, item)
}

func (d *Dynamo) PutItem(ctx context.Context, tableName string, item interface{}) error {
	itemMap, err := attributevalue.MarshalMap(item)
	if err != nil {
		return err
	}
	_, err = d.Client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: &tableName,
		Item:      itemMap,
	})
	return err
}

func (d *Dynamo) List(ctx context.Context, tableName string, obj interface{}) error {
	res, err := d.Client.Scan(ctx, &dynamodb.ScanInput{
		TableName: &tableName,
	})
	if err != nil {
		return err
	}
	return attributevalue.UnmarshalListOfMaps(res.Items, obj)
}

func (d *Dynamo) DeleteItem(ctx context.Context, tableName, id string) error {
	idAttr, err := attributevalue.Marshal(id)
	_, err = d.Client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: &tableName,
		Key: map[string]types.AttributeValue{
			"id": idAttr,
		},
	})
	return err
}
