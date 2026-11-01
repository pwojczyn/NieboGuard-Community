# Data Flow — Protocol Lab

## Unified Data Flow

```
SOURCE → ADAPTER → RAW OBSERVATION → NORMALIZATION → TRACK ENGINE → OUTPUT ADAPTER → API/WSS/UI
```

## Stages

1. **SOURCE**  
   - Real: ADS-B, CAT021, CAT048, Link 16, etc.  
   - Synthetic: Simulator (all profiles)  
   - All inputs are `models.AircraftEvent` or `models.SimulatedObservation`  
   - `simulated=true` and `SimulationNote` preserved  

2. **ADAPTER**  
   - `adapter_adsb_json` (HTTP)  
   - `adapter_http_json` (polling)  
   - `adapter_adsb_beast` (TCP)  
   - `adapter_link16` (simulator)  
   - All adapters output raw data structure  

3. **RAW OBSERVATION**  
   - Raw JSON or binary data  
   - May be incomplete (missing fields)  
   - May have invalid coordinates  
   - May have duplicate entries  
   - May have TTL expiration  

4. **NORMALIZATION**  
   - `NormalizationService`  
   - Validates all required fields  
   - Fills missing required fields with defaults  
   - Converts coordinates to uniform units  
   - Applies time zone normalization (UTC)  
   - Adds `source_id` if missing  
   - Adds `timestamp` if missing  
   - Adds `quality` flag  
   - **Preserves `synthetic` flag and `simulation_note`**  

5. **TRACK ENGINE**  
   - `TrackManager`  
   - Groups observations by `source_id` + `flight`  
   - Updates position/velocity/heading  
   - Manages TTL (30s default)  
   - Detects duplicates (same `observation_id`)  
   - Handles track termination (TTL expired)  
   - Supports `Multi-Sensor` mode (multiple sources)  

6. **OUTPUT ADAPTER**  
   - `wss_` (WebSocket)  
   - `api/` (REST)  
   - `panel/` (Web UI)  
   - All output `models.AircraftEvent` with `synthetic=true` preserved  

7. **API/WSS/UI**  
   - `GET /api/v2/state` → normalized data  
   - `GET /sim/api/` → simulator control  
   - `ws://localhost:8787/sim/` → real-time updates  
   - Panel shows `TEST` tag on all synthetic objects