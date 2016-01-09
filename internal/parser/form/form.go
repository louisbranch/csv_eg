package form

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/luizbranco/waukeen/internal/parser/csv"
	"github.com/luizbranco/waukeen/internal/parser/worker"
	"github.com/luizbranco/waukeen/internal/transaction"
)

const MB = 1 << (10 * 2)

type Example struct {
	Record []string
	ID     string
}

func ParseFile(p worker.Pool, r *http.Request) (Example, error) {
	e := Example{}

	err := r.ParseMultipartForm(5 * MB)
	if err != nil {
		return e, fmt.Errorf("Error reading request %s", err)
	}

	f, h, err := r.FormFile("statement")
	if err != nil {
		return e, fmt.Errorf("Error reading uploaded file %s", err)
	}

	record, err := csv.First(f)
	if err != nil {
		return e, fmt.Errorf("Error parsing csv file %s", err)
	}

	e.Record = record
	e.ID = p.Enqueue(h)

	return e, nil
}

func ParseExample(p worker.Pool, r *http.Request) ([]transaction.Transaction, error) {
	m := csv.Mapping{}
	id := r.FormValue("id")
	if id == "" {
		return nil, fmt.Errorf("Error reading file, please try again")
	}

	date := r.FormValue("date")
	i, err := strconv.Atoi(date)
	if err != nil {
		return nil, fmt.Errorf("Error parsing date column index")
	}
	m.Date = uint(i)

	format := r.FormValue("date_format")
	if format == "" {
		return nil, fmt.Errorf("Error parsing date format")
	}
	m.DateFormat = format

	desc := r.FormValue("description")
	i, err = strconv.Atoi(desc)
	if err != nil {
		return nil, fmt.Errorf("Error parsing description column index")
	}
	m.Description = uint(i)

	value := r.FormValue("value")
	i, err = strconv.Atoi(value)
	if err != nil {
		return nil, fmt.Errorf("Error parsing value column index")
	}
	m.Value = uint(i)

	return p.Retrieve(id, m)
}
