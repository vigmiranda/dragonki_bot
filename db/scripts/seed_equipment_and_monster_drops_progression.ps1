param(
    [string]$DatabaseUrl = $env:DATABASE_URL
)

$ErrorActionPreference = "Stop"

if ([string]::IsNullOrWhiteSpace($DatabaseUrl)) {
    throw "Informe -DatabaseUrl ou defina a variavel de ambiente DATABASE_URL."
}

$scriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$runnerDir = Join-Path $scriptDir "..\..\cmd\tmp_seed_sql_runner"
$runnerPath = Join-Path $runnerDir "main.go"

$runnerSource = @'
package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	dsn := strings.TrimSpace(os.Getenv("DATABASE_URL"))
	if dsn == "" {
		log.Fatal("DATABASE_URL is required")
	}
	if len(os.Args) < 2 {
		log.Fatal("at least one sql file is required")
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()

	for _, arg := range os.Args[1:] {
		content, readErr := os.ReadFile(arg)
		if readErr != nil {
			log.Fatal(readErr)
		}
		if _, execErr := tx.ExecContext(ctx, string(content)); execErr != nil {
			log.Fatalf("execute %s: %v", filepath.Base(arg), execErr)
		}
		fmt.Printf("executed %s\n", filepath.Base(arg))
	}

	if err := tx.Commit(); err != nil {
		log.Fatal(err)
	}
}
'@

if (-not (Test-Path $runnerDir)) {
    New-Item -ItemType Directory -Path $runnerDir | Out-Null
}

try {
    Set-Content -Path $runnerPath -Value $runnerSource -Encoding ASCII

    $sqlFiles = @(
        (Join-Path $scriptDir "seed_shop_low_level_equipments_2_14.sql"),
        (Join-Path $scriptDir "seed_equipment_catalog_1_80.sql"),
        (Join-Path $scriptDir "update_equipment_item_name_prefixes.sql"),
        (Join-Path $scriptDir "seed_monster_drops_equipment_progression.sql")
    )

    $env:DATABASE_URL = $DatabaseUrl
    & go run $runnerPath @sqlFiles
}
finally {
    if (Test-Path $runnerPath) {
        Remove-Item $runnerPath -Force
    }
    if (Test-Path $runnerDir) {
        Remove-Item $runnerDir -Force
    }
}
