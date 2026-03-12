param(
    [switch]$ForceChoco,
    [switch]$ForceWinget
)

$ErrorActionPreference = "Stop"

function Get-CommandSafe {
    param([string]$Name)

    return Get-Command $Name -ErrorAction SilentlyContinue
}

function Assert-SingleInstallerChoice {
    if ($ForceChoco.IsPresent -and $ForceWinget.IsPresent) {
        throw "Use apenas uma das opcoes: -ForceChoco ou -ForceWinget."
    }
}

function Test-K6Installed {
    $k6 = Get-CommandSafe "k6"
    if (-not $k6) {
        return $false
    }

    try {
        $version = & k6 version 2>$null
        if ($LASTEXITCODE -eq 0) {
            Write-Host "k6 ja instalado:"
            Write-Host $version
            return $true
        }
    }
    catch {
    }

    return $false
}

function Install-WithWinget {
    $winget = Get-CommandSafe "winget"
    if (-not $winget) {
        return $false
    }

    Write-Host "Instalando k6 via winget..."
    & winget install k6 --source winget --accept-package-agreements --accept-source-agreements
    if ($LASTEXITCODE -ne 0) {
        throw "Falha ao instalar k6 via winget."
    }
    return $true
}

function Install-WithChoco {
    $choco = Get-CommandSafe "choco"
    if (-not $choco) {
        return $false
    }

    Write-Host "Instalando k6 via Chocolatey..."
    & choco install k6 -y
    if ($LASTEXITCODE -ne 0) {
        throw "Falha ao instalar k6 via Chocolatey."
    }
    return $true
}

function Refresh-Path {
    $machinePath = [System.Environment]::GetEnvironmentVariable("Path", "Machine")
    $userPath = [System.Environment]::GetEnvironmentVariable("Path", "User")
    $env:Path = ($machinePath, $userPath -join ";")
}

Assert-SingleInstallerChoice

if (Test-K6Installed) {
    exit 0
}

$installed = $false

if ($ForceWinget.IsPresent) {
    $installed = Install-WithWinget
} elseif ($ForceChoco.IsPresent) {
    $installed = Install-WithChoco
} else {
    $installed = Install-WithWinget
    if (-not $installed) {
        $installed = Install-WithChoco
    }
}

if (-not $installed) {
    throw "Nem winget nem Chocolatey estao disponiveis nesta maquina. Instale um deles e rode novamente."
}

Refresh-Path

if (-not (Test-K6Installed)) {
    throw "k6 foi instalado, mas nao ficou disponivel no PATH desta sessao. Abra um novo terminal e rode 'k6 version'."
}

Write-Host "Instalacao concluida com sucesso."
