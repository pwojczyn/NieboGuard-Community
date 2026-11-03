; NieboGuard Audio AI Trainer Installer Script
; Version 1.0
; Copyright (c) 2026 NieboGuard Ecosystem

; Configuration
Name "NieboGuard Audio AI Trainer"
OutFile "build/NieboGuard-Audio-Trainer-Setup.exe"
InstallDir "$PROGRAMFILES\NieboGuard\AudioAI"
SetCompressor /SOLID lzma
SetCompressorOverrideList /SOLID lzma

; Include modern UI
!include "MUI2.nsh"
!include "FileDialogs.nsh"
!include "LogicLib.nsh"
!include "WinVer.nsh"

!define APP_NAME "NieboGuard Audio AI Trainer"
!define COMPANY_NAME "NieboGuard Ecosystem"
!define PRODUCT_VERSION "1.0.0"
!define PRODUCT_ID "NieboGuardAudioAI"
!define BUILD_VERSION "1.0.0"

; Language
!insertmacro MUI_LANGUAGE "English"

; Colors
!insertmacro MUI_HEADERIMAGE
!insertmacro MUI_ICON "build/NieboGuard-Audio-Trainer.exe"
!insertmacro MUI_ICONGROUP "NieboGuardAudioAI"

; Language settings
!insertmacro MUI_LANGUAGE "English"

; Global settings
!insertmacro MUI_PAGE_WELCOME
!insertmacro MUI_PAGE_DIRECTORY
!insertmacro MUI_PAGE_INSTFILES
!insertmacro MUI_PAGE_FINISH

!insertmacro MUI_LANGUAGE "English"

; Page configuration
!insertmacro MUI_PAGE_WELCOME
!insertmacro MUI_PAGE_DIRECTORY
!insertmacro MUI_PAGE_INSTFILES
!insertmacro MUI_PAGE_FINISH

; Include NSIS core files
!insertmacro MUI_ICONGROUP "NieboGuardAudioAI"

; Files to install
Section "NieboGuard Audio AI Trainer" SecMain
    SetOutPath "$INSTDIR"
    SetIncludesDir "$INSTDIR"

    ; Main executable
    File "build/NieboGuard-Audio-Trainer.exe" "NieboGuard-Audio-Trainer.exe"
    File "build/NieboGuard-Audio-Trainer.exe" "NieboGuard-Audio-Trainer.exe" ; Copy to executable name

    ; Configuration files
    File "apps/audio-trainer/config.yaml" "config.yaml"
    File "apps/audio-trainer/config.yaml" "config.yaml"

    ; Desktop shortcut
    CreateShortCut "$DESKTOP\NieboGuard Audio AI Trainer.lnk" "$INSTDIR\NieboGuard-Audio-Trainer.exe" "" "$INSTDIR\NieboGuard-Audio-Trainer.exe" 0

    ; Start menu entry
    CreateDirectory "$SMPROGRAMS\NieboGuard"
    CreateShortCut "$SMPROGRAMS\NieboGuard\NieboGuard Audio AI Trainer.lnk" "$INSTDIR\NieboGuard-Audio-Trainer.exe" "" "$INSTDIR\NieboGuard-Audio-Trainer.exe" 0

    ; Documentation
    CreateDirectory "$INSTDIR\docs"
    File "README.md" "docs\README.md"
    File "README.md" "docs\README.md"

    ; License
    CreateDirectory "$INSTDIR\license"
    File "LICENSE" "docs\license\LICENSE"

    ; Create uninstaller
    WriteUninstaller "$INSTDIR\uninstall.exe"
SectionEnd

; Uninstallation section
Section "Uninstall"
    MessageBox MB_OK "Are you sure you want to uninstall NieboGuard Audio AI Trainer?"
    Delete "$INSTDIR\NieboGuard-Audio-Trainer.exe"
    Delete "$INSTDIR\config.yaml"
    Delete "$INSTDIR\docs\README.md"
    Delete "$INSTDIR\license\LICENSE"
    Delete "$INSTDIR\docs\README.md"
    RMDir "$INSTDIR"
    RMDir "$SMPROGRAMS\NieboGuard"
    RMDir "$DESKTOP\NieboGuard Audio AI Trainer.lnk"
    RMDir "$SMPROGRAMS\NieboGuard"
EndSection

; Create installer
!insertmacro MUI_UNPAGE_CONFIRM
!insertmacro MUI_INSTALLCONFIRM
!insertmacro MUI_INSTFILES

; Finalize installer
!insertmacro MUI_PAGE_FINISH
!insertmacro MUI_UNPAGE_CONFIRM