package main

import (
	GroupieFilters "GroupieFilters/logic"
	"fmt"
	"html/template"
	"net/http"
)

var tpl *template.Template
var filterLocations []string

func init() {
	tpl = template.Must(template.ParseGlob("static/*.html"))
}

func main() {
	println("SERVER RUNNING...")
	http.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("./assets"))))

	http.HandleFunc("/", home)
	http.HandleFunc("/search", search)
	http.HandleFunc("/filter", filter)
	http.HandleFunc("/artist/", artist)

	http.ListenAndServe(":8080", nil)
}

func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" && r.URL.Path != "/search" {
		http.NotFound(w, r)
		return
	}

	artistCards, err := GroupieFilters.CreateArtistCards()
	if err != nil {
		http.Error(w, "Failed to fetch data", http.StatusInternalServerError)
		return
	}

	artistData := GroupieFilters.CreateArtistData(artistCards)

	for _, x := range artistData {
		if x.FieldName == "location" {
			var found bool
			for _, filterlocs := range filterLocations {
				if x.FieldValue == filterlocs {
					found = true
				}
			}
			if found == false {
				str := fmt.Sprintf("%v", x.FieldValue)
				filterLocations = append(filterLocations, str)
			}
			found = false
		}
	}

	data := struct {
		Query           string
		ArtistData      []GroupieFilters.ArtistData
		ArtistCards     []GroupieFilters.ArtistCard
		FilterLocations []string
	}{
		Query:           "",
		ArtistData:      artistData,
		ArtistCards:     artistCards,
		FilterLocations: filterLocations,
	}

	err = tpl.ExecuteTemplate(w, "index.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func search(w http.ResponseWriter, r *http.Request) {
	artistCards, err := GroupieFilters.CreateArtistCards()
	if err != nil {
		http.Error(w, "Failed to fetch data", http.StatusInternalServerError)
		return
	}

	artistData := GroupieFilters.CreateArtistData(artistCards)

	query := r.URL.Query().Get("query")

	searchResults := GroupieFilters.SearchArtistCards(query, artistCards)

	data := struct {
		Query           string
		Results         []GroupieFilters.ArtistCard
		ArtistData      []GroupieFilters.ArtistData
		ArtistCards     []GroupieFilters.ArtistCard
		FilterLocations []string
	}{
		Query:           query,
		Results:         searchResults,
		ArtistData:      artistData,
		ArtistCards:     artistCards,
		FilterLocations: filterLocations,
	}

	err = tpl.ExecuteTemplate(w, "index.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func filter(w http.ResponseWriter, r *http.Request) {
	artistCards, err := GroupieFilters.CreateArtistCards()
	if err != nil {
		http.Error(w, "Failed to fetch data", http.StatusInternalServerError)
		return
	}

	artistData := GroupieFilters.CreateArtistData(artistCards)

	CreationDateMin := r.URL.Query().Get("CreationDateMin")
	CreationDateMax := r.URL.Query().Get("CreationDateMax")
	FirstAlbumMin := r.URL.Query().Get("FirstAlbumMin")
	FirstAlbumMax := r.URL.Query().Get("FirstAlbumMax")
	MembersAmountMin := r.URL.Query().Get("MembersAmountMin")
	MembersAmountMax := r.URL.Query().Get("MembersAmountMax")
	SelectedLocations := r.URL.Query()["Location"]

	filterResults := GroupieFilters.FilterArtistCards(CreationDateMin, CreationDateMax, FirstAlbumMin, FirstAlbumMax, MembersAmountMin, MembersAmountMax, SelectedLocations, artistCards)

	data := struct {
		Query           string
		Results         []GroupieFilters.ArtistCard
		ArtistData      []GroupieFilters.ArtistData
		ArtistCards     []GroupieFilters.ArtistCard
		FilterLocations []string
	}{
		Query:           "true",
		Results:         filterResults,
		ArtistData:      artistData,
		ArtistCards:     artistCards,
		FilterLocations: filterLocations,
	}

	err = tpl.ExecuteTemplate(w, "index.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func artist(w http.ResponseWriter, r *http.Request) {
	artistID := r.URL.Path[len("/artist/"):]

	artistCards, err := GroupieFilters.CreateArtistCards()
	if err != nil {
		http.Error(w, "Failed to fetch data", http.StatusInternalServerError)
		return
	}

	artistData, err := GroupieFilters.GetArtistDataByID(artistID, artistCards)
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
