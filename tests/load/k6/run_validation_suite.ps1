param(
    [ValidateSet("telegram_battle", "telegram_big_boss", "mercadopago_webhook")]
    [string]$Scenario = "telegram_battle",
    [string]$BaseUrl = "http://localhost:8080",
    [string]$BotObservabilityUrl = "http://localhost:8080",
    [string]$WebhookObservabilityUrl = "http://localhost:8080",
    [string]$WorkerObservabilityUrl = "http://localhost:9091",
    [string]$MPWebhookToken = "",
    [switch]$RunIntegration
)

$ErrorActionPreference = "Stop"

function Require-Command {
    param([string]$Name)

    if (-not (Get-Command $Name -ErrorAction SilentlyContinue)) {
        throw "Required command '$Name' was not found in PATH."
    }
}

Require-Command "go"
Require-Command "k6"

Write-Host ""
Write-Host "== Go unit suite =="
go test ./...

if ($RunIntegration.IsPresent) {
    Write-Host ""
    Write-Host "== Go integration suite =="
    go test ./... -run Integration -count=1
}

Write-Host ""
Write-Host "== k6 baseline =="

$baselineArgs = @(
    "-ExecutionPolicy", "Bypass",
    "-File", "tests/load/k6/run_baseline.ps1",
    "-Scenario", $Scenario,
    "-BaseUrl", $BaseUrl,
    "-BotObservabilityUrl", $BotObservabilityUrl,
    "-WebhookObservabilityUrl", $WebhookObservabilityUrl,
    "-WorkerObservabilityUrl", $WorkerObservabilityUrl
)

if ($MPWebhookToken -ne "") {
    $baselineArgs += @("-MPWebhookToken", $MPWebhookToken)
}

& powershell @baselineArgs
if ($LASTEXITCODE -ne 0) {
    throw "Baseline de carga falhou com exit code $LASTEXITCODE."
}
