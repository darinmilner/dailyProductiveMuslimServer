package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type postData struct {
	key   string
	value string
}

var theTests = []struct {
	name               string
	url                string
	method             string
	expectedStatusCode int
}{
	{"home", "/", "GET", http.StatusOK},
	{"about", "/about", "GET", http.StatusOK},
	{"notfound", "/abc", "GET", http.StatusNotFound},
	{"Get All hadith", "/hadiths", "GET", http.StatusOK},
	{"Get All ayahs", "/ayahs", "GET", http.StatusOK},
	{"Get All duas", "/duas", "GET", http.StatusOK},
	{"Get One hadith", "/hadiths/1", "GET", http.StatusOK},
	{"Get One ayah", "/ayahs/1", "GET", http.StatusOK},
	{"Get One dua", "/duas/0", "GET", http.StatusOK},
	{"Get One surah", "/surahs/1", "GET", http.StatusOK},
	{"Get all surahs", "/surahs", "GET", http.StatusOK},
	{"Get One hadith with week that does not exist", "/hadiths/1000000", "GET", http.StatusNotFound},
	{"Get One ayah with week that does not exist", "/ayahs/1000000", "GET", http.StatusNotFound},
	{"Get One dua with ID that does not exist", "/duas/1000000", "GET", http.StatusNotFound},
	{"Get One surah with ID that does not exist", "/surahs/1000000", "GET", http.StatusNotFound},
}

func TestHandlers(t *testing.T) {
	routes := GetRoutes()
	ts := httptest.NewTLSServer(routes)

	defer ts.Close()

	for _, e := range theTests {
		if e.method == "GET" {
			resp, err := ts.Client().Get(ts.URL + e.url)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s, expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}

		} else {
			// values := url.Values{}

			// for _, x := range e.params {
			// 	values.Add(x.key, x.value)
			// }

			// resp, err := ts.Client().PostForm(ts.URL+e.url, values)
			// if err != nil {
			// 	t.Log(err)
			// 	t.Fatal(err)
			// }

			// if resp.StatusCode != e.expectedStatusCode {
			// 	t.Errorf("for %s, expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
			// }
		}
	}
}
