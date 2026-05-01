$ErrorActionPreference = "Stop"

$scriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$projectRoot = Resolve-Path (Join-Path $scriptDir "..\..")
$distDir = Join-Path $projectRoot "build\dist"

New-Item -ItemType Directory -Force $distDir | Out-Null

$env:CGO_ENABLED = "1"
$env:CC = "clang --target=x86_64-w64-windows-gnu"
$env:CXX = "clang++ --target=x86_64-w64-windows-gnu"
$env:CGO_CFLAGS = "-O3 -march=x86-64-v3"
$env:CGO_LDFLAGS = "-O3"

Push-Location $projectRoot
try {
    go build -ldflags="-s -w -H=windowsgui -extld=x86_64-w64-mingw32-gcc" -o (Join-Path $distDir "musicalc.exe")
    $exitCode = $LASTEXITCODE
} finally {
    Pop-Location
}

if ($exitCode -ne 0) {
    exit $exitCode
}
