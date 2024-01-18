package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"rest-service/db"
)

type Agify struct {
	Count int    `json:"count"`
	Name  string `json:"name"`
	Age   int64  `json:"age"`
}

type Genderize struct {
	Count       int     `json:"count"`
	Name        string  `json:"name"`
	Gender      string  `json:"gender"`
	Probability float64 `json:"probability"`
}

type Nationalize struct {
	Count   int    `json:"count"`
	Name    string `json:"name"`
	Country []struct {
		CountryID   string  `json:"country_id"`
		Probability float64 `json:"probability"`
	} `json:"country"`
}

func GetPersonInfo(person *db.Person) {
	person.Age = getAge(person.Name)
	person.Gender = getGender(person.Name)
	person.Nationality = getNationality(person.Name)

}
func getAge(name string) int64 {
	endpoint := fmt.Sprintf("https://api.agify.io/?name=%s", name)
	resp, err := http.Get(endpoint)
	if err != nil {
		return 0
	}
	var result Agify
	err2 := json.NewDecoder(resp.Body).Decode(&result)
	if err2 != nil {
		return 0
	}
	return result.Age
}

func getGender(name string) string {
	endpoint := fmt.Sprintf("https://api.genderize.io/?name=%s", name)
	resp, err := http.Get(endpoint)
	if err != nil {
		return ""
	}
	var result Genderize
	err2 := json.NewDecoder(resp.Body).Decode(&result)
	if err2 != nil {
		return ""
	}
	return result.Gender
}

func getNationality(name string) string {
	endpoint := fmt.Sprintf("https://api.nationalize.io/?name=%s", name)
	resp, err := http.Get(endpoint)
	if err != nil {
		return ""
	}
	var result Nationalize
	err2 := json.NewDecoder(resp.Body).Decode(&result)
	if err2 != nil {
		return ""
	}
	if len(result.Country) > 0 {
		return result.Country[0].CountryID
	}
	return ""

}
