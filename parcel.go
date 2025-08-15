package main

import (
	"database/sql"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {
	sqlStmt := `INSERT INTO parcels (client, status, address, created_at) VALUES ($1, $2, $3, $4) RETURNING number`

	var id int
	err := s.db.QueryRow(sqlStmt, p.Client, p.Status, p.Address, p.CreatedAt).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	sqlStmt := `SELECT number, client, status, address, created_at FROM parcels WHERE number = $1 LIMIT 1`

	var p Parcel
	err := s.db.QueryRow(sqlStmt, number).Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err != nil {
		return Parcel{}, err
	}

	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	sqlStmt := `SELECT number, client, status, address, created_at FROM parcels WHERE client = $1`

	rows, err := s.db.Query(sqlStmt, client)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var parcels []Parcel
	for rows.Next() {
		var p Parcel
		err := rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		parcels = append(parcels, p)
	}

	return parcels, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	sqlStmt := `UPDATE parcels SET status = $1 WHERE number = $2`

	_, err := s.db.Exec(sqlStmt, status, number)
	return err
}

func (s ParcelStore) SetAddress(number int, address string) error {
	sqlStmt := `UPDATE parcels SET address = $1 WHERE number = $2 AND status = $3`

	_, err := s.db.Exec(sqlStmt, address, number, ParcelStatusRegistered)
	return err
}

func (s ParcelStore) Delete(number int) error {
	sqlStmt := `DELETE FROM parcels WHERE number = $1 AND status = $2`

	_, err := s.db.Exec(sqlStmt, number, ParcelStatusRegistered)
	return err
}
