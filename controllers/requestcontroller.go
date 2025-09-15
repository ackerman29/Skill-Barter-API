
package controllers

import (
	"context"
	"net/http"
	"time"
	"temp/websocket"
	"fmt"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"temp/config"
	"temp/models"
)

func getRequestCollection() *mongo.Collection {
	return config.DB.Database("temp").Collection("requests")
}

func SendSkillRequest(c *gin.Context) {
	var request models.SkillRequest

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	fromEmail := c.MustGet("email").(string)

	// Get sender's name
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var sender models.User
	err := config.DB.Database("temp").Collection("users").FindOne(ctx, bson.M{"email": fromEmail}).Decode(&sender)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Sender not found"})
		return
	}

	request.ID = primitive.NewObjectID()
	request.FromEmail = fromEmail
	request.FromName = sender.Name
	request.Status = "pending"
	request.CreatedAt = primitive.NewDateTimeFromTime(time.Now())

	_, err = getRequestCollection().InsertOne(ctx, request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not send request"})
		return
	}
	fmt.Println("DEBUG: Skill in request:", request.Skill)

	websocket.SendToUser(request.ToEmail, fmt.Sprintf("Hey! You got a new skill request from %s for: %s", sender.Name, request.Skill))

	c.JSON(http.StatusCreated, gin.H{"message": "Skill request sent"})
}


func RespondToSkillRequest(c *gin.Context) {
	var body struct {
		FromName string `json:"fromName"`
		Status   string `json:"status"` // "accepted" or "rejected"
	}
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	toEmail := c.MustGet("email").(string)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Find the pending request
	filter := bson.M{
		"fromName": body.FromName,
		"toEmail":  toEmail,
		"status":   "pending",
	}

	var req models.SkillRequest
	err := getRequestCollection().FindOne(ctx, filter).Decode(&req)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Request not found"})
		return
	}

	// Update request status
	update := bson.M{"$set": bson.M{"status": body.Status}}
	_, err = getRequestCollection().UpdateOne(ctx, filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update request"})
		return
	}

	if body.Status == "accepted" {
		chatCollection := config.DB.Database("temp").Collection("chats")

		// check if chat already exists
		var existing models.ChatSession
		err := chatCollection.FindOne(ctx, bson.M{"users": bson.M{"$all": []string{req.FromEmail, req.ToEmail}}}).Decode(&existing)

		if err == mongo.ErrNoDocuments {
			newChat := models.ChatSession{
				ID:        primitive.NewObjectID(),
				Users:     []string{req.FromEmail, req.ToEmail},
				CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
			}
			_, err = chatCollection.InsertOne(ctx, newChat)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create chat"})
				return
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Request updated"})
}


