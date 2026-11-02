# APK distribution notice

`nieboguard-0.1.2-beta-release.apk` to publiczna wersja beta Androida
`pl.nieboguard.app`, dystrybuowana poza Google Play do czasu publikacji
w sklepie.

- Jest podpisana kluczem uploadu `NieboGuard Upload`, nie kluczem Play App Signing.
- Nie jest jeszcze wersją z Google Play.
- FCM jest w kliencie, ale backend FCM jest wyłączony do testów.
- Wi-Fi NAN wymaga sprzętu z Wi-Fi Aware.
- APK nie zawiera sekretów serwerowych ani plików `.env`.
- Zawiera komponenty open-source; licencje i atrybucje są w
  [`../THIRD_PARTY_NOTICES.md`](../THIRD_PARTY_NOTICES.md) oraz w zasobach APK.
- Dystrybucja APK nie oznacza licencji MIT na produkcyjny kod NieboGuard.

`nieboguard-dev-0.1.2-beta-debug.apk` jest opcjonalnym buildem dla testerów
(`pl.nieboguard.app.dev`, debug).

Przed instalacją sprawdź SHA-256 z `SHA256SUMS`.
