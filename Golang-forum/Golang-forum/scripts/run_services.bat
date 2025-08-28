@echo off
start "User Service" cmd /k go run cmd\user-service\main.go
start "Chat Service" cmd /k go run cmd\chat-service\main.go