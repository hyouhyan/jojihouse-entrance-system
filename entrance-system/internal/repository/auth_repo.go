package repository

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthRepository struct {
	db *mongo.Database
}

func NewAuthRepository(db *mongo.Database) *AuthRepository {
	return &AuthRepository{db: db}
}
