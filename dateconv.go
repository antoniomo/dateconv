package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/kierdavis/dateparser"
)

func main() {
	var (
		dayFirst  = flag.Bool("df", false, "Disambiguate day-month as day first (like in the US)")
		yearFirst = flag.Bool("yf", false, "Disambiguate 2-digit yeas (yy-mm-dd/dd-mm-yy) as years first")
		fuzzy     = flag.Bool("fuzzy", false, "Allow more resilient fuzzy matching")
		format    = flag.String("format", time.RFC850, "Output format, Golang output format rules")
		config    = flag.String("conf", "~/.dateconv", "Extra config/timezone json mappings")
		// In order of precedence
		toUTC    = flag.Bool("utc", false, "Converts date to UTC format")
		toTs     = flag.Bool("ts", false, "Converts date to Unix timestamp")
		toTsNano = flag.Bool("tsNano", false, "Converts date to Unix timestamp with nanoseconds")
		toTz     = flag.String("tz", "local", "To specific Timezone. Defaults to local.")
	)
	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}
	toParse := flag.Args()[0]

	// Attempt to load config, or fail silently
	var tzInfos map[string]int
	conf, err := ioutil.ReadFile(*config)
	if err == nil {
		err = json.Unmarshal(conf, &tzInfos)
		if err != nil {
			tzInfos = nil
		}
	}

	parser := &dateparser.Parser{
		DayFirst:  *dayFirst,
		YearFirst: *yearFirst,
		Fuzzy:     *fuzzy,
		TZInfos:   tzInfos,
	}
	dt, err := parser.Parse(toParse)
	if err != nil {
		log.Fatal(err)
	}
	// Now just print the output in the date on the desired format
	switch {
	case *toUTC:
		fmt.Println(dt.UTC().Format(*format))
	case *toTs:
		fmt.Println(dt.Unix())
	case *toTsNano:
		fmt.Println(dt.UnixNano())
	default: // To timezone
		if strings.ToLower(*toTz) == "local" {
			fmt.Println(dt.Local().Format(*format))
		} else {
			loc, err := time.LoadLocation(*toTz)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(dt.In(loc).Format(*format))
		}
	}
}
