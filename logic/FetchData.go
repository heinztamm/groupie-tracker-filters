package GroupieFilters

import (
	"encoding/json"
	"io"
	"net/http"
)

func FetchArtists() ([]Artist, error) {
	artistResponse, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		return nil, err
	}
	defer artistResponse.Body.Close()

	artistData, err := io.ReadAll(artistResponse.Body)
	if err != nil {
		return nil, err
	}

	var Artists []Artist
	if err := json.Unmarshal(artistData, &Artists); err != nil {
		return nil, err
	}
	return Artists, nil
}

func FetchLocations() (LocationData, error) {
	locationResponse, err := http.Get("https://groupietrackers.herokuapp.com/api/locations")
	if err != nil {
		return LocationData{}, err
	}

	defer locationResponse.Body.Close()

	var locationData LocationData
	if err := json.NewDecoder(locationResponse.Body).Decode(&locationData); err != nil {
		return LocationData{}, err
	}

	return locationData, nil
}

func FetchDates() (DatesData, error) {
	datesResponse, err := http.Get("https://groupietrackers.herokuapp.com/api/dates")
	if err != nil {
		return DatesData{}, err
	}
	defer datesResponse.Body.Close()

	var datesData DatesData
	if err := json.NewDecoder(datesResponse.Body).Decode(&datesData); err != nil {
		return DatesData{}, err
	}

	return datesData, nil
}
