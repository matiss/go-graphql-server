package migrations

import (
	"github.com/go-pg/pg/v9/orm"
	migrations "github.com/robinjoseph08/go-pg-migrations/v2"
)

func init() {
	up := func(db orm.DB) error {
		_, err := db.Exec(`
			CREATE TABLE users (
				id SERIAL PRIMARY KEY,
				email VARCHAR(255) NOT NULL UNIQUE,
				name VARCHAR(100) NOT NULL,
				password VARCHAR(60) NOT NULL,
				number VARCHAR(50),
				status SMALLINT NOT NULL DEFAULT 0,
				role SMALLINT NOT NULLL DEFAULT 1,
				login_count INT NOT NULL DEFAULT 0,
				login_ip VARCHAR(45),
				login_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
				created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
				deleted_at TIMESTAMPTZ
			);

			CREATE UNIQUE INDEX users_email_idx ON users (email);
		`)
		return err
	}

	down := func(db orm.DB) error {
		_, err := db.Exec(`DROP TABLE users;`)
		return err
	}

	opts := migrations.MigrationOptions{}

	migrations.Register("20200220123508_create_users_table", up, down, opts)
}
