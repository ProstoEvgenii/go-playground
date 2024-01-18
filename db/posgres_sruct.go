package db

import (
	"net/http"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type Person struct {
	gorm.Model
	Name        string
	Surname     string
	Patronymic  string
	Age         int64
	Nationality string
	Gender      string
}

type PersonResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Surname     string    `json:"surname"`
	Patronymic  string    `json:"patronymic"`
	Age         int64     `json:"age"`
	Nationality string    `json:"nationality"`
	Gender      string    `json:"gender"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   time.Time `json:"-"`
}

func FilterData(r *http.Request) *gorm.DB {
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

	query := postgresConnection.Model(&Person{})

	if search != "" {
		query = query.Where("name LIKE ? OR surname LIKE ? OR patronymic LIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	query = query.Limit(limit).Offset((page - 1) * limit)
	return query
}
