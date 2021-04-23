package dbrepo

import (
	"errors"
	"server/everydaymuslimappserver/internal/models"
)

func (m *testDBRepo) AllUsers() bool {
	return true
}

//InsertReservation inserts a reservation to the DB
func (m *testDBRepo) InsertReservation(res models.Reservation) (int, error) {

	if res.CounselingSessionID == 2 {
		return 0, errors.New("An error occurred")
	}
	return 1, nil
}

func (m *testDBRepo) InsertCounselingTimeRestriction(r models.CounselingSessionTimeRestriction) error {
	if r.Restriction.CounselingSessionID == 200_000 {
		return errors.New("An error occurred")
	}
	return nil
}

func (m *testDBRepo) UpdateUser(models.User) error {

	return nil
}

func (m *testDBRepo) Authenticate(email, testPassword string) (int, string, error) {
	return 1, "", nil
}

//AllReservations returns a slice of all reservations
func (m *testDBRepo) AllReservations() ([]models.Reservation, error) {

	var reservations []models.Reservation

	return reservations, nil

}

//AllNewReservations returns a slice of all new reservations
func (m *testDBRepo) AllNewReservations() ([]models.Reservation, error) {

	var reservations []models.Reservation

	return reservations, nil

}

func (m *testDBRepo) GetReservationByID(id int) (models.Reservation, error) {
	var reservations models.Reservation
	return reservations, nil
}

func (m *testDBRepo) UpdateReservation(u models.Reservation) error {
	return nil
}

//DeleteReservation deletes one reservation from the DB
func (m *testDBRepo) DeleteReservation(id int) error {
	return nil
}

func (m *testDBRepo) UpdateProcessedForReservation(id, processed int) error {
	return nil
}
