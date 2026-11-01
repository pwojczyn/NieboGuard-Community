<#
.SYNOPSIS
  Usuwa klienta NieboGuard Compute z Windows.
.DESCRIPTION
  Domyślnie zachowuje konfigurację i checkpointy. -Purge usuwa również dane
  i tokeny. Nie wymaga uprawnień administratora.
#>
[CmdletBinding()]
param([switch]$Purge)
$ErrorActionPreference = "Stop"
$Root = Split-Path -Parent $MyInvocation.MyCommand.Path
$wrapper = Join-Path $Root "bin\nieboguard-compute.cmd"
if (Test-Path $wrapper) {
  & $wrapper --data-dir (Join-Path $Root "data") stop 2>$null | Out-Null
}
Remove-Item -Force -ErrorAction SilentlyContinue (Join-Path $Root "bin\nieboguard-compute.cmd")
Remove-Item -Force -ErrorAction SilentlyContinue (Join-Path $Root "bin\nieboguard-compute.ps1")
Remove-Item -Recurse -Force -ErrorAction SilentlyContinue (Join-Path $Root ".venv")
if ($Purge) {
  Remove-Item -Recurse -Force -ErrorAction SilentlyContinue (Join-Path $Root "data")
  Write-Host "Usunięto program, konfigurację, tokeny i checkpointy."
} else {
  Write-Host "Usunięto program i środowisko Python. Konfiguracja i checkpointy pozostały w $Root\data."
  Write-Host "Aby usunąć wszystko: .\uninstall-windows.ps1 -Purge"
}
