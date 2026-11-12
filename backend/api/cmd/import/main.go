package main

import (
	"context"
	"database/sql"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

func main() {
	dbURL := flag.String("db", os.Getenv("DATABASE_URL"), "PostgreSQL Connection URL")
	tenantID := flag.String("tenant", "", "Tenant UUID")
	csvPath := flag.String("file", "", "Path to CSV file (columns: type[vehicle|driver],identifier[plate|cpf],name_or_model,km)")
	flag.Parse()

	if *dbURL == "" || *tenantID == "" || *csvPath == "" {
		fmt.Println("Usage: go run ./cmd/import -db <url> -tenant <uuid> -file <path.csv>")
		os.Exit(1)
	}

	db, err := sql.Open("postgres", *dbURL)
	if err != nil {
		fmt.Printf("Database connection failed: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	file, err := os.Open(*csvPath)
	if err != nil {
		fmt.Printf("Failed to open CSV file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	// Skip header if present
	header, err := reader.Read()
	if err != nil {
		fmt.Printf("Failed to read CSV header: %v\n", err)
		os.Exit(1)
	}
	_ = header

	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		fmt.Printf("Failed to start transaction: %v\n", err)
		os.Exit(1)
	}
	defer tx.Rollback()

	importedVehicles := 0
	importedDrivers := 0
	importedOdometers := 0

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Printf("CSV read error: %v\n", err)
			os.Exit(1)
		}
		if len(record) < 3 {
			continue
		}

		itemType := strings.ToLower(strings.TrimSpace(record[0]))
		ident := strings.TrimSpace(record[1])
		nameOrModel := strings.TrimSpace(record[2])

		now := time.Now().UTC()

		switch itemType {
		case "vehicle":
			vehicleID := uuid.NewString()
			_, err := tx.ExecContext(ctx, `
				INSERT INTO fleet_vehicles (id, tenant_id, plate, status, created_at, updated_at, version)
				VALUES ($1::uuid, $2::uuid, $3, 'active', $4, $4, 1)
				ON CONFLICT (tenant_id, plate) DO NOTHING
			`, vehicleID, *tenantID, ident, now)
			if err != nil {
				fmt.Printf("Error importing vehicle %s: %v\n", ident, err)
				os.Exit(1)
			}
			importedVehicles++

			// Check if initial KM provided
			if len(record) >= 4 && strings.TrimSpace(record[3]) != "" {
				km, err := strconv.ParseInt(strings.TrimSpace(record[3]), 10, 64)
				if err == nil && km >= 0 {
					_, err = tx.ExecContext(ctx, `
						INSERT INTO fleet_odometer_readings (id, tenant_id, vehicle_id, current_km, reading_at, source, is_correction, reason, created_at)
						VALUES ($1::uuid, $2::uuid, $3::uuid, $4, $5, 'legacy_import', false, 'Carga inicial de dados legados do cliente', $5)
					`, uuid.NewString(), *tenantID, vehicleID, km, now)
					if err == nil {
						importedOdometers++
					}
				}
			}

		case "driver":
			driverID := uuid.NewString()
			_, err := tx.ExecContext(ctx, `
				INSERT INTO drivers (id, tenant_id, name, tax_id, status, created_at, updated_at, version)
				VALUES ($1::uuid, $2::uuid, $3, $4, 'active', $5, $5, 1)
				ON CONFLICT (tenant_id, tax_id) DO NOTHING
			`, driverID, *tenantID, nameOrModel, ident, now)
			if err != nil {
				fmt.Printf("Error importing driver %s: %v\n", nameOrModel, err)
				os.Exit(1)
			}
			importedDrivers++
		}
	}

	if err := tx.Commit(); err != nil {
		fmt.Printf("Failed to commit import transaction: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Importação legada concluída com sucesso:\n- Veículos importados: %d\n- Leituras de hodômetro inicial: %d\n- Motoristas importados: %d\n",
		importedVehicles, importedOdometers, importedDrivers)
}

