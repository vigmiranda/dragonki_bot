param(
    [Parameter(Mandatory = $true)]
    [string]$RunDir
)

$ErrorActionPreference = "Stop"

function Load-JsonFile {
    param([string]$Path)

    if (-not (Test-Path $Path)) {
        throw "File not found: $Path"
    }

    $raw = Get-Content $Path -Raw -Encoding UTF8
    if ([string]::IsNullOrWhiteSpace($raw)) {
        return $null
    }
    return $raw | ConvertFrom-Json -Depth 100
}

function Flatten-Counters {
    param(
        [object]$Node,
        [string]$Prefix = ""
    )

    $result = @{}
    if ($null -eq $Node) {
        return $result
    }

    if ($Node -is [System.Management.Automation.PSCustomObject]) {
        foreach ($prop in $Node.PSObject.Properties) {
            $name = if ($Prefix) { "$Prefix.$($prop.Name)" } else { $prop.Name }
            $child = Flatten-Counters -Node $prop.Value -Prefix $name
            foreach ($key in $child.Keys) {
                $result[$key] = $child[$key]
            }
        }
        return $result
    }

    if ($Node -is [System.Collections.IDictionary]) {
        foreach ($key in $Node.Keys) {
            $name = if ($Prefix) { "$Prefix.$key" } else { [string]$key }
            $child = Flatten-Counters -Node $Node[$key] -Prefix $name
            foreach ($childKey in $child.Keys) {
                $result[$childKey] = $child[$childKey]
            }
        }
        return $result
    }

    if ($Node -is [System.Array]) {
        return $result
    }

    if ($Prefix) {
        $result[$Prefix] = $Node
    }
    return $result
}

function To-NumberOrNull {
    param([object]$Value)

    if ($null -eq $Value) {
        return $null
    }

    $number = 0.0
    if ([double]::TryParse([string]$Value, [ref]$number)) {
        return $number
    }

    return $null
}

function Compare-Snapshot {
    param(
        [string]$Service,
        [string]$BeforePath,
        [string]$AfterPath
    )

    $before = Load-JsonFile $BeforePath
    $after = Load-JsonFile $AfterPath

    $beforeFlat = Flatten-Counters $before
    $afterFlat = Flatten-Counters $after

    $keys = @($beforeFlat.Keys + $afterFlat.Keys | Sort-Object -Unique)
    $rows = foreach ($key in $keys) {
        $beforeValue = $beforeFlat[$key]
        $afterValue = $afterFlat[$key]
        $beforeNumber = To-NumberOrNull $beforeValue
        $afterNumber = To-NumberOrNull $afterValue
        $delta = $null
        if ($null -ne $beforeNumber -and $null -ne $afterNumber) {
            $delta = $afterNumber - $beforeNumber
        }

        [PSCustomObject]@{
            service = $Service
            key = $key
            before = $beforeValue
            after = $afterValue
            delta = $delta
        }
    }

    return $rows
}

$botRows = Compare-Snapshot `
    -Service "bot" `
    -BeforePath (Join-Path $RunDir "bot_debug_vars_before.json") `
    -AfterPath (Join-Path $RunDir "bot_debug_vars_after.json")

$webhookRows = Compare-Snapshot `
    -Service "webhook" `
    -BeforePath (Join-Path $RunDir "webhook_debug_vars_before.json") `
    -AfterPath (Join-Path $RunDir "webhook_debug_vars_after.json")

$workerRows = Compare-Snapshot `
    -Service "worker" `
    -BeforePath (Join-Path $RunDir "worker_debug_vars_before.json") `
    -AfterPath (Join-Path $RunDir "worker_debug_vars_after.json")

$allRows = @($botRows + $webhookRows + $workerRows)

$csvPath = Join-Path $RunDir "debug_vars_diff.csv"
$jsonPath = Join-Path $RunDir "debug_vars_diff.json"

$allRows | Export-Csv -NoTypeInformation -Path $csvPath -Encoding UTF8
$allRows | ConvertTo-Json -Depth 6 | Set-Content -Path $jsonPath -Encoding UTF8

$interesting = $allRows |
    Where-Object { $null -ne $_.delta -and $_.delta -ne 0 } |
    Sort-Object service, key

Write-Host ""
Write-Host "Diff summary for $RunDir"
Write-Host ""
if (-not $interesting) {
    Write-Host "No numeric deltas were detected."
} else {
    $interesting | Format-Table -AutoSize
}
Write-Host ""
Write-Host "Saved:"
Write-Host " - $csvPath"
Write-Host " - $jsonPath"
