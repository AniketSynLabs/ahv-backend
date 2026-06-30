package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// DB is the raw sql.DB pool used by all service queries.
var DB *sql.DB

// gormDB is used only for AutoMigrate.
var gormDB *gorm.DB

// -- GORM model structs (table names match the ahv_worldwide schema) ----------

type SiteSetting struct {
	Key   string `gorm:"primaryKey;column:key"`
	Value string `gorm:"not null;default:''"`
}

func (SiteSetting) TableName() string { return "ahv_worldwide.site_settings" }

type Lead struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	Name      string    `gorm:"not null"`
	Email     string    `gorm:"not null;default:''"`
	Phone     string    `gorm:"not null;default:''"`
	Country   string    `gorm:"not null;default:''"`
	Inquiry   string    `gorm:"not null;default:''"`
	Message   string    `gorm:"not null;default:''"`
	Source    string    `gorm:"not null;default:'contact_form'"`
	Status    string    `gorm:"not null;default:'New'"`
	CreatedAt time.Time `gorm:"not null;autoCreateTime"`
}

func (Lead) TableName() string { return "ahv_worldwide.leads" }

// -- Connection ---------------------------------------------------------------

func Connect(dsn string) error {
	var err error
	DB, err = sql.Open("pgx", dsn)
	if err != nil {
		return fmt.Errorf("open db: %w", err)
	}

	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)
	DB.SetConnMaxLifetime(30 * time.Minute)
	DB.SetConnMaxIdleTime(10 * time.Minute)

	if err = DB.PingContext(context.Background()); err != nil {
		return fmt.Errorf("ping db: %w", err)
	}
	log.Println("✅ Database connected")

	// Open a GORM instance over the same pool for AutoMigrate.
	gormDB, err = gorm.Open(postgres.New(postgres.Config{Conn: DB}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "ahv_worldwide.",
			SingularTable: true,
		},
	})
	if err != nil {
		return fmt.Errorf("gorm open: %w", err)
	}

	return nil
}

// -- Migration ----------------------------------------------------------------

func Migrate() error {
	// Ensure the schema exists first (GORM AutoMigrate won't create schemas).
	if _, err := DB.Exec(`CREATE SCHEMA IF NOT EXISTS ahv_worldwide`); err != nil {
		return fmt.Errorf("create schema: %w", err)
	}

	// AutoMigrate creates / alters tables to match the struct definitions above.
	if err := gormDB.AutoMigrate(&SiteSetting{}, &Lead{}); err != nil {
		return fmt.Errorf("auto migrate: %w", err)
	}

	// Seed default site settings (idempotent).
	if err := seedSettings(); err != nil {
		return fmt.Errorf("seed settings: %w", err)
	}

	log.Println("Schema migrated (ahv_worldwide)")
	return nil
}

func seedSettings() error {
	defaults := []SiteSetting{
		{Key: "brand_name", Value: "AHV Worldwide"},
		{Key: "tagline", Value: "Go Global"},
		{Key: "whatsapp", Value: "+919921481220"},
		{Key: "email", Value: "info@ahvworldwide.com"},
		{Key: "address", Value: "Navi Mumbai, Maharashtra, India"},
		{Key: "phone", Value: "+919921481220"},
	}
	return gormDB.Clauses().Save(&defaults).Error
}
