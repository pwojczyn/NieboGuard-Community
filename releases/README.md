# NieboGuard Mobile — Android APK

## Wersja dla użytkowników

- `nieboguard-0.1.2-beta-release.apk`
- wersja: `0.1.2-beta`
- package: `pl.nieboguard.app`
- rozmiar: ok. 48 MB
- SHA-256: `925419dfa5facedf1be840e93599c5b0a9e6459aed883f496aec502e0034f498`
- podpisana kluczem uploadu `NieboGuard Upload`

Pobierz i zainstaluj na Androidzie. To jest wersja dla normalnych użytkowników
do czasu publikacji w Google Play.

```bash
sha256sum -c SHA256SUMS
adb install -r nieboguard-0.1.2-beta-release.apk
```

## Wersja DEV — tylko dla testerów

- `nieboguard-dev-0.1.2-beta-debug.apk`
- package: `pl.nieboguard.app.dev`
- SHA-256: `3509d64bed17945cea99dd3605951d71a8bff27e0b2c9c42bf26e79be1f5fe18`
- podpisana kluczem debug

DEV może być zainstalowany obok normalnej wersji. Zwykły użytkownik powinien
wybrać plik `nieboguard-0.1.2-beta-release.apk`.

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
