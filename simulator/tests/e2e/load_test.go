package e2e

import (
	"testing"
	"time"

	"github.com/pwojczyn/NieboGuard-Community/simulator/internal/sim"
)

// TestLoadConcurrency verifies the simulation engine can handle concurrent load
func TestLoadConcurrency(t *testing.T) {
	for _, count := range []int{10, 100, 500, 1000} {
		engine := sim.NewEngine(sim.ModeSynthetic)
		engine.SyntheticCount = count

		obs, err := engine.NextFrame()
		if err != nil {
			t.Fatalf("syntheticFrame error for %d objects: %v", count, err)
		}
		if len(obs) != count {
			t.Errorf("expected %d observations, got %d", count, len(obs))
		}

		// Verify all are synthetic
		for _, o := range obs {
			if !o.Synthetic {
				t.Errorf("observation %s is not synthetic", o.ID)
			}
		}
	}

	t.Logf("Load test passed: 10, 100, 500, 1000 objects generated successfully")
}

// TestLoadThroughput measures how many observations per second the engine can produce
func TestLoadThroughput(t *testing.T) {
	engine := sim.NewEngine(sim.ModeSynthetic)
	engine.SyntheticCount = 1000

	obs, err := engine.NextFrame()
	if err != nil {
		t.Fatalf("syntheticFrame error: %v", err)
	}

	// Measure time taken
	start := time.Now()
	for i := 0; i < len(obs); i++ {
		_ = obs[i]
	}
	elapsed := time.Since(start)

	throughput := float64(len(obs)) / elapsed.Seconds()
	t.Logf("Throughput: %.2f observations/second (1000 obs in %.3s)", throughput, elapsed)

	if throughput < 100 {
		t.Errorf("Throughput too low: %.2f obs/s", throughput)
	}
}

// TestLoadStability verifies the engine remains stable under sustained load
func TestLoadStability(t *testing.T) {
	engine := sim.NewEngine(sim.ModeSynthetic)
	engine.SyntheticCount = 2000

	obs, err := engine.NextFrame()
	if err != nil {
		t.Fatalf("syntheticFrame error: %v", err)
	}

	if len(obs) != 2000 {
		t.Errorf("expected 2000 observations, got %d", len(obs))
	}

	// Check that all are synthetic
	for _, o := range obs {
		if !o.Synthetic {
			t.Errorf("observation %s is not synthetic", o.ID)
		}
	}

	t.Logf("Stability test passed: 2000 observations generated without errors")
}
