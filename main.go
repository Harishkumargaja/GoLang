package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/supabase-community/supabase-go"
)

// Note struct to represent a note
type Note struct {
	ID        uuid.UUID  `json:"id"`
	Title     string     `json:"title" binding:"required"`
	Content   *string    `json:"content"` // Make content nullable
	CreatedAt string     `json:"created_at"`
	UserID    *uuid.UUID `json:"user_id"` // Assuming you might want to associate with a user later
}

var client *supabase.Client
var validate *validator.Validate

func main() {
	supabaseURL := "https://rojqpylyzaaryefgalrp.supabase.co"
	supabaseKey := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6InJvanFweWx5emFhcnllZmdhbHJwIiwicm9sZSI6ImFub24iLCJpYXQiOjE3NDIxODcwNjYsImV4cCI6MjA1Nzc2MzA2Nn0.4CbjARN73KTJVma3Yarf2CNT6FyPw2qGW9ENyv_f5Ns"

	if supabaseURL == "" || supabaseKey == "" {
		log.Fatal("SUPABASE_URL and SUPABASE_KEY environment variables must be set.")
		return
	}

	var err error
	client, err = supabase.NewClient(supabaseURL, supabaseKey, nil)
	if err != nil {
		log.Fatalf("Error creating Supabase client: %v", err)
		return
	}

	// user, err := client.Auth.SignIn(context.Background(), supabase.UserCredentials{
	// 	Email:    "user@example.com",
	// 	Password: "yourpassword",
	// })
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("Access Token: %s\n", user.AccessToken)

	validate = validator.New()

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

// getNotes handles GET requests to /notes to retrieve all notes from Supabase
func getNotes(c *gin.Context) {
	var results []Note // We'll unmarshal directly into our Note struct
	data, _, err := client.From("notes").Select("*", "", false).Execute()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch notes from Supabase"})
		log.Printf("Error fetching notes: %v", err)
		return
	}
	err = json.Unmarshal(data, &results)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unmarshal notes data"})
		log.Printf("Error unmarshalling notes: %v", err)
		return
	}
	c.JSON(http.StatusOK, results)
}

// createNote handles POST requests to /notes to create a new note in Supabase
func createNote(c *gin.Context) {
	var newNote Note
	if err := c.BindJSON(&newNote); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Basic validation
	if err := validate.Struct(newNote); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var results []Note
	data, _, err := client.From("notes").Insert([]Note{newNote}, false, "", "", "").Execute()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create note in Supabase"})
		log.Printf("Error creating note: %v", err)
		return
	}
	err = json.Unmarshal(data, &results)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unmarshal created note data"})
		log.Printf("Error unmarshalling created note: %v", err)
		return
	}
	if len(results) > 0 {
		c.JSON(http.StatusCreated, &results[0])
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve created note"})
	}
}

// getNoteByID handles GET requests to /notes/:id to retrieve a specific note from Supabase
func getNoteByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid note ID"})
		return
	}

	var results []Note // Expecting a single result, but unmarshal into a slice
	data, _, err := client.From("notes").Select("*", "", false).Eq("id", id.String()).Execute()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch note from Supabase"})
		log.Printf("Error fetching note by ID: %v", err)
		return
	}
	err = json.Unmarshal(data, &results)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unmarshal note data"})
		log.Printf("Error unmarshalling note: %v", err)
		return
	}
	if len(results) > 0 {
		c.JSON(http.StatusOK, &results[0])
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
	}
}

// updateNote handles PUT requests to /notes/:id to update an existing note in Supabase
func updateNote(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid note ID"})
		return
	}

	var updatedNote Note
	if err := c.BindJSON(&updatedNote); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Basic validation (title is still required)
	if updatedNote.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title is required"})
		return
	}

	var results []Note
	data, _, err := client.From("notes").Update(updatedNote, "*", "").Eq("id", id.String()).Execute()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update note in Supabase"})
		log.Printf("Error updating note: %v", err)
		return
	}
	err = json.Unmarshal(data, &results)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unmarshal updated note data"})
		log.Printf("Error unmarshalling updated note: %v", err)
		return
	}
	if len(results) > 0 {
		c.JSON(http.StatusOK, &results[0])
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
	}
}

// deleteNote handles DELETE requests to /notes/:id to delete a note from Supabase
func deleteNote(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid note ID"})
		return
	}

	_, _, err = client.From("notes").Delete("", "").Eq("id", id.String()).Execute()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete note from Supabase"})
		log.Printf("Error deleting note: %v", err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
