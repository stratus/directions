package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

var (
	origin      = flag.String("origin", "Google San Francisco, San Francisco, CA", "Starting point")
	destination = flag.String("destination", "Google Headquarters, Amphitheatre Parkway, Mountain View, CA", "Destination")
)

const DIR_BASE_URL = "http://maps.googleapis.com/maps/api/directions/json"

type Directions struct {
	Routes []struct {
		Bounds struct {
			Northeast struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"northeast"`
			Southwest struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"southwest"`
		} `json:"bounds"`
		Copyrights string `json:"copyrights"`
		Legs       []struct {
			Distance struct {
				Text  string `json:"text"`
				Value int64  `json:"value"`
			} `json:"distance"`
			Duration struct {
				Text  string `json:"text"`
				Value int64  `json:"value"`
			} `json:"duration"`
			EndAddress  string `json:"end_address"`
			EndLocation struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"end_location"`
			StartAddress  string `json:"start_address"`
			StartLocation struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"start_location"`
			Steps []struct {
				Distance struct {
					Text  string `json:"text"`
					Value int64  `json:"value"`
				} `json:"distance"`
				Duration struct {
					Text  string `json:"text"`
					Value int64  `json:"value"`
				} `json:"duration"`
				EndLocation struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"end_location"`
				HtmlInstructions string `json:"html_instructions"`
				Polyline         struct {
					Points string `json:"points"`
				} `json:"polyline"`
				StartLocation struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"start_location"`
				TravelMode string `json:"travel_mode"`
			} `json:"steps"`
			ViaWaypoint []interface{} `json:"via_waypoint"`
		} `json:"legs"`
		OverviewPolyline struct {
			Points string `json:"points"`
		} `json:"overview_polyline"`
		Summary       string        `json:"summary"`
		Warnings      []interface{} `json:"warnings"`
		WaypointOrder []interface{} `json:"waypoint_order"`
	} `json:"routes"`
	Status string `json:"status"`
}

// decodeJSON decodes JSON out of io.Reader input
func decodeJSON(r io.Reader) (*Directions, error) {
	d := json.NewDecoder(r)
	var ret *Directions
	if err := d.Decode(&ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// encodeURL encodes URL given base URI, origin and destination
func encodeURL(base, origin, destination string) (*url.URL, error) {
	u, err := url.Parse(base)
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Set("origin", origin)
	q.Set("destination", destination)
	q.Set("sensor", "false")
	u.RawQuery = q.Encode()
	return u, nil
}

// getDirections finds out directions given an origin and destination
func getDirections(origin, destination string) (*Directions, error) {
	u, err := encodeURL(DIR_BASE_URL, origin, destination)
	if err != nil {
		return nil, err
	}
	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	d, err := decodeJSON(resp.Body)
	if err != nil {
		return nil, err
	}
	return d, nil
}

// stripHTML strips HTML tags out of given string
func stripHTML(s string) string {
	b := bytes.NewBufferString("")
	tagHere := false
	for _, v := range s {
		switch v {
		case '<':
			tagHere = true
		case '>':
			tagHere = false
		default:
			if !tagHere {
				b.WriteRune(v)
			}
		}
	}
	return b.String()
}

func main() {
	flag.Parse()
	directions, err := getDirections(*origin, *destination)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s -> %s\n\n", *origin, *destination)
	for _, r := range directions.Routes {
		for _, v := range r.Legs {
			fmt.Printf("%s\n", v.Distance.Text)
			for k, s := range v.Steps {
				fmt.Printf("%d - %s (%s)\n", k+1, stripHTML(s.HtmlInstructions), s.Duration.Text)
			}
		}
		fmt.Println()
	}
}
