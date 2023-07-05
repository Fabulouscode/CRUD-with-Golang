package services

import (
	"context"
	"errors"
	"github/fabulousCode/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserServiceImpl struct {
	userCollection *mongo.Collection
	ctx            context.Context
}

func NewUserServiceImpl(userCollection *mongo.Collection, ctx context.Context) *UserServiceImpl {
	return &UserServiceImpl{
		userCollection: userCollection,
		ctx:            ctx,
	}
}

func (s *UserServiceImpl) CreateUser(user *models.User) error {
	_, err := s.userCollection.InsertOne(s.ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserServiceImpl) GetAllUsers() ([]*models.User, error) {
	cursor, err := s.userCollection.Find(s.ctx, nil)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(s.ctx)

	var users []*models.User
	for cursor.Next(s.ctx) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (s *UserServiceImpl) GetUser(userID *string) (*models.User, error) {
	var user models.User
	err := s.userCollection.FindOne(s.ctx, userID).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (s *UserServiceImpl) UpdateUser(user *models.User) error {
	_, err := s.userCollection.ReplaceOne(s.ctx, bson.M{"name": user.Name}, user)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserServiceImpl) DeleteUser(userID *string) error {
	_, err := s.userCollection.DeleteOne(s.ctx, userID)
	if err != nil {
		return err
	}
	return nil
}
