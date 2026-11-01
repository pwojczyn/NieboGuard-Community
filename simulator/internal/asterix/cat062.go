package asterix

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"time"
)

// CAT062Encoder generates simulated system tracks (tracked objects with history).
type CAT062Encoder struct {
	DataSourceID uint16
	TrackCounter uint16
	ActiveTracks map[string]*TrackState
}

type TrackState struct {
	TrackNumber uint16
	Lat         float64
	Lon         float64
	Speed       float64
	Track       float64
	LastUpdate  time.Time
	TTL         int // seconds
	Status      uint8 // 0=new, 1=active, 2=coasting, 3=terminated
}

func (e *CAT062Encoder) Encode(simObs *SimulatedObservation) ([]byte, error) {
	e.init()

	track, exists := e.ActiveTracks[simObs.ID]
	if !exists {
		e.TrackCounter++
		track = &TrackState{
			TrackNumber: e.TrackCounter,
			Status:      0, // new
			TTL:         30,
		}
		e.ActiveTracks[simObs.ID] = track
	}

	// Update track with new observation
	track.Lat = simObs.Lat
	track.Lon = simObs.Lon
	track.Speed = simObs.SpeedMps * 1.944 // m/s to knots
	track.Track = simObs.HeadingDeg
	track.LastUpdate = time.Now()
	track.Status = 1 // active

	buf := new(bytes.Buffer)

	// FSPEC (simplified: position + speed + track + status)
	buf.WriteByte(0x01)

	// Data Source ID
	binary.Write(buf, binary.BigEndian, e.DataSourceID)

	// Track Number
	binary.Write(buf, binary.BigEndian, track.TrackNumber)

	// Latitude
	latScaled := int32(track.Lat * 128 * 60)
	binary.Write(buf, binary.BigEndian, latScaled)

	// Longitude
	lonScaled := int32(track.Lon * 128 * 60)
	binary.Write(buf, binary.BigEndian, lonScaled)

	// Speed
	speedScaled := uint16(track.Speed / 2)
	binary.Write(buf, binary.BigEndian, speedScaled)

	// Track angle
	trackScaled := uint16(track.Track * 90)
	binary.Write(buf, binary.BigEndian, trackScaled)

	// Track Status
	buf.WriteByte(track.Status)

	// Time of Track (3 bytes)
	secOfDay := uint32(time.Now().UTC().Hour()*3600 + time.Now().UTC().Minute()*60 + time.Now().UTC().Second())
	buf.WriteByte(byte(secOfDay >> 16))
	buf.WriteByte(byte(secOfDay >> 8))
	buf.WriteByte(byte(secOfDay))

	return buf.Bytes(), nil
}

func (e *CAT062Encoder) init() {
	if e.ActiveTracks == nil {
		e.ActiveTracks = make(map[string]*TrackState)
	}
}

func (e *CAT062Encoder) CleanupExpired(maxAge time.Duration) {
	for id, t := range e.ActiveTracks {
		if time.Since(t.LastUpdate) > maxAge {
			delete(e.ActiveTracks, id)
		}
	}
}

type SimulatedObservation struct {
	ID         string  `json:"id"`
	Lat        float64 `json:"lat"`
	Lon        float64 `json:"lon"`
	SpeedMps   float64 `json:"speed_mps"`
	HeadingDeg float64 `json:"heading_deg"`
}