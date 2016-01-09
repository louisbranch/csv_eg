package main

import (
	"fmt"
	"net/http"

	"github.com/luizbranco/waukeen/internal/parser/form"
	"github.com/luizbranco/waukeen/internal/parser/worker"
	"github.com/luizbranco/waukeen/internal/views"
)

func main() {
	p := worker.NewPool(5)

	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			views.Render(w, "upload", nil)
		case "POST":
			eg, err := form.ParseFile(p, r)
			if err != nil {
				views.Error(w, err)
			} else {
				views.Render(w, "example", eg)
			}
		default:
			fmt.Fprint(w, http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/parse", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			records, err := form.ParseExample(p, r)
			if err != nil {
				views.Error(w, err)
			} else {
				views.Render(w, "records", records)
			}
		} else {
			fmt.Fprint(w, http.StatusMethodNotAllowed)
		}
	})

	http.ListenAndServe(":8080", nil)
}
