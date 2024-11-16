package database

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/lmtani/learning-client-server-api/internal/entities"
)

func CreateTable(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS quotes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		code TEXT,
		codein TEXT,
		name TEXT,
		high TEXT,
		low TEXT,
		varBid TEXT,
		pctChange TEXT,
		bid TEXT,
		ask TEXT,
		timestamp TEXT,
		createDate TEXT
	)`)
	if err != nil {
		return err
	}
	return nil
}

func AddQuote(db *sql.DB, c *entities.UsdBrl, timeout time.Duration) error {
	stmt, err := db.Prepare(`INSERT INTO quotes (
                    		code,
                    		codein,
                    		name,
                    		high,
                    		low,
                    		varBid,
                    		pctChange,
                    		bid,
                    		ask,
                    		timestamp,
                    		createDate
						) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	_, err = stmt.ExecContext(ctx,
		c.Code,
		c.CodeIn,
		c.Name,
		c.High,
		c.Low,
		c.VarBid,
		c.PctChange,
		c.Bid,
		c.Ask,
		c.Timestamp,
		c.CreateDate,
	)
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return context.DeadlineExceeded
		}

		return err
	}

	return nil
}
