package GroupieSearch

import (
	"strconv"
)

func MaxMemberCount(artistCards []ArtistCard) []int {
	maxNr := 1
	checkboxNrs := []int{}
	for _, card := range artistCards {
		nrMembers := len(card.Members)
		if nrMembers > maxNr {
			maxNr = nrMembers
		}
	}
	for nr := 1; nr <= maxNr; nr++ {
		checkboxNrs = append(checkboxNrs, nr)
	}
	return checkboxNrs
}
func GetMinStartYear(artistCards []ArtistCard) int {
	minYear := artistCards[0].Created
	for _, card := range artistCards {
		if card.Created < minYear {
			minYear = card.Created
		}
	}
	return minYear
}
func GetMaxStartYear(artistCards []ArtistCard) int {
	maxYear := artistCards[0].Created
	for _, card := range artistCards {
		if card.Created > maxYear {
			maxYear = card.Created
		}
	}
	return maxYear
}

func GetMinFirstAlbumYear(artistCards []ArtistCard) int {
	minFirst, _ := strconv.Atoi(artistCards[0].Album[6:])
	for _, card := range artistCards {
		intYear, _ := strconv.Atoi(card.Album[6:])
		if intYear < minFirst {
			minFirst = intYear
		}
	}
	return minFirst
}
func GetMaxFirstAlbumYear(artistCards []ArtistCard) int {
	maxFirst, _ := strconv.Atoi(artistCards[0].Album[6:])
	for _, card := range artistCards {
		intYear, _ := strconv.Atoi(card.Album[6:])
		if intYear > maxFirst {
			maxFirst = intYear
		}
	}
	return maxFirst
}
