@echo off
set MIGRATE_PATH=..\tools\migrate.exe
set MIGRATIONS_DIR=..\migrations\test

%MIGRATE_PATH% -path %MIGRATIONS_DIR% -database "your_database_connection_string" up
pause