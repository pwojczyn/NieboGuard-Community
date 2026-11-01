# Protocol Compatibility Matrix

## Legend

| Status | Meaning |
|--------|---------|
| DOCUMENTED | Public specification documented |
| IMPLEMENTED | Encoder/decoder exists in codebase |
| TESTED | Successfully tested with test data |
| INTEROPERABLE | Confirmed exchange with independent implementation |
| HARDWARE VERIFIED | Tested with real device or source |
| EDUCATIONAL ONLY | Educational model only, no compliance claims |

## Protocols

### ADS-B 1090ES
| Aspect | Status |
|--------|--------|
| Documentation | DOCUMENTED |
| Adapter (HTTP JSON) | IMPLEMENTED |
| Adapter (Beast TCP) | IMPLEMENTED (stub) |
| Decoder | IMPLEMENTED |
| Normalization | IMPLEMENTED |
| Tests | TESTED (fixture-based) |
| Interoperability | INTEROPERABLE (readsb/dump1090 reference) |
| Hardware | NOT VERIFIED (no RTL-SDR available) |
| Compliance | EDUCATIONAL ONLY (no production claim) |

### ASTERIX CAT021
| Aspect | Status |
|--------|--------|
| Documentation | DOCUMENTED |
| Encoder | IMPLEMENTED |
| Decoder | NOT IMPLEMENTED |
| Tests | NOT TESTED |
| Interoperability | NOT VERIFIED |
| Compliance | EDUCATIONAL ONLY (simplified FSPEC) |

### ASTERIX CAT048
| Aspect | Status |
|--------|--------|
| Documentation | DOCUMENTED |
| Encoder | IMPLEMENTED |
| Decoder | NOT IMPLEMENTED |
| Tests | NOT TESTED |
| Interoperability | NOT VERIFIED |
| Compliance | EDUCATIONAL ONLY |

### ASTERIX CAT062
| Aspect | Status |
|--------|--------|
| Documentation | DOCUMENTED |
| Encoder | IMPLEMENTED |
| Decoder | NOT IMPLEMENTED |
| Tests | NOT TESTED |
| Interoperability | NOT VERIFIED |
| Compliance | EDUCATIONAL ONLY |

### Link 16 (STANAG 5516)
| Aspect | Status |
|--------|--------|
| Documentation | DOCUMENTED |
| Simulator | IMPLEMENTED (educational) |
| Tests | TESTED (unit tests) |
| Interoperability | NOT VERIFIED |
| Compliance | EDUCATIONAL ONLY |

### Link 11 / Link 11B
| Aspect | Status |
|--------|--------|
| Documentation | DOCUMENTED |
| Simulator | PLANNED |
| Compliance | EDUCATIONAL ONLY |

### Link 22
| Aspect | Status |
|--------|--------|
| Documentation | DOCUMENTED |
| Simulator | PLANNED |
| Compliance | EDUCATIONAL ONLY |

### JREAP
| Aspect | Status |
|--------|--------|
| Documentation | DOCUMENTED |
| Simulator | PLANNED |
| Compliance | EDUCATIONAL ONLY |

### HLA
| Aspect | Status |
|--------|--------|
| Documentation | DOCUMENTED |
| Simulator | PLANNED |
| Compliance | EDUCATIONAL ONLY |

### DIS
| Aspect | Status |
|--------|--------|
| Documentation | DOCUMENTED |
| Simulator | PLANNED |
| Compliance | EDUCATIONAL ONLY |

## Summary

- **Fully Implemented**: ADS-B (HTTP JSON adapter)
- **Educational Only**: Link 16 (simulator + tests)
- **Encoder Only**: ASTERIX CAT021/048/062 (no decoder, no interoperability test)
- **Documentation Only**: Link 11/22, JREAP, HLA, DIS
- **No Hardware Verification**: All protocols (no RTL-SDR, no military hardware)
- **No Interoperability Verification**: All protocols (no independent decoder available)

## Verification Status

- **ADS-B**: DOCUMENTED, IMPLEMENTED, TESTED, INTEROPERABLE (with readsb reference)
- **ASTERIX CAT021/048/062**: DOCUMENTED, IMPLEMENTED (encoder only), NOT TESTED
- **Link 16**: DOCUMENTED, IMPLEMENTED (educational), TESTED
- **Link 11/22, JREAP, HLA, DIS**: DOCUMENTED only