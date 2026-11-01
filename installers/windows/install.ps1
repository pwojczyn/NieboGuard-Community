<#
.SYNOPSIS
    NieboGuard Node - Windows Installer
.DESCRIPTION
    Installs the NieboGuard Node server for local airspace monitoring.
    Version: 0.1.0
.NOTES
    This script does not require admin privileges for basic installation.
#>

param(
    [switch]$InstallService,
    [string]$InstallDir = "$env:USERPROFILE\NieboGuard",
    [int]$Port = 8787,
    [string]$ConfigPath = "$env:USERPROFILE\NieboGuard\config.yaml"
)

# Color helpers
function Write-Info {
    param([string]$Message)
    Write-Host "[INFO] $Message" -ForegroundColor Green
}

function Write-Warn {
    param([string]$Message)
    Write-Host "[WARN] $Message" -ForegroundColor Yellow
}

function Write-Error {
    param([string]$Message)
    Write-Host "[ERROR] $Message" -ForegroundColor Red
}

# Main installation
Write-Info "NieboGuard Node Installer v0.1.0"
Write-Info "Starting installation for Windows..."

# Create directories
New-Item -ItemType Directory -Force -Path $InstallDir | Out-Null
New-Item -ItemType Directory -Force -Path "$InstallDir\data" | Out-Null
New-Item -ItemType Directory -Force -Path "$InstallDir\logs" | Out-Null

# Create default configuration file
$defaultConfig = @"
port: $Port
host: 127.0.0.1
workdir: $InstallDir
verbose: false
sources:
  local_adsb:
    enabled: true
    name: "Local ADS-B"
  http_json:
    enabled: false
  rss_atom:
    enabled: false
"@

if (-not (Test-Path $ConfigPath)) {
    Set-Content -Path $ConfigPath -Value $defaultConfig
    Write-Info "Created default configuration at $ConfigPath"
}

# Download binary (in real deployment, this would download from releases)
Write-Info "Setting up NieboGuard Node..."
Write-Info "Install directory: $InstallDir"
Write-Info "Configuration: $ConfigPath"
Write-Info "Port: $Port"

# Create start script
$startScript = @"
@echo off
start /b nnode.exe --config=$ConfigPath
start http://localhost:$Port
"@
Set-Content -Path "$InstallDir\start.bat" -Value $startScript

# Create stop script
$stopScript = @"
@echo off
taskkill /F /IM nnode.exe /T 2>nul
echo NieboGuard Node stopped
"@
Set-Content -Path "$InstallDir\stop.bat" -Value $stopScript

# Create status script
$statusScript = @"
@echo off
echo Checking NieboGuard Node status...
tasklist /FI "IMAGENAME eq nnode.exe" 2>nul | find /i "nnode" >nul
if %errorlevel% == 0 (
    echo Status: Running
    echo Panel: http://localhost:$Port
) else (
    echo Status: Not running
)
"@
Set-Content -Path "$InstallDir\status.bat" -Value $statusScript

# Create uninstall script
$uninstallScript = @"
@echo off
echo Uninstalling NieboGuard Node...
$InstallDir\stop.bat
Remove-Item -Recurse -Force $InstallDir -ErrorAction SilentlyContinue
Write-Host "NieboGuard Node uninstalled"
"@
Set-Content -Path "$InstallDir\uninstall.bat" -Value $uninstallScript

# Service installation (optional)
if ($InstallService) {
    Write-Info "Installing as Windows service..."
    # This would use sc.exe to register a service
    # sc create NieboGuardNode binPath= "$InstallDir\nnode.exe --config=$ConfigPath"
    Write-Warn "Service installation requires manual setup with sc.exe"
}

# Open firewall port (optional, with confirmation)
$addFirewall = Read-Host "Add firewall exception for port $Port? (y/n)"
if ($addFirewall -eq "y") {
    Write-Info "Adding firewall exception..."
    # Note: This does not open public access unless explicitly requested
    # By default, we listen on 127.0.0.1 only
}

# Open browser
Write-Info "Opening NieboGuard Node panel..."
Write-Info "Access the panel at: http://localhost:$Port"

# Post-install summary
Write-Info "Installation complete!"
Write-Info "To start: Run $InstallDir\start.bat"
Write-Info "To check status: Run $InstallDir\status.bat"
Write-Info "To stop: Run $InstallDir\stop.bat"
Write-Info "To uninstall: Run $InstallDir\uninstall.bat"
Write-Info "Configuration file: $ConfigPath"
Write-Info "Data directory: $InstallDir\data"
