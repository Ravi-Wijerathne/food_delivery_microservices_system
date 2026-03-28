# Food Delivery Microservices - Cleanup Script
# Stops all services and Docker containers

Write-Host @"

==============================================
  Food Delivery Microservices - Cleanup
==============================================
"@ -ForegroundColor Magenta

# Stop all running .exe files
Write-Host "`nStopping all service processes..." -ForegroundColor Yellow

$processes = @(
    "gateway",
    "auth-service",
    "restaurant-service", 
    "order-service",
    "payment-service",
    "delivery-service",
    "notification-service"
)

foreach ($proc in $processes) {
    $running = Get-Process -Name $proc -ErrorAction SilentlyContinue
    if ($running) {
        Stop-Process -Name $proc -Force -ErrorAction SilentlyContinue
        Write-Host "  Stopped: $proc" -ForegroundColor Green
    }
}

# Kill processes by port
$ports = @(8080, 8081, 8082, 8083, 8084, 8085, 8086, 8087)
Write-Host "`nReleasing ports..." -ForegroundColor Yellow

foreach ($port in $ports) {
    $conns = Get-NetTCPConnection -LocalPort $port -ErrorAction SilentlyContinue
    foreach ($conn in $conns) {
        Stop-Process -Id $conn.OwningProcess -Force -ErrorAction SilentlyContinue
        Write-Host "  Released port: $port" -ForegroundColor Green
    }
}

# Stop Docker containers
Write-Host "`nStopping Docker containers..." -ForegroundColor Yellow

$containers = @("mongodb", "rabbitmq")
foreach ($container in $containers) {
    $running = docker ps --filter "name=$container" --format "{{.Names}}" 2>$null
    if ($running -eq $container) {
        docker stop $container 2>$null | Out-Null
        Write-Host "  Stopped: $container" -ForegroundColor Green
    }
}

# Remove containers (optional - comment out if you want to keep data)
Write-Host "`nRemoving Docker containers..." -ForegroundColor Yellow

foreach ($container in $containers) {
    docker rm -f $container 2>$null | Out-Null
    Write-Host "  Removed: $container" -ForegroundColor Green
}

Write-Host "`nCleanup complete!" -ForegroundColor Green
Write-Host "To restart: .\start.ps1`n" -ForegroundColor Cyan
