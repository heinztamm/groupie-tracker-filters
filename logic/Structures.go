package GroupieSearch

type Artist struct {
	ID       int      `json:"id"`
	ImageURL string   `json:"image"`
	Name     string   `json:"name"`
	Members  []string `json:"members"`
	Created  int      `json:"creationDate"`
	Album    string   `json:"firstAlbum"`
}

type ArtistCard struct {
	ID        int
	ImageURL  string
	Name      string
	Members   []string
	Created   int
	Album     string
	Locations []string
	Dates     []string
}

type ArtistData struct {
	FieldName  string
	FieldValue interface{}
}

type LocationData struct {
	Index []struct {
		ID        int      `json:"id"`
		Locations []string `json:"locations"`
	} `json:"index"`
}

type DatesData struct {
	Index []struct {
		ID    int      `json:"id"`
		Dates []string `json:"dates"`
	} `json:"index"`
}
type FilterValues struct {
	MembersNumbers    []int
	MinStartYear      int
	MaxStartYear      int
	MinFirstAlbumYear int
	MaxFirstAlbumYear int
	LocationSlice     []string
	AllLocations      []string
}
