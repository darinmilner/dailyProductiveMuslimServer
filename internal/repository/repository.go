package repository

import "server/everydaymuslimappserver/internal/models"

type DatabaseRepo interface {
	AllUsers() bool
	UpdateUser(m models.User) error
	Authenticate(email, testPassword string) (int, string, error)

	InsertReservation(res models.Reservation) (int, error)
	InsertCounselingTimeRestriction(r models.CounselingSessionTimeRestriction) error

	AllReservations() ([]models.Reservation, error)
	AllNewReservations() ([]models.Reservation, error)
	GetReservationByID(id int) (models.Reservation, error)

	UpdateReservation(u models.Reservation) error

	DeleteReservation(id int) error
	UpdateProcessedForReservation(id, processed int) error
}
