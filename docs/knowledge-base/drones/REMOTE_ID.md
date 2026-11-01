# ASTM F3411 Remote ID

## Introduction
Remote ID is a requirement for drone identification in regulated airspace. It broadcasts drone identity and position.

## Source
- ASTM F3411-21 / SAE
- FAA: https://www.faa.gov/uas

## Format
- Broadcast over WiFi / Bluetooth / LTE
- Messages include: drone ID, location, altitude, velocity

## Status in NieboGuard
- NOT SUPPORTED yet
- Considered for P1 integration
- Simulator can be extended with synthetic drone data

## Integration Plan
- Adapter for Remote ID receiver
- Synthetic drone generation
- Link to simulator feature
