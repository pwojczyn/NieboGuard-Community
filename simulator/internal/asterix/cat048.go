package asterix

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"time"
)

// CAT048Encoder generates simulated radar observations.
type CAT048Encoder struct {
	RadarLat float64
	RadarLon float64
	RadarAlt float64
	MaxRangeKm float64
}

func (e *CAT048Encoder) Encode(simObs *SimulatedObservation, radarLat, radarLon float64) ([]byte, error) {
	// Simulated radar observation: compute slant range, azimuth
	dlat := simObs.Lat - radarLat
	dlon := simObs.Lon - radarLon
	distKm := math.Sqrt(dlat*dlat + dlon*dlon) * 111.0 // approximate

	if distKm > e.MaxRangeKm {
		return nil, fmt.Errorf("object out of radar range: %.1f km > %.1f km", distKm, e.MaxRangeKm)
	}

	buf := new(bytes.Buffer)
	// FSPEC (simplified)
	buf.WriteByte(0x01)

	// Data Source ID
	binary.Write(buf, binary.BigEndian, uint16(1))

	// Target Address (4 bytes, simulated)
	binary.Write(buf, binary.BigEndian, uint32(0x1001))

	// Slant Range (2 bytes, nautical miles * 4 approximation)
	slantNM := uint16(distKm * 0.539957) // km to NM
	binary.Write(buf, binary.BigEndian, slantNM)

	// Azimuth (2 bytes, degrees * 2)
	az := math.Atan2(dlon, dlat) * 180 / math.Pi
	if az < 0 {
		az += 360
	}
	azScaled := uint16(az * 2)
	binary.Write(buf, binary.BigEndian, azScaled)

	// Track Number (2 bytes)
	binary.Write(buf, binary.BigEndian, uint16(1))

	// Mode A Code (2 bytes)
	binary.Write(buf, binary.BigEndian, uint16(1234))

	// Mode S Code (4 bytes simulated)
	binary.Write(buf, binary.BigEndian, uint32(0xABCDEF01))

	// Time of occurrence (3 bytes, sec of day)
	now := time.Now().UTC()
	secOfDay := uint32(now.Hour()*3600 + now.Minute()*60 + now.Second())
	buf.WriteByte(byte(secOfDay >> 16))
	buf.WriteByte(byte(secOfDay >> 8))
	buf.WriteByte(byte(secOfDay))

	return buf.Bytes(), nil
}

type SimulatedObservation struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}
