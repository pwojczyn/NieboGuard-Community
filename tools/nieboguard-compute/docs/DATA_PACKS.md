# Małe paczki danych akademickich

Nie wolno zmuszać wolontariusza do pobierania wielogigabajtowego archiwum.
Paczka jest podzbiorem zatwierdzonych shardów datasetu, z własnym manifestem i
SHA-256. Nie zawiera kodu.

## Klasy

| Klasa | Zakres | Domyślny target |
|---|---|---|
| `MINI` | 25–50 MB | 40 MB |
| `STANDARD` | 100–200 MB | 150 MB |
| `EXTENDED` | 250–500 MB | 375 MB |
| `GPU` | 500–1024 MB | 750 MB |

Rozmiar modelu i bibliotek jest **osobny** od rozmiaru paczki. Pola
`model_bytes` i `library_bytes` są w manifescie paczki.

Wartości są konfigurowalne w panelu (`compute_settings.packs`) i walidowane
przez Worker. Test E2E używa małych wartości wyłącznie w lokalnym D1.

## Budowanie

Admin uploaduje shardy datasetu, tworzy dataset (`POST /admin/compute/datasets`),
a następnie:

```http
POST /api/v2/admin/compute/datasets/:id/:version/build-packs
{"tiers":["MINI","STANDARD"],"model_bytes":5242880,"library_bytes":31457280}
```

Builder:

1. sprawdza, że każdy shard istnieje w tej samej wersji i ma SHA-256,
2. odrzuca duplikaty SHA-256,
3. przypisuje grupę sesji i split deterministycznie (70/15/15), jeśli brak,
4. nie pozwala, aby jedna grupa była w dwóch splitach,
5. pakuje shardy w bin-packing do targetu, bez rozcinania sesji,
6. łączy zbyt małe biny, jeśli nadal mieszczą się w limicie,
7. zapisuje manifest `compute/packs/<dataset>/v<version>/<tier>/<index>/manifest.json`,
8. zapisuje rekord `compute_data_packs` z SHA-256 i licznikami klas.

Paczka poniżej `min_mb` lub powyżej `max_mb` nie jest oznaczana `READY`.
`skipped` w odpowiedzi wyjaśnia, dlaczego nie powstała.

## Przydział

- Zadanie ma `resource.pack_tier`.
- Urządzenie deklaruje `max_data_mb` i `pack_tiers` (zgoda użytkownika).
- `resourceMatches` odrzuca zadanie, jeśli paczka przekracza limit lub klasa
  nie jest dozwolona.
- `next_task` wybiera paczkę najrzadziej używaną, preferując taką, której dane
  urządzenie jeszcze nie dostało.
- Manifest zadania zawiera `pack` z ID, klasą, bajtami, SHA-256, liczbą sesji
  i osobnymi rozmiarami modelu/bibliotek.
- Klient sprawdza SHA-256 każdego shardu, limit `max_data_mb` i klasę paczki.
- Ukończenie zapisuje `pack_id` w przypisaniu; agregacja nie łączy wyników z
  różnych paczek bez zgodności task/round/base/architektury.

## Deduplikacja i split

- klucz grupy: `group_key` / `session_id` / nazwa sharda,
- ten sam shard SHA-256 nie może wystąpić dwa razy w datasecie,
- grupa nie może wystąpić w `train` i `validation/test`,
- builder nie rozcina sharda; shard powinien odpowiadać jednej sesji lub
  spójnemu fragmentowi,
- paczki tej samej klasy nie dzielą grup (każdy shard trafia do jednej paczki).

## Czego nie robimy

- nie publikujemy całych archiwów Zenodo/GitHub,
- nie dzielimy jednego nagrania na wiele paczek bez metadanej granicy,
- nie używamy danych bez prawa redystrybucji,
- nie traktujemy rozmiaru paczki jako gwarancji jakości lub zbalansowania klas,
- nie wysyłamy mowy ludzkiej ani danych osobowych.
