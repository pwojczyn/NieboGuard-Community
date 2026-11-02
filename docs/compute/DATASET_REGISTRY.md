# NieboGuard Dataset Registry

## Overview

The Dataset Registry is a centralized system for managing all audio datasets used in the NieboGuard Audio AI Trainer. It handles:

- Dataset metadata and licensing
- Preparation status tracking
- Version control and history
- Classification into packs (MINI/STANDARD/EXTENDED)
- Integration with R2 storage
- Validation and preparation pipeline

## Dataset Metadata

Each dataset has the following metadata fields:

| Field | Type | Description | Example |
|-------|------|-------------|---------|
| `id` | string | Unique dataset identifier | `svanstrom-audio-v1` |
| `name` | string | Human-readable name | `Svanstrom Audio Dataset v1.0` |
| `version` | string | Version number | `v1.0` |
| `source` | string | Origin/source information | `Freesound` |
| `recordings` | integer | Number of audio files | `45000` |
| `classes` | integer | Number of distinct classes | `40` |
| `size_mb` | number | Dataset size in MB | `1800.5` |
| `license` | string | License terms | `CC-BY-SA 3.0` |
| `license_url` | string | URL to license text | `https://creativecommons.org/licenses/by-sa/3.0/` |
| `distribution_rights` | boolean | Can dataset be distributed | `true` |
| `training_rights` | boolean | Can dataset be used for training | `true` |
| `import_date` | date | Date dataset was imported | `2025-06-15` |
| `preparation_status` | string | Current preparation state | `ready` |
| `approval_status` | string | Administrator approval status | `approved` |
| `packs` | object | Available packs with metadata | `{ "mini": {...}, "standard": {...}, "extended": {...} }` |

## Dataset Preparation Pipeline

The preparation pipeline processes datasets through several stages:

1. **Import** - Dataset uploaded by administrator via AI Lab
2. **Format Validation** - WAV files, metadata, class consistency
3. **Audio Processing** - Resampling to 16kHz, normalization, format conversion
4. **Data Cleaning** - Remove duplicates, fix corrupted files
5. **Label Verification** - Ensure labels match class count
6. **Split Preparation** - Create train/validation/test splits
7. **Shard Creation** - Divide into manageable chunks
8. **Hash Calculation** - SHA-256 for integrity verification
9. **Publication** - Store in R2 buckets

## Pack System

Datasets are divided into training packs with different sizes:

### MINI Pack
- Size: 25-50 MB
- Classes: 5-10
- Recordings: 500-2000
- Use case: Quick training, evaluation, experimentation

### STANDARD Pack
- Size: 100-200 MB
- Classes: 10-30
- Recordings: 2000-5000
- Use case: Standard training runs

### EXTENDED Pack
- Size: 250-500 MB
- Classes: 30-50
- Recordings: 5000-15000
- Use case: Large-scale training, model refinement

### GPU Pack
- Size: 500 MB+
- Classes: 50+
- Recordings: 15000+
- Note: For GPU instances, not available in CPU Edition

## Status Enums

### PreparationStatus
- `pending` - Awaiting administrator review
- `in_progress` - Currently being prepared
- `processing` - Audio processing stage
- `splitting` - Data split preparation
- `hashing` - Integrity hash calculation
- `publishing` - Publishing to R2
- `ready` - Dataset ready for training
- `failed` - Preparation failed
- `cancelled` - Preparation cancelled

### ApprovalStatus
- `pending` - Awaiting administrator approval
- `approved` - Administrator approved
- `rejected` - Administrator rejected
- `requires_revision` - Needs changes
- `expired` - Approval period expired

### SplitStage
- `pending` - Split not yet calculated
- `calculating` - Split calculation in progress
- `ready` - Splits calculated
- `failed` - Split calculation failed

## Package JSON Structure

Each pack has the following metadata:

```json
{
  "pack_type": "mini",
  "size_mb": 32,
  "classes": 8,
  "recordings": 1500,
  "files": [
    {
      "id": "svanstrom_0001.wav",
      "path": "compute/datasets/svanstrom/v1.0/shards/mini/part-000000.wav",
      "sha256": "a1b2c3d4e5f6...",
      "size_bytes": 1048576,
      "sample_rate": 16000,
      "duration_seconds": 5.0,
      "verified": true,
      "labels": ["car", "motorcycle", "engine", ...]
    }
  ],
  "split": {
    "train_percentage": 70,
    "val_percentage": 20,
    "test_percentage": 10,
    "train_recordings": 1050,
    "val_recordings": 300,
    "test_recordings": 150
  }
}
```

