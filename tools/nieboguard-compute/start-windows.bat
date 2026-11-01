@echo off
setlocal
set "ROOT=%~dp0"
if not exist "%ROOT%bin\nieboguard-compute.cmd" (
  echo Brak instalacji. Uruchom najpierw install-windows.ps1
  exit /b 1
)
echo NieboGuard Compute - uruchomienie za zgoda uzytkownika. Zatrzymanie: Ctrl+C.
call "%ROOT%bin\nieboguard-compute.cmd" start %*
exit /b %ERRORLEVEL%
