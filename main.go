package main

import (
	"fmt"
	"net/http"

	"github.com/luizbranco/waukeen/internal/parser/csv"
	"github.com/luizbranco/waukeen/internal/views"
)

const MB = 1 << (10 * 2)

func main() {
	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			views.Render(w, "upload", nil)
		case "POST":
			err := r.ParseMultipartForm(5 * MB)
			if err != nil {
				fmt.Fprintf(w, "Error reading uploaded file %s", err)
				return
			}
			f, _, err := r.FormFile("statement")
			if err != nil {
				fmt.Fprintf(w, "Error reading uploaded file %s", err)
				return
			}
			record, err := csv.First(f)
			if err != nil {
				fmt.Fprintf(w, "Error reading uploaded file %s", err)
				return
			}
			fmt.Fprintf(w, "Record %s", record)
		default:
			fmt.Fprint(w, http.StatusMethodNotAllowed)
		}
	})
	http.ListenAndServe(":8080", nil)
}
