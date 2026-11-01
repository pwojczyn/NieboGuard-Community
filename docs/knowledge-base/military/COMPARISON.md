# Comparison with Other Standards

## Link 16 vs ADS-B vs ASTERIX

| Feature | Link 16 | ADS-B | ASTERIX CAT021 |
|---------|---------|--------|---------------|
| Purpose | Military data link | Civil aviation surveillance | Military position reporting |
| Standard | NATO STANAG 5516 | FAA/ICAO | EUROCONTROL SPEC-014 |
| Data Format | Binary | JSON/Beast | Binary |
| Security | Encrypted | Basic | None |
| Range | 200-400 km | 200-400 km | Depends on source |
| Real-time | Yes | Yes | Yes |
| Simulation | Educational only | Educational | Educational |
| Simulation in NieboGuard | ✅ Link 16 Simulator | ✅ ADS-B Simulator | ✅ ASTERIX Simulator |

## Link 16 vs Link 11 vs Link 22

| Feature | Link 11 | Link 16 | Link 22 |
|---------|---------|---------|---------|
| Data Rate | 1200-2400 bps | 50-100 kbps | Up to 1.2 Mbps |
| Security | Basic | Basic | Enhanced |
| Coalition Support | Limited | Limited | Full |
| Maturity | Established | Mature | Next-gen |
| Complexity | Medium | Low | High |
| Simulation | Educational | Educational | Educational |

## ADS-B vs ASTERIX

| Feature | ADS-B | ASTERIX CAT021 |
|---------|--------|---------------|
| Purpose | Civil aviation | Military position reporting |
| Standard | FAA/ICAO | EUROCONTROL |
| Format | JSON/Beast | Binary |
| Security | Basic | None |
| Real-time | Yes | Yes |
| Simulation | ✅ | ✅ |

## Military vs Civilian Standards

| Aspect | Military (Link 16, DIS) | Civilian (ADS-B, ASTERIX) |
|--------|-------------------------|---------------------------|
| Purpose | Tactical operations | Surveillance |
| Security | High | Low |
| Standardization | NATO/STANAG | FAA/ICAO |
| Interoperability | High (military) | Limited (civilian) |
| Simulation | ✅ Educational | ✅ Educational |
| Real-time | ✅ | ✅ |

## Implementation Status

- **Link 16** — ✅ Simulator implemented (educational)
- **Link 11** — ✅ Documentation created
- **Link 22** — ✅ Documentation created
- **JREAP** — ✅ Documentation created
- **DIS** — ✅ Documentation created
- **HLA** — ✅ Documentation created
- **SISO Simulation** — ✅ Documentation and implementation completed

## Status Summary

- **Documentation** — ✅ Complete for all standards
- **Simulator** — ✅ Link 16 educational simulator implemented
- **Integration** — Planned for future
- **Testing** — In progress

## Status Summary

- **Documentation** — ✅ Complete
- **Simulator** — ✅ Link 16 educational simulator implemented
- **Integration** — Planned
- **Testing** — In progress