package csv

import (
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"strconv"
	"time"

	"github.com/luizbranco/csv_eg/internal/transaction"
)

type Mapping struct {
	Date        uint
	DateFormat  string
	Description uint
	Value       uint
}

func First(f io.Reader) ([]string, error) {
	r := csv.NewReader(f)
	return r.Read()
}

func Parse(in io.Reader, m Mapping) ([]transaction.Transaction, error) {
	r := csv.NewReader(in)
	r.LazyQuotes = true
	t := []transaction.Transaction{}
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			txt, _ := ioutil.ReadAll(in)
			return nil, fmt.Errorf("Error reading line %s (%s)", txt, err)
		}

		d, err := time.Parse(m.DateFormat, record[m.Date])
		if err != nil {
			err = fmt.Errorf("Error parsing record date %s (%s) ",
				record[m.Date], err)
			return nil, err
		}

		v, err := strconv.ParseFloat(record[m.Value], 64)
		if err != nil {
			err = fmt.Errorf("Error parsing record value %s (%s) ",
				record[m.Value], err)
			return nil, err
		}

		t = append(t, transaction.Transaction{
			Date:        d,
			Description: record[m.Description],
			Type:        transaction.Debit,
			Value:       uint(math.Abs(v * 100)),
		})
	}

	return t, nil
}
