package models

import "time"

type User struct {
	ID          int       `json:"id"`
	FirstName   string    `json:"firstName"`
	LastName    string    `json:"lastName"`
	Email       string    `json:"email"`
	Password    string    `json:"-"`
	Gender      string    `json:"gender"`
	AccessLevel int       `json:"accessLevel"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type Signup struct {
	FirstName string
	LastName  string
	Email     string
}

type UserRegistration struct {
	FirstName string
	LastName  string
	Email     string
}

type HijriDate struct {
	Day   int
	Month string
	Year  int
}

//MailData model
type MailData struct {
	To       string
	From     string
	Subject  string
	Content  string
	Template string
}

//Reservation is the reservations Model
type Reservation struct {
	ID                  int
	FirstName           string
	LastName            string
	Email               string
	StartTime           time.Time
	EndTime             time.Time
	Date                time.Time
	CounselingSessionID int
	CreatedAt           time.Time
	UpdatedAt           time.Time
	CounselingSession   CounselingSession
	Processed           int
}

//CounselingSession struct has the data about counseling sessions
type CounselingSession struct {
	ID            int
	CounselorName string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

//Restriction is the room DB model
type Restriction struct {
	ID              int
	RestrictionName string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

//CounselingSessionRestriction is the couseling session time restriction DB model
type CounselingSessionTimeRestriction struct {
	ID            int
	StartTime     time.Time
	EndTime       time.Time
	Date          time.Time
	ReservationID int
	RestrictionID int
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Session       CounselingSession
	Reservation   Reservation
	Restriction   Reservation
}

type CounselingRegistration struct {
	ID        int    `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Gender    string `json:"gender"`
	Reason    string `json:"reason"`
}
