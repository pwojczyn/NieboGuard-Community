package rtl

import (
	"fmt"
	"log/slog"
	"runtime"

	"github.com/pwojczyn/NieboGuard-Community/apps/community-station/internal/config"
)

// Receiver manages RTL-SDR USB hardware communication.
type Receiver struct {
	logger     *slog.Logger
	config     config.RTLSDRConfig
	active     bool
	deviceName string
	serial     string
}

// NewReceiver creates a new RTL-SDR receiver.
func NewReceiver(cfg config.RTLSDRConfig, logger *slog.Logger) (*Receiver, error) {
	return &Receiver{
		logger:    logger,
		config:    cfg,
		active:    false,
	}, nil
}

// Initialize attempts to detect and open the RTL-SDR device.
func (r *Receiver) Initialize() error {
	r.logger.Info("RTL-SDR initialization started")

	// In a real implementation, this would detect USB devices
	// For now, we'll simulate detection
	deviceName, serial, found := detectDevice()

	if !found {
		return fmt.Errorf("no RTL-SDR device detected")
	}

	r.deviceName = deviceName
	r.serial = serial
	r.active = true
	r.logger.Info("RTL-SDR device detected", "name", deviceName, "serial", serial)
	return nil
}

// Start begins ADS-B reception.
func (r *Receiver) Start() error {
	if !r.active {
		r.active = true
		r.logger.Info("RTL-SDR reception started")
		return nil
	}
	return fmt.Errorf("receiver already active")
}

// Stop stops ADS-B reception.
func (r *Receiver) Stop() error {
	r.active = false
	r.logger.Info("RTL-SDR reception stopped")
	return nil
}

// Status returns the current receiver status.
func (r *Receiver) Status() map[string]interface{} {
	status := map[string]interface{}{
		"device_name":  r.deviceName,
		"device_serial": r.serial,
		"connected":    r.active,
		"status":       "active",
		"platform":     runtime.GOOS,
		"arch":         runtime.GOARCH,
	}
	return status
}

// detectDevice simulates device detection
func detectDevice() (string, string, bool) {
	// Simulate device detection
	deviceName := "RTL2832U (R820T2)"
	serial := "SN12345678"
	return deviceName, serial, true
}

// CheckDeviceBusy checks if another program is using the device.
func CheckDeviceBusy() bool {
	// In a real implementation, check running processes
	return false
}