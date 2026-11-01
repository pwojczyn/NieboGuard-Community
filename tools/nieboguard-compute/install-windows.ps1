<#
.SYNOPSIS
  NieboGuard Compute — instalator dla Windows 10/11 x64.
.DESCRIPTION
  Tworzy prywatne środowisko Python w katalogu klienta. Nie wymaga
  administratora, nie pobiera skryptów z internetu i nie uruchamia kodu z
  backendu. Nie tworzy autostartu ani ukrytej usługi.
#>
[CmdletBinding()]
param(
  [string]$Python = "",
  [switch]$WithOptional,
  [switch]$SkipBackendCheck,
  [switch]$NonInteractive,
  [string]$Api = "https://api.neboguard.pl",
  [string]$ManifestKey = ""
)

$ErrorActionPreference = "Stop"
$Root = Split-Path -Parent $MyInvocation.MyCommand.Path
function Say($m) { Write-Host $m }
function Fail($m) { Write-Error $m; exit 1 }

if ($env:OS -ne "Windows_NT") { Fail "Ten instalator obsługuje tylko Windows." }
$arch = $env:PROCESSOR_ARCHITECTURE
if ($arch -notin @("AMD64", "x86_64")) { Say "UWAGA: architektura $arch nie jest jeszcze obsługiwana przez backend." }

$os = [System.Environment]::OSVersion.Version
Say ("Windows {0}.{1} build {2} · {3}" -f $os.Major, $os.Minor, $os.Build, $arch)
$ramGb = [math]::Round((Get-CimInstance Win32_ComputerSystem).TotalPhysicalMemory / 1GB, 1)
$cpu = Get-CimInstance Win32_Processor | Select-Object -First 1
Say ("RAM: {0} GB · CPU: {1} rdzeni ({2})" -f $ramGb, $cpu.NumberOfLogicalProcessors, $cpu.Name)
if ($ramGb -lt 2) { Fail "Wymagane co najmniej 2 GB RAM (zalecane 4 GB)." }
if ($ramGb -lt 4) { Say "UWAGA: mniej niż 4 GB RAM. Zadania Vision mogą zostać odrzucone." }
$gpu = Get-CimInstance Win32_VideoController | Select-Object -First 1
if ($gpu) { Say ("GPU: {0}" -f $gpu.Name) } else { Say "GPU: nie wykryto. Zadania CPU nadal działają." }

$candidates = @()
if ($Python) { $candidates += $Python }
$candidates += @("py", "python", "python3")
$pyCmd = $null
foreach ($c in $candidates) {
  if (-not $c) { continue }
  $cmd = Get-Command $c -ErrorAction SilentlyContinue
  if ($cmd) {
    try {
      $ver = & $c -c "import sys; print('%d.%d' % sys.version_info[:2])" 2>$null
      if ([version]$ver -ge [version]"3.10") { $pyCmd = $c; break }
    } catch { }
  }
}
if (-not $pyCmd) { Fail "Brak Pythona 3.10+. Zainstaluj Python z python.org (bez uprawnień administratora) i uruchom ponownie." }
Say ("Python: {0}" -f (& $pyCmd -c "import sys; print(sys.version.split()[0])"))

$Venv = Join-Path $Root ".venv"
if (-not (Test-Path (Join-Path $Venv "Scripts\python.exe"))) {
  Say "Tworzę izolowane środowisko Python w $Venv"
  & $pyCmd -m venv $Venv
  if ($LASTEXITCODE -ne 0) { Fail "Nie udało się utworzyć venv." }
}
$vpy = Join-Path $Venv "Scripts\python.exe"
& $vpy -m pip --version | Out-Null
if ($LASTEXITCODE -ne 0) { Fail "venv nie zawiera pip." }
& $vpy -m pip install --disable-pip-version-check --quiet --upgrade "pip>=24,<26"
if ($WithOptional) {
  $lock = Join-Path $Root "requirements-optional.lock"
  if (-not (Test-Path $lock)) { Fail "Brak requirements-optional.lock." }
  Say "Instaluję opcjonalne zależności z weryfikacją SHA-256…"
  & $vpy -m pip install --disable-pip-version-check --require-hashes -r $lock
  if ($LASTEXITCODE -ne 0) { Fail "Instalacja opcjonalnych zależności nie powiodła się." }
}

