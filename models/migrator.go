package models

import (
	"database/sql"
	"if_inventory/errorhandler"
)

type Migrator struct{
	migrations map[string]bool /// used as a set (i.e. bool will always be true if it exists)
}

func NewMigrator(db *sql.DB)*Migrator{

	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS migrations (name varchar(30))`)
	errorhandler.PanicOnErr(err)

	///TODO read table into map
	migs := make(map[string]bool)
	return &Migrator{
		migrations: migs,
	}
}

func (p *Migrator)CreateTable(tablename string, migrationname string,
	cols map[string]string, indexs []string, primes []string){
	///TODO ... check the migration ... create table if it doesn't exist
	/// update the migration map and the database

}


