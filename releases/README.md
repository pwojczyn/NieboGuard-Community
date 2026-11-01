# NieboGuard Mobile — test build

## Plik

- `nieboguard-mobile-0.1.0-dev-debug.apk`
- wersja: `0.1.0-dev`
- package: `pl.nieboguard.mobile`
- rozmiar: ok. 53 MB
- SHA-256: `e4857a79fbf9a640fbc06b1a50bedeabde194d08d021a480f6d461ffba9d33a6`

```bash
sha256sum -c releases/SHA256SUMS
adb install -r releases/nieboguard-mobile-0.1.0-dev-debug.apk
```

## To jest build testowy

- podpisany kluczem **debug**, nie release,
- nie jest przeznaczony do szerokiej dystrybucji produkcyjnej,
- może zawierać błędy i niedokończone funkcje,
- wymaga Androida zgodnego z Capacitor 8 / SDK projektu,
- nie zawiera FCM; alerty lokalne działają, gdy aplikacja może odpytać feed.

## Bezpieczeństwo i licencje

- APK nie zawiera `.env`, tokenów Cloudflare, kluczy SSH ani prywatnych kluczy,
  zgodnie z audytem paczki.
- Zawiera m.in. YAMNet (Apache-2.0), Capacitor/AndroidX (Apache-2.0) i zasoby
  NieboGuard/Strażnik. Obowiązuje [`../THIRD_PARTY_NOTICES.md`](../THIRD_PARTY_NOTICES.md)
  oraz licencje dołączone w APK.
- To nie jest oficjalny system ostrzegania. Nie zastępuje RCB, syren ani 112.

## Czego tu nie ma

- `Straznik.apk` — starszy artefakt z nieaktualną konfiguracją Firebase i
  niejasnym pochodzeniem; nie został dołączony.
- kodu źródłowego APK — to repozytorium nie jest repozytorium produkcyjnym.
- podpisu release, keystore i sekretów.
