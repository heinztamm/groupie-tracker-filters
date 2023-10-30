package main

import (
	GroupieSearch "GroupieSearch/logic"
	"html/template"
	"net/http"
	"slices"
	"strconv"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("static/*.html"))
}

func main() {

	http.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("./assets"))))

	http.HandleFunc("/", home)
	http.HandleFunc("/search", search)

	http.HandleFunc("/artist/", artist)

	http.ListenAndServe(":8080", nil)
}

func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" && r.URL.Path != "/search" {
		http.NotFound(w, r)
		return
	}

	artistCards, err := GroupieSearch.CreateArtistCards()
	if err != nil {
		http.Error(w, "Failed to fetch data", http.StatusInternalServerError)
		return
	}

	artistData := GroupieSearch.CreateArtistData(artistCards)

	allLocations := []string{}

	for _, act := range artistCards {
		for _, location := range act.Locations {
			if !(slices.Contains(allLocations, location)) {
				allLocations = append(allLocations, location)
			}
		}
	}

	var filterValues GroupieSearch.FilterValues

	checkboxNrs := GroupieSearch.MaxMemberCount(artistCards)

	var searchResults []GroupieSearch.ArtistCard
	data := struct {
		Query             string
		ArtistData        []GroupieSearch.ArtistData
		ArtistCards       []GroupieSearch.ArtistCard
		FilterValuesSlice []GroupieSearch.FilterValues
		Results           []GroupieSearch.ArtistCard
		CheckboxNrs       []int
		MembersNumbers    []int
		LocationSlice     []string
	}{
		Query:             "",
		ArtistData:        artistData,
		ArtistCards:       artistCards,
		FilterValuesSlice: []GroupieSearch.FilterValues{},
		Results:           searchResults,
		CheckboxNrs:       checkboxNrs,
		MembersNumbers:    filterValues.MembersNumbers,
		LocationSlice:     allLocations,
	}

	err = tpl.ExecuteTemplate(w, "index.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func search(w http.ResponseWriter, r *http.Request) {
	artistCards, err := GroupieSearch.CreateArtistCards()
	if err != nil {
		http.Error(w, "Failed to fetch data", http.StatusInternalServerError)
		return
	}

	var filterValues GroupieSearch.FilterValues

	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		for _, str := range r.Form["nr_members"] {
			intValue, _ := strconv.Atoi(str)
			filterValues.MembersNumbers = append(filterValues.MembersNumbers, intValue)
		}
		for _, location := range r.Form["location"] {
			filterValues.LocationSlice = append(filterValues.LocationSlice, location)
		}

		filterValues.MinStartYear, _ = strconv.Atoi(r.FormValue("minStart"))
		filterValues.MaxStartYear, _ = strconv.Atoi(r.FormValue("maxStart"))
		filterValues.MinFirstAlbumYear, _ = strconv.Atoi(r.FormValue("minFirst"))
		filterValues.MaxFirstAlbumYear, _ = strconv.Atoi(r.FormValue("maxFirst"))
	}

	artistData := GroupieSearch.CreateArtistData(artistCards)

	query := r.URL.Query().Get("query")

	searchResults := GroupieSearch.SearchArtistCards(query, filterValues, artistCards)

	checkboxNrs := GroupieSearch.MaxMemberCount(artistCards)

	data := struct {
		Query             string
		Results           []GroupieSearch.ArtistCard
		ArtistData        []GroupieSearch.ArtistData
		ArtistCards       []GroupieSearch.ArtistCard
		FilterValuesSlice []GroupieSearch.FilterValues
		CheckboxNrs       []int
		MembersNumbers    []int
		LocationSlice     []string
	}{
		Query:             query,
		Results:           searchResults,
		ArtistData:        artistData,
		ArtistCards:       artistCards,
		FilterValuesSlice: []GroupieSearch.FilterValues{filterValues},
		CheckboxNrs:       checkboxNrs,
		MembersNumbers:    filterValues.MembersNumbers,
		LocationSlice:     filterValues.LocationSlice,
	}

	err = tpl.ExecuteTemplate(w, "index.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func artist(w http.ResponseWriter, r *http.Request) {
	artistID := r.URL.Path[len("/artist/"):]

	artistCards, err := GroupieSearch.CreateArtistCards()
	if err != nil {
		http.Error(w, "Failed to fetch data", http.StatusInternalServerError)
		return
	}

	artistData, err := GroupieSearch.GetArtistDataByID(artistID, artistCards)
	if err != nil {
		http.Error(w, "Page not found", http.StatusNotFound)
		return
	}

	templateFilePath := "static/artist.html"

	tmpl, err := template.ParseFiles(templateFilePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, artistData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
