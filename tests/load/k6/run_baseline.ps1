param(
    [ValidateSet("telegram_battle", "telegram_big_boss", "mercadopago_webhook")]
    [string]$Scenario = "telegram_battle",
    [string]$BaseUrl = "http://localhost:8080",
    [string]$BotObservabilityUrl = "http://localhost:8080",
    [string]$WebhookObservabilityUrl = "http://localhost:8080",
    [string]$WorkerObservabilityUrl = "http://localhost:9091",
    [string]$OutputDir = "tests/load/k6/results",
    [string]$MPWebhookToken = ""
)

$ErrorActionPreference = "Stop"

function Require-Command {
    param([string]$Name)

    if (-not (Get-Command $Name -ErrorAction SilentlyContinue)) {
        throw "Required command '$Name' was not found in PATH."
    }
}

function Save-DebugVarsSnapshot {
    param(
        [string]$Url,
        [string]$Path
    )

    try {
        $response = Invoke-WebRequest -Uri $Url -UseBasicParsing -TimeoutSec 10
        Set-Content -Path $Path -Value $response.Content -Encoding UTF8
    }
    catch {
        $payload = @{
            url = $Url
            error = $_.Exception.Message
            captured_at = (Get-Date).ToUniversalTime().ToString("o")
        } | ConvertTo-Json -Depth 4
        Set-Content -Path $Path -Value $payload -Encoding UTF8
    }
}

Require-Command "k6"

$scenarioFile = Join-Path "tests/load/k6/scenarios" ($Scenario + ".js")
if (-not (Test-Path $scenarioFile)) {
    throw "Scenario file not found: $scenarioFile"
}

$timestamp = Get-Date -Format "yyyyMMdd-HHmmss"
$runDir = Join-Path $OutputDir ($Scenario + "-" + $timestamp)
New-Item -ItemType Directory -Path $runDir -Force | Out-Null

$botVarsBefore = Join-Path $runDir "bot_debug_vars_before.json"
$botVarsAfter = Join-Path $runDir "bot_debug_vars_after.json"
$webhookVarsBefore = Join-Path $runDir "webhook_debug_vars_before.json"
$webhookVarsAfter = Join-Path $runDir "webhook_debug_vars_after.json"
$workerVarsBefore = Join-Path $runDir "worker_debug_vars_before.json"
$workerVarsAfter = Join-Path $runDir "worker_debug_vars_after.json"
$k6Summary = Join-Path $runDir "k6_summary.json"
$k6Stdout = Join-Path $runDir "k6_stdout.txt"

Save-DebugVarsSnapshot -Url ($BotObservabilityUrl.TrimEnd("/") + "/debug/vars") -Path $botVarsBefore
Save-DebugVarsSnapshot -Url ($WebhookObservabilityUrl.TrimEnd("/") + "/debug/vars") -Path $webhookVarsBefore
Save-DebugVarsSnapshot -Url ($WorkerObservabilityUrl.TrimEnd("/") + "/debug/vars") -Path $workerVarsBefore

$env:BASE_URL = $BaseUrl
if ($MPWebhookToken -ne "") {
    $env:MP_WEBHOOK_TOKEN = $MPWebhookToken
}

$k6Command = @(
    "run",
    "--summary-export", $k6Summary,
    $scenarioFile
)

$previousErrorActionPreference = $ErrorActionPreference
$ErrorActionPreference = "Continue"
& k6 @k6Command 2>&1 | Tee-Object -FilePath $k6Stdout
$k6ExitCode = $LASTEXITCODE
$ErrorActionPreference = $previousErrorActionPreference

Save-DebugVarsSnapshot -Url ($BotObservabilityUrl.TrimEnd("/") + "/debug/vars") -Path $botVarsAfter
Save-DebugVarsSnapshot -Url ($WebhookObservabilityUrl.TrimEnd("/") + "/debug/vars") -Path $webhookVarsAfter
Save-DebugVarsSnapshot -Url ($WorkerObservabilityUrl.TrimEnd("/") + "/debug/vars") -Path $workerVarsAfter

Write-Host "Baseline run saved to $runDir"
if ($k6ExitCode -ne 0) {
    throw "k6 terminou com exit code $k6ExitCode. Veja $k6Stdout e $k6Summary."
}
