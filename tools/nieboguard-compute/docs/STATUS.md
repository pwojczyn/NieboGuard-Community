# NieboGuard Compute — status publicznego modułu

Aktualizacja: 2026-11-01. To status publicznego klienta. Backend produkcyjny
pozostaje w prywatnym repozytorium; jego status nie jest publikowany tutaj.

## Działa i zostało sprawdzone

- [x] instalator Linux i prywatne `.venv` bez administratora,
- [x] wykrywanie systemu, architektury, RAM, CPU i NVIDIA GPU,
- [x] rejestracja z zaproszeniem i jawną zgodą,
- [x] weryfikacja podpisanego manifestu Ed25519,
- [x] pobieranie małych paczek danych z SHA-256 i limitem rozmiaru,
- [x] trening Audio AI na CPU z checkpointem i wznowieniem,
- [x] trening głowy Vision AI na CPU z checkpointem,
- [x] lokalny agent start/stop/resume tylko na loopbacku,
- [x] limity CPU/GPU, harmonogram bezczynności i limit danych,
- [x] lokalne przygotowanie paczek akademickich z SHA-256 i bezpiecznym rozpakowaniem,
- [x] safetensors, bez pickle i wykonywalnych checkpointów,
- [x] testy jednostkowe klienta (35 testów).

## Częściowe / wymaga potwierdzenia

- [ ] **Windows runtime** — skrypty `.ps1` mają poprawną składnię, ale nie
  uruchomiono ich na Windows 10/11.
- [ ] **GPU NVIDIA** — opcjonalny PyTorch/CUDA nie jest domyślnie instalowany;
  brak testu na fizycznej karcie.
- [ ] **Rzeczywiste paczki akademickie** — źródła mają status `DRAFT`; brak
  manifestów i SHA-256 plików.
- [ ] **Podpisane instalatory** EXE/MSI/DEB/RPM.
- [ ] **UI E2E** strony `/compute/`.

## Zablokowane

- redystrybucja YOLOBirDrone bez pisemnej zgody,
- publiczna beta przed audytem bezpieczeństwa,
- federacyjny wariant z prywatnymi zdjęciami/nagraniami,
- twierdzenie o pełnej anonimowości aktualizacji wag.

## Pliki

- `src/nieboguard_compute/` — klient,
- `install-*`, `start-*`, `uninstall-*` — Windows/Linux,
- `tests/` — testy jednostkowe,
- `docs/` — bezpieczeństwo, paczki, źródła akademickie, API.
