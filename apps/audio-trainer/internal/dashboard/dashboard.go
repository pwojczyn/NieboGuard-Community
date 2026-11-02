package dashboard

import (
	"encoding/json"
	"log"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/pwojczyn/NieboGuard-Community/apps/audio-trainer/internal/config"
	"github.com/pwojczyn/NieboGuard-Community/apps/audio-trainer/internal/compute"
	"github.com/pwojczyn/NieboGuard-Community/apps/audio-trainer/internal/dataset"
	"github.com/pwojczyn/NieboGuard-Community/apps/audio-trainer/internal/model"
	"github.com/pwojczyn/NieboGuard-Community/apps/audio-trainer/internal/r2"
)

type SystemMetrics struct {
	CPU        float64 `json:"cpu"`
	RAM        float64 `json:"ram"`
	Uptime      string  `json:"uptime"`
	Goroutines  int     `json:"goroutines"`
	Temperature float64 `json:"temperature"`
}

type TrainingEvent struct {
	Type      string                 `json:"type"`
	Timestamp time.Time              `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
}

type WebSocketMessage struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type Dashboard struct {
	cfg           *config.Configuration
	computeClient *compute.Client
	datasetManager *dataset.DatasetManager
	modelManager   *model.ModelManager
	r2DownloadMgr *r2.DownloadManager
	historyMutex sync.RWMutex
	systemMetrics  SystemMetrics
	upTime         time.Time
	wsConnections  map[interface{}]bool
	wsMutex       sync.RWMutex
	trainingEvents chan TrainingEvent
	logger        *log.Logger
}

func New(cfg *config.Configuration, computeClient *compute.Client, datasetManager *dataset.DatasetManager, modelManager *model.ModelManager, r2DownloadMgr *r2.DownloadManager) (*Dashboard, error) {
	d := &Dashboard{
		cfg:           cfg,
		computeClient: computeClient,
		datasetManager: datasetManager,
		modelManager:  modelManager,
		r2DownloadMgr: r2DownloadMgr,
		wsConnections: make(map[interface{}]bool),
		upTime:        time.Now(),
		trainingEvents: make(chan TrainingEvent, 100),
		logger:        log.New(os.Stdout, "[DASHBOARD] ", log.LstdFlags),
	}

	// Start background processes
	go d.CollectMetrics()
	go d.ProcessTrainingEvents()

	return d, nil
}

func (d *Dashboard) CollectMetrics() {
	for {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		d.historyMutex.Lock()
		d.systemMetrics.CPU = getCPUUsage()
		d.systemMetrics.RAM = float64(m.Alloc) / (1024 * 1024) // MB
		d.systemMetrics.Uptime = time.Since(d.upTime).String()
		d.systemMetrics.Goroutines = runtime.NumGoroutine()
		d.historyMutex.Unlock()

		time.Sleep(5 * time.Second)
	}
}

func (d *Dashboard) ProcessTrainingEvents() {
	for {
		select {
		case evt := <-d.trainingEvents:
			d.logger.Printf("Training event: %s - %v", evt.Type, evt.Data)
		case <-time.After(10 * time.Second):
		}
	}
}

func (d *Dashboard) RegisterRoutes(r interface{}) {
	// Routes registered via main router integration
}

func getCPUUsage() float64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return float64(m.Alloc) / (1024 * 1023)
}