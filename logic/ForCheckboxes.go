package GroupieSearch

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
