package asterix

import (
	"bytes"
	"encoding/binary"
	"time"

	"github.com/pwojczyn/NieboGuard-Community/simulator/pkg/models"
)

// CAT062Encoder generates simulated system tracks (tracked objects with history).
type CAT062Encoder struct {
	DataSourceID  uint16
	TrackCounter  uint16
	ActiveTracks  map[string]*TrackState
}

type TrackState struct {
	TrackNumber uint16
	Lat         float64
	Lon         float64
	Speed       float64
	Track       float64
	LastUpdate  time.Time
	TTL         int
	Status      uint8 // 0=new, 1=active, 2=coasting, 3=terminated
}

// Encode takes a unified AircraftEvent and produces a CAT062 record.
func (e *CAT062Encoder) Encode(event *models.AircraftEvent) ([]byte, error) {
	e.init()

	trackKey := event.SourceID + ":" + event.Flight
	track, exists := e.ActiveTracks[trackKey]
	if !exists {
		e.TrackCounter++
		track = &TrackState{
			TrackNumber: e.TrackCounter,
			Status:      0, // new
			TTL:         30,
		}
		e.ActiveTracks[trackKey] = track
	}

	// Update track with new observation
	track.Lat = event.Lat
	track.Lon = event.Lon
	track.Speed = event.Gs * 1.944 // m/s to knots
	track.Track = event.HeadingDeg
	track.LastUpdate = time.Now()
	track.Status = 1 // active

	buf := new(bytes.Buffer)

	// FSPEC (simplified: position + speed + track + status)
	buf.WriteByte(0x01)

	// Data Source ID
	binary.Write(buf, binary.BigEndian, e.DataSourceID)

	// Track Number
	binary.Write(buf, binary.BigEndian, track.TrackNumber)

	// Latitude (scaled 1/128 minutes)
	latScaled := int32(track.Lat * 128 * 60)
	binary.Write(buf, binary.BigEndian, latScaled)

	// Longitude
	lonScaled := int32(track.Lon * 128 * 60)
	binary.Write(buf, binary.BigEndian, lonScaled)

	// Speed (knots / 2)
	speedScaled := uint16(track.Speed / 2)
	binary.Write(buf, binary.BigEndian, speedScaled)

	// Track angle (degrees * 90)
	trackScaled := uint16(track.Track * 90)
	binary.Write(buf, binary.BigEndian, trackScaled)

	// Track Status
	buf.WriteByte(track.Status)

	// Time of Track (3 bytes, sec of day)
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