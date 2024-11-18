package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

type Item struct {
	Barcode   string  `json:"barcode"`
	Item      string  `json:"item"`
	Category  string  `json:"category"`
	Price     int     `json:"price"`
	Discount  float64 `json:"discount"`
	Available int     `json:"available"`
}

type Response struct {
	Page       int    `json:"page"`
	PerPage    int    `json:"per_page"`
	Total      int    `json:"total"`
	TotalPages int    `json:"total_pages"`
	Data       []Item `json:"data"`
}

func main() {
	http.HandleFunc("/filter", func(w http.ResponseWriter, r *http.Request) {
		category := r.URL.Query().Get("category")

		minPrice, err := strconv.Atoi(r.URL.Query().Get("minPrice"))
		if err != nil {
			http.Error(w, "cannot convert into int minPrice", http.StatusBadRequest)
			return
		}
		maxPrice, err := strconv.Atoi(r.URL.Query().Get("maxPrice"))
		if err != nil {
			http.Error(w, "cannot convert into int maxPrice", http.StatusBadRequest)
			return
		}

		resp, err := http.Get("https://jsonmock.hackerrank.com/api/inventory")
		if err != nil {
			http.Error(w, "error fetching data from API", http.StatusBadRequest)
			return
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, "error reading api response", http.StatusInternalServerError)
			return
		}

		var apiResponse Response
		err = json.Unmarshal(body, &apiResponse)
		if err != nil {
			http.Error(w, "error parsing api response", http.StatusInternalServerError)
			return
		}

		var filteredItems []Item
		for _, item := range apiResponse.Data {
			if item.Available == 1 && item.Price >= minPrice && item.Price <= maxPrice && item.Category == category {
				filteredItems = append(filteredItems, item)
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(filteredItems)
		// json.NewEncoder(w).Encode(len(filteredItems))
	})

	fmt.Println("Server started to work")
	http.ListenAndServe(":8000", nil)
}
