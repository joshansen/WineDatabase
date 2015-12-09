package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	geocodeURL = "https://maps.googleapis.com/maps/api/geocode/json?address="
)

type geocodingResults struct {
	Results []struct {
		Geometry struct {
			Location struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"location"`
		} `json:"geometry"`
	} `json:"results"`
}

func geocode(address string) (lat float64, lng float64, err error) {
	//can add api key with the following  + "&key=" + apiKey
	//fmt.Println(geocodeURL + url.QueryEscape(address))
	resp, err := http.Get(geocodeURL + url.QueryEscape(address))

	if err != nil {
		return 0, 0, fmt.Errorf("Error geocoding address: <%v>", err)
	}

	defer resp.Body.Close()

	var result geocodingResults

	body, err := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &result)

	//err = json.NewDecoder(resp.Body).Decode(&result)
	//fmt.Println(result)

	if err != nil {
		return 0, 0, fmt.Errorf("Error decoding geocoding result: <%v>", err)
	}

	if len(result.Results) > 0 {
		lat = result.Results[0].Geometry.Location.Lat
		lng = result.Results[0].Geometry.Location.Lng
	}

	return lat, lng, nil
}

// testing
// func main() {
// 	lat, lng, err := geocode("11130 Madison St. NE, Blaine, MN 55434")

// 	if err != nil {
// 		fmt.Println(err)
// 	} else {
// 		fmt.Printf("Lat: %v\nLng: %v\n", lat, lng)
// 	}

// }
