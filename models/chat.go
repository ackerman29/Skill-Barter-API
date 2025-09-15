package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type ChatSession struct {
    ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    Users     []string           `bson:"users" json:"users"` // [userA, userB]
    CreatedAt primitive.DateTime `bson:"createdAt" json:"createdAt"`
}
