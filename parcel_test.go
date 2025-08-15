package main

import (
	"database/sql"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	_ "modernc.org/sqlite"
)

var (
	// randSource источник псевдо случайных чисел.
	// Для повышения уникальности в качестве seed
	// используется текущее время в unix формате (в виде числа)
	randSource = rand.NewSource(time.Now().UnixNano())
	// randRange использует randSource для генерации случайных чисел
	randRange = rand.New(randSource)
)

// getTestParcel возвращает тестовую посылку
func getTestParcel() Parcel {
	return Parcel{
		Client:    1000,
		Status:    ParcelStatusRegistered,
		Address:   "test",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

// TestAddGetDelete проверяет добавление, получение и удаление посылки
func TestAddGetDelete(t *testing.T) {
	// настройте подключение к БД
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	// Добавляем посылку
	id, err := store.Add(parcel)
	require.NoError(t, err)
	require.Greater(t, id, 0)

	// Получаем добавленную посылку
	got, err := store.Get(id)
	require.NoError(t, err)
	//assert.Equal(t, parcel.Number, got.Number)
	assert.Equal(t, parcel.Client, got.Client)
	assert.Equal(t, parcel.Status, got.Status)
	assert.Equal(t, parcel.Address, got.Address)
	assert.Equal(t, parcel.CreatedAt, got.CreatedAt)

	// Удаляем посылку
	err = store.Delete(id)
	require.Nil(t, err)

	// Повторно пытаемся получить удалённую посылку
	_, err = store.Get(id)
	require.Error(t, err)
}

// TestSetAddress проверяет обновление адреса
func TestSetAddress(t *testing.T) {
	// настройте подключение к БД
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	// Добавляем посылку
	id, err := store.Add(parcel)
	require.Nil(t, err)

	// Обновляем адрес
	newAddress := "new test address"
	err = store.SetAddress(id, newAddress)
	assert.Nil(t, err)

	// Проверяем обновление
	updatedParcel, err := store.Get(id)
	assert.Nil(t, err)
	assert.Equal(t, newAddress, updatedParcel.Address)
}

// TestSetStatus проверяет обновление статуса
func TestSetStatus(t *testing.T) {
	// настройте подключение к БД
	db, err := sql.Open("sqlite", "tracker.db")
	require.Nil(t, err)
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	// Добавляем посылку
	id, err := store.Add(parcel)
	require.Nil(t, err)

	// Обновляем статус
	err = store.SetStatus(id, ParcelStatusSent)
	require.Nil(t, err)

	// Проверяем обновление
	updatedParcel, err := store.Get(id)
	assert.Nil(t, err)
	assert.Equal(t, ParcelStatusSent, updatedParcel.Status)
}

// TestGetByClient проверяет получение посылок по идентификатору клиента
func TestGetByClient(t *testing.T) {
	// настройте подключение к БД
	db, err := sql.Open("sqlite", "tracker.db")
	require.Nil(t, err)
	defer db.Close()

	store := NewParcelStore(db)

	parcels := []Parcel{
		getTestParcel(),
		getTestParcel(),
		getTestParcel(),
	}

	// Назначаем одинаковые клиенты
	client := randRange.Intn(10_000_000)
	for i := range parcels {
		parcels[i].Client = client
	}

	// Добавляем посылки
	for i, p := range parcels {
		id, err := store.Add(p)
		require.Nil(t, err)
		parcels[i].Number = id
	}

	// Получаем посылки по клиенту
	gotParcels, err := store.GetByClient(client)
	require.NoError(t, err)
	assert.Len(t, gotParcels, len(parcels))

	// Проверяем соответствие посылок
	assert.ElementsMatch(t, parcels, gotParcels)
}
