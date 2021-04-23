package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"server/everydaymuslimappserver/internal/config"
	"server/everydaymuslimappserver/internal/driver"
	"server/everydaymuslimappserver/internal/forms"
	"server/everydaymuslimappserver/internal/helpers"
	"server/everydaymuslimappserver/internal/models"
	"server/everydaymuslimappserver/internal/render"
	"server/everydaymuslimappserver/internal/repository"
	"server/everydaymuslimappserver/internal/repository/dbrepo"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-chi/chi"
	"github.com/hablullah/go-hijri"
)

var Repo *Repository

//Repository is the repository type struct
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

//NewRepo creates a new repository
func NewRepo(app *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: app,
		DB:  dbrepo.NewPostgresRepo(db.SQL, app),
	}
}

//NewTestRepo creates a new test repository
func NewTestRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

//NewHandlers sets the repository for the handlers function
func NewHandlers(r *Repository) {
	Repo = r
}

//GetHijiriCalendarDay returns the date on the Hijri Calender
func GetHijriCalendarDay() hijri.HijriDate {
	today := time.Now()
	hijriDate, _ := hijri.CreateHijriDate(today, hijri.Default)
	fmt.Printf("%s %04d-%02d-%02d \n",
		today.Format("2006-01-02"),
		hijriDate.Year,
		hijriDate.Month,
		hijriDate.Day)

	return hijriDate
}

//Home page function
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {

	data := GetHijriCalendarDay()

	var hijriDate models.HijriDate

	var today int64
	today = data.Day
	month := data.Month

	hijriDate.Day = int(today)

	if month == 1 {
		hijriDate.Month = "Muharram"
	} else if month == 2 {
		hijriDate.Month = "Safar"
	} else if month == 3 {
		hijriDate.Month = "Rabbi alAwwal"
	} else if month == 4 {
		hijriDate.Month = "Rabbi alThani"
	} else if month == 5 {
		hijriDate.Month = "Jumada alAwwal"
	} else if month == 6 {
		hijriDate.Month = "Jumada alThani"
	} else if month == 7 {
		hijriDate.Month = "Rajab"
	} else if month == 8 {
		hijriDate.Month = "Shaban"
	} else if month == 9 {
		hijriDate.Month = "Ramadan"
	} else if month == 10 {
		hijriDate.Month = "Shawwal"
	} else if month == 11 {
		hijriDate.Month = "Dhu alQi'dah"
	} else if month == 12 {
		hijriDate.Month = "Dhu alHijjah"
	}

	fmt.Print("HijriDay ", hijriDate.Day)
	fmt.Print("hijriMonth ", hijriDate.Month)

	log.Print(data)

	render.Templates(w, r, "home.page.html", &models.TemplateData{
		Day:   hijriDate.Day,
		Month: hijriDate.Month,
	})
}

//About page function
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	render.Templates(w, r, "about.page.html", &models.TemplateData{})
}

func (m *Repository) DoesNotExistPage(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	//http.Error(w, "Page Not found", http.StatusNotFound)
	render.Templates(w, r, "404.page.html", &models.TemplateData{})
}

func createHeader(w http.ResponseWriter) {
	w.Header().Add("content-type", "application/json")
}

//enableCors enables Cors
func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

//Hadith struct
type Hadith struct {
	Week int    `json:"Week"`
	Text string `json:"Text"`
}

type Dua struct {
	Name        string `json:"Name"`
	Text        string `json:"Text"`
	Translation string `json:"Translation"`
}

//Surah is the struct holding data about surahs in the Quran API
type Surah struct {
	SurahName     string `json:"surahName"`
	Juz           string `json:"juz"`
	NumberOfAyahs int    `json:"numberOfAyahs"`
	Location      string `json:"location"`
	Description   string `json:"description"`
}

