$Env:GOOS = "linux"
$Env:GOARCH = "arm64"
go build -o ./bin/ server.go
# Reset GOOS and GOARCH
$Env:GOOS = "windows"
$Env:GOARCH = "amd64"
