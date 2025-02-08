package psql

import (
	"context"
	"database/sql"

	"github.com/DexScen/VKtestTask/backend/internal/domain"
)

type Containers struct {
	db *sql.DB
}

func NewContainers(db *sql.DB) *Containers {
	return &Containers{db}
}

func (c *Containers) GetContainers(ctx context.Context, list *domain.ListContainer) error {
	rows, err := c.db.Query("SELECT ip_address, ping_time, last_successful_ping_date FROM pings")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var cont domain.Container
		var successDate sql.NullTime

		err := rows.Scan(&cont.IP, &cont.PingTime, &successDate)
		if err != nil {
			return err
		}

		if successDate.Valid {
			cont.SuccessDate = successDate.Time
		}
		*list = append(*list, cont)
	}
	return err
}

func (c *Containers) PostContainers(ctx context.Context, list *domain.ListContainer) error {
	tr, err := c.db.Begin()
	if err != nil {
		return err
	}

	_, err = tr.Exec("DELETE FROM pings")
	if err != nil {
		tr.Rollback()
		return err
	}

	statement, err := tr.Prepare("INSERT INTO pings (ip_address, ping_time, last_successful_ping_date) VALUES ($1, $2, $3)")
	if err != nil {
		tr.Rollback()
		return err
	}
	defer statement.Close()

	for _, cont := range *list {
		_, err = statement.Exec(cont.IP, cont.PingTime, cont.SuccessDate)
		if err != nil {
			tr.Rollback()
			return err
		}
	}

	return tr.Commit()
}
