@echo off
setlocal EnableDelayedExpansion

:: Get the DB_DATABASE value from test.env
for /f "tokens=2 delims==" %%i in ('findstr "^DB_DATABASE=" test.env') do set "DB_DATABASE=%%i"

:: Drop the public schema
psql -d "%DB_DATABASE%" -c "DROP SCHEMA public CASCADE;"

:: Create the public schema
psql -d "%DB_DATABASE%" -c "CREATE SCHEMA public;"

endlocal
