package web

import (
	"github.com/gin-gonic/gin"
	"text/template"
)

func Render(w gin.ResponseWriter, templateMain string) {
	templates, err := template.ParseGlob("./data/templates/*")
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error while parsing tmepalates :: " + err.Error()))
		return
	}
	err = templates.ExecuteTemplate(w, templateMain, nil)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error while executing tmepalates :: " + err.Error()))
		return

	}
	return

}
