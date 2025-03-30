package repository

import (
	"fmt"
	"os"
	"strconv"
	"task/config"
	"task/models"

	"github.com/jinzhu/gorm"
)

var db *gorm.DB

func LoadDB() {
	var err error
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")

	connectionString := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable", user, password, dbname, host, port)

	db, err = gorm.Open("postgres", connectionString)
	if err != nil {
		config.Logger.Fatal("Ошибка при подключении к базе данных: ", err)
	}

	db.AutoMigrate(&models.Person{})

	config.Logger.Debug("Успешно подключено к базе данных PostgreSQL")
}

func GetPeople(idStr, name, surname, ageStr, gender, nationality string, limit, offset int) ([]models.Person, error) {
	reqdb := db
	if idStr != "" {
		id, err := strconv.Atoi(idStr)
		if err == nil {
			if id <= 0 {
				return nil, fmt.Errorf("некорректный ID: %d", id)
			}
			reqdb = reqdb.Where("id = ?", id)
		}
	}
	if name != "" {
		reqdb = reqdb.Where("name LIKE ?", "%"+name+"%")
	}
	if surname != "" {
		reqdb = reqdb.Where("surname LIKE ?", "%"+surname+"%")
	}
	if ageStr != "" {
		age, err := strconv.Atoi(ageStr)
		if err == nil {
			if age <= 0 {
				return nil, fmt.Errorf("некорректный возраст: %d", age)
			}
			reqdb = reqdb.Where("age = ?", age)
		}
	}
	if gender != "" {
		reqdb = reqdb.Where("gender LIKE ?", "%"+gender+"%")
	}
	if nationality != "" {
		reqdb = reqdb.Where("nationality LIKE ?", "%"+nationality+"%")
	}

	reqdb = reqdb.Offset(offset).Limit(limit)

	var people []models.Person
	err := reqdb.Find(&people).Error
	if err != nil {
		return nil, err
	}

	return people, nil
}

func CreatePerson(person models.Person) error {
	if err := db.Create(&person).Error; err != nil {
		return fmt.Errorf("ошибка при создании пользователя: %v", err)
	}
	return nil
}

func UpdatePerson(person models.Person) error {
	if err := db.Model(&person).Updates(models.Person{
		Name:        person.Name,
		Surname:     person.Surname,
		Patronymic:  person.Patronymic,
		Age:         person.Age,
		Gender:      person.Gender,
		Nationality: person.Nationality,
	}).Error; err != nil {
		return fmt.Errorf("ошибка при обновлении пользователя: %v", err)
	}
	return nil
}

func DeletePerson(id int) error {

	if id <= 0 {
		return fmt.Errorf("некорректный ID: %d", id)
	}

	err := db.Where("id = ?", id).Delete(&models.Person{}).Error
	if err != nil {
		return err
	}
	return nil
}