## Versioning

Datasets use semantic versioning:
- **MAJOR** - Breaking changes
- **MINOR** - New features, backward compatible
- **PATCH** - Bug fixes, minor improvements

Version updates increment the whole number, and SHA256 hashes change for all files.

## License Validation

The system validates licenses:

- **CC-BY-SA 3.0** - Free for commercial use, requires attribution
- **CC-BY 4.0** - Free for commercial use, requires attribution
- **CC0 1.0** - Public domain, no attribution required
- **MIT** - Free software license, compatible with many projects
- **Apache 2.0** - Free software license, patent grant

Datasets with non-permissive licenses cannot be used for public training tasks.

## Integration with R2

Dataset files are stored in the following R2 paths:

- `compute/datasets/{dataset_id}/{version}/pack/{pack_type}/{shard}/` - Training packs
- `compute/datasets/{dataset_id}/{version}/manifest.json` - Dataset manifest
- `compute/datasets/{dataset_id}/{version}/metadata.yaml` - Metadata
- `compute/datasets/{dataset_id}/{version}/logs/` - Processing logs
- `compute/datasets/{dataset_id}/{version}/temp/` - Temporary files

## Validation Requirements

Before a dataset is marked as `ready`:

1. **Audio Quality** - All files are valid WAV files
2. **Sample Rate** - All files are resampled to 16000 Hz
3. **Normalization** - Audio is normalized
4. **Class Consistency** - All labels match declared classes
5. **File Integrity** - SHA256 hashes match
6. **Split Validity** - Splits are balanced and stratified
7. **Access Rights** - License permits distribution and training

## API Endpoints

### Dataset Registry APIs

| Method | Path | Description |
|--------|------|-------------| |
| GET | `/api/v2/datasets` | List all datasets |
| GET | `/api/v2/datasets/{id}` | Get dataset metadata |
| GET | `/api/v2/datasets/{id}/versions` | Get all versions |
| GET | `/api/v2/datasets/{id}/{version}/pack/{pack_type}` | Get dataset pack |
| POST | `/api/v2/datasets` | Create new dataset |
| PUT | `/api/v2/datasets/{id}` | Update dataset |
| DELETE | `/api/v2/datasets/{id}` | Delete dataset |
| POST | `/api/v2/datasets/{id}/{version}/prepare` | Start preparation |
| GET | `/api/v2/datasets/{id}/{version}/status` | Get preparation status |
| GET | `/api/v2/datasets/{id}/{version}/packs` | List available packs |

## Integration Points

### With AI Lab
- Dataset creation, approval, and status management
- Dataset selection for task creation

### With R2
- Upload of raw dataset files
- Download of prepared packs for Windows app

### With Windows App
- Fetch list of available datasets
- Query dataset preparation status
- Download required packs
- Verify SHA256 hashes

### With Training Pipeline
- Dataset validation before training
- Monitoring of dataset preparation progress
- Retry mechanism for failed validations

## Monitoring & Logging

The system tracks:

- Dataset preparation progress
- File processing times
- Memory usage during processing
- Disk usage
- Network transfer times
- Failed file hashes
- Preparation error logs

## Cache System

Dataset packs are cached:

- **SHA256 Cache** - Maps file hashes to R2 paths
- **Version Cache** - Tracks dataset versions
- **Status Cache** - Caches preparation status
- **Pack Cache** - Tracks available packs per dataset

## Security

- All dataset uploads require authorization
- License validation is performed before publication
- Dataset metadata includes usage rights
- Access to sensitive metadata is controlled

## Migration

When moving datasets:

1. Update all references
2. Re-hash all files
3. Update version number
4. Invalidate old caches
5. Update references in R2

## Status

- **Status**: In Progress
- **Dependencies**: R2 bucket setup, Dataset preparation pipeline
- **Next Steps**: Implement dataset creation API, R2 integration

## TODO

- [ ] Implement dataset registry API endpoints
- [ ] Create dataset preparation pipeline
- [ ] Implement R2 bucket integration
- [ ] Add license validation
- [ ] Create dataset management UI
- [ ] Implement cache system
- [ ] Add monitoring and logging