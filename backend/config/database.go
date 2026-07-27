package config

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

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
	host := getEnv("DB_HOST", "localhost")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "postgres")
	dbname := getEnv("DB_NAME", "hrm_db")
	port := getEnv("DB_PORT", "5432")

	var dsn string

	if databaseURL != "" {
		if u, parseErr := url.Parse(databaseURL); parseErr == nil {
			dsn = u.String()
		} else {
			dsn = databaseURL
		}
		log.Println("Connecting to PostgreSQL using DATABASE_URL environment variable...")
	} else {
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=require TimeZone=UTC",
			host, user, password, dbname, port)
		log.Printf("Connecting to PostgreSQL at %s:%s (DB: %s)...", host, port, dbname)
	}

	// Use PreferSimpleProtocol: true to disable prepared statement caching (Required for Supabase PgBouncer / Transaction Pooler)
	db, err = gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	// Retry with sslmode=disable for local fallback if sslmode=require fails
	if err != nil && databaseURL == "" {
		dsnDisable := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
			host, user, password, dbname, port)
		db, err = gorm.Open(postgres.New(postgres.Config{
			DSN:                  dsnDisable,
			PreferSimpleProtocol: true,
		}), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
	}

	// Auto-creation attempt for local postgres database if database is missing
	if err != nil && databaseURL == "" && strings.Contains(err.Error(), "does not exist") {
		log.Printf("Direct connect to %s failed. Attempting auto-creation via default 'postgres' database...", dbname)
		adminDSN := fmt.Sprintf("host=%s user=%s password=%s dbname=postgres port=%s sslmode=disable TimeZone=UTC",
			host, user, password, port)
		adminDB, adminErr := gorm.Open(postgres.New(postgres.Config{
			DSN:                  adminDSN,
			PreferSimpleProtocol: true,
		}), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})

		if adminErr == nil {
			log.Printf("Creating PostgreSQL database '%s'...", dbname)
			adminDB.Exec(fmt.Sprintf("CREATE DATABASE %s;", dbname))
			db, err = gorm.Open(postgres.New(postgres.Config{
				DSN:                  dsn,
				PreferSimpleProtocol: true,
			}), &gorm.Config{
				Logger: logger.Default.LogMode(logger.Info),
			})
		}
	}

	// Fallback to SQLite ONLY if running locally without DATABASE_URL
	if err != nil {
		if databaseURL != "" {
			log.Fatalf("Fatal: Failed to connect to cloud PostgreSQL database (%v). Please verify DATABASE_URL or DB credentials.", err)
		}
		log.Printf("PostgreSQL unavailable (%v). Operating on local SQLite database (hrm.db)...", err)
		db, err = gorm.Open(sqlite.Open("hrm.db"), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
		if err != nil {
			log.Fatalf("Failed to initialize local database connection: %v", err)
		}
	}

	log.Println("Database connection established successfully.")

	// Auto Migration
	log.Println("Running AutoMigrate for GORM models...")
	err = db.AutoMigrate(
		&models.Role{},
		&models.Department{},
		&models.Employee{},
		&models.RefreshToken{},
		&models.AuditLog{},
		&models.LeaveType{},
		&models.LeaveRequest{},
		&models.LeaveBalance{},
		&models.AttendanceLog{},
	)
	if err != nil {
		log.Fatalf("Database AutoMigrate failed: %v", err)
	}

	DB = db

	// Seed Initial Data
	seedDatabase(db)

	return db
}

func seedDatabase(db *gorm.DB) {
	var empCount int64
	db.Model(&models.Employee{}).Count(&empCount)
	if empCount == 0 {
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
	}

	// Seed Leave Types & Initial Leave Balances if missing
	var leaveTypeCount int64
	db.Model(&models.LeaveType{}).Count(&leaveTypeCount)
	if leaveTypeCount == 0 {
		log.Println("Seeding initial PeoplesHR Leave Types...")
		leaveTypes := []models.LeaveType{
			{Name: "Annual Leave", MaxDaysPerYear: 14},
			{Name: "Casual Leave", MaxDaysPerYear: 7},
			{Name: "Sick Leave", MaxDaysPerYear: 7},
			{Name: "Maternity Leave", MaxDaysPerYear: 84},
		}

		for i := range leaveTypes {
			db.Create(&leaveTypes[i])
		}

		// Ensure all existing employees have LeaveBalances initialized
		var allEmps []models.Employee
		db.Find(&allEmps)

		var createdLeaveTypes []models.LeaveType
		db.Find(&createdLeaveTypes)

		for _, emp := range allEmps {
			for _, lt := range createdLeaveTypes {
				db.Create(&models.LeaveBalance{
					EmployeeID:    emp.ID,
					LeaveTypeID:   lt.ID,
					AllocatedDays: lt.MaxDaysPerYear,
					UsedDays:      0,
					RemainingDays: lt.MaxDaysPerYear,
				})
			}
		}
		log.Println("Leave Types & Employee Balances seeded successfully!")
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
