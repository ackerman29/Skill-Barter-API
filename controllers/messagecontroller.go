package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"temp/config"
	"temp/models"
	"temp/websocket"
)

// getMessageCollection returns the messages collection
func getMessageCollection() *mongo.Collection {
	return config.DB.Database("temp").Collection("messages")
}

// SendMessage handles sending a message in a chat
func SendMessage(c *gin.Context) {
	var body struct {
		ChatID  string `json:"chatId"`
		Content string `json:"content"`
	}

	if err := c.BindJSON(&body); err != nil || body.Content == "" || body.ChatID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "chatId and content are required"})
		return
	}

	chatID, err := primitive.ObjectIDFromHex(body.ChatID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chatId"})
		return
	}

	sender := c.MustGet("email").(string)

	msg := models.Message{
		ID:        primitive.NewObjectID(),
		ChatID:    chatID,
		Sender:    sender,
		Content:   body.Content,
		CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Insert message into DB
	_, err = getMessageCollection().InsertOne(ctx, msg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send message"})
		return
	}

	// Fetch chat to know participants
	chatCollection := config.DB.Database("temp").Collection("chats")
	var chat models.ChatSession
	err = chatCollection.FindOne(ctx, bson.M{"_id": chatID}).Decode(&chat)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Chat not found"})
		return
	}

	// Push real-time via WebSocket
	websocket.SendChatMessage(chat.Users, fmt.Sprintf("%s: %s", sender, body.Content))

	c.JSON(http.StatusCreated, msg)
}

// GetMessages fetches all messages for a chat
func GetMessages(c *gin.Context) {
	chatID := c.Param("chatId")
	objID, err := primitive.ObjectIDFromHex(chatID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chatId"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Find all messages for this chat
	cursor, err := getMessageCollection().Find(ctx, bson.M{"chatId": objID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch messages"})
		return
	}
	defer cursor.Close(ctx)

	var messages []models.Message
	if err := cursor.All(ctx, &messages); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse messages"})
		return
	}

	c.JSON(http.StatusOK, messages)
}
