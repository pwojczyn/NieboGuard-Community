# NieboGuard R2 Architecture — Automatic Data and Model Distribution

## Overview

The NieboGuard Audio AI Trainer uses Cloudflare R2 for data and model distribution. The system follows a distributed architecture with:

1. **Cloudflare R2** - Object storage for datasets, models, and artifacts
2. **Cloudflare D1** - Database for task management, user permissions, and manifest validation
3. **Cloudflare Workers** - Serverless functions handling API requests and data transfers
4. **NieboGuard AI Lab** - Web-based admin interface for task creation and monitoring
5. **NieboGuard Audio AI Trainer (Windows)** - Local application that downloads data, trains models, and uploads results

## Data Structure

### Logical Prefixes

| Prefix | Description | Example Path |
|--------|-------------|--------------|
| `compute/datasets/` | Dataset registry and source data | `compute/datasets/admin/` |
| `compute/packs/` | Training data packs (MINI/STANDARD/EXTENDED) | `compute/packs/mini/` |
| `compute/base-models/` | Base models for transfer learning | `compute/base-models/yamnet/` |
| `compute/quarantine/updates/` | Quarantined updates for review | `compute/quarantine/updates/2026/` |
| `compute/aggregated/` | Aggregated results from multiple users | `compute/aggregated/round1/` |
| `compute/reports/` | Generated reports and metrics | `compute/reports/round1-report.json` |

## Data Flow

1. **Administrator** creates dataset and model in AI Lab
2. **AI Lab** stores metadata in D1 and data in R2
3. **Worker** validates manifest signatures (Ed25519)
4. **Windows app** downloads required data and model via signed URLs
5. **App** downloads missing shards from R2
6. **App** trains model locally on CPU
7. **App** uploads results to R2
8. **AI Lab** validates results, aggregates, and evaluates
9. **Administrator** approves or rejects model for production

## Security & Validation

- All data transfers use signed URLs with short expiration
- Manifests are signed with Ed25519 keys
- Dataset registration requires license verification
- Model registry requires architecture compatibility check
- All API endpoints require OAuth2 token validation
- No direct access to R2 buckets - all access through signed URLs

## Dataset Registry

Administrator must be able to view:
- Dataset name and ID
- Version number
- Source of origin
- Number of recordings
- Number of classes
- Data size
- License information
- Distribution rights
- Training rights
- Import date
- Preparation status
- Approval status
- Available packs

## Model Registry

- Versioning of base models
- Compatibility checks with datasets
- Architecture and tensor shape validation
- Versioned checkpoints
- Compatibility with YAMNet and other architectures

## API Contracts

The system uses the following API contracts (existing or to be implemented):

- `POST /api/v2/compute/enroll` - Register new task
- `POST /api/v2/compute/token` - Get authorization token
- `GET /api/v2/compute/tasks/next` - Get next task
- `POST /api/v2/compute/tasks/:id/heartbeat` - Report task progress
- `POST /api/v2/compute/tasks/:id/upload-url` - Upload data to R2
- `POST /api/v2/compute/tasks/:id/complete` - Mark task as complete
- `POST /api/v2/compute/tasks/:id/fail` - Mark task as failed
- `GET /api/v2/compute/objects/:uploadId` - Get object status
- `PUT /api/v2/compute/objects/:uploadId` - Update object status

## Cache System

- Cache based on SHA-256 hash and version
- Reuse identical packs between tasks
- Automatic detection of corrupted files
- Configurable disk usage limit
- Checkpoint persistence
- Separate storage for training data and results

## Network Resilience

- Retry mechanism with exponential backoff
- Timeout handling for network operations
- Error handling for 401/403/429/5xx responses
- Internet connectivity detection
- URL expiration detection and renewal
- Large transfer resumption (when supported)
- Secure temporary storage before atomic commit

## Windows CPU Edition Requirements

- Local CPU-only training (no GPU)
- CPU throttling controls
- RAM limits
- Checkpoint resume capability
- Offline operation support
- Windows 10/11 x64 compatibility
- No Electron or Node.js dependencies

## Implementation Status

- [x] Basic R2 architecture documentation created
- [ ] Dataset registry implementation
- [ ] Model registry implementation
- [ ] API contract implementation
- [ ] Windows app integration with R2
- [ ] Checkpoint and resume functionality
- [ ] Upload pipeline with R2
- [ ] AI Lab integration
- [ ] Test suite for all components

## Status

- **Status**: In Progress
- **Priority**: High
- **Dependencies**: Compute API contracts, R2 bucket setup
- **Next Steps**: Implement dataset registry and model download logic