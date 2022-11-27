// Package timescalewrapper provides a wrapper for timescale.
package timescalewrapper

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

// Database is the database.
type Database struct {
	ctx  context.Context
	pool *pgxpool.Pool
}

// SensorData is the data from the sensor.
type SensorData struct {
	Mac   string    `json:"name"`
	Idx   int       `json:"reading_idx"`
	Temps []float64 `json:"temps"`
}

// NewDatabase creates a new database.
func NewDatabase(ctx context.Context, connStr string) (*Database, error) {
	dbpool, err := pgxpool.Connect(ctx, connStr)
	if err != nil {
		return nil, fmt.Errorf("connecting to tsdb: %w", err)
	}

	db := &Database{
		ctx:  ctx,
		pool: dbpool,
	}

	return db, nil
}

// Shutdown closes all the connections to the database.
func (db *Database) Shutdown() {
	db.pool.Close()
}

// InsertData inserts the data into the database.
func (db *Database) InsertData(reading SensorData) error {
	queryInsertData := `INSERT INTO conditions (time, mac, temp_delta, raw_temp0, raw_temp1) VALUES ($1, $2, $3, $4, $5)`

	_, err := db.pool.Exec(context.Background(), queryInsertData, time.Now(), reading.Mac, reading.Temps[1]-reading.Temps[0], reading.Temps[0], reading.Temps[1])
	if err != nil {
		return fmt.Errorf("inserting data: %w", err)
	}

	return nil
}
