:: Mac
SET CGO_ENABLED=0
SET GOOS=darwin
SET GOARCH=amd64
go build -o dockser_darwin_amd64 main.go

:: Linux
SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go build -o dockser_linux_amd64 main.go

:: Windows
SET CGO_ENABLED=0
SET GOOS=windows
SET GOARCH=amd64
go build -o dockser_windows_amd64.exe main.go