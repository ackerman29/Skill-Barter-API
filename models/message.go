package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Message struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ChatID    primitive.ObjectID `bson:"chatId" json:"chatId"`
	Sender    string             `bson:"sender" json:"sender"`       // sender email
	Content   string             `bson:"content" json:"content"`     // message text
	CreatedAt primitive.DateTime `bson:"createdAt" json:"createdAt"` // timestamp
}
