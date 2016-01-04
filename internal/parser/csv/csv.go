package csv

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/luizbranco/waukeen/internal/transaction"
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

func Parse(f io.Reader, m Mapping) ([]transaction.Transaction, error) {
	r := csv.NewReader(f)
	t := []transaction.Transaction{}
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		d, err := time.Parse(m.DateFormat, record[m.Date])
		if err != nil {
			err = fmt.Errorf("error parsing record date %s (%s) ",
				record[m.Date], err)
			return nil, err
		}

		v, err := strconv.ParseFloat(record[m.Value], 64)
		if err != nil {
			err = fmt.Errorf("error parsing record value %s (%s) ",
				record[m.Value], err)
			return nil, err
		}

		t = append(t, transaction.Transaction{
			Date:        d,
			Description: record[m.Description],
			Type:        transaction.Debit,
			Value:       uint(v * 100),
		})
	}

	return t, nil
}
