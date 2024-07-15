package storage

import (
	"bytes"
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type S3 struct {
	Client *s3.Client
}

const region = "us-west-1"
const bucket = "shoppinglistapi"

func NewS3(profile string) (*S3, error) {
	client, err := session(profile)
	if err != nil {
		return nil, err
	}
	err = VerifyCollections(client, []string{"users", "collections"})
	if err != nil {
		return nil, err
	}
	return &S3{
		Client: client,
	}, nil
}

func session(profile string) (*s3.Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithSharedConfigProfile(profile),
		config.WithRegion(region),
	)
	if err != nil {
		return nil, err
	}
	svc := s3.NewFromConfig(cfg)
	return svc, nil
}

func (s *S3) Get(collection, query string) ([]byte, error) {
	out, err := s.Client.SelectObjectContent(context.TODO(), &s3.SelectObjectContentInput{
		Bucket:         aws.String(bucket),
		Key:            &collection,
		Expression:     aws.String(query),
		ExpressionType: "SQL",
		InputSerialization: &types.InputSerialization{
			JSON: &types.JSONInput{
				Type: types.JSONTypeDocument,
			},
		},
		OutputSerialization: &types.OutputSerialization{
			JSON: &types.JSONOutput{
				RecordDelimiter: aws.String("\n"),
			},
		},
	})
	if err != nil {
		return nil, err
	}

	stream := out.GetStream()
	defer stream.Close()
	for event := range stream.Events() {
		if err = stream.Err(); err != nil {
			return nil, err
		}
		switch ev := event.(type) {
		case *types.SelectObjectContentEventStreamMemberRecords:
			return ev.Value.Payload, nil
		}
	}
	return nil, err
}

func (s *S3) List(collection, query string, f func([]byte) error) error {
	out, err := s.Client.SelectObjectContent(context.TODO(), &s3.SelectObjectContentInput{
		Bucket:         aws.String(bucket),
		Key:            &collection,
		Expression:     aws.String(query),
		ExpressionType: "SQL",
		InputSerialization: &types.InputSerialization{
			JSON: &types.JSONInput{
				Type: types.JSONTypeDocument,
			},
		},
		OutputSerialization: &types.OutputSerialization{
			JSON: &types.JSONOutput{
				RecordDelimiter: aws.String("\n"),
			},
		},
	})
	if err != nil {
		return err
	}

	stream := out.GetStream()
	defer stream.Close()
	for event := range stream.Events() {
		if err = stream.Err(); err != nil {
			return err
		}
		switch ev := event.(type) {
		case *types.SelectObjectContentEventStreamMemberRecords:
			if err = f(ev.Value.Payload); err != nil {
				return err
			}
		}
	}
	return err
}

func (s *S3) Insert(collection string, object []byte) error {
	_, err := s.Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(collection),
		Body:   bytes.NewReader(object),
	})
	return err
}

func VerifyCollections(client *s3.Client, collections []string) error {
	for _, collection := range collections {
		_, err := client.HeadObject(context.TODO(), &s3.HeadObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(fmt.Sprintf("%s.json", collection)),
		})
		if err != nil {
			_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
				Bucket: aws.String(bucket),
				Key:    aws.String(fmt.Sprintf("%s.json", collection)),
				Body:   bytes.NewReader([]byte(fmt.Sprintf(`{"%s":[]}`, collection))),
			})
			return err
		}
	}
	return nil
}
