@echo off
if "%PROCESSOR_ARCHITECTURE%"=="AMD64" (
    set GOARCH=amd64
    set ARCH_NAME=amd64
) else if "%PROCESSOR_ARCHITECTURE%"=="x86" (
    set GOARCH=386
    set ARCH_NAME=386
) else if "%PROCESSOR_ARCHITECTURE%"=="ARM64" (
    set GOARCH=arm64
    set ARCH_NAME=arm64
) else (
    set GOARCH=amd64
    set ARCH_NAME=amd64
)

set /p VERSION="Enter version number: "
set GOOS=windows
go build -ldflags="-s -w" -o job-visualizer-%VERSION%-windows-%ARCH_NAME%.exe ./cmd/app
