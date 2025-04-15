package postgres

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/AsaHero/whereismycity/internal/inerr"
	"github.com/AsaHero/whereismycity/pkg/config"
	"github.com/AsaHero/whereismycity/pkg/utility"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func New(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.DB.Host, cfg.DB.User, cfg.DB.Password, cfg.DB.Name, cfg.DB.Port, cfg.DB.Sslmode)

	// Set up log level
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,  // Slow SQL threshold
			LogLevel:                  logger.Error, // Log level
			Colorful:                  true,         // Disable color
			IgnoreRecordNotFoundError: true,
		},
	)

	// Open connection
	connection, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return nil, err
	}
	return connection, nil
}

func Error[T any](err error, operation string, entity T) error {
	switch err {
	case gorm.ErrRecordNotFound:
		return inerr.NewErrNotFound(utility.GetTypeName(entity))
	case gorm.ErrDuplicatedKey, gorm.ErrForeignKeyViolated, gorm.ErrCheckConstraintViolated:
		return inerr.NewErrConflict(utility.GetTypeName(entity))
	default:
		if err.Error() == "no rows affected" {
			return inerr.NewErrNoChanges(utility.GetTypeName(entity))
		}
		log.Println(err)
		return fmt.Errorf("failed to %s entity %s: \n %s", operation, utility.GetTypeName(entity), utility.FormatStruct(entity))
	}
}
