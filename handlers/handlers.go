package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"rest-service/db"
	"strconv"

	"github.com/gorilla/mux"
)

func GetPeople(rw http.ResponseWriter, r *http.Request) {
	postgres := db.GetPostgres()
	queryValues := r.URL.Query()
	search := queryValues.Get("seach")
	page, _ := strconv.Atoi(queryValues.Get("page"))
	limit, _ := strconv.Atoi(queryValues.Get("limit"))
	if page == 0 {
		page = 1
	}
	if limit == 0 {
		limit = 10
	}

	var people []db.PersonResponse
	var count int64

	query := postgres.Model(&db.Person{})

	if search != "" {
		query = query.Where("name LIKE ? OR surname LIKE ? OR patronymic LIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	query.Count(&count)
	query = query.Limit(limit).Offset((page - 1) * limit)

	query.Find(&people)

	response := map[string]interface{}{
		"result": people,
		"count":  count,
		"page":   page,
		"limit":  limit,
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(response)
}

func GetPerson(rw http.ResponseWriter, r *http.Request) {
	postgres := db.GetPostgres()
	params := mux.Vars(r)
	personID, err := strconv.Atoi(params["id"])
	if err != nil {
		SendErr(rw, "Invalid person ID")
		return
	}

	var person db.Person

	if err := postgres.First(&person, personID).Error; err != nil {
		SendErr(rw, "Person not found")
		return
	}
	response := map[string]interface{}{
		"result": db.PersonResponse{
			ID:          person.ID,
			Name:        person.Name,
			Surname:     person.Surname,
			Patronymic:  person.Patronymic,
			Age:         person.Age,
			Nationality: person.Nationality,
			Gender:      person.Gender,
			CreatedAt:   person.CreatedAt,
			UpdatedAt:   person.UpdatedAt,
			DeletedAt:   person.DeletedAt.Time,
		},
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(&response)
	fmt.Println(personID)
}

func CreatePerson(rw http.ResponseWriter, r *http.Request) {
	postgres := db.GetPostgres()
	var person db.Person
	err := json.NewDecoder(r.Body).Decode(&person)
	if err != nil {
		SendErr(rw, "Invalid body")
		return
	}

	if person.Surname == "" {
		SendErr(rw, "Surname is required")
		return
	}
	if person.Name == "" {
		SendErr(rw, "Name is required")
		return
	}

	GetPersonInfo(&person)

	if err := postgres.Create(&person).Error; err != nil {
		SendErr(rw, "Error creating person")
		return
	}
	response := map[string]interface{}{
		"result": fmt.Sprintf("Created person with ID: %d", person.ID),
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(&response)
}

func UpdatePerson(rw http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	personID, err := strconv.Atoi(params["id"])
	if err != nil {
		SendErr(rw, "Invalid person ID")
		return
	}
	postgres := db.GetPostgres()

	var existingPerson db.Person

	if err := postgres.First(&existingPerson, personID).Error; err != nil {
		SendErr(rw, "Person not found")
		return
	}

	var personUpdate db.Person

	err = json.NewDecoder(r.Body).Decode(&personUpdate)
	if err != nil {
		SendErr(rw, "Invalid request body")
		return
	}
	if personUpdate.Name != "" {
		existingPerson.Name = personUpdate.Name
	}
	if personUpdate.Surname != "" {
		existingPerson.Surname = personUpdate.Surname
	}
	if personUpdate.Patronymic != "" {
		existingPerson.Patronymic = personUpdate.Patronymic
	}
	if personUpdate.Age != 0 {
		existingPerson.Age = personUpdate.Age
	}
	if personUpdate.Gender != "" {
		existingPerson.Gender = personUpdate.Gender
	}
	if personUpdate.Nationality != "" {
		existingPerson.Nationality = personUpdate.Nationality
	}
	if err := postgres.Save(&existingPerson).Error; err != nil {
		SendErr(rw, "Error updating person")
		return
	}
	response := map[string]interface{}{
		"result": fmt.Sprintf("Updated person with ID: %d", existingPerson.ID),
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(&response)
}

func DeletePerson(rw http.ResponseWriter, r *http.Request) {
	postgres := db.GetPostgres()
	param := mux.Vars(r)
	personID, err := strconv.Atoi(param["id"])
	if err != nil {
		SendErr(rw, "Invalid person ID")
		return
	}
	var person db.Person
	if err := postgres.First(&person, personID).Error; err != nil {
		SendErr(rw, fmt.Sprintf("Person with ID %d does not exists.", personID))
		return
	}
	if err := postgres.Delete(&person).Error; err != nil {
		SendErr(rw, "Error deleting person")
		return
	}

	response := map[string]interface{}{
		"result": fmt.Sprintf("Deleted person with ID: %d", person.ID),
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(&response)
}

func SendErr(rw http.ResponseWriter, message string) {
	response := map[string]interface{}{
		"error": message,
	}

	json.NewEncoder(rw).Encode(response)
}
