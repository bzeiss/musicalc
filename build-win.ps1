$ErrorActionPreference = "Stop"

$env:CGO_ENABLED = "1"
$env:CC = "clang --target=x86_64-w64-windows-gnu"
$env:CXX = "clang++ --target=x86_64-w64-windows-gnu"
$env:CGO_CFLAGS = "-O3 -march=x86-64-v3"
$env:CGO_LDFLAGS = "-O3"

go build -ldflags="-s -w -H=windowsgui -extld=x86_64-w64-mingw32-gcc" -o musicalc.exe
