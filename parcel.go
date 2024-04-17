package main

import (
	"database/sql"
	"fmt"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {
	// реализуйте добавление строки в таблицу parcel, используйте данные из переменной p
	res, err := s.db.Exec("INSERT INTO parcel (client, status, address, created_at) VALUES (:client, :status, :address, :createdAt)",
		sql.Named("client", p.Client),
		sql.Named("status", p.Status),
		sql.Named("address", p.Address),
		sql.Named("createdAt", p.CreatedAt))
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	// реализуйте чтение строки по заданному number
	// здесь из таблицы должна вернуться только одна строка
	parcel := Parcel{}

	row := s.db.QueryRow("SELECT client, status, address, created_at FROM parcel WHERE number = :number", sql.Named("number", number))
	err := row.Scan(&parcel.Client, &parcel.Status, &parcel.Address, &parcel.CreatedAt)

	return parcel, err
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	// реализуйте чтение строк из таблицы parcel по заданному client
	// здесь из таблицы может вернуться несколько строк

	// заполните срез Parcel данными из таблицы
	var res []Parcel

	rows, err := s.db.Query("SELECT client, status, address, created_at FROM parcel WHERE client = ?", client)
	if err != nil {
		return []Parcel{}, err
	}
	for rows.Next() {
		var parcel Parcel
		if err := rows.Scan(&parcel.Client, &parcel.Status, &parcel.Address, &parcel.CreatedAt); err != nil {
			return res, err
		}
		res = append(res, parcel)
	}
	if err = rows.Err(); err != nil {
		return res, err
	}

	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	// реализуйте обновление статуса в таблице parcel
	_, err := s.db.Exec("UPDATE parcel SET status = :status WHERE number = :number",
		sql.Named("status", status),
		sql.Named("number", number))
	if err != nil {
		return err
	}

	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {
	parcel, err := s.Get(number)
	if err != nil {
		return fmt.Errorf("не удалось найти посылку с id %d", number)
	}

	// менять адрес можно только если значение статуса registered
	if parcel.Status != ParcelStatusRegistered {
		return fmt.Errorf("изменять адрес посылки только если значение статуса registered. Статус посылки с id = %d: %s", number, parcel.Status)
	}

	_, err = s.db.Exec("UPDATE parcel SET address = :address WHERE number = :number",
		sql.Named("address", address),
		sql.Named("number", number))
	if err != nil {
		return err
	}

	return nil
}

func (s ParcelStore) Delete(number int) error {
	parcel, err := s.Get(number)
	if err != nil {
		return fmt.Errorf("не удалось найти посылку с id %d", number)
	}

	// удалять строку можно только если значение статуса registered
	if parcel.Status != ParcelStatusRegistered {
		return fmt.Errorf("удалять строку можно только если значение статуса registered. Статус посылки с id = %d: %s", number, parcel.Status)
	}

	_, err = s.db.Exec("DELETE FROM parcel WHERE number = :number", sql.Named("number", number))
	if err != nil {
		return err
	}

	return nil
}
