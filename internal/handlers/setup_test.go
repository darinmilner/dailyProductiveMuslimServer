package handlers

import (
	"encoding/gob"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"server/everydaymuslimappserver/internal/config"
	"server/everydaymuslimappserver/internal/models"
	"server/everydaymuslimappserver/internal/render"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/justinas/nosurf"
)

var app config.AppConfig
var session *scs.SessionManager

var functions = template.FuncMap{}

const pathToTemplates = "./../../templates"

func TestMain(m *testing.M) {
	//put into the session
	gob.Register(models.User{})
	//Change to true when in production
	app.InProduction = false

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	app.InfoLog = infoLog

	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour

	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction //True in Production

	//test mail chan without sending mail
	mailChan := make(chan models.MailData)
	app.MailChan = mailChan

	defer close(mailChan)

	listenForMail()

	app.Session = session
	tc, err := CreateTestTemplateCache()
	if err != nil {
		log.Fatal("Can not create template cache", err)

	}

	app.TemplateCache = tc
	app.UseCache = true

	repo := NewTestRepo(&app)

	NewHandlers(repo)

	render.NewRenderer(&app)

	os.Exit(m.Run())
}

func listenForMail() {
	go func() {
		for {
			_ = <-app.MailChan
		}
	}()
}

func GetRoutes() http.Handler {

	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(WriteToConsole)
	//mux.Use(NoSurf)
	mux.Use(SessionLoad)

	hadithHandler := NewHadithHandlers()
	ayahHandler := NewAyahsHandlers()
	duaHandler := NewDuaHandlers()
	surahHandler := NewSurahHandlers()

	mux.Get("/", Repo.Home)
	mux.Get("/about", Repo.About)

	mux.Get("/hadiths", hadithHandler.GetHadith)
	mux.Get("/hadiths/{id}", hadithHandler.GetHadith)

	mux.Get("/duas", duaHandler.GetDuas)
	mux.Get("/duas/{id}", duaHandler.GetDuas)

	mux.Get("/ayahs", ayahHandler.GetAyahs)
	mux.Get("/ayahs/{id}", ayahHandler.GetAyahs)

	mux.Get("/surahs", surahHandler.GetSurahs)
	mux.Get("/surahs/{id}", surahHandler.GetSurahs)

	mux.Get("/*", Repo.DoesNotExistPage)

	fileServer := http.FileServer(http.Dir("./static/"))

	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}

//WriteToConsole middleware--USELESS
func WriteToConsole(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Hit the page")
		next.ServeHTTP(w, r)
	})
}

//NoSurf adds CSRF to all POST requests
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.InProduction,
		SameSite: http.SameSiteLaxMode,
	})
	return csrfHandler
}

//SessionLoad middleware loads and saves the session on each request
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

func CreateTestTemplateCache() (map[string]*template.Template, error) {

	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.html", pathToTemplates))

	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.html", pathToTemplates))
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.html", pathToTemplates))
			if err != nil {
				return myCache, err
			}
		}
		myCache[name] = ts
	}
	return myCache, nil
}
