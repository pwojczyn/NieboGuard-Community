package errorgen

import (
	"fmt"
	"math/rand"
	"time"
)

type Scenario int

const (
	ScenarioNone Scenario = iota
	ScenarioPacketLoss5
	ScenarioPacketLoss10
	ScenarioPacketLoss25
	ScenarioPacketLoss50
	ScenarioDelay
	ScenarioDuplicates
	ScenarioOutOfOrder
	ScenarioTCPDrop
	ScenarioReconnect
	ScenarioMissingFields
	ScenarioOutOfRange
	ScenarioInvalidBinary
	ScenarioHighLoad
	ScenarioUTCDateChange
)

type ErrorGenerator struct {
	Scenario Scenario
	LossPct  float64
	DelayMs  int
	Duplicates bool
	OutOfOrder bool
	MissingFields bool
	OutOfRange bool
	InvalidBinary bool
	HighLoad bool
	UTCDateChange bool

	// State tracking
	lastSeq int
	dropped int
	delayed int
	dups    int
}

func New(scenario Scenario) *ErrorGenerator {
	eg := &ErrorGenerator{Scenario: scenario}
	switch scenario {
	case ScenarioPacketLoss5:
		eg.LossPct = 5
	case ScenarioPacketLoss10:
		eg.LossPct = 10
	case ScenarioPacketLoss25:
		eg.LossPct = 25
	case ScenarioPacketLoss50:
		eg.LossPct = 50
	case ScenarioDelay:
		eg.DelayMs = 200
	case ScenarioDuplicates:
		eg.Duplicates = true
	case ScenarioOutOfOrder:
		eg.OutOfOrder = true
	case ScenarioTCPDrop:
		// marked externally
	case ScenarioReconnect:
		// marked externally
	case ScenarioMissingFields:
		eg.MissingFields = true
	case ScenarioOutOfRange:
		eg.OutOfRange = true
	case ScenarioInvalidBinary:
		eg.InvalidBinary = true
	case ScenarioHighLoad:
		eg.HighLoad = true
	case ScenarioUTCDateChange:
		eg.UTCDateChange = true
	}
	return eg
}

func (e *ErrorGenerator) ShouldDrop() bool {
	if e.LossPct <= 0 {
		return false
	}
	return rand.Float64()*100 < e.LossPct
}

func (e *ErrorGenerator) GetDelay() int {
	if e.DelayMs > 0 {
		return e.DelayMs
	}
	return 0
}

func (e *ErrorGenerator) Apply(data interface{}) (interface{}, error) {
	// Return modified data based on scenario
	if e.Duplicates && e.lastSeq < 100 {
		e.lastSeq++
		// Return duplicate marker - caller should duplicate last item
		return data, fmt.Errorf("duplicate_marker")
	}
	if e.OutOfOrder {
		// reorder - caller shuffles
	}
	if e.InvalidBinary {
		return nil, fmt.Errorf("invalid_binary")
	}
	return data, nil
}

func (e *ErrorGenerator) GenerateInvalidBinary() ([]byte, error) {
	return []byte{0xFF, 0x00, 0xDE, 0xAD, 0xBE, 0xEF}, nil
}

func (e *ErrorGenerator) Stats() map[string]interface{} {
	return map[string]interface{}{
		"dropped": e.dropped,
		"delayed": e.delayed,
		"dups":    e.dups,
		"scenario": e.Scenario,
	}
}

var _ = time.Second // placeholder for potential time usage
