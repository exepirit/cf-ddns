package repository

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
	"log"
)

type sqliteRepository struct {
	db *sql.DB
}

func NewSqlite(filename string) (BindingRepository, error) {
	repo := &sqliteRepository{}
	var err error

	repo.db, err = sql.Open("sqlite3", filename)
	if err != nil {
		return nil, err
	}

	return repo, nil
}

func (repo sqliteRepository) Init() error {
	sqlStmt := `
	create table if not exists ddns_records (
		'record_id' integer primary key autoincrement not null,
		'domain' varchar(255) not null,
		'period' integer not null
	);`
	_, err := repo.db.Exec(sqlStmt)
	return err
}

func (repo sqliteRepository) Reset() error {
	sqlStmt := `delete * from ddns_record;`
	_, err := repo.db.Exec(sqlStmt)
	return err
}

func (repo sqliteRepository) GetAll() []DnsBinding {
	rows, err := repo.db.Query(`select domain, period from ddns_records;`)
	if err != nil {
		log.Panicln(err)
		return []DnsBinding{}
	}

	records := make([]DnsBinding, 0)
	var tmpRecord DnsBinding
	for rows.Next() {
		err := rows.Scan(&tmpRecord.Domain, &tmpRecord.UpdatePeriod)
		if err != nil {
			return []DnsBinding{}
		}
		records = append(records, tmpRecord)
	}
	return records
}

func (repo sqliteRepository) Add(record DnsBinding) error {
	sqlStmt := `insert into ddns_records ('domain', 'period') values (?, ?);`
	_, err := repo.db.Exec(sqlStmt, record.Domain, record.UpdatePeriod)
	return err
}

func (repo sqliteRepository) Update(binding DnsBinding) error {
	sqlStmt := `update ddns_records set (period = ?) where domain = ?;`
	rows, err := repo.db.Exec(sqlStmt, binding.UpdatePeriod, binding.Domain)
	if err != nil {
		return errors.WithMessage(err, "failed to query database")
	}

	if affected, err := rows.RowsAffected(); affected == 0 || err != nil {
		return errors.New("binding not found")
	}
	return nil
}
