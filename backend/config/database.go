package config

import (
	"fmt"
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"hrm-backend/models"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
	var db *gorm.DB
	var err error

	databaseURL := os.Getenv("DATABASE_URL")
	var dsn string

	if databaseURL != "" {
		dsn = databaseURL
		log.Println("Using DATABASE_URL environment variable for PostgreSQL connection...")
	} else {
		host := getEnv("DB_HOST", "localhost")
		user := getEnv("DB_USER", "postgres")
		password := getEnv("DB_PASSWORD", "postgres")
		dbname := getEnv("DB_NAME", "hrm_db")
		port := getEnv("DB_PORT", "5432")

		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
			host, user, password, dbname, port)
		log.Printf("Connecting to PostgreSQL at %s:%s (DB: %s)...", host, port, dbname)
	}

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	// If database missing, attempt to connect to 'postgres' system DB and CREATE DATABASE hrm_db automatically
	if err != nil && databaseURL == "" {
		host := getEnv("DB_HOST", "localhost")
		user := getEnv("DB_USER", "postgres")
		password := getEnv("DB_PASSWORD", "postgres")
		dbname := getEnv("DB_NAME", "hrm_db")
		port := getEnv("DB_PORT", "5432")

		log.Printf("Direct connect to %s failed (%v). Attempting auto-creation via default 'postgres' database...", dbname, err)
		adminDSN := fmt.Sprintf("host=%s user=%s password=%s dbname=postgres port=%s sslmode=disable TimeZone=UTC",
			host, user, password, port)
		adminDB, adminErr := gorm.Open(postgres.Open(adminDSN), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})

		if adminErr == nil {
			log.Printf("Creating PostgreSQL database '%s'...", dbname)
			adminDB.Exec(fmt.Sprintf("CREATE DATABASE %s;", dbname))
			// Retry connecting to newly created database
			db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
				Logger: logger.Default.LogMode(logger.Info),
			})
		}
	}

	// Fallback to SQLite if PostgreSQL service is unavailable locally
	if err != nil {
		log.Printf("PostgreSQL unavailable (%v). Operating on local SQLite database (hrm.db)...", err)
		db, err = gorm.Open(sqlite.Open("hrm.db"), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
		if err != nil {
			log.Fatalf("Failed to initialize database connection: %v", err)
		}
	}

	log.Println("Database connection established successfully.")

	// Auto Migration
	log.Println("Running AutoMigrate for GORM models...")
	err = db.AutoMigrate(&models.Role{}, &models.Department{}, &models.Employee{}, &models.RefreshToken{}, &models.AuditLog{})
	if err != nil {
		log.Fatalf("Database AutoMigrate failed: %v", err)
	}

	DB = db

	// Seed Initial Data
	seedDatabase(db)

	return db
}

func seedDatabase(db *gorm.DB) {
	var roleCount int64
	db.Model(&models.Role{}).Count(&roleCount)
	if roleCount > 0 {
		log.Println("Database already contains roles data. Skipping seeder.")
		return
	}

	log.Println("Seeding initial Roles, Departments & Organizational Hierarchy...")

	// 1. Seed Roles
	roles := []models.Role{
		{Name: "Admin", AccessLevel: 100},
		{Name: "HR", AccessLevel: 80},
		{Name: "Manager", AccessLevel: 50},
		{Name: "Employee", AccessLevel: 10},
	}

	for i := range roles {
		db.Create(&roles[i])
	}

	// 2. Seed Departments
	departments := []models.Department{
		{Name: "Executive", Description: "Executive Leadership & Board"},
		{Name: "People Ops", Description: "Human Resources, Talent Acquisition & Payroll"},
		{Name: "Engineering", Description: "Software Engineering & Infrastructure"},
		{Name: "Sales", Description: "Enterprise Sales & Account Management"},
		{Name: "Finance", Description: "Corporate Finance & Disbursals"},
	}

	for i := range departments {
		db.Create(&departments[i])
	}

	// Reference IDs
	adminRole := roles[0]
	hrRole := roles[1]
	managerRole := roles[2]
	employeeRole := roles[3]

	execDept := departments[0]
	hrDept := departments[1]
	engDept := departments[2]

	defaultPasswordHash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	// CEO / Top Manager
	ceo := models.Employee{
		FirstName:    "Alex",
		LastName:     "Vance",
		Email:        "admin@company.com",
		PasswordHash: string(defaultPasswordHash),
		RoleID:       adminRole.ID,
		DepartmentID: &execDept.ID,
		Status:       "Active",
		ManagerID:    nil,
	}
	db.Create(&ceo)

	// HR Director
	hrDir := models.Employee{
		FirstName:    "Sarah",
		LastName:     "Connor",
		Email:        "hr@company.com",
		PasswordHash: string(defaultPasswordHash),
		RoleID:       hrRole.ID,
		DepartmentID: &hrDept.ID,
		Status:       "Active",
		ManagerID:    &ceo.ID,
	}
	db.Create(&hrDir)

	// Engineering Manager
	engManager := models.Employee{
		FirstName:    "Marcus",
		LastName:     "Brody",
		Email:        "manager@company.com",
		PasswordHash: string(defaultPasswordHash),
		RoleID:       managerRole.ID,
		DepartmentID: &engDept.ID,
		Status:       "Active",
		ManagerID:    &ceo.ID,
	}
	db.Create(&engManager)

	// Direct Reports under Engineering Manager
	dev1 := models.Employee{
		FirstName:    "Elena",
		LastName:     "Rostova",
		Email:        "employee@company.com",
		PasswordHash: string(defaultPasswordHash),
		RoleID:       employeeRole.ID,
		DepartmentID: &engDept.ID,
		Status:       "Active",
		ManagerID:    &engManager.ID,
	}
	db.Create(&dev1)

	dev2 := models.Employee{
		FirstName:    "David",
		LastName:     "Kim",
		Email:        "david@company.com",
		PasswordHash: string(defaultPasswordHash),
		RoleID:       employeeRole.ID,
		DepartmentID: &engDept.ID,
		Status:       "Active",
		ManagerID:    &engManager.ID,
	}
	db.Create(&dev2)

	log.Println("Database seeder completed successfully! Sample Users & Departments Inserted:")
	log.Println(" 1. Admin:    admin@company.com    / password123 (Access Level 100)")
	log.Println(" 2. HR:       hr@company.com       / password123 (Access Level 80)")
	log.Println(" 3. Manager:  manager@company.com  / password123 (Access Level 50)")
	log.Println(" 4. Employee: employee@company.com / password123 (Access Level 10)")
	log.Println(" 5. Employee: david@company.com    / password123 (Access Level 10)")
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
