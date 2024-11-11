#!/usr/bin/env bash
version="24.6.4"
build="2406041"
hash="$(git show --format="%h" --no-patch)"
go build -ldflags="-X 'github.com/intiqo/app-platform/internal/version.BuildVersion=$version' -X 'github.com/intiqo/app-platform/internal/version.BuildNumber=$build' -X 'github.com/intiqo/app-platform/internal/version.CommitHash=$hash'" -o ./bin/app github.com/intiqo/app-platform/cmd
