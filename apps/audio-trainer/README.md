# NieboGuard Audio AI Trainer

## Overview

NieboGuard Audio AI Trainer is a Windows CPU-only application for audio AI model training, built as part of the NieboGuard ecosystem. It provides:

- **CPU-only training** without GPU requirements
- **Data distribution** via Cloudflare R2 storage
- **Model management** with base models and transfer learning
- **Windows 10/11** application with no external dependencies
- **Real-time monitoring** of training progress and system metrics

## Features

### Core Functionality

- **Dataset Management** - Import, prepare, and distribute audio datasets
- **Model Training** - Train audio classification models locally on CPU
- **Checkpointing** - Resume interrupted training
- **Result Upload** - Upload trained models to R2
- **Web Dashboard** - Real-time monitoring and control

### Training Modes

#### CPU Lite
- Lightweight model training
- Pre-extracted features (log-mel, embeddings)
- Small datasets (25-50MB)
- Fast training for experimentation

#### CPU Standard
- Full model training
- Optional YAMNet integration
- Medium datasets (100-200MB)
- Production-ready training

### System Features

- **CPU Throttling** - Control CPU usage (25%, 50%, 75%, 100%)
- **RAM Limits** - Configurable memory limits
- **Cache System** - Smart caching of datasets and models
- **Offline Support** - Works without internet (resume functionality)
- **Error Handling** - Detailed error reporting and recovery

## Architecture

```
┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│  Cloudflare    │     │  Cloudflare    │     │  Cloudflare    │
│     R2         │     │      D1        │     │   Workers      │
└─────────────────┘     └─────────────────┘     └─────────────────┘
         │                       │                       │
         │                       │                       │
         └───────────────────────┼───────────────────────┘
                                 │
                                 ▼
                      ┌─────────────────┐
                      │ AI Lab (Admin)  │
                      │                 │
                      │  - Dataset     │
                      │    Registry    │
                      │  - Task Queue  │
                      │  - Model Hub   │
                      │  - Monitoring  │
                      └─────────────────┘
                                 │
                                 ▼
                      ┌─────────────────┐
                      │ NieboGuard      │
                      │ Audio AI Trainer│
                      │ (Windows EXE)   │
                      │                 │
                      │  - Dataset     │
                      │    Downloader  │
                      │  - Model       │
                      │    Manager     │
                      │  - Trainer     │
                      │  - Checkpoint  │
                      │  - Uploader    │
                      └─────────────────┘
                                 │
                                 ▼
                      ┌─────────────────┐
                      │   Cloudflare   │
                      │       R2        │
                      │                 │
                      │  - Datasets    │
                      │  - Base Models │
                      │  - Checkpoints │
                      │  - Results     │
                      └─────────────────┘
```

## Technical Requirements

### System Requirements

- **OS**: Windows 10/11 (x64)
- **CPU**: Intel/AMD with SSE4.2+ (no GPU required)
- **RAM**: Minimum 8GB (16GB+ recommended)
- **Disk Space**: 500MB+ (depending on dataset size)
- **Internet**: Required for initial setup and model download

### Software Requirements

- **Go** (for building the application)
- **No external dependencies** required at runtime
- **No Node.js, Electron, or Python** required

## Installation

### Building from Source

```bash
# Clone the repository
cd /path/to/NieboGuard-Community

# Build for Windows
mkdir -p build
go build -o build/NieboGuard-Audio-Trainer.exe ./apps/audio-trainer/cmd/audio-trainer/

# The executable will be available at:
# build/NieboGuard-Audio-Trainer.exe
```

### Running

```bash
# First run - will open config wizard
build/NieboGuard-Audio-Trainer.exe

# Or with specific config
build/NieboGuard-Audio-Trainer.exe --config C:\MyConfig.yaml
```

## Configuration

### Default Configuration

The application comes with sensible defaults:

