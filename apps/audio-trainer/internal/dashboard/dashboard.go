package dashboard

import (
	"context"
	"runtime"
	"github.com/gorilla/websocket"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/mux"
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
}

type TrainingEvent struct {
	Type      string          `json:"type"`
	Timestamp time.Time       `json:"timestamp"`
	Payload   json.RawMessage `json:"payload"`
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
	}

	// Start background metrics collection
	go d.collectMetrics()

	// Start processing training events
	go d.processTrainingEvents()

	return d, nil
}

func (d *Dashboard) collectMetrics() {
	for {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
			d.historyMutex.Lock()
			d.systemMetrics.CPU = getCPUUsage() // Implement if needed
			d.systemMetrics.RAM = float64(m.Alloc) / (1024 * 1024) // MB
			d.systemMetrics.Uptime = time.Since(d.upTime).String()
			d.systemMetrics.Goroutines = runtime.NumGoroutine()
			d.historyMutex.Unlock()

		time.Sleep(5 * time.Second)
	}
}

func getCPUUsage() float64 {
	// Try to get actual CPU usage from /proc/stat on Linux
	// On Windows, this would use different APIs
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return m.TotalAlloc / (1024 * 1024) // MB
}

func (d *Dashboard) processTrainingEvents() {
	for {
		select {
		case event := <-d.trainingEvents:
			d.handleTrainingEvent(event)
		default:
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func (d *Dashboard) handleTrainingEvent(event TrainingEvent) {
	switch event.Type {
	case "training_start":
			d.logger.Printf("Training started")
	case "training_pause":
			d.logger.Printf("Training paused")
	case "training_complete":
			d.logger.Printf("Training complete")
	}
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
	}

	// Start background processes
	ctx, cancel := context.WithCancel(context.Background())
	go dashboard.collectMetrics()
	go dashboard.processTrainingEvents()

	return d, nil
}

func (d *Dashboard) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/", d.handleIndex).Methods("GET")
	r.HandleFunc("/api/v1/health", d.handleHealth).Methods("GET")
	r.HandleFunc("/api/v1/ready", d.handleReady).Methods("GET")
	r.HandleFunc("/api/v1/training/status", d.handleTrainingStatus).Methods("GET")
	r.HandleFunc("/api/v1/training/start", d.handleTrainingStart).Methods("POST")
	r.HandleFunc("/api/v1/training/pause", d.handleTrainingPause).Methods("POST")
	r.HandleFunc("/api/v1/dataset/history", d.handleDatasetHistory).Methods("GET")
	r.HandleFunc("/api/v1/dataset/list", d.handleDatasetList).Methods("GET")
	r.HandleFunc("/api/v1/model/status", d.handleModelStatus).Methods("GET")
	r.HandleFunc("/api/v1/system/metrics", d.handleSystemMetrics).Methods("GET")
	r.HandleFunc("/api/v1/compute/enroll", d.handleComputeEnroll).Methods("POST")
	r.HandleFunc("/api/v1/compute/tasks/next", d.handleNextTask).Methods("GET")
	r.HandleFunc("/api/v1/compute/tasks/:id/heartbeat", d.handleHeartbeat).Methods("POST")
	r.HandleFunc("/api/v1/r2/download", d.handleDownload).Methods("GET")
	r.HandleFunc("/api/v1/config", d.handleConfig).Methods("GET")
	r.HandleFunc("/ws", d.handleWebSocket).Methods("GET")
}

func (d *Dashboard) handleIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, filepath.Join("web", "index.html"))
}

func (d *Dashboard) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}

func (d *Dashboard) handleReady(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "{\"status\": \"ready\"}")
}

func (d *Dashboard) handleTrainingStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": d.modelManager.Status(),
		"info":   d.modelManager.GetInfo(),
	})
}

