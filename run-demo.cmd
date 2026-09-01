@echo off
title Lunas - AI Collections Agent
cd /d "%~dp0"

echo ================================================
echo   Lunas - starting...
echo ================================================
echo.

if not exist "web\dist\index.html" (
  echo [1/2] Building the web app first ^(one time only^)...
  cd web
  call bun install || goto :error
  call bun run build || goto :error
  cd ..
)

echo [1/2] Web app   : ready (web\dist)
echo [2/2] Starting server...

echo.
echo   On this laptop : http://localhost:8080
echo   On your phone  : http://192.168.1.5:8080  (same WiFi/network)
echo.
echo   Press Ctrl+C here to stop.
echo.

go run ./cmd/lunas
goto :eof

:error
echo.
echo Build failed - check the messages above.
pause
