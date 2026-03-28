# Food Delivery Microservices - Startup Automation Script
# This script checks prerequisites, starts infrastructure, and runs all services

param(
    [switch]$SkipBuild,
    [switch]$UseDockerCompose,
    [switch]$DeployToK8s,
    [switch]$Cleanup
)

$ErrorActionPreference = "Continue"
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
Set-Location $ScriptDir

# Colors for output
function Write-Success { param($msg) Write-Host "[OK] $msg" -ForegroundColor Green }
function Write-Info { param($msg) Write-Host "[INFO] $msg" -ForegroundColor Cyan }
function Write-Warn { param($msg) Write-Host "[WARN] $msg" -ForegroundColor Yellow }
function Write-Fail { param($msg) Write-Host "[FAIL] $msg" -ForegroundColor Red }

Write-Host @"

==============================================
  Food Delivery Microservices - Startup Script
==============================================
"@ -ForegroundColor Magenta

# ============================================
# Phase 1: Check Prerequisites
# ============================================
Write-Host "`n[1/5] Checking Prerequisites..." -ForegroundColor Yellow

$prereqs = @{
    "Go" = @{ cmd = "go"; args = "version"; pattern = "go" }
    "Docker" = @{ cmd = "docker"; args = "--version"; pattern = "Docker" }
}

$allGood = $true

foreach ($name in $prereqs.Keys) {
    $p = $prereqs[$name]
    try {
        $output = & $p.cmd $p.args 2>&1
        if ($output -match $p.pattern) {
            Write-Success "$name installed - $($output -split '`n')[0]"
        } else {
            Write-Fail "$name not found"
            $allGood = $false
        }
    } catch {
        Write-Fail "$name not found"
        $allGood = $false
    }
}

if (-not $allGood) {
    Write-Fail "Please install missing prerequisites"
    Write-Host @"
Install with:
  - Go: https://go.dev/dl/
  - Docker: https://www.docker.com/
  - kubectl: choco install kubernetes-cli
  - kind: choco install kind
"@ -ForegroundColor Yellow
    exit 1
}

# ============================================
# Phase 2: Start Infrastructure (Docker)
# ============================================
Write-Host "`n[2/5] Checking Infrastructure..." -ForegroundColor Yellow

# Check Docker is running
try {
    docker info 2>&1 | Out-Null
    if ($LASTEXITCODE -ne 0) { throw }
    Write-Success "Docker is running"
} catch {
    Write-Fail "Docker is not running. Please start Docker Desktop."
    exit 1
}

# Start MongoDB
Write-Info "Checking MongoDB..."
$mongoRunning = docker ps --filter "name=mongodb" --format "{{.Names}}" 2>$null
if ($mongoRunning -ne "mongodb") {
    Write-Info "Starting MongoDB..."
    docker run -d --name mongodb -p 27017:27017 mongo:latest 2>$null | Out-Null
    if ($?) {
        Write-Success "MongoDB started on port 27017"
        Start-Sleep 3
    } else {
        # Try to start existing container
        docker start mongodb 2>$null | Out-Null
        Write-Success "MongoDB started"
    }
} else {
    Write-Success "MongoDB is already running"
}

# Start RabbitMQ
Write-Info "Checking RabbitMQ..."
$rabbitRunning = docker ps --filter "name=rabbitmq" --format "{{.Names}}" 2>$null
if ($rabbitRunning -ne "rabbitmq") {
    Write-Info "Starting RabbitMQ..."
    docker run -d --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:3-management 2>$null | Out-Null
    if ($?) {
        Write-Success "RabbitMQ started on ports 5672, 15672"
        Start-Sleep 5
    } else {
        docker start rabbitmq 2>$null | Out-Null
        Write-Success "RabbitMQ started"
    }
} else {
    Write-Success "RabbitMQ is already running"
}

# Verify RabbitMQ is ready
Write-Info "Waiting for RabbitMQ to be ready..."
$rabbitReady = $false
for ($i = 0; $i -lt 30; $i++) {
    try {
        $response = Invoke-RestMethod -Uri "http://localhost:15672/api/overview" -UserAgent "PowerShell" -ErrorAction SilentlyContinue
        if ($response) {
            $rabbitReady = $true
            break
        }
    } catch {}
    Start-Sleep 1
}
if ($rabbitReady) {
    Write-Success "RabbitMQ management UI ready"
} else {
    Write-Warn "RabbitMQ may not be fully ready yet"
}

