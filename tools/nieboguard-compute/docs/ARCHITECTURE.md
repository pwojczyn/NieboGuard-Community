# Architektura NieboGuard Compute

## Przepływ

```text
Administrator AI Lab
  │  zaproszenie, dataset, model bazowy, zadanie
  ▼
Worker / D1 / R2 ── podpisany manifest (Ed25519) ──► klient Windows/Linux
  │                                                     │
  │  capability URL (download)                          ├─ weryfikacja podpisu
  │◄────────────────────────────────────────────────────┤
  │                                                     ├─ lokalny trening / checkpoint
  │  capability URL (upload do kwarantanny)             │
  │◄────────────────────────────────────────────────────┤
  │                                                     ▼
  │  safetensors + metryki                        kwarantanna R2
  │                                                     │
  ├─ walidacja integralności, nagłówka, zakresów, dryfu ─┘
  │
  ├─ niezależny evaluator (zaufany runner, HMAC)
  │
  ├─ agregacja tylko ACCEPTED i zgodnych architektur
  │
  ├─ kandydat modelu (CANDIDATE)
  │
  └─ ręczna publikacja / rollback (admin)
```

## Granice odpowiedzialności

### Worker

- rejestr urządzeń i zaproszeń,
- krótkotrwałe tokeny i cofanie dostępu,
- podpisywanie manifestów,
- leasing zadań,
- capability URL do R2,
- walidacja safetensors i limitów,
- rejestr rund, aktualizacji, ewaluacji i modeli.

Worker **nie** trenuje, nie uruchamia Pythona i nie publikuje modelu
automatycznie.

### D1

Tabele `compute_*`: urządzenia, tokeny, zaproszenia, datasety, zadania, rundy,
przypisania, uploady, aktualizacje, ewaluacje, modele, audyt, limity.
Migracja jest addytywna i nie zmienia tabel produkcyjnych.

### R2 `MEDIA`

Prefiksy:

- `compute/datasets/` — manifesty danych,
- `compute/base-models/` — artefakty bazowe,
- `compute/quarantine/updates/` — surowe wyniki wolontariuszy,
- `compute/aggregated/` — modele po agregacji,
- `compute/reports/` — raporty evaluatora.

Kwarantanna jest logicznie oddzielona. Worker nie serwuje jej publicznie.

### Klient

- weryfikuje podpis manifestu kluczem publicznym,
- pobiera wyłącznie publiczne shardy z prawem redystrybucji,
- trenuje lokalnie i zapisuje checkpointy safetensors,
- respektuje limity CPU/GPU i czas zadania,
- wysyła tylko wagi i metryki.

### Evaluator

Zaufany runner poza Workerem uruchamia testy na niezależnym holdoucie i
podpisuje raport HMAC. Worker nie ufa metrykom klienta; używa ich tylko jako
kontekstu.

## Protokół

1. `POST /api/v2/compute/enroll` — zaproszenie + zgoda → refresh token.
2. `POST /api/v2/compute/token` — refresh → access token 20 min.
3. `GET /api/v2/compute/tasks/next` — dopasowanie zasobów, leasing.
4. Manifest zawiera podpis, dataset, model bazowy, architekturę,
   hiperparametry, limit czasu, upload i progi walidacji.
5. Klient pobiera shardy, trenuje, wysyła safetensors.
6. `POST /compute/tasks/:id/complete` — serwer sprawdza R2, SHA-256, nagłówek,
   tensor spec, zakresy i dryf; tworzy `PENDING_EVAL`.
7. Evaluator wysyła raport; serwer stosuje progi i oznacza
   `ACCEPTED`/`REJECTED`/`MANUAL_REVIEW`.
8. Administrator agreguje tylko `ACCEPTED`, zgodne co do task/round/base i
   tensorów; wynik jest `CANDIDATE`.
9. Evaluator ocenia kandydata na holdoucie; administrator publikuje lub
   cofa.

## Dopasowanie zasobów

Zadanie ma `resource_json`: system, architektura, min. rdzenie, min. RAM,
wymaganie GPU, min. VRAM, max urządzeń i modalność. Urządzenie deklaruje
kubełki, nie pełną specyfikację. Nie zbieramy nazwy hosta, adresów ani
numerów seryjnych.

## Agregacja

- tylko ta sama runda, task, wersja bazowa i hash architektury,
- te same nazwy, dtype i kształty tensorów,
- unikalne SHA-256,
- wagi według `sample_count` z nagłówka safetensors,
- ważony **trimmed mean** (dla ≥5 aktualizacji odrzucane są skrajne wartości),
- brak publikacji automatycznej.

## Dlaczego nie Flower w MVP

Flower wymaga serwera koordynującego, własnego transportu i dodatkowej
powierzchni ataku. Obecny MVP ma podpisane manifesty, capability URL i D1/R2,
które już istnieją w projekcie. Prostszy protokół jest łatwiejszy do audytu.
Flower pozostaje opcją dla przyszłego wariantu federacyjnego.
