@echo off
:: DNS-SHIELD Windows Installer
:: Run as Administrator

setlocal EnableDelayedExpansion

set "INSTALL_DIR=C:\Program Files\DNS-SHIELD"
set "CONFIG_DIR=C:\ProgramData\DNS-SHIELD"
set "DATA_DIR=C:\ProgramData\DNS-SHIELD\blocklists"
set "LOG_DIR=C:\ProgramData\DNS-SHIELD\logs"
set "BIN=%INSTALL_DIR%\dns-shield.exe"
set "REPO=https://github.com/24kode-team/DNS-SHIELD"

echo.
echo  ================================================
echo   DNS-SHIELD Installer - Privacy-first DNS
echo  ================================================
echo.

:: Check admin
net session >nul 2>&1
if %errorLevel% neq 0 (
    echo [ERROR] Please run as Administrator
    echo Right-click install.bat and select "Run as administrator"
    pause
    exit /b 1
)

:: Check architecture
if "%PROCESSOR_ARCHITECTURE%"=="AMD64" (
    set "ARCH=amd64"
) else if "%PROCESSOR_ARCHITEW6432%"=="AMD64" (
    set "ARCH=amd64"
) else (
    set "ARCH=386"
)

echo [1/6] Detected Windows/%ARCH%

:: Get latest version via GitHub API
echo [2/6] Fetching latest version...
for /f "usebackq tokens=2 delims=:, " %%v in (
    `curl -fsSL "https://api.github.com/repos/dns-shield/shield/releases/latest" ^| findstr "tag_name"`
) do set "VERSION=%%~v"
set "VERSION=%VERSION:"=%"
if "%VERSION%"=="" set "VERSION=v1.0.0"
echo     Version: %VERSION%

:: Create directories
echo [3/6] Creating directories...
if not exist "%INSTALL_DIR%" mkdir "%INSTALL_DIR%"
if not exist "%CONFIG_DIR%\tls" mkdir "%CONFIG_DIR%\tls"
if not exist "%DATA_DIR%" mkdir "%DATA_DIR%"
if not exist "%LOG_DIR%" mkdir "%LOG_DIR%"

:: Download binary
echo [4/6] Downloading dns-shield.exe...
set "DL_URL=https://github.com/24kode-team/DNS-SHIELD/releases/download/%VERSION%/dns-shield_windows_%ARCH%.zip"
curl -fsSL "%DL_URL%" -o "%TEMP%\dns-shield.zip"
if %errorLevel% neq 0 (
    echo [ERROR] Download failed. Check internet connection.
    pause & exit /b 1
)
powershell -Command "Expand-Archive -Force '%TEMP%\dns-shield.zip' '%TEMP%\dns-shield-extract'"
copy /Y "%TEMP%\dns-shield-extract\dns-shield.exe" "%BIN%" >nul

:: Write config
echo [5/6] Writing config...
if not exist "%CONFIG_DIR%\shield.yaml" (
    (
        echo resolver:
        echo   listen_addr: "0.0.0.0:53"
        echo   dot_addr: "0.0.0.0:853"
        echo   doh_path: "/dns-query"
        echo   read_timeout: 5s
        echo   upstreams:
        echo     - "9.9.9.9:53"
        echo     - "149.112.112.112:53"
        echo     - "1.1.1.1:53"
        echo filter:
        echo   block_page: "0.0.0.0"
        echo   categories:
        echo     - phishing
        echo     - malware
        echo     - scam
        echo     - porn
        echo     - gambling
        echo     - predatory
        echo     - deepfake
        echo     - sextortion
        echo   allowlist:
        echo     - "canada.ca"
        echo     - "gc.ca"
        echo api:
        echo   listen_addr: "0.0.0.0:8080"
        echo blocklists:
        echo   data_dir: "C:\\ProgramData\\DNS-SHIELD\\blocklists"
        echo   refresh_every: 24h
        echo   feeds:
        echo     - name: "PhishingArmy"
        echo       url: "https://phishing.army/download/phishing_army_blocklist_extended.txt"
        echo       category: phishing
        echo       format: domains
        echo       enabled: true
        echo     - name: "URLhaus"
        echo       url: "https://urlhaus.abuse.ch/downloads/hostfile/"
        echo       category: malware
        echo       format: hosts
        echo       enabled: true
        echo     - name: "StevenBlack Adult"
        echo       url: "https://raw.githubusercontent.com/StevenBlack/hosts/master/alternates/porn/hosts"
        echo       category: porn
        echo       format: hosts
        echo       enabled: true
    ) > "%CONFIG_DIR%\shield.yaml"
)

:: Generate admin token
for /f %%i in ('powershell -Command "[System.Guid]::NewGuid().ToString('N') + [System.Guid]::NewGuid().ToString('N')"') do set "TOKEN=%%i"
echo SHIELD_ADMIN_TOKEN=%TOKEN% > "%CONFIG_DIR%\.env"

:: Install as Windows Service
echo [6/6] Installing Windows service...
sc query dns-shield >nul 2>&1
if %errorLevel% equ 0 (
    sc stop dns-shield >nul 2>&1
    sc delete dns-shield >nul 2>&1
    timeout /t 2 >nul
)

sc create dns-shield ^
    binPath= "\"%BIN%\"" ^
    DisplayName= "DNS-SHIELD Filtering Service" ^
    Description= "Privacy-first DNS filtering. Blocks phishing, malware, adult content." ^
    start= auto ^
    obj= LocalSystem

sc failure dns-shield reset= 60 actions= restart/5000/restart/5000/restart/10000

:: Set environment variable for service
reg add "HKLM\SYSTEM\CurrentControlSet\Services\dns-shield\Parameters" ^
    /v "SHIELD_ADMIN_TOKEN" /t REG_SZ /d "%TOKEN%" /f >nul

sc start dns-shield

:: Wait for health
echo.
echo Waiting for service to start...
timeout /t 4 >nul
curl -fsSL http://127.0.0.1:8080/health >nul 2>&1
if %errorLevel% equ 0 (
    echo [OK] DNS-SHIELD is running!
) else (
    echo [WARN] Service may still be starting. Check: sc query dns-shield
)

echo.
echo  ================================================
echo   DNS-SHIELD installed successfully!
echo  ================================================
echo.
echo   DNS:        127.0.0.1:53
echo   DoT:        127.0.0.1:853
echo   Dashboard:  http://127.0.0.1:8080
echo   Admin Token: %TOKEN%
echo.
echo   Config: %CONFIG_DIR%\shield.yaml
echo   Logs:   %LOG_DIR%
echo.
echo   Quick test (in cmd):
echo     nslookup canada.ca 127.0.0.1
echo.
echo  ================================================
echo.
pause
