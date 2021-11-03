package app

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

type listing struct {
	id    int    `json:"id"`
	title string `json:"title"`
}

var listings []listing

func (app *App) Initialize() {
	app.Router = mux.NewRouter()
	app.Router.HandleFunc("/listing", app.getListings).Methods("GET")
	app.Router.HandleFunc("/listing", app.createListing).Methods("POST")
	// app.DB = sql.OpenDB("")
}

func (app *App) Run() {
	err := http.ListenAndServe(":"+os.Getenv("PORT"), app.Router)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func (app *App) getListings(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(listings)
}

func (app *App) createListing(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	new_listing := listing{id: 1, title: "Listing"}
	_ = json.NewEncoder(w).Encode(&new_listing)
	listings = append(listings, new_listing)
	json.NewEncoder(w).Encode(listings)
}
