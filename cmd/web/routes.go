package main

import (
	"net/http"
	"server/everydaymuslimappserver/internal/config"
	"server/everydaymuslimappserver/internal/handlers"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func routes(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)
	mux.Use(WriteToConsole)
	//Session middleware
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	hadithHandler := handlers.NewHadithHandlers()
	ayahHandler := handlers.NewAyahsHandlers()
	duaHandler := handlers.NewDuaHandlers()
	surahHandler := handlers.NewSurahHandlers()

	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)

	mux.Get("/hadiths", hadithHandler.GetHadith)
	mux.Get("/hadiths/{id}", hadithHandler.GetHadith)

	mux.Get("/ayahs", ayahHandler.GetAyahs)
	mux.Get("/ayahs/{id}", ayahHandler.GetAyahs)

	mux.Get("/duas", duaHandler.GetDuas)
	mux.Get("/duas/{id}", duaHandler.GetDuas)

	mux.Get("/surahs", surahHandler.GetSurahs)
	mux.Get("/surahs/{id}", surahHandler.GetSurahs)

	mux.Post("/search-availability-json", handlers.Repo.AvailabilityJSON)

	mux.Get("/signup", handlers.Repo.NewsLetterSignup)
	mux.Post("/signup", handlers.Repo.PostNewsLetterSignUp)

	mux.Get("/signup-success", handlers.Repo.SignupSuccess)

	mux.Get("/create-user", handlers.Repo.UserRegistration)
	mux.Post("/create-user", handlers.Repo.PostUserRegistration)
	mux.Get("/user-created-success", handlers.Repo.RegistrationSignupSuccess)

	mux.Get("/login", handlers.Repo.ShowLogin)
	mux.Post("/login", handlers.Repo.PostShowLogin)

	mux.Get("/make-reservation", handlers.Repo.Reservation)

	mux.Get("/counseling-reservation", handlers.Repo.CounselingSessionRegistration)

	mux.Post("/make-session-reservation", handlers.Repo.PostCounselingReservation)
	mux.Get("/logout", handlers.Repo.Logout)

	// mux.Get("/admin/dashboard", handlers.Repo.AdminDashboard)
	mux.Route("/admin", func(mux chi.Router) {

		mux.Use(Auth)
		mux.Get("/dashboard", handlers.Repo.AdminDashboard)
		mux.Get("/all-reservations", handlers.Repo.AdminAllReservations)
		mux.Get("/new-reservations", handlers.Repo.AdminNewReservations)
		mux.Get("/calender", handlers.Repo.AdminReservationsCalendar)

		mux.Get("/process-reservation/{src}/{id}", handlers.Repo.AdminProcessReservation)
		mux.Get("/reservations/{src}/{id}", handlers.Repo.AdminShowReservation)
		mux.Post("/reservations/{src}/{id}", handlers.Repo.AdminPostShowReservation)

	})
	mux.Get("/*", handlers.Repo.DoesNotExistPage)

	fileServer := http.FileServer(http.Dir("./static/"))

	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
