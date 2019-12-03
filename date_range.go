package daterange

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// DateRange ...
type DateRange struct {
	Start time.Time
	End   time.Time
}

// Constructor ...
type Constructor struct {
	Pattern string
	Regexp  *regexp.Regexp
	Handler func([]string) *DateRange
}

// Parser ...
type Parser struct {
	constructors []*Constructor
	monthRegexps []*regexp.Regexp
}

func NewParser() *Parser {
	parser := Parser{}

	monthPatterns := []string{
		"jan(?:uary)?",
		"feb(?:ruary)?",
		"mar(?:ch)?",
		"apr(?:il)?",
		"may",
		"jun(?:e)?",
		"jul(?:y)?",
		"aug(?:ust)?",
		"sep(?:tember)?",
		"oct(?:ober)?",
		"nov(?:ember)?",
		"dec(?:ember)?",
	}

	monthRegexps := make([]*regexp.Regexp, len(monthPatterns))

	for i, p := range monthPatterns {
		monthRegexps[i] = regexp.MustCompile("(?i)" + p)
	}

	parser.monthRegexps = monthRegexps

	monthAnyPattern := strings.Join(monthPatterns, "|")

	parser.constructors = []*Constructor{
		&Constructor{
			Pattern: fmt.Sprintf("(?i)(\\d{1,2})\\s*(%[1]s)\\s*(\\d{4})\\s*[–-]\\s*(\\d{1,2})\\s*(%[1]s)\\s*(\\d{4})", monthAnyPattern),
			Handler: func(m []string) *DateRange {
				return &DateRange{Start: parser.startDate(m[3], m[2], m[1]), End: parser.endDate(m[6], m[5], m[4])}
			},
		},
		&Constructor{
			Pattern: fmt.Sprintf("(?i)(\\d{1,2})\\s*(%[1]s)\\s*[–-]\\s*(\\d{1,2})\\s*(%[1]s)\\s*(\\d{4})", monthAnyPattern),
			Handler: func(m []string) *DateRange {
				return &DateRange{Start: parser.startDate(m[5], m[2], m[1]), End: parser.endDate(m[5], m[4], m[3])}
			},
		},
		&Constructor{
			Pattern: fmt.Sprintf("(?i)(\\d{1,2})\\s*[–-]\\s*(\\d{1,2})\\s*(%[1]s)\\s*(\\d{4})", monthAnyPattern),
			Handler: func(m []string) *DateRange {
				return &DateRange{Start: parser.startDate(m[4], m[3], m[1]), End: parser.endDate(m[4], m[3], m[2])}
			},
		},
	}

	for _, c := range parser.constructors {
		c.Regexp = regexp.MustCompile(c.Pattern)
	}

	return &parser
}

// Parse returns a date range from a given string
func (p Parser) Parse(text string) (*DateRange, error) {
	var m []string

	for _, c := range p.constructors {
		m = c.Regexp.FindStringSubmatch(text)

		if len(m) > 0 {
			return c.Handler(m), nil
		}
	}

	return nil, errors.New("could not parse date range")
}

func (p Parser) startDate(y, m, d string) time.Time {
	return time.Date(p.parseYear(y), p.parseMonth(m), p.parseDay(d), 0, 0, 0, 0, time.UTC)
}

func (p Parser) endDate(y, m, d string) time.Time {
	return time.Date(p.parseYear(y), p.parseMonth(m), p.parseDay(d), 23, 59, 59, 999999999, time.UTC)
}

func (p Parser) parseDay(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Fatal(err)
	}

	return i
}

func (p Parser) parseMonth(s string) time.Month {
	for i, re := range p.monthRegexps {
		m := re.MatchString(s)
		if m {
			return time.Month(i + 1)
		}
	}

	return 0
}

func (p Parser) parseYear(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Fatal(err)
	}

	return i
}
