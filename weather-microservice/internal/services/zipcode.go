package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type Location struct {
	City string `json:"localidade"`
}

func GetLocationByZipcode(zipcode string) (Location, error) {
	url := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", zipcode)
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != 200 {
		return Location{}, errors.New("could not fetch location")
	}
	defer resp.Body.Close()

	var location Location
	if err := json.NewDecoder(resp.Body).Decode(&location); err != nil || location.City == "" {
		return Location{}, errors.New("invalid location data")
	}

	return location, nil
}
