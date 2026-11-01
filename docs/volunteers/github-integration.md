# Projekt integracji GitHub z kontem NieboGuard

## Cel

Wolontariusz może opcjonalnie połączyć istniejące konto NieboGuard z GitHubem,
aby przypisać zadania, PR i odznaki. Nie tworzymy nowego systemu kont.

## Przepływ

1. Użytkownik loguje się do NieboGuard istniejącym mechanizmem.
2. W profilu wybiera „Połącz GitHub”.
3. Backend generuje `state` i przekierowuje do GitHub OAuth.
4. GitHub zwraca `code`; backend wymienia go na token.
5. Zapisujemy minimalny identyfikator GitHub, login i link do profilu.
6. Użytkownik może odłączyć konto; token jest unieważniany/ usuwany.

## Zakres OAuth

- minimalny scope: publiczny profil (`read:user`),
- bez dostępu do repozytoriów, e-maili, organizacji i prywatnych danych,
- token wyłącznie po stronie serwera, nigdy w publicznym repo/APK,
- brak automatycznego publikowania w imieniu użytkownika.

## Model danych (projekt)

```text
community_github_links(
  user_id, github_id, login, avatar_url, linked_at, revoked_at
)
```

- `github_id` unikalny,
- brak e-maila i tokenu w tabeli publicznej,
- token (jeśli potrzebny) w osobnym, szyfrowanym magazynie,
- możliwość usunięcia na żądanie.

## Zadania i reputacja

- issues z etykietami `good-first-issue`, `help-wanted`,
- PR przypisany do konta przez login GitHub,
- po merge moderator może przyznać XP/odznakę w istniejącym systemie,
- brak automatycznego XP za sam commit/PR bez review,
- status zadania: `open`, `in_progress`, `review`, `done`, `blocked`.

## Nadużycia

- limit łączenia kont,
- weryfikacja własności GitHub przez OAuth,
- moderacja i możliwość cofnięcia odznak,
- brak publicznego ujawniania e-maila,
- log audytu operacji administracyjnych.

## Status

To projekt. Nie wdrożono OAuth ani nowych tabel w produkcji w ramach tego
zadania.
