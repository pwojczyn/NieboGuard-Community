# Współtworzenie materiałów Knowledge Hub

## Zgłaszanie poprawek

- Błąd merytoryczny: otwórz issue z linkiem do pliku i źródłem korekty.
- Literówka lub link: może być bezpośredni PR.
- Zmiana parametru technicznego: podaj źródło i datę weryfikacji.
- Materiał niebezpieczny lub niejawny: nie publikuj go w issue; skontaktuj
  się prywatnie według `SECURITY.md`.

## Nowy tutorial

1. Sprawdź, czy temat nie istnieje w `docs/knowledge-base/INDEX.md` i w
   Centrum Wiedzy.
2. Utwórz plik w odpowiednim katalogu.
3. Dodaj H1, krótki wstęp, sekcje H2, przykłady i ograniczenia.
4. Dodaj źródła oraz informację o licencji.
5. Uruchom testy repozytorium:

```bash
python -m unittest discover -s tests -v
node --test tests/*.test.mjs
```

6. Otwórz PR i zaznacz, czy materiał jest kandydatem do Centrum Wiedzy.

## Statusy

- `candidate` — może być opracowany,
- `review` — wymaga sprawdzenia redakcyjnego,
- `published` — decyzja administratora NieboGuard, nie zdarzenie GitHub.

## Aktualizacje i wycofanie

Zmiana pliku nie publikuje treści automatycznie. Usunięcie pliku oznacza
weryfikację, nie natychmiastowe usunięcie artykułu. Historia wersji i źródło
pozostają zachowane.

## Licencja i prawa

Nowe pliki są udostępniane na MIT, chyba że plik zawiera inną licencję.
Nie kopiuj kodu ani tekstu o niejasnej licencji. Zewnętrzne standardy mogą być
płatne lub zastrzeżone — linkuj do legalnego źródła, nie kopiuj całości.
