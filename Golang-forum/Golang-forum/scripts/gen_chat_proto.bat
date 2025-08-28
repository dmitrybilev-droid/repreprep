@echo off
setlocal

set PROTO_PATH=api
set OUTPUT_PATH=internal\chat\transport\grpc\gen

protoc --go_out=%OUTPUT_PATH% --go_opt=paths=source_relative ^
       --go-grpc_out=%OUTPUT_PATH% --go-grpc_opt=paths=source_relative ^
       %PROTO_PATH%\chat_service.proto

endlocal