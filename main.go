package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Character struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Alias    string `json:"alias"`
	Superpower string `json:"superpower"`
}

var characters []Character
var nextID = 1

func main() {
	r := mux.NewRouter()

	// Inisialisasi data dummy
	characters = append(characters, Character{ID: nextID, Name: "Tony Stark", Alias: "Iron Man", Superpower: "Genius-level intellect"})
	nextID++
	characters = append(characters, Character{ID: nextID, Name: "Steve Rogers", Alias: "Captain America", Superpower: "Enhanced strength and agility"})
	nextID++

	// Endpoint REST API CRUD
	r.HandleFunc("/characters", getCharacters).Methods("GET")
	r.HandleFunc("/characters/{id}", getCharacterByID).Methods("GET")
	r.HandleFunc("/characters", createCharacter).Methods("POST")
	r.HandleFunc("/characters/{id}", updateCharacter).Methods("PUT")
	r.HandleFunc("/characters/{id}", deleteCharacter).Methods("DELETE")

	log.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

// GET /characters
func getCharacters(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(characters)
}

// GET /characters/{id}
func getCharacterByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	for _, c := range characters {
		if c.ID == id {
			json.NewEncoder(w).Encode(c)
			return
		}
	}
	http.Error(w, "Character not found", http.StatusNotFound)
}

// POST /characters
func createCharacter(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var newCharacter Character
	if err := json.NewDecoder(r.Body).Decode(&newCharacter); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	newCharacter.ID = nextID
	nextID++
	characters = append(characters, newCharacter)
	json.NewEncoder(w).Encode(newCharacter)
}

// PUT /characters/{id}
func updateCharacter(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	var updatedCharacter Character
	if err := json.NewDecoder(r.Body).Decode(&updatedCharacter); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	for i, c := range characters {
		if c.ID == id {
			updatedCharacter.ID = id
			characters[i] = updatedCharacter
			json.NewEncoder(w).Encode(updatedCharacter)
			return
		}
	}
	http.Error(w, "Character not found", http.StatusNotFound)
}

// DELETE /characters/{id}
func deleteCharacter(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	for i, c := range characters {
		if c.ID == id {
			characters = append(characters[:i], characters[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	http.Error(w, "Character not found", http.StatusNotFound)
}