//jsonResponse struct
type jsonResponse struct {
	OK        bool   `json:"ok"`
	Message   string `json:"message"`
	Date      string `json:"Date"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
}

//Hadiths type
type Hadiths []Hadith

//Dua type
type Duas []Dua

//Surah type
type Surahs []Surah

//HadithHandlers struct
type hadithHandlers struct {
	sync.Mutex
	hadiths Hadiths
}

//duaHandlers struct
type duaHandlers struct {
	sync.Mutex
	duas Duas
}

//Ayahs type
type Ayahs []Ayah

//Ayah struct
type Ayah struct {
	Day  int    `json:"Day"`
	Text string `json:"Text"`
}
type ayahHandlers struct {
	ayahs Ayahs
}

type surahHandlers struct {
	surahs Surahs
}

func NewAyahsHandlers() *ayahHandlers {
	return &ayahHandlers{
		ayahs: Ayahs{
			Ayah{1, "Ayah One"},
			Ayah{2, "Ayah Two"},
		},
	}
}

func NewSurahHandlers() *surahHandlers {
	return &surahHandlers{
		surahs: Surahs{
			Surah{"Annass", "Juz Amma", 6, "Mecca", "Surah Annas text"},
			Surah{"Affalaq", "Juz Amma", 5, "Mecca", "Surah Annas text"},
			Surah{"Iklas", "Juz Amma", 4, "Mecca", "Surah Annas text"},
		},
	}
}

//GetAyahs function sends ayahs as JSON
func (h *ayahHandlers) GetAyahs(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	id, err := idFromUrl(r)
	log.Println(id)
	if err != nil {
		respondWithJSON(w, http.StatusOK, h.ayahs)
		return
	}

	if id >= len(h.ayahs) || id < 0 {
		http.Error(w, "Ayah Not Found", http.StatusNotFound)
		return
	}

	respondWithJSON(w, http.StatusOK, h.ayahs[id])

}

//GetDuas function sends duas as JSON
func (h *duaHandlers) GetDuas(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	id, err := idFromUrl(r)
	log.Println(id)
	if err != nil {
		respondWithJSON(w, http.StatusOK, h.duas)
		return
	}

	if id >= len(h.duas) || id < 0 {
		http.Error(w, "Dua Not Found", http.StatusNotFound)
		return
	}

	respondWithJSON(w, http.StatusOK, h.duas[id])

}

//GetDuas function sends duas as JSON
func (s *surahHandlers) GetSurahs(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	id, err := idFromUrl(r)
	log.Println(id)
	if err != nil {
		respondWithJSON(w, http.StatusOK, s.surahs)
		return
	}

	if id >= len(s.surahs) || id < 0 {
		http.Error(w, "Surah Not Found", http.StatusNotFound)
		return
	}

	respondWithJSON(w, http.StatusOK, s.surahs[id])

}

func NewHadithHandlers() *hadithHandlers {
	return &hadithHandlers{
		hadiths: Hadiths{
			Hadith{1, "Hadith One"},
			Hadith{2, "Hadith Two"},
		},
	}
}

func NewDuaHandlers() *duaHandlers {
	return &duaHandlers{
		duas: Duas{
			Dua{"Morning dua", "La illah il Allah", "There is no god besides Allah"},
		},
	}
}

//idFromUrl returns the id from the req.params
func idFromUrl(r *http.Request) (int, error) {
	parts := strings.Split(r.URL.String(), "/")
	if len(parts) != 3 {
		return 0, errors.New("Week Number Not Found")
	}

	id, err := strconv.Atoi(parts[len(parts)-1])
	if err != nil {
		return 0, errors.New("Not Found")
	}

	return id, nil
}

func (h *hadithHandlers) GetHadith(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	defer h.Unlock()
	h.Lock()
	id, err := idFromUrl(r)
	log.Println(id)
	//Returns all Hadith
	if err != nil {
		respondWithJSON(w, http.StatusOK, h.hadiths)
	}

	if id >= len(h.hadiths) || id < 0 {
		http.Error(w, "Hadith Not Found", http.StatusNotFound)
		return
	}

	respondWithJSON(w, http.StatusOK, h.hadiths[id])
}

func respondWithJSON(w http.ResponseWriter, code int, data interface{}) {
	jsonBytes, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(code)
	w.Write(jsonBytes)
}

//AvailabilityJSON return JSON for available time slots
func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		resp := jsonResponse{
			OK:      false,
			Message: "Internal Server Error",
		}
		out, _ := json.MarshalIndent(resp, "", "/t")
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
		return
	}

	//TODO: change to search DB for available times
	available := true

	if err != nil {
		resp := jsonResponse{
			OK:      false,
			Message: "Error connecting to Database",
		}
		out, _ := json.MarshalIndent(resp, "", "/t")
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
		return
	}

	resp := jsonResponse{
		//TODO: Add available times
		OK:      available,
		Message: "This time is available!",
	}

	//Manually construct JSON--Won't throw an error
	out, _ := json.MarshalIndent(resp, "", "     ")

	w.Header().Set("Content-Type", "application/json")

	w.Write(out)

}

func (m *Repository) SignupSuccess(w http.ResponseWriter, r *http.Request) {
	signup, ok := m.App.Session.Get(r.Context(), "signup").(models.Signup)
	if !ok {
		m.App.ErrorLog.Println("Could not get signup model from the session")
		m.App.Session.Put(r.Context(), "error", "Could not get signup from context")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	m.App.Session.Remove(r.Context(), "signup")
	data := make(map[string]interface{})
	data["signup"] = signup
	render.Templates(w, r, "signup-success.page.html", &models.TemplateData{
		Data: data,
	})
}

func (m *Repository) RegistrationSignupSuccess(w http.ResponseWriter, r *http.Request) {
	signup, ok := m.App.Session.Get(r.Context(), "user-signup").(models.UserRegistration)
	if !ok {
		m.App.ErrorLog.Println("Could not get user registration model from the session")
		m.App.Session.Put(r.Context(), "error", "Could not get user registration from context")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	m.App.Session.Remove(r.Context(), "user-signup")
	data := make(map[string]interface{})
	data["user-signup"] = signup
	render.Templates(w, r, "registered-newuser-success.page.html", &models.TemplateData{
		Data: data,
	})
}

//NewsLetterSignup page function
func (m *Repository) NewsLetterSignup(w http.ResponseWriter, r *http.Request) {
	var emptySignupForm models.Signup

	data := make(map[string]interface{})
	data["signup"] = emptySignupForm

	render.Templates(w, r, "signup.page.html", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

func (m *Repository) PostNewsLetterSignUp(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	signup := models.Signup{
		FirstName: r.Form.Get("first-name"),
		LastName:  r.Form.Get("last-name"),
		Email:     r.Form.Get("email"),
	}

	log.Print("User Data from form")
	log.Print(signup.FirstName)
	log.Print(signup.LastName)
	log.Print(signup.Email)

	form := forms.New(r.PostForm)

	form.Required("first-name", "last-name", "email")

	form.MinLength("first-name", 3)

	form.MinLength("last-name", 3)

	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["signup"] = signup
		render.Templates(w, r, "signup.page.html", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return

	}

	var user models.User
	user.FirstName = r.Form.Get("first-name")
	user.LastName = r.Form.Get("last-name")
	user.Email = r.Form.Get("email")

	log.Print(user)
	//TODO: CreateUserInDB(w, r, user)
	htmlMessage := fmt.Sprintf(`
		<strong>Thank You for Signing Up</strong><br>
		Dear %s %s, <br>
		You have successfully signed up for our bimonthly newsletter and email list for app updates.
		Please watch your inbox for our newsletters.  
		JazakAllahu Khairun

	`, signup.FirstName, signup.LastName)

	msg := models.MailData{
		To:       signup.Email,
		From:     "productivedailymuslim@aaaaaaa.com",
		Subject:  "Newsletter Signup Confirmation",
		Content:  htmlMessage,
		Template: "basic.html",
	}

	m.App.MailChan <- msg

	//Send email to User and Admin

	//Add session
	m.App.Session.Put(r.Context(), "signup", signup)

	http.Redirect(w, r, "/signup-success", http.StatusSeeOther)
}

//UserRegstration allosw users to register an account
func (m *Repository) UserRegistration(w http.ResponseWriter, r *http.Request) {
	var emptySignupForm models.UserRegistration

	data := make(map[string]interface{})
	data["user-signup"] = emptySignupForm

	render.Templates(w, r, "create-user.page.html", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

//PostUserRegistration handles the posting of a new user registration
func (m *Repository) PostUserRegistration(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	signup := models.UserRegistration{
		FirstName: r.Form.Get("first-name"),
		LastName:  r.Form.Get("last-name"),
		Email:     r.Form.Get("email"),
	}

	log.Print("User Data from registration form")
	log.Print(signup.FirstName)
	log.Print(signup.LastName)
	log.Print(signup.Email)

	form := forms.New(r.PostForm)

	form.Required("first-name", "last-name", "email")

	form.MinLength("first-name", 3)

	form.MinLength("last-name", 3)

	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["user-signup"] = signup
		render.Templates(w, r, "create-user.page.html", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return

	}

	var newUser models.UserRegistration
	newUser.FirstName = r.Form.Get("first-name")
	newUser.LastName = r.Form.Get("last-name")
	newUser.Email = r.Form.Get("email")

	log.Print(newUser)
	//TODO: CreateUserInDB(w, r, user)
	//TODO: change url on signup link
	htmlMessage := fmt.Sprintf(`
		<strong>Thank You for Registering your account with the Productive Muslim App</strong><br>
		Dear %s %s, <br>
		You have successfully signed up for the productive Muslim app. May it bring you many rewards
		and benefit.
		Please consider signing up for our newsletters. https://localhost:8001/signup  
		JazakAllahu Khairun

	`, signup.FirstName, signup.LastName)

	msg := models.MailData{
		To:       signup.Email,
		From:     "productivedailymuslim@aaaaaaa.com",
		Subject:  "Newsletter Signup Confirmation",
		Content:  htmlMessage,
		Template: "basic.html",
	}

	m.App.MailChan <- msg

	//Send email to User and Admin

	//Add session
	m.App.Session.Put(r.Context(), "user-signup", signup)

	http.Redirect(w, r, "/user-created-success", http.StatusSeeOther)
}

//CouncilingSessionRegistration allows users to register an account
func (m *Repository) CounselingSessionRegistration(w http.ResponseWriter, r *http.Request) {
	var emptySignupForm models.UserRegistration

	data := make(map[string]interface{})
	data["user-signup"] = emptySignupForm

	render.Templates(w, r, "counciling-registration.page.html", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

func (m *Repository) PostCounselingReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	signup := models.CounselingRegistration{
		FirstName: r.Form.Get("first-name"),
		LastName:  r.Form.Get("last-name"),
		Email:     r.Form.Get("email"),
		Gender:    r.Form.Get("gender"),
	}

	log.Print("counseling data from registration form")
	log.Print(signup.FirstName)
	log.Print(signup.LastName)
	log.Print(signup.Email)

	form := forms.New(r.PostForm)

	form.Required("first-name", "last-name", "email", "gender")

	form.MinLength("first-name", 3)

	form.MinLength("last-name", 3)

	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["user-signup"] = signup
		render.Templates(w, r, "counceling-signup.page.html", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return

	}

	var newSession models.CounselingRegistration
	newSession.FirstName = r.Form.Get("first-name")
	newSession.LastName = r.Form.Get("last-name")
	newSession.Email = r.Form.Get("email")

	log.Print(newSession)
	//TODO: CreateCouncelingRegistrationInDB(w, r, user)
	//TODO: change url on signup link
	htmlMessage := fmt.Sprintf(`
		<strong>Thank You for Requesting a counseling session</strong><br>
		Dear %s %s, <br>
		You have successfully signed up for a session. 
		Someone will be contacting you shortly about the time and link for the session.
		May it bring you many rewards
		and benefit.
		Please consider signing up for our newsletters. https://localhost:8001/signup  
		JazakAllahu Khairun

	`, signup.FirstName, signup.LastName)

	msg := models.MailData{
		To:       signup.Email,
		From:     "productivedailymuslim@aaaaaaa.com",
		Subject:  "Newsletter Signup Confirmation",
		Content:  htmlMessage,
		Template: "basic.html",
	}

	m.App.MailChan <- msg

	//Send email to User and Admin

	//Add session
	m.App.Session.Put(r.Context(), "counseling-signup", signup)

	http.Redirect(w, r, "/counceling-signup-success", http.StatusSeeOther)
}

//Logout logs a user out
func (m *Repository) Logout(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.Destroy(r.Context())
	_ = m.App.Session.RenewToken(r.Context())

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

//ShowLogin shows the login page
func (m *Repository) ShowLogin(w http.ResponseWriter, r *http.Request) {
	render.Templates(w, r, "login.page.html", &models.TemplateData{
		Form: forms.New(nil),
	})
}

func (m *Repository) PostShowLogin(w http.ResponseWriter, r *http.Request) {
	//Renew token when logging in and out
	_ = m.App.Session.RenewToken(r.Context())

	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	form := forms.New(r.PostForm)
	form.Required("email", "password")
	form.IsEmail("email")

	log.Println("email ", email, "password: ", password)
	if !form.Valid() {
		//Take user back
		render.Templates(w, r, "login.page.html", &models.TemplateData{
			Form: form,
		})
		return
	}

	log.Print("Authenticate in DB")
	id, _, err := m.DB.Authenticate(email, password)
	if err != nil {
		log.Println("Auth Error", err)
		m.App.Session.Put(r.Context(), "error", "Invalid login credentials")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	m.App.Session.Put(r.Context(), "userId", id)

	log.Println("logged in")
	m.App.Session.Put(r.Context(), "flash", "Logged in successfully")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

//Reservation route handler
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {

	//Pull reservation out of the session
	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		//helpers.ServerError(w, errors.New("can not get reservation from session."))
		m.App.Session.Put(r.Context(), "error", "can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	//room, err := m.DB.GetRoomByID(res.RoomID)

	// if err != nil {
	// 	m.App.Session.Put(r.Context(), "error", "can't find a room")
	// 	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	// 	return
	// }

	res.CounselingSession.CounselorName = "Session1"

	m.App.Session.Put(r.Context(), "reservation", res)

	sd := res.Date.Format("2006-01-02")
	//ed := res.EndDate.Format("2006-01-02")

	stringMap := make(map[string]string)
	stringMap["date"] = sd
	//stringMap["end-date"] = ed

	data := make(map[string]interface{})
	data["reservation"] = res

	render.Templates(w, r, "make-reservation.page.html", &models.TemplateData{
		Form:      forms.New(nil),
		Data:      data,
		StringMap: stringMap,
	})
}

func (m *Repository) AdminDashboard(w http.ResponseWriter, r *http.Request) {
	render.Templates(w, r, "admin.dashboard.page.html", &models.TemplateData{})
}
func (m *Repository) AdminNewReservations(w http.ResponseWriter, r *http.Request) {
	reservations, err := m.DB.AllNewReservations()

	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	log.Print(reservations)

	data := make(map[string]interface{})
	data["reservations"] = reservations

	log.Print(data["reservations"])

	render.Templates(w, r, "admin.new-reservations.page.html", &models.TemplateData{
		Data: data,
	})
}

func (m *Repository) AdminAllReservations(w http.ResponseWriter, r *http.Request) {
	reservations, err := m.DB.AllReservations()

	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})
	data["reservations"] = reservations

	render.Templates(w, r, "admin.all-reservations.page.html", &models.TemplateData{
		Data: data,
	})
}

//AdminShowReservation shows the reservation details
func (m *Repository) AdminShowReservation(w http.ResponseWriter, r *http.Request) {

	exploded := strings.Split(r.RequestURI, "/")

	id, err := strconv.Atoi(exploded[4])

	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	log.Print(id)

	src := exploded[3]

	stringMap := make(map[string]string)
	stringMap["src"] = src

	//Get reservation from the Database
	res, err := m.DB.GetReservationByID(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})
	data["reservations"] = res

	render.Templates(w, r, "admin.reservations.show.page.html", &models.TemplateData{
		StringMap: stringMap,
		Data:      data,
		Form:      forms.New(nil),
	})
}

//AdminPostShowReservation shows the reservation details
func (m *Repository) AdminPostShowReservation(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	exploded := strings.Split(r.RequestURI, "/")

	id, err := strconv.Atoi(exploded[4])

	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	log.Print(id)

	src := exploded[3]

	stringMap := make(map[string]string)
	stringMap["src"] = src

	//Get reservation from the Database
	res, err := m.DB.GetReservationByID(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	res.FirstName = r.Form.Get("first-name")
	res.LastName = r.Form.Get("last-name")
	res.Email = r.Form.Get("email")
	//res.Phone = r.Form.Get("phone")

	err = m.DB.UpdateReservation(res)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	m.App.Session.Put(r.Context(), "flash", "changes saved")
	http.Redirect(w, r, fmt.Sprintf("/admin/%s-reservations", src), http.StatusSeeOther)
}

//AdminReservationsCalender display the reservation calender
func (m *Repository) AdminReservationsCalendar(w http.ResponseWriter, r *http.Request) {
	render.Templates(w, r, "admin.reservations.calendar.page.html", &models.TemplateData{})
}

//AdminProcessReservation processes the reservation
func (m *Repository) AdminProcessReservation(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	src := chi.URLParam(r, "src")

	_ = m.DB.UpdateProcessedForReservation(id, 1)

	m.App.Session.Put(r.Context(), "flash", "Reservation marked as complete")
	http.Redirect(w, r, fmt.Sprintf("/admin/%s-reservations", src), http.StatusSeeOther)
}
