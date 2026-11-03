# NieboGuard Mobile — Android APK

## Wersja dla użytkowników

- `nieboguard-0.1.3-beta-release.apk`
- wersja: `0.1.3-beta`
- package: `pl.nieboguard.app`
- rozmiar: ok. 48 MB
- SHA-256: `3632cb15b126bace800717a0bb94db4fed60d1095e8993167c11058e39e45951`
- podpisana kluczem uploadu `NieboGuard Upload`

Pobierz i zainstaluj na Androidzie. To jest wersja dla normalnych użytkowników
do czasu publikacji w Google Play.

```bash
sha256sum -c SHA256SUMS
adb install -r nieboguard-0.1.3-beta-release.apk
```

## Wersja DEV — tylko dla testerów

- `nieboguard-dev-0.1.3-beta-debug.apk`
- package: `pl.nieboguard.app.dev`
- SHA-256: `59058f346ab0af71debaa4ce0e4e0c0cd7efeab23436feff3d13ac0f486e204e`
- podpisana kluczem debug

DEV może być zainstalowany obok normalnej wersji. Zwykły użytkownik powinien
wybrać plik `nieboguard-0.1.3-beta-release.apk`.

## Ważne

- To wersja beta dystrybuowana poza Google Play.
- Normalny APK jest podpisany kluczem uploadu, a nie kluczem Play App Signing.
  Jeżeli później zainstalujesz wersję z Play podpisaną innym kluczem, Android
  może wymagać odinstalowania APK z GitHub.
- Po publikacji w Play przejdź na wersję ze sklepu; nie instaluj obu wersji
  naprzemiennie.
- FCM jest w kliencie, ale backend FCM jest jeszcze wyłączony
  (`FCM_ENABLED=0`); zdalny push nie działa. Powiadomienia lokalne działają.
- Wi-Fi NAN Remote ID wymaga telefonu z Wi-Fi Aware; test HW jest wymagany.
- Wi-Fi Beacon jest best-effort i zależy od chipsetu.
- To nie jest oficjalny system ostrzegania. Nie zastępuje RCB, syren ani 112.

## Bezpieczeństwo i licencje

- APK nie zawiera `.env`, tokenów Cloudflare, kluczy SSH, kluczy podpisu ani
  service account Firebase.
- Zawiera m.in. YAMNet (Apache-2.0), Capacitor/AndroidX (Apache-2.0) i zasoby
  NieboGuard. Obowiązuje [`../THIRD_PARTY_NOTICES.md`](../THIRD_PARTY_NOTICES.md)
  oraz licencje dołączone w APK.
- Kod źródłowy aplikacji nie jest publikowany w tym repozytorium.
