; Inno Setup Script for MusiCalc
; Music Calculator Application

#define MyAppName "MusiCalc"
#define MyAppVersion "0.8.0"
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
; Uncomment the following line to run in non administrative install mode (install for current user only.)
;PrivilegesRequired=lowest
OutputDir=installer
OutputBaseFilename=MusiCalc-Setup-{#MyAppVersion}
SetupIconFile=icons\appicon.ico
Compression=lzma
SolidCompression=yes
WizardStyle=modern
UninstallDisplayIcon={app}\{#MyAppExeName}

[Languages]
Name: "english"; MessagesFile: "compiler:Default.isl"
Name: "german"; MessagesFile: "compiler:Languages\German.isl"

[Tasks]
Name: "desktopicon"; Description: "{cm:CreateDesktopIcon}"; GroupDescription: "{cm:AdditionalIcons}"; Flags: unchecked

[Files]
Source: "musicalc.exe"; DestDir: "{app}"; Flags: ignoreversion
; NOTE: Don't use "Flags: ignoreversion" on any shared system files

[Icons]
Name: "{autoprograms}\{#MyAppName}"; Filename: "{app}\{#MyAppExeName}"
Name: "{autodesktop}\{#MyAppName}"; Filename: "{app}\{#MyAppExeName}"; Tasks: desktopicon

[Run]
Name: "{autorun}\{#MyAppName}"; Filename: "{app}\{#MyAppExeName}"; Description: "{cm:AutoStartProgram,{#MyAppName}}"; Flags: runascurrentuser nowait postinstall skipifsilent
