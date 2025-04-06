package internal

import (
	"encoding/json"
	"net/http"
	"os"
	"sort"
	"task/models"
)

func EnrichPerson(name string) (int, models.Gender, string, error) {

	agifyURL := os.Getenv("AGIFY_URL") + "/?name=" + name
	resp, err := http.Get(agifyURL)
	if err != nil {
		return 0, 0, "", err
	}
	defer resp.Body.Close()

	var ageData models.PersonWihtAge
	if err := json.NewDecoder(resp.Body).Decode(&ageData); err != nil {
		return 0, 0, "", err
	}

	genderizeURL := os.Getenv("GENDERIZE_URL") + "/?name=" + name
	resp, err = http.Get(genderizeURL)
	if err != nil {
		return 0, 0, "", err
	}
	defer resp.Body.Close()

	var genderData models.PersonWihtGender
	if err := json.NewDecoder(resp.Body).Decode(&genderData); err != nil {
		return 0, 0, "", err
	}

	nationalizeURL := os.Getenv("NATIONALIZE_URL") + "/?name=" + name
	resp, err = http.Get(nationalizeURL)
	if err != nil {
		return 0, 0, "", err
	}
	defer resp.Body.Close()

	var natData models.PersonWihtNationality
	if err := json.NewDecoder(resp.Body).Decode(&natData); err != nil {
		return 0, 0, "", err
	}

	nationality := ""
	if len(natData.Nationality) > 0 {
		sort.Slice(natData.Nationality, func(i, j int) bool {
			return natData.Nationality[i].Probability > natData.Nationality[j].Probability
		})
		nationality = natData.Nationality[0].CountryID
	}

	gender := models.Unknown
	switch genderData.Gender {
	case "male":
		gender = models.Male
	case "female":
		gender = models.Female
	}

	return ageData.Age, gender, nationality, nil
}
