//go:build !solution

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
)

type MedalCount struct {
	Gold   int `json:"gold"`
	Silver int `json:"silver"`
	Bronze int `json:"bronze"`
	Total  int `json:"total"`
}

type Athlete1 struct {
	Name         string                `json:"athlete"`
	Country      string                `json:"country"`
	Medals       MedalCount            `json:"medals"`
	MedalsByYear map[string]MedalCount `json:"medals_by_year"`
}

var athletes1 []Athlete1

type Athlete2 struct {
	Name         string                `json:"athlete"`
	Country      string                `json:"country"`
	Medals       MedalCount            `json:"medals"`
	MedalsByYear map[string]MedalCount `json:"medals_by_year"`
	Sport        string                `json:"sport"`
}

type Country struct {
	Country string `json:"country"`
	Bronze  int    `json:"bronze"`
	Gold    int    `json:"gold"`
	Silver  int    `json:"silver"`
	Total   int    `json:"total"`
	Year    int    `json:"year"`
}

var athletes2 []Athlete2
var countries []Country

type AthleteRecord struct {
	Name    string `json:"athlete"`
	Country string `json:"country"`
	Year    int    `json:"year"`
	Date    string `json:"date"`
	Sport   string `json:"sport"`
	Gold    int    `json:"gold"`
	Silver  int    `json:"silver"`
	Bronze  int    `json:"bronze"`
	Total   int    `json:"total"`
}

var athleteRecords []AthleteRecord

func convertToAthletes1(records []AthleteRecord) []Athlete1 {
	athletesMap := make(map[string]*Athlete1)

	for _, record := range records {
		athlete, exists := athletesMap[record.Name]
		if !exists {
			athlete = &Athlete1{
				Name:         record.Name,
				Country:      record.Country,
				MedalsByYear: make(map[string]MedalCount),
				Medals:       MedalCount{},
			}
			athletesMap[record.Name] = athlete
		}

		yearKey := fmt.Sprintf("%d", record.Year)

		medalsByYear := athlete.MedalsByYear[yearKey]
		medalsByYear.Gold += record.Gold
		medalsByYear.Silver += record.Silver
		medalsByYear.Bronze += record.Bronze
		medalsByYear.Total += record.Total

		medals := athlete.Medals
		medals.Gold += record.Gold
		medals.Silver += record.Silver
		medals.Bronze += record.Bronze
		medals.Total += record.Total
		athlete.MedalsByYear[yearKey] = medalsByYear
		athlete.Medals = medals
	}

	athletes := make([]Athlete1, 0, len(athletesMap))
	for _, athlete := range athletesMap {
		athletes = append(athletes, *athlete)
	}

	return athletes
}

func convertToAthletes2(records []AthleteRecord) []Athlete2 {
	athletesMap := make(map[string]*Athlete2)

	for _, record := range records {
		athlete, exists := athletesMap[record.Name+record.Sport]
		if !exists {
			athlete = &Athlete2{
				Name:         record.Name,
				Country:      record.Country,
				MedalsByYear: make(map[string]MedalCount),
				Medals:       MedalCount{},
				Sport:        record.Sport,
			}
			athletesMap[record.Name+record.Sport] = athlete
		}

		yearKey := fmt.Sprintf("%d", record.Year)

		medalsByYear := athlete.MedalsByYear[yearKey]
		medalsByYear.Gold += record.Gold
		medalsByYear.Silver += record.Silver
		medalsByYear.Bronze += record.Bronze
		medalsByYear.Total += record.Total

		medals := athlete.Medals
		medals.Gold += record.Gold
		medals.Silver += record.Silver
		medals.Bronze += record.Bronze
		medals.Total += record.Total
		athlete.MedalsByYear[yearKey] = medalsByYear
		athlete.Medals = medals
	}

	athletes := make([]Athlete2, 0, len(athletesMap))
	for _, athlete := range athletesMap {
		athletes = append(athletes, *athlete)
	}

	return athletes
}

func convertToCountries(records []AthleteRecord) []Country {
	countryMap := make(map[string]*Country)

	for _, record := range records {
		country, exists := countryMap[record.Country+fmt.Sprintf("%d", record.Year)]
		if !exists {
			country = &Country{
				Country: record.Country,
				Year:    record.Year,
			}
		}
		country.Gold += record.Gold
		country.Silver += record.Silver
		country.Bronze += record.Bronze
		country.Total += record.Total
		countryMap[record.Country+fmt.Sprintf("%d", record.Year)] = country
	}

	countries_ := make([]Country, 0, len(countryMap))
	for _, country := range countryMap {
		countries_ = append(countries_, *country)
	}

	return countries_
}

func loadData(filePath string) {
	data, _ := os.ReadFile(filePath)

	err := json.Unmarshal(data, &athleteRecords)
	if err != nil {
		return
	}
	athletes1 = convertToAthletes1(athleteRecords)

	athletes2 = convertToAthletes2(athleteRecords)
	countries = convertToCountries(athleteRecords)
}

func athleteInfoHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "name is missing", http.StatusBadRequest)
		return
	}

	for _, athlete := range athletes1 {
		if athlete.Name == name {
			athleteWithoutSport := struct {
				Name         string                `json:"athlete"`
				Country      string                `json:"country"`
				Medals       MedalCount            `json:"medals"`
				MedalsByYear map[string]MedalCount `json:"medals_by_year"`
			}{
				Name:         athlete.Name,
				Country:      athlete.Country,
				Medals:       athlete.Medals,
				MedalsByYear: athlete.MedalsByYear,
			}

			w.Header().Set("Content-Type", "application/json")
			err := json.NewEncoder(w).Encode(athleteWithoutSport)
			if err != nil {
				return
			}
			return
		}
	}

	http.Error(w, "athlete not found", http.StatusNotFound)
}

func topAthletesInSportHandler(w http.ResponseWriter, r *http.Request) {
	sport := r.URL.Query().Get("sport")
	if sport == "" {
		http.Error(w, "sport is missing", http.StatusBadRequest)
		return
	}

	limitStr := r.URL.Query().Get("limit")

	limit := 3
	if limitStr != "" {
		limitInt, err := strconv.Atoi(limitStr)
		if err != nil || limitInt <= 0 {
			http.Error(w, "limit is invalid", http.StatusBadRequest)
			return
		}
		limit = limitInt
	}

	var filteredAthletes []Athlete2
	for _, athlete := range athletes2 {
		if athlete.Sport == sport {
			filteredAthletes = append(filteredAthletes, athlete)
		}
	}

	if len(filteredAthletes) == 0 {
		http.Error(w, "no athletes found for the given sport", http.StatusNotFound)
		return
	}

	sort.Slice(filteredAthletes, func(i, j int) bool {
		a, b := filteredAthletes[i], filteredAthletes[j]
		if a.Medals.Gold != b.Medals.Gold {
			return a.Medals.Gold > b.Medals.Gold
		}
		if a.Medals.Silver != b.Medals.Silver {
			return a.Medals.Silver > b.Medals.Silver
		}
		if a.Medals.Bronze != b.Medals.Bronze {
			return a.Medals.Bronze > b.Medals.Bronze
		}
		return strings.Compare(a.Name, b.Name) < 0
	})

	type AthleteWithoutSport struct {
		Name         string                `json:"athlete"`
		Country      string                `json:"country"`
		Medals       MedalCount            `json:"medals"`
		MedalsByYear map[string]MedalCount `json:"medals_by_year"`
	}

	w.Header().Set("Content-Type", "application/json")

	filteredWithoutSport := make([]AthleteWithoutSport, 0, len(filteredAthletes))
	for _, athlete := range filteredAthletes[0:limit] {
		filteredWithoutSport = append(filteredWithoutSport, AthleteWithoutSport{
			Name:         athlete.Name,
			Country:      athlete.Country,
			Medals:       athlete.Medals,
			MedalsByYear: athlete.MedalsByYear,
		})
	}

	err := json.NewEncoder(w).Encode(filteredWithoutSport)
	if err != nil {
		return
	}

}
func topCountriesInYearHandler(w http.ResponseWriter, r *http.Request) {
	year := r.URL.Query().Get("year")
	if year == "" {
		http.Error(w, "year is missing", http.StatusBadRequest)
		return
	}

	limitStr := r.URL.Query().Get("limit")
	limit := 3
	if limitStr != "" {
		limitInt, err := strconv.Atoi(limitStr)
		if err != nil || limitInt <= 0 {
			http.Error(w, "limit is invalid", http.StatusBadRequest)
			return
		}
		limit = limitInt
	}

	var filteredCountries []Country
	for _, country := range countries {
		yearInt, err := strconv.Atoi(year)
		if err != nil || yearInt <= 0 {
			http.Error(w, "year is invalid", http.StatusBadRequest)
			return
		}
		if country.Year == yearInt {
			filteredCountries = append(filteredCountries, country)
		}
	}

	if len(filteredCountries) == 0 {
		http.Error(w, "no countries found for the given year", http.StatusNotFound)
		return
	}

	sort.Slice(filteredCountries, func(i, j int) bool {
		a, b := filteredCountries[i], filteredCountries[j]
		if a.Gold != b.Gold {
			return a.Gold > b.Gold
		}
		if a.Silver != b.Silver {
			return a.Silver > b.Silver
		}
		if a.Bronze != b.Bronze {
			return a.Bronze > b.Bronze
		}
		return strings.Compare(a.Country, b.Country) < 0
	})

	if limit > len(filteredCountries) {
		limit = len(filteredCountries)
	}

	type CountryWithoutYear struct {
		Country string `json:"country"`
		Bronze  int    `json:"bronze"`
		Gold    int    `json:"gold"`
		Silver  int    `json:"silver"`
		Total   int    `json:"total"`
	}

	filteredWithoutYear := make([]CountryWithoutYear, 0, len(filteredCountries))
	for _, country := range filteredCountries[0:limit] {
		filteredWithoutYear = append(filteredWithoutYear, CountryWithoutYear{
			Country: country.Country,
			Bronze:  country.Bronze,
			Gold:    country.Gold,
			Silver:  country.Silver,
			Total:   country.Total,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(filteredWithoutYear)
	if err != nil {
		return
	}
}

func main() {
	port := flag.String("port", "6029", "port")
	dataPath := flag.String("data", "", "data")
	flag.Parse()

	loadData(*dataPath)

	http.HandleFunc("/athlete-info", athleteInfoHandler)
	http.HandleFunc("/top-athletes-in-sport", topAthletesInSportHandler)
	http.HandleFunc("/top-countries-in-year", topCountriesInYearHandler)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
