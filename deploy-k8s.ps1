# Food Delivery Microservices - Kubernetes Deployment Script

param(
    [string]$ClusterName = "food-delivery",
    [switch]$DeleteCluster
)

$ErrorActionPreference = "Stop"
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
Set-Location $ScriptDir

function Write-Success { param($msg) Write-Host "[OK] $msg" -ForegroundColor Green }
function Write-Info { param($msg) Write-Host "[INFO] $msg" -ForegroundColor Cyan }
function Write-Warn { param($msg) Write-Host "[WARN] $msg" -ForegroundColor Yellow }
function Write-Fail { param($msg) Write-Host "[FAIL] $msg" -ForegroundColor Red }

Write-Host @"

==============================================
  Food Delivery - Kubernetes Deployment
==============================================
"@ -ForegroundColor Magenta

# Check prerequisites
Write-Host "`nChecking prerequisites..." -ForegroundColor Yellow
try {
    kubectl version --client 2>$null | Out-Null
    Write-Success "kubectl installed"
} catch {
    Write-Fail "kubectl not found"
    exit 1
}

try {
    kind version 2>$null | Out-Null
    Write-Success "kind installed"
} catch {
    Write-Fail "kind not found"
    exit 1
}

# Delete cluster if requested
if ($DeleteCluster) {
    Write-Host "`nDeleting cluster: $ClusterName" -ForegroundColor Yellow
    kind delete cluster --name $ClusterName 2>$null
    Write-Success "Cluster deleted"
    exit 0
}

# Create cluster if not exists
Write-Host "`nChecking Kind cluster..." -ForegroundColor Yellow
$clusterExists = kind get clusters 2>$null | Where-Object { $_ -eq $ClusterName }

if (-not $clusterExists) {
    Write-Info "Creating Kind cluster: $ClusterName"
    kind create cluster --name $ClusterName
    Write-Success "Cluster created"
} else {
    Write-Success "Cluster already exists"
}

# Build and load Docker images
Write-Host "`nBuilding Docker images..." -ForegroundColor Yellow

$services = @(
    "services/auth-service",
    "services/restaurant-service",
    "services/order-service",
    "services/payment-service",
    "services/delivery-service",
    "services/notification-service",
    "gateway"
)

foreach ($svc in $services) {
    Write-Info "Building $svc..."
    Set-Location "$ScriptDir\$svc"
    
    # Build for Linux (required for K8s)
    $imageName = ($svc -split "/")[-1] -replace "service", "-service"
    if ($imageName -eq "gateway") { $imageName = "gateway" }
    
    docker build -t $imageName:latest . 2>$null | Out-Null
    
    if ($LASTEXITCODE -eq 0) {
        Write-Success "Built: $imageName"
        
        # Load to Kind cluster
        Write-Info "Loading to Kind cluster..."
        kind load docker-image $imageName:latest --name $ClusterName 2>$null
        Write-Success "Loaded: $imageName"
    } else {
        Write-Fail "Failed to build: $svc"
    }
    
    Set-Location $ScriptDir
}

# Deploy infrastructure
Write-Host "`nDeploying infrastructure..." -ForegroundColor Yellow
kubectl apply -f "$ScriptDir/deployments/mongodb.yaml"
kubectl apply -f "$ScriptDir/deployments/rabbitmq.yaml"

Write-Info "Waiting for infrastructure..."
Start-Sleep 10

# Deploy services
Write-Host "`nDeploying services..." -ForegroundColor Yellow

$serviceFiles = Get-ChildItem "$ScriptDir/deployments" -Filter "*.yaml"
foreach ($file in $serviceFiles) {
    Write-Info "Deploying $($file.Name)..."
    kubectl apply -f $file.FullName
}

# Wait for pods to be ready
Write-Host "`nWaiting for pods to be ready..." -ForegroundColor Yellow
kubectl wait --for=condition=ready pod -l app=gateway --timeout=120s 2>$null
Write-Success "Gateway ready"

# Show status
Write-Host "`n==============================================
           Deployment Status
==============================================" -ForegroundColor Green

kubectl get pods -o wide
kubectl get services

Write-Host @"

==============================================
           Access Information
==============================================

Gateway:          kubectl port-forward svc/gateway 8080:80
RabbitMQ:         kubectl port-forward svc/rabbitmq-service 15672:15672
MongoDB:          kubectl port-forward svc/mongodb-service 27017:27017

Test Gateway:
  curl http://localhost:8080/api/health

"@ -ForegroundColor Cyan
