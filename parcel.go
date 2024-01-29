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
	newParcel, err := s.db.Exec("INSERT INTO parcel (client, status, address, created_at) VALUES (:client, :status, :address, :created_at)",
		sql.Named("client", p.Client),
		sql.Named("status", p.Status),
		sql.Named("address", p.Address),
		sql.Named("created_at", p.CreatedAt))
	if err != nil {
		return 0, err
	}

	lastId, err := newParcel.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(lastId), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	p := Parcel{}

	getParcel := s.db.QueryRow("SELECT * FROM parcel WHERE number = :number",
		sql.Named("number", number))
	err := getParcel.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err != nil {
		return p, err
	}

	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	var res []Parcel

	getParcels, err := s.db.Query("SELECT * FROM parcel WHERE client = :client",
		sql.Named("client", client))
	if err != nil {
		return res, err
	}
	defer getParcels.Close()

	for getParcels.Next() {
		p := Parcel{}

		err := getParcels.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
		if err != nil {
			return res, err
		}

		res = append(res, p)
	}

	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	_, err := s.db.Exec("UPDATE parcel SET status = :new_status WHERE number = :number",
		sql.Named("new_status", status),
		sql.Named("number", number))
	if err != nil {
		return err
	}

	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {
	// реализуйте обновление адреса в таблице parcel
	// менять адрес можно только если значение статуса registered

	_, err := s.db.Exec("UPDATE parcel SET address = :new_address WHERE number = :number AND status = :registered",
		sql.Named("new_address", address),
		sql.Named("number", number),
		sql.Named("registered", ParcelStatusRegistered))
	if err != nil {
		return err
	}

	return nil
}

func (s ParcelStore) Delete(number int) error {
	// реализуйте удаление строки из таблицы parcel
	// удалять строку можно только если значение статуса registered

	_, err := s.db.Exec("DELETE FROM parcel WHERE number = :number AND status = :registered",
		sql.Named("number", number),
		sql.Named("registered", ParcelStatusRegistered))
	if err != nil {
		return err
	}

	return nil
}
