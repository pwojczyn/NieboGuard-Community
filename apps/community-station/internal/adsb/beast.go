package adsb

import (
	"bufio"
	"fmt"
	"log/slog"
	"net"
	"os"
	"strings"
	"time"
)

// BeastReader reads Beast TCP data from a port or connection
type BeastReader struct {
	conn      net.Conn
	logger    *slog.Logger
	enabled   bool
	readRate  time.Duration
}

// NewBeastReader creates a new Beast TCP reader
func NewBeastReader(logger *slog.Logger) *BeastReader {
	return &BeastReader{logger: logger, readRate: 100 * time.Millisecond}
}

// Connect connects to Beast TCP on the given port
func (b *BeastReader) Connect(port int) error {
	addr := fmt.Sprintf("localhost:%d", port)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return fmt.Errorf("connecting to Beast TCP on %s: %w", addr, err)
	}
	b.conn = conn
	b.enabled = true
	b.logger.Info("Connected to Beast TCP", "port", port)
	return nil
}

// ReadFrame reads one aircraft frame from Beast TCP
func (b *BeastReader) ReadFrame() (map[string]string, error) {
	if !b.enabled || b.conn == nil {
		return nil, fmt.Errorf("Beast reader not connected")
	}

	// Read with timeout
	// beast frames end with \n
	reader := bufio.NewReaderTimeout(b.conn, 5*time.Second)
	line, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	line = strings.TrimSpace(line)
	if line == "" {
		return nil, nil
	}

	// Parse Beast format: hex,lat,lon,alt,gs,track,seen
	// Example: A024E8,52.1,19.2,35000,450,90,30
	fields := strings.Split(line, ",")
	if len(fields) < 7 {
		b.logger.Warn("Invalid Beast frame, too few fields", "fields", len(fields))
		return nil, nil
	}

	hex := fields[0]
	lat := fields[1]
	lon := fields[2]
	alt := fields[3]
	gs := fields[4]
	track := fields[5]
	seen := fields[6]

	result := map[string]string{
		"hex":     hex,
		"lat":     lat,
		"lon":     lon,
		"alt":     alt,
		"gs":      gs,
		"track":   track,
		"seen":    seen,
	}

	return result, nil
}

// Close closes the connection
func (b *BeastReader) Close() error {
	if b.conn != nil {
		return b.conn.Close()
	}
	return nil
}

// ReadLoop continuously reads frames in a goroutine
func (b *BeastReader) ReadLoop(ch chan<- map[string]string) {
	defer close(ch)
	for {
		frame, err := b.ReadFrame()
		if err != nil {
			b.logger.Error("Error reading Beast frame", "error", err)
			time.Sleep(b.readRate)
			continue
		}
		if frame != nil {
			ch <- frame
		}
		time.Sleep(b.readRate)
	}
}