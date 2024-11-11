@echo off
setlocal EnableDelayedExpansion

:: Check if test.env file exists, create if it does not
if not exist "test.env" (
    echo. > test.env
)

:: Load environment variables from test.env file
for /f "tokens=* delims=" %%i in (test.env) do (
    set "%%i"
)

:: Print running migrations message
echo Running migrations for tests

:: Construct the Postgres URL from environment variables
set "POSTGRES_URL=host=%DB_HOST% port=%DB_PORT% user=%DB_USERNAME% password=%DB_PASSWORD% dbname=%DB_DATABASE_NAME% sslmode=disable"

:: Run migrations using goose
goose --dir "./internal/database/migrations" postgres "%POSTGRES_URL%" up

endlocal