```yaml
port: 8788
host: 127.0.0.1
workdir: "%LOCALAPPDATA%\NieboGuard\AudioTrainer"
language: "pl"
theme: "dark"

# Training settings
training:
  enabled: false
  mode: "cpu-lite"
  epochs: 50
  batch_size: 32
  learning_rate: 0.001
  max_threads: 4
  cpu_limit_percent: 70

# Dataset settings
dataset:
  import_dir: "imports"
  manifest_path: "manifest.yaml"
  split_train: 70
  split_val: 20
  split_test: 10
  resample_hz: 16000
  mono: true
  normalize: true
```

### Configuration Options

- `port` - Port for web dashboard
- `host` - Host for web dashboard (127.0.0.1 for local only)
- `workdir` - Application data directory
- `language` - UI language (pl/en)
- `theme` - UI theme (light/dark)
- `training.*` - Training parameters
- `dataset.*` - Dataset processing parameters
- `compute.*` - Compute integration settings
- `consent.*` - User consent settings

## Usage

### Dashboard

The web dashboard is available at:
```
http://127.0.0.1:8788/
```

**Main Features:**

1. **Task Management**
   - Check for new tasks
   - Start, pause, or stop training
   - View task status and progress

2. **Dataset Management**
   - Import audio files
   - Verify dataset integrity
   - View dataset statistics

3. **Model Management**
   - Load pre-trained models
   - Save checkpoints
   - View model information

4. **System Monitoring**
   - CPU and RAM usage
   - Training progress
   - Network activity
   - Log monitoring

### Dataset Import

1. Click "Import Dataset" button
2. Select audio files (WAV format)
3. Wait for processing and verification
4. Dataset becomes available for training

### Training

1. Click "Check Tasks" to see available tasks
2. Select a task and click "Start Training"
3. Monitor progress in real-time
4. Use "Pause" or "Stop" as needed

### Checkpointing

Training automatically saves checkpoints:

- After each epoch
- When interrupted
- On application shutdown

Training can be resumed by selecting the appropriate checkpoint.

## System Features

### CPU Throttling

Control CPU usage from the dashboard:

- **25%** - Minimal CPU usage, slower training
- **50%** - Balanced performance
- **75%** - Better performance, higher CPU usage
- **100%** - Maximum performance

### Cache System

The application uses a smart caching system:

- **Dataset Cache** - Stores downloaded datasets locally
- **Model Cache** - Stores trained models
- **Automatic Cleanup** - Removes unused caches
- **Integrity Verification** - Checks file hashes before use

### Error Handling

The application provides detailed error messages:

- Network errors with retry suggestions
- Disk space warnings
- File format validation
- License and permission errors
- System resource limits

## Advanced Features

### Command Line Arguments

```bash
build/NieboGuard-Audio-Trainer.exe [options]

Options:
  --config <path>     Path to config file (default: use workdir)
  --port <port>       Override default port
  --host <host>       Override default host
  --nopanel          Do not open browser automatically
```

### API

The application exposes a REST API:

- `GET /api/v1/health` - Health check
- `GET /api/v1/ready` - Readiness check
- `GET /api/v1/training/status` - Get training status
- `POST /api/v1/training/start` - Start training
- `POST /api/v1/training/pause` - Pause training
- `GET /api/v1/dataset/history` - Get dataset import history
- `POST /api/v1/dataset/import` - Import dataset
- `POST /api/v1/dataset/verify` - Verify dataset
- `GET /api/v1/model/status` - Get model status
- `GET /api/v1/system/metrics` - Get system metrics

### WebSocket

For real-time updates:

- Connect to `/ws`
- Receive training progress updates
- Receive system notifications

## Data Distribution

The application integrates with Cloudflare R2 for:

### Model Distribution

- Download base models from R2
- Upload checkpoints and final models
- Retrieve trained models for validation

### Dataset Distribution

- Download training datasets from R2
- Upload verified datasets to R2
- Maintain dataset versioning

### Security

- All R2 operations use signed URLs
- API keys are never exposed
- Data is encrypted during transfer

## Troubleshooting

### Common Issues

#### "Cannot bind to port"

**Solution:** Port 8788 is already in use.