# ============================================
# Phase 3: Build Services
# ============================================
if (-not $SkipBuild) {
    Write-Host "`n[3/5] Building Services..." -ForegroundColor Yellow
    
    $services = @(
        @{ name = "Auth Service"; path = "services/auth-service"; port = 8081 },
        @{ name = "Restaurant Service"; path = "services/restaurant-service"; port = 8082 },
        @{ name = "Order Service"; path = "services/order-service"; port = 8083 },
        @{ name = "Payment Service"; path = "services/payment-service"; port = 8085 },
        @{ name = "Delivery Service"; path = "services/delivery-service"; port = 8086 },
        @{ name = "Notification Service"; path = "services/notification-service"; port = 8087 },
        @{ name = "Gateway"; path = "gateway"; port = 8080 }
    )

    foreach ($svc in $services) {
        Write-Info "Building $($svc.name)..."
        Set-Location "$ScriptDir\$($svc.path)"
        
        # Clean old binary
        $exeName = ($svc.path -split "/")[-1]
        Remove-Item "$exeName.exe" -ErrorAction SilentlyContinue
        
        # Build
        go build -o "$exeName.exe" .
        if ($LASTEXITCODE -eq 0) {
            Write-Success "$($svc.name) built"
        } else {
            Write-Fail "$($svc.name) build failed"
        }
    }
    
    Set-Location $ScriptDir
}

# ============================================
# Phase 4: Start Services
# ============================================
Write-Host "`n[4/5] Starting Services..." -ForegroundColor Yellow

# Kill existing processes on ports
$ports = @(8081, 8082, 8083, 8084, 8085, 8086, 8087, 8080)
foreach ($port in $ports) {
    $proc = Get-NetTCPConnection -LocalPort $port -ErrorAction SilentlyContinue | 
            Select-Object -ExpandProperty OwningProcess -First 1
    if ($proc) {
        Stop-Process -Id $proc -Force -ErrorAction SilentlyContinue
    }
}

# Start services in order
$servicesToStart = @(
    @{ name = "Auth Service"; path = "services/auth-service"; exe = "auth-service.exe"; port = 8081 },
    @{ name = "Restaurant Service"; path = "services/restaurant-service"; exe = "restaurant-service.exe"; port = 8082 },
    @{ name = "Order Service"; path = "services/order-service"; exe = "order-service.exe"; port = 8083 },
    @{ name = "Payment Service"; path = "services/payment-service"; exe = "payment-service.exe"; port = 8085 },
    @{ name = "Delivery Service"; path = "services/delivery-service"; exe = "delivery-service.exe"; port = 8086 },
    @{ name = "Notification Service"; path = "services/notification-service"; exe = "notification-service.exe"; port = 8087 },
    @{ name = "Gateway"; path = "gateway"; exe = "gateway.exe"; port = 8080 }
)

foreach ($svc in $servicesToStart) {
    Write-Info "Starting $($svc.name)..."
    $fullPath = "$ScriptDir\$($svc.path)\$($svc.exe)"
    
    if (Test-Path $fullPath) {
        Start-Process -FilePath $fullPath -WorkingDirectory "$ScriptDir\$($svc.path)" -WindowStyle Hidden
        Start-Sleep 1
        Write-Success "$($svc.name) started on port $($svc.port)"
    } else {
        Write-Fail "$($svc.exe) not found at $fullPath"
    }
}

# Wait for services to be ready
Write-Host "`nWaiting for services to be ready..." -ForegroundColor Cyan
Start-Sleep 3

# ============================================
# Phase 5: Verify Services
# ============================================
Write-Host "`n[5/5] Verifying Services..." -ForegroundColor Yellow

$healthEndpoints = @{
    "Auth Service" = "http://localhost:8081/health"
    "Restaurant Service" = "http://localhost:8082/health"
    "Order Service" = "http://localhost:8083/health"
    "Payment Service" = "http://localhost:8085/health"
    "Delivery Service" = "http://localhost:8086/health"
    "Notification Service" = "http://localhost:8087/health"
    "API Gateway" = "http://localhost:8080/api/health"
}

$allHealthy = $true

foreach ($name in $healthEndpoints.Keys) {
    try {
        $response = Invoke-RestMethod -Uri $healthEndpoints[$name] -TimeoutSec 3 -ErrorAction Stop
        if ($response.status -eq "healthy") {
            Write-Success "$name - Healthy"
        }
    } catch {
        Write-Warn "$name - Not responding"
        $allHealthy = $false
    }
}

# ============================================
# Summary
# ============================================
Write-Host @"

==============================================
           Services Ready!
==============================================

API Gateway:     http://localhost:8080
Auth Service:    http://localhost:8081
Restaurant:     http://localhost:8082
Order Service:  http://localhost:8083
Payment:        http://localhost:8085
Delivery:       http://localhost:8086
Notification:   http://localhost:8087

RabbitMQ UI:    http://localhost:15672 (guest/guest)

Test the API:
  curl -X POST http://localhost:8080/api/auth/register -H "Content-Type: application/json" -d '{\"email\":\"test@test.com\",\"password\":\"pass123\",\"name\":\"Test\"}'

"@ -ForegroundColor Green

# Quick test
Write-Host "Running quick API test..." -ForegroundColor Cyan
try {
    $testResp = Invoke-RestMethod -Uri "http://localhost:8080/api/restaurants" -TimeoutSec 5
    if ($testResp) {
        Write-Success "API Gateway is working - restaurants endpoint responding"
    }
} catch {
    Write-Warn "Could not test API Gateway automatically"
}

Write-Host "`nAll done! Happy coding!`n" -ForegroundColor Magenta
