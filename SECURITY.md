# Security policy

## Zgłaszanie luk

Nie otwieraj publicznego issue dla podatności. Wyślij prywatnie:

**kontakt@nieboguard.pl** z tematem `SECURITY — NieboGuard Community`.

Podaj:

- opis i wpływ,
- minimalną reprodukcję bez sekretów i danych osobowych,
- wersję/commit i środowisko,
- propozycję poprawki, jeśli ją masz.

Nie testuj publicznych PR ani infrastruktury produkcyjnej bez zgody. Nie
publikuj exploitów, tokenów, danych użytkowników ani współrzędnych prywatnych
instalacji.

## Zakres

- SDK i narzędzia w tym repozytorium,
- szablony i workflow Community,
- dokumentacja, jeśli błąd może wprowadzić użytkownika w błąd w sprawach
  bezpieczeństwa.

## Poza zakresem

- produkcyjne API i infrastruktura NieboGuard (zgłaszaj przez kontakt),
- dane dostawców zewnętrznych,
- urządzenia i sieci, których nie jesteś właścicielem.

## Zasady

- brak sekretów w issue/PR,
- brak danych osobowych,
- brak automatycznego publikowania kodu produkcyjnego,
- publiczne workflow nie mają dostępu do sekretów prywatnego repo.

Dziękujemy za odpowiedzialne zgłoszenia.

## NieboGuard Compute

Klient Compute przyjmuje wyłącznie podpisane manifesty Ed25519, safetensors bez
pickle i pliki z limitem rozmiaru oraz SHA-256. Nie uruchamia kodu z backendu.
Wyniki są kwarantanną po stronie prywatnego backendu, nie automatyczną
produkcją.

Zgłaszaj prywatnie:

- obejście weryfikacji podpisu manifestu,
- path traversal lub pobranie spoza dozwolonego origin,
- wykonanie kodu z pliku zadania,
- ujawnienie tokenu urządzenia,
- obejście limitu CPU/GPU/danych,
- błąd w local agent (nasłuch poza loopbackiem, obcy Origin).

Nie publikuj tokenów, kluczy, prywatnych danych ani exploitów. Lokalny agent
Compute jest opcjonalny, działa na `127.0.0.1:8765` i nie może być wystawiany
na LAN.