- Check running processes: `netstat -an | findstr 8788`
- Kill the process: `taskkill /F /PID <PID>`
- Or use a different port in config

#### "Disk space full"

**Solution:** Free up disk space or adjust cache settings.

#### "Audio format not supported"

**Solution:** Convert audio to WAV format with 16kHz sample rate.

### Logs

Logs are stored in the workdir:

```
%LOCALAPPDATA%\NieboGuard\AudioTrainer\logs\
```

Common log files:

- `app.log` - Application logs
- `training.log` - Training progress
- `dataset.log` - Dataset processing
- `system.log` - System events

## Testing

### Unit Tests

```bash
cd apps/audio-trainer
go test ./...
```

### Integration Tests

Integration tests verify end-to-end functionality:

- Dataset import and processing
- Model training and checkpointing
- Result upload and validation
- Dashboard functionality

## Development

### Building

```bash
# Build for current OS/architecture
make build

# Build for Windows
make build-windows

# Build for Linux
make build-linux
```

### Dependencies

All dependencies are managed via Go modules. No external dependencies required at runtime.

### Code Structure

```
apps/audio-trainer/
├── cmd/
│   └── audio-trainer/
│       └── main.go           # Application entry point
├── internal/
│   ├── config/              # Configuration
│   ├── dashboard/           # Web dashboard
│   ├── dataset/             # Dataset management
│   ├── model/               # Model management
│   ├── training/            # Training logic
│   └── ...
├── web/                     # HTML/CSS/JS frontend
├── go.mod                  # Go module file
└── go.work                 # Work file for multiple modules
```

## Status

This is an early release. The application is functional but may have:

- Stability issues
- Limited feature set
- Missing error handling
- Incomplete testing

## Known Issues

1. **Port conflicts** - May conflict with other NieboGuard applications
2. **Cache management** - Cache cleanup needs improvement
3. **Error recovery** - Some error scenarios need better handling
4. **Documentation** - Some areas need more detailed documentation

## Future Work

### High Priority

1. **R2 Integration** - Complete integration with Cloudflare R2
2. **AI Lab Integration** - Connect with AI Lab for task management
3. **Model Hub** - Implement model registry and distribution
4. **Dashboard Improvements** - Enhance UI/UX
5. **Performance Optimization** - Improve training speed

### Medium Priority

1. **Offline Mode** - Full offline support
2. **Advanced Features** - Multi-model support, hyperparameter optimization
3. **Scalability** - Support for larger datasets
4. **Security** - Enhanced security features

### Low Priority

1. **Platform Support** - Add macOS support
2. **Cloud Integration** - Integration with cloud providers
3. **Mobile Support** - Mobile app development

## Support

For issues and support:

1. Check the documentation
2. Review logs in the workdir
3. File an issue in the repository
4. Contact the development team

## License

This project is part of the NieboGuard ecosystem. See the LICENSE file for details.

## Acknowledgments

- [YAMNet](https://github.com/tensorflow/models/tree/master/research/audioset) - Audio feature extractor
- [Go](https://golang.org/) - Programming language
- [Gorilla/mux](https://github.com/gorilla/mux) - HTTP router
- [Gorilla/websocket](https://github.com/gorilla/websocket) - WebSocket library
- [Cloudflare Workers](https://workers.cloudflare.com/) - Serverless platform
- [Cloudflare D1](https://developers.cloudflare.com/d1/) - Database
- [Cloudflare R2](https://developers.cloudflare.com/r2/) - Object storage

## Version History

### Version 0.1.0 (Development)

- Initial implementation
- Basic functionality
- Web dashboard
- Dataset import
- CPU-only training

### Future Versions

- Version 0.2.0 - R2 integration
- Version 0.3.0 - AI Lab integration
- Version 0.4.0 - Production features
- Version 1.0.0 - Stable release

## Contributing

See CONTRIBUTING.md for contribution guidelines.

---

**Built with Go, following NieboGuard coding standards and practices.**

---

*Page last updated: 2026-11-02*
