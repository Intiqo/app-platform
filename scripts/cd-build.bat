@echo off
setlocal EnableDelayedExpansion

:: Set version and build variables
set "version=24.6.4"
set "build=2406041"

:: Get the current git commit hash
for /f "tokens=*" %%i in ('git show --format="%%h" --no-patch') do set "hash=%%i"

:: Set environment variables for cross-compilation
set "GOOS=linux"
set "GOARCH=arm64"

:: Build the Go application with the specified ldflags
set "LDFLAGS=-X \"github.com/intiqo/app-platform/internal/version.BuildVersion=!version!\" -X \"github.com/intiqo/app-platform/internal/version.BuildNumber=!build!\" -X \"github.com/intiqo/app-platform/internal/version.CommitHash=!hash!\""
call go build -ldflags=!LDFLAGS! -o ./bin/app github.com/intiqo/app-platform/cmd

endlocal
