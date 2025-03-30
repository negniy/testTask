package internal

import (
	"encoding/json"
	"net/http"
	"task/models"
)

func EnrichPerson(name string) (int, string, string, error) {

	agifyURL := "https://api.agify.io/?name=" + name
	resp, err := http.Get(agifyURL)
	if err != nil {
		return 0, "", "", err
	}
	defer resp.Body.Close()

	var ageData models.PersonWihtAge
	if err := json.NewDecoder(resp.Body).Decode(&ageData); err != nil {
		return 0, "", "", err
	}

	genderizeURL := "https://api.genderize.io/?name=" + name
	resp, err = http.Get(genderizeURL)
	if err != nil {
		return 0, "", "", err
	}
	defer resp.Body.Close()

	var genderData models.PersonWihtGender
	if err := json.NewDecoder(resp.Body).Decode(&genderData); err != nil {
		return 0, "", "", err
	}

	nationalizeURL := "https://api.nationalize.io/?name=" + name
	resp, err = http.Get(nationalizeURL)
	if err != nil {
		return 0, "", "", err
	}
	defer resp.Body.Close()

	var natData models.PersonWihtNationality
	if err := json.NewDecoder(resp.Body).Decode(&natData); err != nil {
		return 0, "", "", err
	}

	nationality := ""
	if len(natData.Nationality) > 0 {
		nationality = natData.Nationality[0].CountryID
	}

	return ageData.Age, genderData.Gender, nationality, nil
}
