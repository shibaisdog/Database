@echo off
go build -o server.exe ./server/server.go
server.exe "../@swap" ":8080"