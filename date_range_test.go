package daterange_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/theseanything/daterange-parser"
)

func TestParseError(t *testing.T) {
	_, err := daterange.NewParser().Parse("")

	if err == nil {
		t.Error("No error was raised")
	}
}

func TestParse(t *testing.T) {
	p := daterange.NewParser()

	var tests = []struct {
		text                   string
		sY, sM, sD, eY, eM, eD int
	}{
		{"1 January 1900 - 31 December 2000", 1900, 1, 1, 2000, 12, 31},
		{"9 December 2020 – 28 February 2021", 2020, 12, 9, 2021, 2, 28},
		{" 1 January 1900 - 31 December 2000 ", 1900, 1, 1, 2000, 12, 31},
		{" 1 January 1900\n -\n 31 December 2000 ", 1900, 1, 1, 2000, 12, 31},
		{"1January1900-31December2000", 1900, 1, 1, 2000, 12, 31},
		{"1 Jan 1900 - 31 Dec 2000", 1900, 1, 1, 2000, 12, 31},
		{"1 January - 31 December 2000", 2000, 1, 1, 2000, 12, 31},
		{"1 Jan - 31 Dec 2000", 2000, 1, 1, 2000, 12, 31},
		{" 1  Jan  -  31  Dec  2000 ", 2000, 1, 1, 2000, 12, 31},
		{"1Jan-31Dec2000", 2000, 1, 1, 2000, 12, 31},
		{"1 - 31 December 2000", 2000, 12, 1, 2000, 12, 31},
		{" 1  -  31  December  2000 ", 2000, 12, 1, 2000, 12, 31},
		{"1-31December2000", 2000, 12, 1, 2000, 12, 31},
	}

	dateFormat := "2006 January 02"
	for _, tt := range tests {
		testname := fmt.Sprintf("%s", tt.text)

		t.Run(testname, func(t *testing.T) {
			ans, err := p.Parse(tt.text)
			want := daterange.DateRange{
				time.Date(tt.sY, time.Month(tt.sM), tt.sD, 0, 0, 0, 0, time.UTC),
				time.Date(tt.eY, time.Month(tt.eM), tt.eD, 23, 59, 59, 999999999, time.UTC),
			}

			if err != nil {
				t.Error(err)
			} else {
				if ans.Start != want.Start {
					t.Errorf(
						"Start Date: got %s, want %s",
						ans.Start.Format(dateFormat),
						want.Start.Format(dateFormat),
					)
				}

				if ans.End != want.End {
					t.Errorf(
						"End Date: got %s, want %s",
						ans.End.Format(dateFormat),
						want.End.Format(dateFormat),
					)
				}
			}
		})
	}
}
