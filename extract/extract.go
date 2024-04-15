package extract

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"
)

var defaultExtractor *Extractor

func init() {
	defaultExtractor = &Extractor{}
	defaultExtractor.RegisterParser(
		ZapEventParser{},
		RawParser{},
	)
}

func Extract(line string) (Row, error) {
	return defaultExtractor.Extract(line)
}

type Extractor struct {
	parsers []Parser
}

func (ex *Extractor) RegisterParser(parsers ...Parser) {
	ex.parsers = append(ex.parsers, parsers...)
}

func (ex *Extractor) Extract(line string) (Row, error) {
	for _, v := range ex.parsers {
		row := v.Parse(line)
		if row != nil {
			return row, nil
		}
	}
	return nil, fmt.Errorf("no parser match")
}

var _ Parser = (*ZapEventParser)(nil)
var _ Parser = (*RawParser)(nil)
var _ Parser = (*JsonParser)(nil)

type Parser interface {
	Parse(line string) Row // returns nil if it does not match the Parser rule
}

type ZapEventParser struct {
}

var zapEventRegex = regexp.MustCompile(`(?U)^(\S+)\s+(\S+)\s+(\S+)\s+(\S+)\s+(.*)({.*})*$`)

func (z ZapEventParser) Parse(line string) Row {
	r := zapEventRegex.FindAllStringSubmatch(line, 6)
	if len(r) == 0 || len(r[0]) < 7 {
		return nil
	}
	if r[0][1] == "" || r[0][2] == "" {
		return nil
	}
	ts, err := time.Parse("2006-01-02T15:04:05.999Z0700", r[0][1])
	if err != nil {
		return nil
	}
	row := new(ZapEventRow)
	row.Ts = ts
	row.Level = r[0][2]
	row.Module = r[0][3]
	row.Location = r[0][4]
	row.Message = strings.TrimSpace(r[0][5])
	row.Fields = strings.TrimSpace(r[0][6])
	return row
}

type JsonParser struct {
}

func (j JsonParser) Parse(line string) Row {
	return nil
}

type RawParser struct {
}

func (l RawParser) Parse(line string) Row {
	return &Raw{
		message: line,
	}
}

type Row interface {
	Output() string
}

type ZapEventRow struct {
	Ts       time.Time
	Level    string
	Module   string
	Location string
	Message  string
	Fields   string
}

func (z ZapEventRow) Output() string {
	d, _ := json.Marshal(z)
	return string(d)
}

type Raw struct {
	message string
}

func (r Raw) Output() string {
	return r.message
}
