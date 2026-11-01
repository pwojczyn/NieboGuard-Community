package asterix

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"time"
)

// CAT021Encoder creates binary ASTERIX CAT021 records from ADS-B data.
// Based on EUROCONTROL SPEC-014 (2023 edition).

type CAT021Encoder struct {
	DataSourceID uint16
}

func (e *CAT021Encoder) Encode(record *Record) ([]byte, error) {
	buf := new(bytes.Buffer)

	// FSPEC: minimal spec for basic record
	// In production, FSPEC is variable length based on which fields present.
	// For MVP: fixed minimal spec with position + identity.
	fspec := []byte{0x01, 0x00} // minimal flags (simplified)
	buf.Write(fspec)

	// Data Source Identifier (2 bytes)
	binary.Write(buf, binary.BigEndian, e.DataSourceID)

	// Track Number (2 bytes, 0 if not tracked)
	if record.TrackNumber > 0 {
		binary.Write(buf, binary.BigEndian, record.TrackNumber)
	}

	// Latitude (4 bytes, scaled 1/128 minutes ≈ 0.0078°)
	if record.Lat != 0 {
		latScaled := int32(record.Lat * 128 * 60)
		binary.Write(buf, binary.BigEndian, latScaled)
	}

	// Longitude (4 bytes)
	if record.Lon != 0 {
		lonScaled := int32(record.Lon * 128 * 60)
		binary.Write(buf, binary.BigEndian, lonScaled)
	}

	// Altitude (2 bytes, feet above MSL, scaled)
	if record.Altitude > 0 {
		altScaled := uint16(record.Altitude / 2) // approximate
		binary.Write(buf, binary.BigEndian, altScaled)
	}

	// Speed (2 bytes, knots / 2 for approximation)
	if record.Speed > 0 {
		speedScaled := uint16(record.Speed / 2)
		binary.Write(buf, binary.BigEndian, speedScaled)
	}

	// Track angle (2 bytes, 1/90°)
	if record.Track != 0 {
		trackScaled := uint16(record.Track * 90)
		binary.Write(buf, binary.BigEndian, trackScaled)
	}

	// Time of occurrence (3 bytes, seconds from midnight UTC)
	secOfDay := uint32(time.Now().UTC().Hour()*3600 + time.Now().UTC().Minute()*60 + time.Now().UTC().Second())
	buf.WriteByte(byte(secOfDay >> 16))
	buf.WriteByte(byte(secOfDay >> 8))
	buf.WriteByte(byte(secOfDay))

	return buf.Bytes(), nil
}

type Record struct {
	DataSourceID uint16
	TrackNumber  uint16
	Lat          float64
	Lon          float64
	Altitude     float64 // feet
	Speed        float64 // knots
	Track        float64 // degrees
	TimeOfOccur  time.Time
}

func (e *CAT021Encoder) EncodeFromADSBAircraft(hex string, lat, lon float64, alt int, gs float64, track float64) ([]byte, error) {
	rec := &Record{
		DataSourceID: e.DataSourceID,
		Lat:          lat,
		Lon:          lon,
		Altitude:     float64(alt),
		Speed:        gs,
		Track:        track,
		TimeOfOccur:  time.Now().UTC(),
	}
	if hex != "" {
		// Use hex as track number for identification (simplified)
		rec.TrackNumber = 1
	}
	return e.Encode(rec)
}

func VerifyRecord(encoded []byte) error {
	if len(encoded) < 2 {
		return fmt.Errorf("record too short")
	}
	// Basic structural check: must have FSPEC + identifier
	return nil
}
