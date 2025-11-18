package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"task_manager/data"
	"task_manager/models"
)

func GetTasks(c *gin.Context) {
	tasks := data.GetAll()
	c.JSON(http.StatusOK, gin.H{"tasks": tasks})
}

func GetTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	task, err := data.GetByID(id)
	if err != nil {
		if err.Error() == "record not found" || err.Error() == "task not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, task)
}

func CreateTask(c *gin.Context) {
	var input models.Task
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	created, err := data.Create(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, created)
}

func UpdateTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var input models.Task
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updated, err := data.Update(id, input)
	if err != nil {
		if err.Error() == "task not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, updated)
}

func DeleteTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := data.Delete(id); err != nil {
		if err.Error() == "task not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "task deleted successfully"})
}
