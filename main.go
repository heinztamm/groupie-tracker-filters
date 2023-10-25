package main

import (
	GroupieSearch "GroupieSearch/logic"
	"html/template"
	"net/http"
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

	data := struct {
		Query       string
		ArtistData  []GroupieSearch.ArtistData
		ArtistCards []GroupieSearch.ArtistCard
	}{
		Query:       "",
		ArtistData:  artistData,
		ArtistCards: artistCards,
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

	artistData := GroupieSearch.CreateArtistData(artistCards)

	query := r.URL.Query().Get("query")

	searchResults := GroupieSearch.SearchArtistCards(query, artistCards)

	data := struct {
		Query       string
		Results     []GroupieSearch.ArtistCard
		ArtistData  []GroupieSearch.ArtistData
		ArtistCards []GroupieSearch.ArtistCard
	}{
		Query:       query,
		Results:     searchResults,
		ArtistData:  artistData,
		ArtistCards: artistCards,
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
