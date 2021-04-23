package models

import "server/everydaymuslimappserver/internal/forms"

//TemplateData struct holds the template settings
type TemplateData struct {
	StringMap       map[string]string
	IntMap          map[string]int
	FloatMap        map[string]float32
	Data            map[string]interface{}
	CSRFToken       string
	Flash           string
	Warning         string
	Error           string
	Day             int
	Month           string
	Form            *forms.Form
	IsAuthenticated int
}
