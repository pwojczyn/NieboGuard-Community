# Military Data Link Simulation Lab — TODO

## P0 — Core Implementation

| ID | Task | Status | Priority | Dependencies |
|----|------|--------|----------|--------------|
| T1 | Complete Link 11 simulation (educational) | ⚠️ PARTIAL | P0 | T1, T2 |
| T2 | Implement Link 22 simulation (educational) | ⚠️ PARTIAL | P0 | T1, T2 |
| T3 | Complete HLA integration (if needed) | ⚠️ PLANNED | P0 | T1, T2 |
| T4 | Implement DIS simulation (educational) | ⚠️ PLANNED | P0 | T1, T2 |
| T5 | Complete integration with Sensor Hub core engine | ⚠️ PARTIAL | P0 | T1, T2, T3, T4 |

## P1 — Panel Integration

| ID | Task | Status | Priority | Dependencies |
|----|------|--------|----------|--------------|
| P1 | Implement military panel tab | ⚠️ PLANNED | P1 | T1, T2, T4, P1 |
| P2 | Add visualization for Link 16 tracks | ⚠️ PLANNED | P1 | T1, T2, T4, P1 |
| P3 | Add visualization for Link 11/22 tracks | ⚠️ PLANNED | P1 | T1, T2, T4, P1 |
| P4 | Add scenario selection (formation, patrol, etc.) | ⚠️ PLANNED | P1 | T1, T2, T4, P1 |
| P5 | Add export functionality (JSON, CSV) | ⚠️ PLANNED | P1 | T1, T2, T4, P1 |

## P2 — Testing

| ID | Task | Status | Priority | Dependencies |
|----|------|--------|----------|--------------|
| T6 | Unit tests for Link16Simulator | ⚠️ PLANNED | P1 | T1, T2, T3 |
| T7 | Integration tests with Sensor Hub | ⚠️ PLANNED | P1 | T1, T2, T3, T4 |
| T8 | Performance testing (10, 100, 500, 1000 objects) | ⚠️ PLANNED | P1 | T1, T2, T3 |
| T9 | UI tests with Playwright (desktop/mobile) | ⚠️ PLANNED | P1 | T1, T2, T4, P1 |
| T10 | Regression testing for all military protocols | ⚠️ PLANNED | P1 | All |

## P2 — Documentation

| ID | Task | Status | Priority | Dependencies |
|----|----|--------|----------|--------------|
| T11 | Complete LINK11.md with full details | ⚠️ PLANNED | P1 | T1, T2 |
| T12 | Complete LINK22.md with full details | ⚠️ PLANNED | P1 | T1, T2 |
| T13 | Complete JREAP.md with implementation examples | ⚠️ PLANNED | P1 | T1, T2 |
| T14 | Complete DIS.md with examples | ⚠️ PLANNED | P1 | T1, T2 |
| T15 | Complete HLA.md with implementation guide | ⚠️ PLANNED | P1 | T1, T2 |
| T16 | Complete SOURCES.md with all references | ⚠️ PLANNED | P1 | T1, T2 |
| T17 | Complete IMPLEMENTATION_LOG.md | ⚠️ PLANNED | P1 | T1, T2 |
| T18 | Update README.md with military section | ✅ DONE | P1 | - |

## P3 — Integration

| ID | Task | Status | Priority | Dependencies |
|----|------|--------|----------|--------------|
| I1 | Integrate Link16Simulator with Sensor Hub core | ⚠️ PLANNED | P1 | T1, T2, T3, T4 |
| I2 | Implement WebSocket bridge for military data | ⚠️ PLANNED | P1 | T1, T2, T3, I1 |
| I3 | Implement REST API endpoints for military data | ⚠️ PLANNED | P1 | T1, T2, T3, I1 |
| I4 | Implement panel integration with military tab | ⚠️ PLANNED | P1 | P1, I1, I2, I3 |

## P3 — Integration with NieboGuard

| ID | Task | Status | Priority | Dependencies |
|----|------|--------|----------|--------------|
| I5 | Integrate with core normalization engine | ⚠️ PLANNED | P1 | T1, T2, T3, T4 |
| I6 | Implement military-specific normalization rules | ⚠️ PLANNED | P1 | I5 |
| I7 | Implement military data validation | ⚠️ PLANNED | P1 | I5, I6 |
| I8 | Implement military-specific error handling | ⚠️ PLANNED | P1 | I5, I6, I7 |

## P3 — Testing

| ID | Task | Status | Priority | Dependencies |
|----|------|--------|----------|--------------|
| T11 | Unit tests for Link16Simulator | ⚠️ PLANNED | P2 | T1, T2 |
| T12 | Integration tests with Sensor Hub | ⚠️ PLANNED | P2 | T1, T2, T3, I1 |
| T13 | Performance testing (10, 100, 500, 1000 objects) | ⚠️ PLANNED | P2 | T1, T2, T3 |
| T14 | UI tests with Playwright (desktop/mobile) | ⚠️ PLANNED | P2 | T1, T2, T4, P1 |
| T15 | Regression testing for all military protocols | ⚠️ PLANNED | P2 | All |

## P3 — Documentation

| ID | Task | Status | Priority | Dependencies |
|----|------|--------|----------|--------------|
| T16 | Complete IMPLEMENTATION_LOG.md | ⚠️ PLANNED | P1 | T1, T2, T3 |
| T17 | Update all documentation with military section | ⚠️ PLANNED | P1 | T1, T2 |
| T18 | Update ROADMAP.md with military milestones | ⚠️ PLANNED | P1 | T1, T2 |
| T19 | Update CHANGELOG.md with military updates | ⚠️ PLANNED | P1 | T1, T2 |
| T20 | Update README.md with military section | ✅ DONE | P1 | - |

## P3 — Integration with NieboGuard

| ID | Task | Status | Priority | Dependencies |
|----|------|--------|----------|--------------|
| I9 | Implement military-specific normalization rules | ⚠️ PLANNED | P1 | I5, I6 |
| I10 | Implement military-specific validation rules | ⚠️ PLANNED | P1 | I5, I6 |
| I11 | Implement military-specific error handling | ⚠️ PLANNED | P1 | I5, I6, I7 |
| I11 | Implement military-specific security model (educational only) | ⚠️ PLANNED | P1 | I5, I6, I7 |