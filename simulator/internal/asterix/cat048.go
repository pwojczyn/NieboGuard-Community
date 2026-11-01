package asterix

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/pwojczyn/NieboGuard-Community/simulator/pkg/models"
)

// CAT048Encoder generates simulated radar observations (CAT048).
// Educational simulation only - not real PSR/SSR radar data.
type CAT048Encoder struct {
	RadarLat       float64   // Radar station latitude (degrees)
	RadarLon       float64   // Radar station longitude (degrees)
	RadarAlt       float64   // Radar station altitude (meters MSL)
	MaxRangeKm     float64   // Maximum radar range (km)
	UpdateInterval time.Duration // Observation interval
	ErrorPct       float64   // Position error percentage (0-1)
	LossPct        float64   // Packet loss percentage (0-100)
	Rng            *rand.Rand // Random source for error simulation
}

// NewCAT048Encoder creates a CAT048 encoder with default configuration.
func NewCAT048Encoder(radarLat, radarLon, radarAlt float64) *CAT048Encoder {
	return &CAT048Encoder{
		RadarLat:        radarLat,
		RadarLon:        radarLon,
		RadarAlt:        radarAlt,
		MaxRangeKm:      400.0,      // Typical radar range
		UpdateInterval:  5 * time.Second,
		ErrorPct:        0.05,       // 5% position error
		LossPct:         0.0,        // No loss by default
		Rng:             rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// SetLossPercentage configures simulated packet loss.
func (e *CAT048Encoder) SetLossPercentage(pct float64) {
	e.LossPct = pct
}

// SetErrorPercentage configures simulated position measurement error.
func (e *CAT048Encoder) SetErrorPercentage(pct float64) {
	e.ErrorPct = pct
}

// Encode takes a unified AircraftEvent and produces a CAT048 radar plot record.
func (e *CAT048Encoder) Encode(event *models.AircraftEvent) ([]byte, error) {
	if event == nil {
		return nil, fmt.Errorf("nil event")
	}

	// Apply simulated position error
	lat := event.Lat
	lon := event.Lon
	if e.ErrorPct > 0 && e.Rng != nil {
		latErr := (e.Rng.Float64()*2 - 1) * e.ErrorPct * 0.001 // Small position error
		lonErr := (e.Rng.Float64()*2 - 1) * e.ErrorPct * 0.001
		lat += latErr
		lon += lonErr
	}

	// Apply simulated loss
	if e.LossPct > 0 && e.Rng != nil {
		if e.Rng.Float64()*100 < e.LossPct {
			return nil, fmt.Errorf("simulated packet loss (%.1f%%)", e.LossPct)
		}
	}

	// Calculate slant range from radar station
	dlat := lat - e.RadarLat
	dlon := lon - e.RadarLon
	distKm := math.Sqrt(dlat*dlat+dlon*dlon) * 111.0 // Approximate km per degree

	if distKm > e.MaxRangeKm {
		return nil, fmt.Errorf("object out of radar range: %.1f km > %.1f km", distKm, e.MaxRangeKm)
	}

	// Calculate slant range (including altitude difference)
	altDiffKm := math.Abs(event.AltM - e.RadarAlt) / 1000.0 // meters to km
	slantRangeKm := math.Sqrt(distKm*distKm + altDiffKm*altDiffKm)

	// Convert to nautical miles (1 NM = 1.852 km)
	slantNM := uint16(slantRangeKm / 1.852 * 4) // CAT048 scaling: NM * 4

	// Calculate azimuth from radar station
	az := math.Atan2(dlon, dlat) * 180 / math.Pi
	if az < 0 {
		az += 360
	}
	azScaled := uint16(az * 2) // CAT048 scaling: degrees * 2

	buf := new(bytes.Buffer)

	// FSPEC (simplified: position + track number + mode + time)
	buf.WriteByte(0x01)

	// Data Source ID (2 bytes)
	binary.Write(buf, binary.BigEndian, uint16(1))

	// Target Address (4 bytes, simulated)
	binary.Write(buf, binary.BigEndian, uint32(0x1001))

	// Slant Range (2 bytes, NM * 4 approximation)
	binary.Write(buf, binary.BigEndian, slantNM)

	// Azimuth (2 bytes, degrees * 2)
	binary.Write(buf, binary.BigEndian, azScaled)

	// Track Number (2 bytes)
	binary.Write(buf, binary.BigEndian, uint16(1))

	// Mode A Code (2 bytes, simulated)
	binary.Write(buf, binary.BigEndian, uint16(1234))

	// Mode S Code (4 bytes simulated)
	binary.Write(buf, binary.BigEndian, uint32(0xABCDEF01))

	// Time of occurrence (3 bytes, sec of day UTC)
	now := time.Now().UTC()
	secOfDay := uint32(now.Hour()*3600 + now.Minute()*60 + now.Second())
	buf.WriteByte(byte(secOfDay >> 16))
	buf.WriteByte(byte(secOfDay >> 8))
	buf.WriteByte(byte(secOfDay))

	return buf.Bytes(), nil
}

// IsInRange checks if a position is within radar range.
func (e *CAT048Encoder) IsInRange(lat, lon float64) bool {
	dlat := lat - e.RadarLat
	dlon := lon - e.RadarLon
	distKm := math.Sqrt(dlat*dlat+dlon*dlon) * 111.0
	return distKm <= e.MaxRangeKm
}

// GetRangeInfo returns radar range and azimuth information for a position.
func (e *CAT048Encoder) GetRangeInfo(lat, lon float64) (rangeKm, azimuthDeg float64, inRange bool) {
	dlat := lat - e.RadarLat
	dlon := lon - e.RadarLon
	rangeKm = math.Sqrt(dlat*dlat+dlon*dlon) * 111.0
	az := math.Atan2(dlon, dlat) * 180 / math.Pi
	if az < 0 {
		az += 360
	}
	inRange = rangeKm <= e.MaxRangeKm
	return rangeKm, az, inRange
}

// Stats returns encoder statistics.
func (e *CAT048Encoder) Stats() map[string]interface{} {
	return map[string]interface{}{
		"radar_lat":       e.RadarLat,
		"radar_lon":       e.RadarLon,
		"radar_alt":       e.RadarAlt,
		"max_range_km":    e.MaxRangeKm,
		"update_interval": e.UpdateInterval.String(),
		"error_pct":       e.ErrorPct,
		"loss_pct":        e.LossPct,
		"simulated":       true,
	}
}

// Verification: check if CAT048 record decodes correctly.
// Educational test - not real radar interoperability.
func VerifyCAT048Record(data []byte) error {
	if len(data) < 15 {
		return fmt.Errorf("CAT048 record too short: %d bytes", len(data))
	}

	// Check FSPEC
	if data[0] != 0x01 {
		return fmt.Errorf("unexpected FSPEC: %x", data[0])
	}

	// Verify all required fields present (FSPEC + SourceID + TargetAddr + Range + Azimuth + Track + ModeA + ModeS + Time)
	expectedLen := 1 + 2 + 4 + 2 + 2 + 2 + 2 + 4 + 3
	if len(data) < expectedLen {
		return fmt.Errorf("CAT048 record incomplete: %d bytes, expected >= %d", len(data), expectedLen)
	}

	return nil
}
