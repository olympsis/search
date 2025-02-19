package search

import (
	"context"

	"github.com/olympsis/models"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Service struct {
	AuthCol *mongo.Collection
	UserCol *mongo.Collection
	Log     *logrus.Logger
}

func NewSearchService(l *logrus.Logger, a *mongo.Collection, u *mongo.Collection) *Service {
	return &Service{Log: l, AuthCol: a, UserCol: u}
}

func (s *Service) SearchUserByUUID(uuid string) (models.UserData, error) {

	// context/filter
	ctx := context.Background()
	filter := bson.M{"uuid": uuid}
	opts := options.FindOneOptions{}

	// find and decode auth user data
	var auth models.AuthUser
	err := s.AuthCol.FindOne(ctx, filter).Decode(&auth)
	if err != nil {
		return models.UserData{}, err
	}

	// find and decode user metadata
	var user models.User
	err = s.UserCol.FindOne(ctx, filter, &opts).Decode(&user)
	if err != nil {
		return models.UserData{}, err
	}

	imageURL := ""
	if user.ImageURL != nil {
		imageURL = *user.ImageURL
	}

	// create user data object
	userData := models.UserData{
		UUID:                   *auth.UUID,
		Username:               user.UserName,
		FirstName:              *auth.FirstName,
		LastName:               *auth.LastName,
		ImageURL:               imageURL,
		Visibility:             user.Visibility,
		NotificationDevices:    user.NotificationDevices,
		NotificationPreference: user.NotificationPreference,
	}

	// if user visibility is public display this data if not then don't
	if user.Visibility == "public" {
		userData.Bio = user.Bio
		userData.Clubs = user.Clubs
		userData.Sports = user.Sports
	}
	return userData, nil
}

func (s *Service) SearchUserByUsername(name string) (models.UserData, error) {

	// context/filter
	ctx := context.Background()
	filter := bson.M{"username": name}
	opts := options.FindOneOptions{}

	// find and decode user metadata
	var user models.User
	err := s.UserCol.FindOne(ctx, filter, &opts).Decode(&user)
	if err != nil {
		return models.UserData{}, err
	}

	filter = bson.M{"uuid": user.UUID}

	// return only uuid, first name and last name

	// find and decode auth user data
	var auth models.AuthUser
	err = s.AuthCol.FindOne(ctx, filter).Decode(&auth)
	if err != nil {
		return models.UserData{}, err
	}

	imageURL := ""
	if user.ImageURL != nil {
		imageURL = *user.ImageURL
	}

	// create user data object
	userData := models.UserData{
		UUID:                   *auth.UUID,
		Username:               user.UserName,
		FirstName:              *auth.FirstName,
		LastName:               *auth.LastName,
		ImageURL:               imageURL,
		Visibility:             user.Visibility,
		NotificationDevices:    user.NotificationDevices,
		NotificationPreference: user.NotificationPreference,
	}

	// if user visibility is public display this data if not then don't
	if user.Visibility == "public" {
		userData.Bio = user.Bio
		userData.Clubs = user.Clubs
		userData.Sports = user.Sports
	}
	return userData, nil
}
