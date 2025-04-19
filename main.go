package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Note struct to represent a note
type Note struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

// In-memory storage for notes (for demo purposes only)
var notes = []Note{}
var nextID = 1

func main() {
	router := gin.Default()

	// Define routes
	router.GET("/notes", getNotes)
	router.POST("/notes", createNote)
	router.GET("/notes/:id", getNoteByID)
	router.PUT("/notes/:id", updateNote)
	router.DELETE("/notes/:id", deleteNote)

	// Run the server
	router.Run(":8080")
}

// getNotes handles GET requests to /notes to retrieve all notes
func getNotes(c *gin.Context) {
	c.JSON(http.StatusOK, notes)
}

// createNote handles POST requests to /notes to create a new note
func createNote(c *gin.Context) {
	var newNote Note
	if err := c.BindJSON(&newNote); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newNote.ID = nextID
	notes = append(notes, newNote)
	nextID++
	c.JSON(http.StatusCreated, newNote)
}

// getNoteByID handles GET requests to /notes/:id to retrieve a specific note
func getNoteByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid note ID"})
		return
	}

	for _, note := range notes {
		if note.ID == id {
			c.JSON(http.StatusOK, note)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
}

// updateNote handles PUT requests to /notes/:id to update an existing note
func updateNote(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid note ID"})
		return
	}

	var updatedNote Note
	if err := c.BindJSON(&updatedNote); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for i, note := range notes {
		if note.ID == id {
			updatedNote.ID = id // Ensure the ID remains the same
			notes[i] = updatedNote
			c.JSON(http.StatusOK, updatedNote)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
}

// deleteNote handles DELETE requests to /notes/:id to delete a note
func deleteNote(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid note ID"})
		return
	}

	for i, note := range notes {
		if note.ID == id {
			notes = append(notes[:i], notes[i+1:]...)
			c.JSON(http.StatusNoContent, nil) // 204 No Content for successful deletion
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
}
