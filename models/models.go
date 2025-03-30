package models

type PersonWihtAge struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type PersonWihtGender struct {
	Name   string `json:"name"`
	Gender string `json:"gender"`
}

type PersonWihtNationality struct {
	Name        string `json:"name"`
	Nationality []struct {
		CountryID   string  `json:"country_id"`
		Probability float64 `json:"probability"`
	} `json:"country"`
}

type Person struct {
	Id          int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string `json:"name" gorm:"type:varchar(100);not null"`
	Surname     string `json:"surname" gorm:"type:varchar(100);not null"`
	Patronymic  string `json:"patronymic,omitempty" gorm:"type:varchar(100)"`
	Age         int    `json:"age,omitempty" gorm:"default:0"`
	Gender      string `json:"gender,omitempty" gorm:"type:varchar(50)"`
	Nationality string `json:"nationality,omitempty" gorm:"type:varchar(50)"`
}

type UpdatePerson struct {
	Name        *string `json:"name,omitempty"`
	Surname     *string `json:"surname,omitempty"`
	Patronymic  *string `json:"patronymic,omitempty"`
	Age         *int    `json:"age,omitempty"`
	Gender      *string `json:"gender,omitempty"`
	Nationality *string `json:"nationality,omitempty"`
}
