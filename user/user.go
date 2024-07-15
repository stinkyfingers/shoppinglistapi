package user

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/stinkyfingers/shoppinglistapi/storage"
)

type User struct {
	ID      string `json:"id" dynamodbav:"id"`
	Name    string `json:"name" dynamodbav:"name"`
	Email   string `json:"email" dynamodbav:"email"`
	Picture string `json:"picture,omitempty" dynamodbav:"picture"`
}

type UsersCollection struct {
	Users []User `json:"users"`
}

const (
	usersCollection = "sl_users"
)

var (
	ErrUserNotFound  = fmt.Errorf("user not found")
	ErrMultipleUsers = fmt.Errorf("multiple users found")
)

func (u *User) Insert(ctx context.Context, store storage.Storage) error {
	u.ID = uuid.New().String()
	fmt.Println(u)
	err := store.PutItem(ctx, usersCollection, u)
	return err
}

func (u *User) Get(ctx context.Context, store storage.Storage) error {
	return store.GetItem(ctx, usersCollection, u.ID, u)
}
func ListUsers(ctx context.Context, store storage.Storage) ([]User, error) {
	var users []User
	err := store.List(ctx, usersCollection, &users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (u *User) FindByEmail(ctx context.Context, store storage.Storage) error {
	users, err := ListUsers(ctx, store)
	if err != nil {
		return err
	}
	var foundUser *User
	for _, user := range users {
		if user.Email == u.Email {
			if foundUser != nil {
				return ErrMultipleUsers
			}
			foundUser = &user
		}
	}
	if foundUser == nil {
		return ErrUserNotFound
	}
	*u = *foundUser
	return nil
}
