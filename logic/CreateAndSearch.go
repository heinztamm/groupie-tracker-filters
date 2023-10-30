package GroupieSearch

import (
	"errors"
	"slices"
	"strconv"
	"strings"
)

func CreateArtistCards() ([]ArtistCard, error) {
	artists, err := FetchArtists()
	if err != nil {
		return nil, err
	}

	locations, err := FetchLocations()
	if err != nil {
		return nil, err
	}

	dates, err := FetchDates()
	if err != nil {
		return nil, err
	}

	var artistCards []ArtistCard

	// create maps to look up locations and dates
	locationsMap := make(map[int][]string)
	datesMap := make(map[int][]string)

	for _, artistLocations := range locations.Index {
		locationsMap[artistLocations.ID] = artistLocations.Locations
	}

	for _, artistDates := range dates.Index {
		datesMap[artistDates.ID] = artistDates.Dates
	}

	// loop over artists to create ArtistCards
	for _, artist := range artists {

		// remove the "*" from dates
		var modifiedDates []string
		for _, date := range datesMap[artist.ID] {
			if strings.HasPrefix(date, "*") {
				date = strings.TrimPrefix(date, "*")
			}
			modifiedDates = append(modifiedDates, date)
		}

		artistCard := ArtistCard{
			ID:        artist.ID,
			ImageURL:  artist.ImageURL,
			Name:      artist.Name,
			Members:   artist.Members,
			Created:   artist.Created,
			Album:     artist.Album,
			Locations: locationsMap[artist.ID],
			Dates:     modifiedDates,
		}

		artistCards = append(artistCards, artistCard)
	}

	return artistCards, nil
}

func CreateArtistData(artistCards []ArtistCard) []ArtistData {
	var artistData []ArtistData

	for _, card := range artistCards {
		artistData = append(artistData,
			ArtistData{"artist/band", card.Name},
			ArtistData{"creation date", card.Created},
			ArtistData{"first album", card.Album},
		)
		for _, member := range card.Members {
			artistData = append(artistData,
				ArtistData{"member", member},
			)
		}
		for _, location := range card.Locations {
			artistData = append(artistData,
				ArtistData{"location", location},
			)
		}
		for _, date := range card.Dates {
			artistData = append(artistData,
				ArtistData{"date", date},
			)
		}
	}
	return artistData
}

func SearchArtistCards(query string, filterValues FilterValues, artistCards []ArtistCard) []ArtistCard {
	matchingArtists := []ArtistCard{}
	query = strings.ToLower(query)
	MembersNrMatch := false
	LocationMatch := false
	for _, artistCard := range artistCards {
		intFirstAlbumYear, _ := strconv.Atoi(artistCard.Album[6:])
		MembersNrMatch = false
		LocationMatch = false
		for _, nr := range filterValues.MembersNumbers {
			if len(artistCard.Members) == nr {
				MembersNrMatch = true
			}
		}
		for _, location := range filterValues.LocationSlice {
			if slices.Contains(artistCard.Locations, location) {
				LocationMatch = true
			}
		}
		if !(artistCard.Created >= filterValues.MinStartYear && artistCard.Created <= filterValues.MaxStartYear) {
			continue
		}
		if !(intFirstAlbumYear >= filterValues.MinFirstAlbumYear && intFirstAlbumYear <= filterValues.MaxFirstAlbumYear) {
			continue
		}
		if !MembersNrMatch && len(filterValues.MembersNumbers) != 0 {
			continue
		}
		if !LocationMatch && len(filterValues.LocationSlice) != 0 {
			continue
		}
		if strings.Contains(strings.ToLower(artistCard.Name), query) {
			matchingArtists = append(matchingArtists, artistCard)
		} else if strings.Contains(strings.ToLower(artistCard.Album), query) {
			matchingArtists = append(matchingArtists, artistCard)
		} else if strings.Contains(strings.ToLower(strconv.Itoa(artistCard.Created)), query) {
			matchingArtists = append(matchingArtists, artistCard)
		} else {
			for _, member := range artistCard.Members {
				if strings.Contains(strings.ToLower(member), query) {
					matchingArtists = append(matchingArtists, artistCard)
					break
				}
			}

			for _, location := range artistCard.Locations {
				if strings.Contains(strings.ToLower(location), query) {
					matchingArtists = append(matchingArtists, artistCard)
					break
				}
			}

			for _, date := range artistCard.Dates {
				if strings.Contains(strings.ToLower(date), query) {
					matchingArtists = append(matchingArtists, artistCard)
					break
				}
			}
		}
	}
	return matchingArtists
}

func GetArtistDataByID(artistID string, artistCards []ArtistCard) (ArtistCard, error) {
	id, err := strconv.Atoi(artistID)
	if err != nil {
		return ArtistCard{}, err
	}

	for _, card := range artistCards {
		if card.ID == id {
			return card, nil
		}
	}

	return ArtistCard{}, errors.New("Artist not found")
}
