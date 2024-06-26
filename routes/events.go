package routes

import (
	"net/http"
	"strconv"

	"example.com/rest-api/models"
	"github.com/gin-gonic/gin"
)

// Get all events
func getEvents(context *gin.Context) {
	events, err := models.GetAllEvents()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"Could not fetch events": err.Error()})
		return
	}
	context.JSON(http.StatusOK, events)
}

// Get one event
func getEvent(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"Could not parse event id": err.Error()})
		return
	}
	event, err := models.GetEventByID(id)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"Could not fetch event": err.Error()})
		return
	}
	context.JSON(http.StatusOK, event)

}

// Create an event
func createEvent(context *gin.Context) {

	var event models.Event

	err := context.BindJSON(&event)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"Could not parse request data": err.Error()})
		return
	}
	userId := context.GetInt64("user_id")
	event.UserID = userId

	err = event.Save()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"Could not create event": err.Error()})
		return
	}
	context.JSON(http.StatusCreated, gin.H{"Event was created:": event})
}

// Update an event
func updateEvent(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"Could not parse event id": err.Error()})
		return
	}

	userId := context.GetInt64("user_id")

	event, err := models.GetEventByID(id)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"Could not fetch event": err.Error()})
		return
	}
	if event.UserID != userId {
		context.JSON(http.StatusUnauthorized, gin.H{"Unauthorized": "You are not authorized to update this event"})
		return
	}

	var updatedEvent models.Event
	err = context.ShouldBindJSON(&updatedEvent)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"Could not parse request data": err.Error()})
		return
	}
	updatedEvent.ID = id
	err = updatedEvent.Update()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"Could not update event": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"Event was updated:": updatedEvent})
}

// Delete an event
func deleteEvent(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"Could not parse event id": err.Error()})
		return
	}

	userId := context.GetInt64("user_id")
	event, err := models.GetEventByID(id)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"Could not fetch event": err.Error()})
		return
	}
	if event.UserID != userId {
		context.JSON(http.StatusUnauthorized, gin.H{"Unauthorized": "You are not authorized to delete this event"})
		return
	}

	err = event.Delete()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"Could not delete event": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"Event was deleted": event})
}
