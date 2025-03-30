package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"task/config"
	"task/internal"
	"task/models"
	"task/repository"
)

func response(w http.ResponseWriter, code int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if data != nil {
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			config.Logger.Debug(err)
		}
	}
}

func responseError(w http.ResponseWriter, code int, err error) {
	response(w, code, map[string]string{"error :": err.Error()})
}

func GetPeople(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	name := r.URL.Query().Get("name")
	surname := r.URL.Query().Get("surname")
	ageStr := r.URL.Query().Get("age")
	gender := r.URL.Query().Get("gender")
	nationality := r.URL.Query().Get("nationality")
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 10
	offset := 0

	if limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	if offsetStr != "" {
		parsedOffset, err := strconv.Atoi(offsetStr)
		if err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	people, err := repository.GetPeople(id, name, surname, ageStr, gender, nationality, limit, offset)
	if err != nil {
		responseError(w, http.StatusInternalServerError, err)
		return
	}

	response(w, http.StatusOK, people)
}

func CreatePerson(w http.ResponseWriter, r *http.Request) {
	var input models.Person

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		config.Logger.Debug("Ошибка парсинга входных данных: ", err)
		responseError(w, http.StatusBadRequest, err)
	}
	defer r.Body.Close()

	age, gender, nationality, err := internal.EnrichPerson(input.Name)
	if err != nil {
		config.Logger.Debug("Ошибка обогащения: ", err)
		responseError(w, http.StatusInternalServerError, err)
	}

	input.Age = age
	input.Gender = gender
	input.Nationality = nationality

	err = repository.CreatePerson(input)
	if err != nil {
		config.Logger.Debug("Ошибка сохранения данных в БД: ", err)
		responseError(w, http.StatusInternalServerError, err)
	}

	response(w, http.StatusCreated, input)
}

func UpdatePerson(w http.ResponseWriter, r *http.Request) {
	input := new(models.UpdatePerson)

	err := json.NewDecoder(r.Body).Decode(input)
	if err != nil {
		config.Logger.Debug("Ошибка парсинга входных данных: ", err)
		responseError(w, http.StatusBadRequest, err)
	}
	defer r.Body.Close()

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		config.Logger.Debug("Ошибка парсинга id: ", err)
		responseError(w, http.StatusBadRequest, err)
	}

	existed, err := repository.GetPeople(strconv.Itoa(id), "", "", "", "", "", 0, 1)
	if err != nil {
		config.Logger.Debug("Ошибка поиска: ", err)
		responseError(w, http.StatusNotFound, err)
	}

	exist := existed[0]

	switch {
	case input.Name != nil:
		exist.Name = *input.Name
	case input.Surname != nil:
		exist.Surname = *input.Surname
	case input.Patronymic != nil:
		exist.Patronymic = *input.Patronymic
	case input.Age != nil:
		exist.Age = *input.Age
	case input.Gender != nil:
		exist.Gender = *input.Gender
	case input.Nationality != nil:
		exist.Nationality = *input.Nationality
	}

	err = repository.UpdatePerson(exist)
	if err != nil {
		config.Logger.Debug("Ошибка обновления данных в БД: ", err)
		responseError(w, http.StatusInternalServerError, err)
	}

	response(w, http.StatusOK, nil)
}

func DeletePerson(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		config.Logger.Debug("Ошибка парсинга id: ", err)
		responseError(w, http.StatusBadRequest, err)
	}

	err = repository.DeletePerson(id)
	if err != nil {
		config.Logger.Debug("Ошибка удаления: ", err)
		responseError(w, http.StatusNotFound, err)
	}

	response(w, http.StatusNoContent, nil)
}
