package db

import (
	"time"

	"github.com/cenkalti/backoff/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func New(dsn string, maxOpen, maxIdle, maxLifetimeSec int) (*gorm.DB, error) {
	var db *gorm.DB
	operation := func() error {
		var err error
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		return err
	}

	b := backoff.NewExponentialBackOff()
	b.MaxElapsedTime = 30 * time.Second
	if err := backoff.Retry(operation, b); err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(maxOpen)
	sqlDB.SetMaxIdleConns(maxIdle)
	sqlDB.SetConnMaxLifetime(time.Duration(maxLifetimeSec) * time.Second)

	sqlDB.SetConnMaxIdleTime(30 * time.Second)
	return db, nil
}
