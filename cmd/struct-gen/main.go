package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gen"
	"gorm.io/gorm"
)

func main() {
	// Load .env file if exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, reading from system environment")
	}

	DBURL := os.Getenv("PSQL_DB_URL")
	gormdb, err := gorm.Open(postgres.Open(DBURL), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	// Initialize a *gorm.DB instance

	// Use the above `*gorm.DB` instance to initialize the generator,
	// which is required to generate structs from db when using `GenerateModel/GenerateModelAs`

	generator := gen.NewGenerator(gen.Config{
		OutPath:      "internal/pkg/query",
		ModelPkgPath: "models/gormmodels",
		Mode:         gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode
	})
	generator.UseDB(gormdb)
	generator.ApplyBasic(
		// Generate structs from all tables of current database
		generator.GenerateAllTable()...,
	)
	// Generate the code
	generator.Execute()

	// Execute the generator
	generator.Execute()
}
