Write-Host "Testing /api/services/execute endpoint" -ForegroundColor Green
Write-Host "=====================================" -ForegroundColor Green

Write-Host "`n1. Testing text.process service with uppercase operation:" -ForegroundColor Yellow
$body1 = @{
    service = "text.process"
    payload = @{
        text = "hello world"
        operation = "uppercase"
    }
} | ConvertTo-Json -Depth 3

try {
    $response1 = Invoke-RestMethod -Uri "http://localhost:8080/api/services/execute" -Method POST -Body $body1 -ContentType "application/json"
    Write-Host "Response:" -ForegroundColor Cyan
    $response1 | ConvertTo-Json -Depth 3
} catch {
    Write-Host "Error: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host "`n2. Testing text.process service with lowercase operation:" -ForegroundColor Yellow
$body2 = @{
    service = "text.process"
    payload = @{
        text = "HELLO WORLD"
        operation = "lowercase"
    }
} | ConvertTo-Json -Depth 3

try {
    $response2 = Invoke-RestMethod -Uri "http://localhost:8080/api/services/execute" -Method POST -Body $body2 -ContentType "application/json"
    Write-Host "Response:" -ForegroundColor Cyan
    $response2 | ConvertTo-Json -Depth 3
} catch {
    Write-Host "Error: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host "`n3. Testing text.process service with reverse operation:" -ForegroundColor Yellow
$body3 = @{
    service = "text.process"
    payload = @{
        text = "hello"
        operation = "reverse"
    }
} | ConvertTo-Json -Depth 3

try {
    $response3 = Invoke-RestMethod -Uri "http://localhost:8080/api/services/execute" -Method POST -Body $body3 -ContentType "application/json"
    Write-Host "Response:" -ForegroundColor Cyan
    $response3 | ConvertTo-Json -Depth 3
} catch {
    Write-Host "Error: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host "`n4. Testing echo service:" -ForegroundColor Yellow
$body4 = @{
    service = "echo"
    payload = @{
        message = "Hello from PowerShell API!"
    }
} | ConvertTo-Json -Depth 3

try {
    $response4 = Invoke-RestMethod -Uri "http://localhost:8080/api/services/execute" -Method POST -Body $body4 -ContentType "application/json"
    Write-Host "Response:" -ForegroundColor Cyan
    $response4 | ConvertTo-Json -Depth 3
} catch {
    Write-Host "Error: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host "`n5. Testing available services:" -ForegroundColor Yellow
try {
    $services = Invoke-RestMethod -Uri "http://localhost:8080/api/services" -Method GET
    Write-Host "Available services:" -ForegroundColor Cyan
    $services | ConvertTo-Json -Depth 3
} catch {
    Write-Host "Error: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host "`nTesting complete!" -ForegroundColor Green
