package main

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/user/sport-booking/internal/config"
	"github.com/user/sport-booking/internal/repo"
)

func main() {
	_ = godotenv.Load()
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Config load error: %v\n", err)
		os.Exit(1)
	}

	db, err := repo.NewDB(cfg.SupabaseURL, cfg.SupabaseAnonKey)
	if err != nil {
		fmt.Printf("DB init error: %v\n", err)
		os.Exit(1)
	}

	ctx := context.Background()

	// Check if facility already exists to avoid duplicates
	facilities, err := db.ListFacilities(ctx)
	if err != nil {
		fmt.Printf("ListFacilities error: %v\n", err)
		os.Exit(1)
	}

	if len(facilities) > 0 {
		fmt.Println("Facilities already exist. Skipping seed.")
		return
	}

	fmt.Println("Creating facility...")
	err = db.CreateFacility(ctx, "Badminton Court 1", "badminton")
	if err != nil {
		fmt.Printf("CreateFacility error: %v\n", err)
		os.Exit(1)
	}

	// Need to get the ID of the created facility.
	facilities, _ = db.ListFacilities(ctx)
	if len(facilities) == 0 {
		fmt.Println("Failed to retrieve created facility")
		os.Exit(1)
	}
	facID := facilities[0].ID

	fmt.Println("Creating resource unit...")
	err = db.CreateResourceUnit(ctx, facID, "Court A")
	if err != nil {
		fmt.Printf("CreateResourceUnit error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Seed complete!")
}
