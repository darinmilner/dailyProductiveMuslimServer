package dbrepo

import (
	"context"
	"errors"
	"log"
	"server/everydaymuslimappserver/internal/models"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func (m *postgresDBRepo) AllUsers() bool {
	return true
}

//GetUserByID gets a user by ID from the DB
func (m *postgresDBRepo) GetUserByID(id int) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	query := `
	select id, first_name, last_name, email, password, access_level, created_at ,updated_at
	from users where id=$1
	`

	row := m.DB.QueryRowContext(ctx, query, id)

	var u models.User

	err := row.Scan(
		&u.ID,
		&u.FirstName,
		&u.LastName,
		&u.Email,
		&u.Password,
		&u.AccessLevel,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if err != nil {
		return u, err
	}

	return u, nil
}

//UpdateUser updates user in a DB
func (m *postgresDBRepo) UpdateUser(u models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	query := `
		update users set first_name = $1,last_name = $2, email =$3, access_level = $4, updated_at =$5
	`

	_, err := m.DB.ExecContext(ctx, query,
		u.FirstName, u.LastName, u.Email, u.AccessLevel, time.Now(),
	)

	if err != nil {
		return err
	}

	return nil
}

//Authenticate authenticates a user
func (m *postgresDBRepo) Authenticate(email, testPassword string) (int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	log.Println("Got Context")
	defer cancel()

	var id int
	var hashedPassword string

	row := m.DB.QueryRowContext(ctx, "select id, password from users where email = $1", email)

	//gets id and password from DB
	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		return id, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(testPassword))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, "", errors.New("incorrect Password")
	} else if err != nil {
		return 0, "", err
	}

	return id, hashedPassword, nil
}

//InsertReservation inserts a reservation to the DB
func (m *postgresDBRepo) InsertReservation(res models.Reservation) (int, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	var newID int

	stmt := `insert into reservations (first_name, last_name, email, date,
		start_time, end_time, counseling_session_id, created_at, updated_at)
		values($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id`

	err := m.DB.QueryRowContext(ctx, stmt,
		res.FirstName,
		res.LastName,
		res.Email,
		res.Date,
		res.StartTime,
		res.EndTime,
		res.CounselingSessionID,
		time.Now(),
		time.Now(),
	).Scan(&newID)

	if err != nil {
		return 0, err
	}

	return newID, nil
}

//InsertRoomRestriction inserts a room restriction into the DB
func (m *postgresDBRepo) InsertCounselingTimeRestriction(r models.CounselingSessionTimeRestriction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	stmt := `insert into counseling_time_restrictions (start_time, end_time, date,
			reservation_id, created_at, updated_at, restriction_id)
			values
			($1, $2, $3, $4, $5, $6, $7)`

	_, err := m.DB.ExecContext(
		ctx, stmt,
		r.StartTime,
		r.EndTime,
		r.Date,
		r.ReservationID,
		time.Now(),
		time.Now(),
		r.RestrictionID,
	)

	if err != nil {
		return err
	}
	return nil
}

//AllReservations returns a slice of all reservations
func (m *postgresDBRepo) AllReservations() ([]models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	var reservations []models.Reservation

	query := `
		select r.id, r.first_name, r.last_name, r.email,
		r.start_time, r.end_time, r.date, r.counseling_session_id, 
		r.created_at, r.updated_at, r.processed,
		cs.id, cs.counselor_name
		from reservations r 
		left join counseling_session cs on (r.counseling_session_id = cs.id)
		order by r.date asc
	`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return reservations, err
	}

	defer rows.Close()

	for rows.Next() {
		var i models.Reservation
		err := rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.Email,
			&i.StartTime,
			&i.EndTime,
			&i.Date,
			&i.CounselingSessionID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Processed,
			&i.CounselingSession.ID,
			&i.CounselingSession.CounselorName,
		)
		if err != nil {
			return reservations, err
		}

		reservations = append(reservations, i)
	}

	if err = rows.Err(); err != nil {
		return reservations, err
	}

	return reservations, nil
}

//AllNewReservations returns a slice of all new reservations
func (m *postgresDBRepo) AllNewReservations() ([]models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	var reservations []models.Reservation

	query := `
		select r.id, r.first_name, r.last_name, r.email,
		r.start_time, r.end_time, r.date, r.counseling_session_id, 
		r.created_at, r.updated_at, r.processed
		from reservations r 
		left join counseling_session cs on (r.counseling_session_id = cs.id)
		where processed = 0
		order by r.date asc
	`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return reservations, err
	}

	defer rows.Close()

	for rows.Next() {
		var i models.Reservation
		err := rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.Email,
			&i.StartTime,
			&i.EndTime,
			&i.Date,
			&i.CounselingSessionID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Processed,
			// &i.CounselingSession.ID,
			// &i.CounselingSession.CounselorName,
		)
		if err != nil {
			return reservations, err
		}

		reservations = append(reservations, i)
	}

	if err = rows.Err(); err != nil {
		return reservations, err
	}
	log.Print("Reservations in DB: ", reservations)
	return reservations, nil
}

//GetReservationByID gets on reservation by ID
func (m *postgresDBRepo) GetReservationByID(id int) (models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	var res models.Reservation

	query := `
		select r.id, r.first_name, r.last_name, r.email,
		r.start_time, r.end_time, r.date, 
		r.created_at, r.updated_at, r.processed, r.counseling_session_id, 
		cs.id, cs.counselor_name
		from reservations r
		left join counseling_session cs on (r.counseling_session_id = cs.id)
		where r.id = $1
	`

	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&res.ID,
		&res.FirstName,
		&res.LastName,
		&res.Email,
		&res.StartTime,
		&res.EndTime,
		&res.Date,
		&res.CreatedAt,
		&res.UpdatedAt,
		&res.Processed,
		&res.CounselingSessionID,
		&res.CounselingSession.ID,
		&res.CounselingSession.CounselorName,
	)

	if err != nil {
		return res, err
	}

	return res, nil
}

//UpdateUser updates user in a DB
func (m *postgresDBRepo) UpdateReservation(u models.Reservation) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	query := `
		update reservations set first_name = $1,last_name = $2, email =$3, date = $4, start_time = $5, end_time = $6, updated_at = $7
		where id = $8
	`

	_, err := m.DB.ExecContext(ctx, query,
		u.FirstName, u.LastName, u.Email, u.Date, u.StartTime, u.EndTime, time.Now(), u.ID,
	)

	if err != nil {
		return err
	}

	return nil
}

//DeleteReservation deletes one reservation from the DB
func (m *postgresDBRepo) DeleteReservation(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		delete from reservation where id = $1
	`

	_, err := m.DB.ExecContext(ctx, query, id)

	if err != nil {
		return err
	}

	return nil
}

//UpdateProcessedForReservation updates if the reservation has been processed
func (m *postgresDBRepo) UpdateProcessedForReservation(id, processed int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		update reservations set processed = $1 where id = $2
	`

	_, err := m.DB.ExecContext(ctx, query, processed, id)

	if err != nil {
		return err
	}

	return nil
}
