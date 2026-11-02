# Knowledge Hub — materiały dla NieboGuard

Ten katalog opisuje, jak dokumentacja z NieboGuard Community może trafić do
publicznego Centrum Wiedzy NieboGuard.

## Co może trafić do Centrum Wiedzy

Domyślnie brane są pod uwagę pliki Markdown z:

- `docs/tutorials/`
- `docs/knowledge-base/aviation/`
- `docs/knowledge-base/transport/`
- `docs/knowledge-base/drones/`
- `docs/knowledge-base/tutorials/`
- `docs/volunteers/`
- `docs/api/`
- `docs/architecture/`

Nie każdy plik jest publikowany. System najpierw tworzy dokument źródłowy i
szkic, a artykuł wymaga review administratora.

Materiały wojskowe (`docs/knowledge-base/military/`), symulator i wewnętrzne
notatki nie są domyślnie synchronizowane. Link 16 i podobne tematy wymagają
legalnych, jawnych źródeł oraz kontroli redakcyjnej.

## Zalecany nagłówek dokumentu

Nagłówek jest opcjonalny, ale pomaga klasyfikacji:

```yaml
---
title: Krótki tytuł materiału
summary: Jedno lub dwa zdania opisu.
category: adsb
difficulty: beginner
lang: pl
license: MIT
review: required
---
```

Dozwolone kategorie: `adsb`, `standardy-i-protokoly`, `nieboguard-diy`, `ai`,
`spolecznosc`, `dokumentacja-programistow`, `poradniki`.

## Zasady

- Nie dodawaj sekretów, tokenów, konfiguracji VPS, danych użytkowników ani
  prywatnych lokalizacji.
- Nie dodawaj APK ani dużych plików binarnych.
- Podawaj źródła liczb, standardów i parametrów.
- Nie przedstawiaj symulacji ani danych syntetycznych jako rzeczywistych.
- Nie publikuj instrukcji operacyjnych ani materiałów niejawnych.
- Zachowaj jasne ostrzeżenia bezpieczeństwa.

## Jak zaproponować materiał

1. Dodaj plik Markdown w jednym z katalogów powyżej.
2. Uzupełnij źródła i licencję.
3. Otwórz PR z krótkim opisem i informacją, czy materiał może być kandydatem
   do Centrum Wiedzy.
4. Maintainer sprawdza treść i oznacza materiał do review.
5. Po synchronizacji administrator NieboGuard decyduje o publikacji.

## Aktualizacje

- Zmiana pliku jest wykrywana przez SHA.
- Nowa wersja trafia do review; publiczny artykuł nie zmienia się
  automatycznie.
- Usunięcie pliku nie usuwa artykułu; tworzy zadanie weryfikacji.
- Historia wersji pozostaje po stronie NieboGuard.

## Licencje

Kod i dokumentacja w tym repozytorium są udostępniane na licencji MIT, o ile
plik nie mówi inaczej. Zewnętrzne standardy, zrzuty, obrazy i dane mają własne
warunki — dodaj atrybucję i sprawdź licencję przed PR.

Nie dodawaj treści, do której nie masz prawa redystrybucji.
