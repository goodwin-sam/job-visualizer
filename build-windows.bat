@echo off
echo Building Job Visualizer for Windows...

REM Prompt user for version number
set /p VERSION="Enter version number (e.g., v0.1.0): "

echo Building version: %VERSION%

REM Set build environment for Windows 64-bit
set GOOS=windows
set GOARCH=amd64

REM Build the application with optimizations
echo Building executable...
go build -ldflags="-s -w" -o job-visualizer-%VERSION%-windows-amd64.exe ./cmd/app

if %ERRORLEVEL% EQU 0 (
    echo.
    echo Build successful!
    echo Executable created: job-visualizer-%VERSION%-windows-amd64.exe
    echo.
    echo You can now distribute this single file to Windows users.
    echo They can run it by double-clicking the .exe file.
) else (
    echo.
    echo Build failed! Please check the error messages above.
)

pause
