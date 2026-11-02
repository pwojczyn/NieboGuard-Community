# NieboGuard Mobile — test build

## Plik

- `nieboguard-dev-0.1.2-beta-debug.apk`
- wersja: `0.1.2-beta-dev`
- package: `pl.nieboguard.app.dev`
- rozmiar: ok. 54 MB
- SHA-256: `3509d64bed17945cea99dd3605951d71a8bff27e0b2c9c42bf26e79be1f5fe18`

```bash
sha256sum -c releases/SHA256SUMS
adb install -r releases/nieboguard-dev-0.1.2-beta-debug.apk
```

## To jest build testowy

- podpisany kluczem **debug**, nie release,
- nie jest przeznaczony do szerokiej dystrybucji produkcyjnej,
- może zawierać błędy i niedokończone funkcje,
- wymaga Androida zgodnego z Capacitor 8 / SDK projektu,
- zawiera klienta FCM; backend FCM jest produkcyjnie wyłączony, więc nie
  oczekuj zdalnego pushu — powiadomienia lokalne działają,
- zawiera odbiornik Wi-Fi NAN Remote ID; wymaga telefonu z Wi-Fi Aware,
  test na fizycznym urządzeniu jest nadal wymagany,
- Wi-Fi Beacon jest best-effort i zależy od chipsetu,
- to nie jest oficjalny system ostrzegania.

## Bezpieczeństwo i licencje

- APK nie zawiera `.env`, tokenów Cloudflare, kluczy SSH, kluczy podpisu ani
  service account Firebase.
- Zawiera m.in. YAMNet (Apache-2.0), Capacitor/AndroidX (Apache-2.0) i zasoby
  NieboGuard. Obowiązuje [`../THIRD_PARTY_NOTICES.md`](../THIRD_PARTY_NOTICES.md)
  oraz licencje dołączone w APK.
- To nie jest oficjalny system ostrzegania. Nie zastępuje RCB, syren ani 112.

## Czego tu nie ma

- kodu źródłowego APK — to repozytorium nie jest repozytorium produkcyjnym,
- podpisu release, keystore i sekretów,
- potwierdzonego testu FCM/Wi-Fi NAN na fizycznym telefonie.
