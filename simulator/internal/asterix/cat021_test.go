package asterix

import (
	"testing"
)

func TestCAT021Encoder_Encode(t *testing.T) {
	enc := &CAT021Encoder{DataSourceID: 1}
	rec := &Record{DataSourceID: 1, Lat: 52.1, Lon: 19.2, Altitude: 32000, Speed: 450, Track: 90}
	data, err := enc.Encode(rec)
	if err != nil {
		t.Fatalf("Encode error: %v", err)
	}
	if len(data) < 4 {
		t.Errorf("Record too short: %d bytes", len(data))
	}
}

func TestCAT021Encoder_EncodeFromADSBAircraft(t *testing.T) {
	enc := &CAT021Encoder{DataSourceID: 1}
	data, err := enc.EncodeFromADSBAircraft("abc123", 52.1, 19.2, 32000, 450.0, 90.0)
	if err != nil {
		t.Fatalf("Encode error: %v", err)
	}
	if len(data) < 4 {
		t.Error("Data too short")
	}
}

func TestCAT021Encoder_EncodeZeroValues(t *testing.T) {
	enc := &CAT021Encoder{DataSourceID: 1}
	rec := &Record{DataSourceID: 1, Lat: 0, Lon: 0, Altitude: 0, Speed: 0, Track: 0}
	data, err := enc.Encode(rec)
	if err != nil {
		t.Fatalf("Encode zero error: %v", err)
	}
	if len(data) < 2 {
		t.Error("Zero value record too short")
	}
}

func TestVerifyRecord(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		wantErr bool
	}{
		{"too short", []byte{0x01}, true},
		{"valid", []byte{0x01, 0x00, 0xFF, 0x01}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := VerifyRecord(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("VerifyRecord error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}
