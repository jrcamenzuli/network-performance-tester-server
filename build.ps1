$Env:GOOS = "windows"
$Env:GOARCH = "amd64"
go build -v
$Env:GOOS = "darwin"
$Env:GOARCH = "amd64"
go build -v