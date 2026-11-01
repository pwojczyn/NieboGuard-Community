# NieboGuard Compute — uruchomienie w pierwszym planie.
# Nie tworzy usługi ani autostartu. Wymaga jawnej zgody z rejestracji.
$ErrorActionPreference = "Stop"
$Root = Split-Path -Parent $MyInvocation.MyCommand.Path
$wrapper = Join-Path $Root "bin\nieboguard-compute.ps1"
if (-not (Test-Path $wrapper)) {
  Write-Error "Brak instalacji. Uruchom install-windows.ps1"
  exit 1
}
Write-Host "NieboGuard Compute — uruchomienie za zgodą użytkownika. Zatrzymanie: Ctrl+C."
& $wrapper start @args
exit $LASTEXITCODE
