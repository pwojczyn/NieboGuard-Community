package collectors

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// ADSBJSONAdapter reads aircraft.json from dump1090/reads.
type ADSBJSONAdapter struct {
	URL     string
	Timeout time.Duration
}

func (a *ADSBJSONAdapter) Fetch() ([]AircraftRecord, error) {
	client := &http.Client{Timeout: a.Timeout}
	resp, err := client.Get(a.URL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	var data struct {
		Aircraft []json.RawMessage `json:"aircraft"`
		Now      float64           `json:"now"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	var records []AircraftRecord
	for _, raw := range data.Aircraft {
		var a struct {
			Hex      string  `json:"hex"`
			Flight   string  `json:"flight"`
			Lat      float64 `json:"lat"`
			Lon      float64 `json:"lon"`
			AltBaro  int     `json:"alt_baro"`
			GS       float64 `json:"gs"`
			Track    float64 `json:"track"`
			Seen     float64 `json:"seen"`
			Military bool    `json:"military"`
			Category string  `json:"category"`
		}
		if err := json.Unmarshal(raw, &a); err != nil {
			continue
		}
		if a.Hex == "" || a.Lat == 0 || a.Lon == 0 {
			continue
		}
		records = append(records, AircraftRecord{
			Hex:     a.Hex,
			Flight:  strings.TrimSpace(a.Flight),
			Lat:     a.Lat,
			Lon:     a.Lon,
			AltBaro: a.AltBaro,
			GS:      a.GS,
			Track:   a.Track,
			Seen:    a.Seen,
		})
	}
	return records, nil
}

type AircraftRecord struct {
	Hex     string
	Flight  string
	Lat     float64
	Lon     float64
	AltBaro int
	GS      float64
	Track   float64
	Seen    float64
}
