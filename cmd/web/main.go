package main

import (
	"context"
	"encoding/gob"
	"log"
	"net/http"
	"os"
	"os/signal"
	"server/everydaymuslimappserver/internal/config"
	"server/everydaymuslimappserver/internal/driver"
	"server/everydaymuslimappserver/internal/handlers"
	"server/everydaymuslimappserver/internal/helpers"
	"server/everydaymuslimappserver/internal/models"
	"server/everydaymuslimappserver/internal/render"
	"time"

	"github.com/alexedwards/scs/v2"
)

//const portNumber = ":8001"

var session *scs.SessionManager
var app config.AppConfig

var infoLog *log.Logger
var errorLog *log.Logger

func main() {

	portNumber := os.Getenv("PORT")

	if portNumber == "" {
		portNumber = "8001"
	}

	db, err := run()

	if err != nil {
		log.Fatal(err)
	}

	defer db.SQL.Close()

	defer close(app.MailChan)

	log.Println("Starting Email listener...")
	listenForMail()

	log.Println("Server running on port: ", portNumber)
	srv := &http.Server{
		Addr:        ":" + portNumber,
		Handler:     routes(&app),
		IdleTimeout: 120 * time.Second,
	}

	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	log.Println("Server terminate request received, gracefully shutting down", sig)

	timeOutCtx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	srv.Shutdown(timeOutCtx)
}

func run() (*driver.DB, error) {

	//Put model into the session
	gob.Register(models.User{})
	gob.Register(models.Signup{})
	gob.Register(models.Reservation{})
	gob.Register(models.UserRegistration{})
	gob.Register(models.CounselingSession{})

	mailChan := make(chan models.MailData)
	app.MailChan = mailChan

	app.InProduction = false

	//Info log
	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour

	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction //True in Production

	app.Session = session

	log.Println("Connecting to database")
	db, err := driver.ConnectSQL("host=localhost port=5432 dbname=db user=db password=")
	if err != nil {
		log.Fatal("Can not connect to DB", err)
		return nil, err
	}

	log.Println("Connected to DB")

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Can not create template cache", err)
		return nil, err
	}

	app.TemplateCache = tc
	app.UseCache = true

	repo := handlers.NewRepo(&app, db)

	handlers.NewHandlers(repo)

	render.NewRenderer(&app)

	helpers.NewHelpers(&app)

	return db, nil
}