func (d *Dashboard) handleTrainingStart(w http.ResponseWriter, r *http.Request) {
	var req struct {
		TaskID string `json:"task_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	task, err := d.computeClient.NextTask()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	d.modelManager.SetStatus(model.StatusTraining)

	d.trainingEvents <- TrainingEvent{
		Type:      "training_start",
		Timestamp: time.Now(),
		Payload: json.RawMessage(fmt.Sprintf(`{"task_id": "%s"`, task.ID)),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "training_started", "task_id": task.ID})
}

func (d *Dashboard) handleTrainingPause(w http.ResponseWriter, r *http.Request) {
	if err := d.modelManager.Pause(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	d.trainingEvents <- TrainingEvent{
		Type:      "training_pause",
		Timestamp: time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "training_paused"})
}

func (d *Dashboard) handleDatasetImport(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	tempPath := filepath.Join(d.cfg.WorkDir, "imported.wav")
	outFile, err := os.Create(tempPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer outFile.Close()

	if _, err := io.Copy(outFile, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	wavInfo, err := d.datasetManager.ImportWAV(tempPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	d.trainingEvents <- TrainingEvent{
		Type:      "wav_imported",
		Timestamp: time.Now(),
		Payload: json.RawMessage(fmt.Sprintf(`{"path": "%s", "duration": %.2f`, tempPath, wavInfo.Duration)),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "imported", "duration": fmt.Sprintf("%.2f", wavInfo.Duration), "labels": strings.Join(wavInfo.Labels, ", ")})
}

func (d *Dashboard) handleDatasetHistory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(d.datasetManager.GetHistory())
}

func (d *Dashboard) handleDatasetList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	datasets := []map[string]string{
		{"id": "svanstrom-audio", "name": "Svanstrom Audio Dataset"},
		{"id": "other-dataset", "name": "Other Dataset"},
	}
	json.NewEncoder(w).Encode(datasets)
}

func (d *Dashboard) handleModelStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": d.modelManager.Status(),
		"info":   d.modelManager.GetInfo(),
	})
}

func (d *Dashboard) handleSystemMetrics(w http.ResponseWriter, r *http.Request) {
	d.historyMutex.Lock()
	defer d.historyMutex.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(d.systemMetrics)
}

func (d *Dashboard) handleConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(d.cfg)
}

func (d *Dashboard) handleComputeEnroll(w http.ResponseWriter, r *http.Request) {
	var req struct {
		StationID string `json:"station_id"`
		Version   string `json:"version"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := d.computeClient.Enroll(compute.EnrollRequest{
		StationID: req.StationID,
		Version:   req.Version,
		Platform:  "windows-cpu-edition",
		CPU:       4,
		RAM:       8,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (d *Dashboard) handleNextTask(w http.ResponseWriter, r *http.Request) {
	task, err := d.computeClient.NextTask()
	if err != nil {
		if err == compute.ErrNoTasks {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNoContent)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func (d *Dashboard) handleHeartbeat(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID := vars["id"]

	var req struct {
		Progress float64 `json:"progress"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := d.computeClient.Heartbeat(taskID, req.Progress)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "heartbeat_sent"})
}

func (d *Dashboard) handleDownload(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	objectPath := vars["object_path"]

	if objectPath == "" {
		http.Error(w, "missing object_path parameter", http.StatusBadRequest)
		return
	}

	url, err := d.r2DownloadMgr.GenerateSignedURL(objectPath, 24*time.Hour)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"download_url": url,
		"object_path":  objectPath,
	})
}

func (d *Dashboard) handleConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(d.cfg)
}

func (d *Dashboard) HandleWebSocket(w http.ResponseWriter, r *http.Request, upgrader *websocket.Upgrader) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	d.wsMutex.Lock()
	d.wsConnections[conn] = true
	d.wsMutex.Unlock()

	defer func() {
		d.wsMutex.Lock()
		delete(d.wsConnections, conn)
		d.wsMutex.Unlock()
		conn.Close()
	}()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}

		var wsMsg WebSocketMessage
		if err := json.Unmarshal(msg, &wsMsg); err != nil {
			continue
		}

		switch wsMsg.Type {
		case "subscribe_training":
			d.logger.Printf("Subscribed to training updates")
		case "request_status":
			status := map[string]interface{}{
				"training_status": d.modelManager.Status(),
				"system_metrics":  d.systemMetrics,
			}
			payload, _ := json.Marshal(status)
			conn.WriteMessage(websocket.TextMessage, payload)
		default:
			d.logger.Printf("Unknown WebSocket message type: %s", wsMsg.Type)
		}
	}
}