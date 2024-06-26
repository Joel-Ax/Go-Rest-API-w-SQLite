package routes

import (
	"net/http"
	"strconv"

	"example.com/rest-api/models"
	"github.com/gin-gonic/gin"
)

func registerForEvent(context *gin.Context) {
	userId := context.GetInt64("user_id")
	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"Could not parse event id": err.Error()})
		return
	}
	event, err := models.GetEventByID(eventId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"Could not fetch event": err.Error()})
		return
	}

	err = event.RegisterUser(userId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"Could not register user": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"User registered for event": event})
}

func cancelResgistration(context *gin.Context) {
	userId := context.GetInt64("user_id")
	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"Could not parse event id": err.Error()})

		var event models.Event
		event.ID = eventId

		err = event.CancelRegistration(userId)
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"Could not cancel registration": err.Error()})
			return
		}
		context.JSON(http.StatusOK, gin.H{"User canceled registration for event": event})
	}
}
