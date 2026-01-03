; Inno Setup Script for MusiCalc
; Music Calculator Application - Universal Architecture Version

#define MyAppName "MusiCalc"
#define MyAppVersion "0.8.2"
#define MyAppPublisher "B. Zeiss"
#define MyAppURL "https://github.com/bzeiss/musicalc"
#define MyAppExeName "musicalc.exe"

[Setup]
; NOTE: The value of AppId uniquely identifies this application.
AppId={{8F9A3D2E-5B6C-4A7E-9D8F-1C2E3A4B5C6D}
AppName={#MyAppName}
AppVersion={#MyAppVersion}
AppPublisher={#MyAppPublisher}
AppPublisherURL={#MyAppURL}
AppSupportURL={#MyAppURL}
AppUpdatesURL={#MyAppURL}
DefaultDirName={autopf}\{#MyAppName}
DisableProgramGroupPage=yes
; PrivilegesRequired=lowest
OutputDir=installer
OutputBaseFilename=MusiCalc-Setup-{#MyAppVersion}
SetupIconFile=icons\appicon.ico
Compression=lzma
SolidCompression=yes
WizardStyle=modern
UninstallDisplayIcon={app}\{#MyAppExeName}

; --- Architecture Logic ---
; Allow installation on x64 (including Arm64 emulation) and native Arm64
ArchitecturesAllowed=x64compatible arm64
; Enable 64-bit install mode (native Program Files) on both architectures
ArchitecturesInstallIn64BitMode=x64compatible arm64

[Languages]
Name: "english"; MessagesFile: "compiler:Default.isl"
Name: "german"; MessagesFile: "compiler:Languages\German.isl"

[Tasks]
Name: "desktopicon"; Description: "{cm:CreateDesktopIcon}"; GroupDescription: "{cm:AdditionalIcons}"; Flags: unchecked

[Files]
; 1. Install AMD64 version on x64 systems
Source: "dist\musicalc_x64.exe"; DestDir: "{app}"; DestName: "{#MyAppExeName}"; Check: IsX64; Flags: 64bit ignoreversion

; 2. Install ARM64 version on ARM64 systems
Source: "dist\musicalc_arm64.exe"; DestDir: "{app}"; DestName: "{#MyAppExeName}"; Check: IsArm64; Flags: 64bit ignoreversion

; Common files
Source: "icons\appicon.ico"; DestDir: "{app}"; Flags: ignoreversion
Source: "icons\appicon.png"; DestDir: "{app}\icons"; Flags: ignoreversion

[Icons]
Name: "{autoprograms}\{#MyAppName}"; Filename: "{app}\{#MyAppExeName}"; IconFilename: "{app}\appicon.ico"
Name: "{autodesktop}\{#MyAppName}"; Filename: "{app}\{#MyAppExeName}"; IconFilename: "{app}\appicon.ico"; Tasks: desktopicon

[Run]
Filename: "{app}\{#MyAppExeName}"; Description: "{cm:LaunchProgram,{#StringChange(MyAppName, '&', '&&')}}"; Flags: nowait postinstall skipifsilent

[Code]
// Helper functions for the [Files] section check
function IsArm64: Boolean;
begin
  Result := (ProcessorArchitecture = paArm64);
end;

function IsX64: Boolean;
begin
  Result := (ProcessorArchitecture = paX64);
end;