$Bin = Join-Path $Root "bin"
New-Item -ItemType Directory -Force -Path $Bin | Out-Null
$cmdPath = Join-Path $Bin "nieboguard-compute.cmd"
@"
@echo off
setlocal
set "ROOT=%~dp0.."
set "PYTHONPATH=%ROOT%\src;%PYTHONPATH%"
"%ROOT%\.venv\Scripts\python.exe" -m nieboguard_compute %*
exit /b %ERRORLEVEL%
"@ | Set-Content -Encoding ASCII -Path $cmdPath
$psPath = Join-Path $Bin "nieboguard-compute.ps1"
@"
`$ErrorActionPreference = "Stop"
`$Root = Split-Path -Parent (Split-Path -Parent `$MyInvocation.MyCommand.Path)
`$env:PYTHONPATH = Join-Path `$Root "src"
& (Join-Path `$Root ".venv\Scripts\python.exe") -m nieboguard_compute @args
exit `$LASTEXITCODE
"@ | Set-Content -Encoding UTF8 -Path $psPath

$Data = Join-Path $Root "data"
New-Item -ItemType Directory -Force -Path $Data | Out-Null
$cfg = Join-Path $Data "config.json"
$wrapper = Join-Path $Bin "nieboguard-compute.cmd"
if (-not (Test-Path $cfg)) {
  & $wrapper --data-dir $Data config set api_base $Api | Out-Null
  if ($ManifestKey) { & $wrapper --data-dir $Data config set manifest_public_key $ManifestKey | Out-Null }
}

if (-not $NonInteractive) {
  Say ""
  Say "Wybór sposobu użycia komputera (możesz zmienić później)."
  $modality = Read-Host "Model [audio/vision/both] (domyślnie audio)"
  if ($modality -notin @("audio", "vision", "both")) { $modality = "audio" }
  if ($modality -eq "both") { $mods = "audio,vision" } else { $mods = $modality }
  $cpuPct = Read-Host "Limit CPU % (5-100, domyślnie 50)"
  if (-not $cpuPct) { $cpuPct = "50" }
  $gpuPct = Read-Host "Limit GPU % (0-100, domyślnie 60)"
  if (-not $gpuPct) { $gpuPct = "60" }
  $sched = Read-Host "Harmonogram [manual/idle/always] (domyślnie manual)"
  if ($sched -notin @("manual", "idle", "always")) { $sched = "manual" }
  $dataMb = Read-Host "Maks. paczka danych MB [50/200/500/1024] (domyślnie 200)"
  if ($dataMb -notin @("50", "200", "500", "1024")) { $dataMb = "200" }
  switch ($dataMb) { "50" { $tiers = "MINI" } "500" { $tiers = "MINI,STANDARD,EXTENDED" } "1024" { $tiers = "GPU" } default { $tiers = "MINI,STANDARD" } }
  & $wrapper --data-dir $Data config set modalities $mods | Out-Null
  & $wrapper --data-dir $Data config set max_data_mb $dataMb | Out-Null
  & $wrapper --data-dir $Data config set pack_tiers $tiers | Out-Null
  & $wrapper --data-dir $Data config set max_cpu_percent $cpuPct | Out-Null
  & $wrapper --data-dir $Data config set max_gpu_percent $gpuPct | Out-Null
  & $wrapper --data-dir $Data config set schedule $sched | Out-Null
}

if (-not $SkipBackendCheck) {
  & $wrapper --data-dir $Data doctor --json | Out-Null
  if ($LASTEXITCODE -ne 0) {
    Say "UWAGA: backend jest niedostępny lub nieobsługiwany. Instalacja lokalna jest gotowa."
  } else {
    Say "Backend odpowiada."
  }
}

Say ""
Say "Instalacja zakończona: $Root"
Say "Uruchamianie:   $Root\start-windows.bat"
Say "Diagnostyka:   $Root\bin\nieboguard-compute.cmd doctor"
Say "Rejestracja:   $Root\bin\nieboguard-compute.cmd enroll --invite <kod> --consent"
Say ""
Say "Trening nie uruchamia się w tle. Program startuje wyłącznie po Twojej komendzie."
