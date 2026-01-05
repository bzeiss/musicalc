$env:CGO_FLAGS="-O3 -flto=auto -march=x86-64-v3"
$env:CGO_LDFLAGS="-O3 -flto=auto"
$env:CGO_ENABLED=1
go build -ldflags="-s -w" -o musicalc.exe
