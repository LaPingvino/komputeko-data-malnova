package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	data "github.com/komputeko/komputeko-data"
	"os"
	"regexp"
	"strings"
)

func ridx(s string) string {
	r := regexp.MustCompile("(.*) /.*").FindStringSubmatch(s)
	if len(r) > 1 {
		return r[1]
	} else {
		return s
	}
}

func konverti(de string) (entry data.Entry, err error) {
	sourceline := strings.Split(de, "\t")
	if len(sourceline) != 36 {
		err = fmt.Errorf("Input file not valid")
		return
	}

	// Getting the word type
	if sourceline[2] != "" {
		entry.Wordtype = sourceline[2]
	} else {
		// Get word type from English definition
		match := strings.TrimSpace(regexp.MustCompile(`\([a-z]+\.?\)\ *$`).FindString(sourceline[3]))
		switch match {
		case "(subst.)", "(verbo)", "(mallon.)", "(adj.)":
			entry.Wordtype = strings.Trim(match, "()")
			sourceline[3] = strings.TrimSuffix(strings.TrimSpace(sourceline[3]), match)
		}
	}

	// Getting the English definitions
	if sourceline[3] != "" {
		word := data.Word{
			Written: strings.TrimSpace(regexp.MustCompile("^[^(]+").FindString(sourceline[3])),
		}
		words := []data.Word{word}
		if sourceline[4] != "" {
			words = append(words, data.Word{Written: sourceline[4]})
		} else {
			match := regexp.MustCompile(`\([^)]+\)`).FindString(sourceline[3])
			if match != "" {
				matchsplit := strings.Split(match, ",")
				for i := 0; i < len(matchsplit); i++ {
					words = append(words, data.Word{Written: strings.Trim(matchsplit[i], "() ")})
				}
			}
		}
		english := data.Translation{"en", words}
		entry.Translations = append(entry.Translations, english)
	}

	// Getting the Esperanto definitions
	if sourceline[5] != "" {
		var sources []string
		var sources2 []string
		var sources3 []string
		for _, source := range sourceline[6:10] {
			if source != "" {
				sources = append(sources, source)
			}
		}
		word := data.Word{
			Written: ridx(sourceline[5]),
			Sources: sources,
		}
		words := []data.Word{word}
		if sourceline[10] != "" {
			for _, source := range sourceline[11:13] {
				if source != "" {
					sources2 = append(sources2, source)
				}
			}
			words = append(words, data.Word{Written: ridx(sourceline[10]), Sources: sources2})
		}
		if sourceline[13] != "" {
			if sourceline[14] != "" {
				sources3 = []string{sourceline[14]}
			}
			words = append(words, data.Word{Written: ridx(sourceline[13]), Sources: sources3})
		}
		esperanto := data.Translation{"eo", words}
		entry.Translations = append(entry.Translations, esperanto)
	}

	// Getting the Dutch definitions
	if sourceline[16] != "" {
		var extrainfo string
		var sourcesnl []string
		var sourcesnl2 []string
		for _, source := range sourceline[18:21] {
			if source != "" {
				sourcesnl = append(sourcesnl, source)
			}
		}
		if sourceline[15] != "" {
			extrainfo = ", " + sourceline[15]
		}
		word := data.Word{
			Written: sourceline[16] + extrainfo,
			Sources: sourcesnl,
		}
		words := []data.Word{word}
		if sourceline[21] != "" {
			if sourceline[22] != "" {
				sourcesnl2 = []string{sourceline[22]}
			}
			words = append(words, data.Word{Written: sourceline[21], Sources: sourcesnl2})
		}
		dutch := data.Translation{"nl", words}
		entry.Translations = append(entry.Translations, dutch)
	}

	// Getting the French definitions
	if sourceline[24] != "" {
		var extrainfofr string
		var sourcesfr []string
		var sourcesfr2 []string
		for _, source := range sourceline[26:29] {
			if source != "" {
				sourcesfr = append(sourcesfr, source)
			}
		}
		if sourceline[23] != "" {
			extrainfofr = ", " + sourceline[23]
		}
		word := data.Word{
			Written: sourceline[24] + extrainfofr,
			Sources: sourcesfr,
		}
		words := []data.Word{word}
		if sourceline[29] != "" {
			if sourceline[30] != "" {
				sourcesfr2 = []string{sourceline[30]}
			}
			words = append(words, data.Word{Written: sourceline[29], Sources: sourcesfr2})
		}
		french := data.Translation{"fr", words}
		entry.Translations = append(entry.Translations, french)
	}

	// Getting the German definitions
	if sourceline[32] != "" {
		var extrainfode string
		var sourcesde []string
		var sourcesde2 []string
		if sourceline[31] != "" {
			extrainfode = ", " + sourceline[31]
		}
		if sourceline[33] != "" {
			sourcesde = []string{sourceline[33]}
		}
		word := data.Word{
			Written: sourceline[32] + extrainfode,
			Sources: sourcesde,
		}
		words := []data.Word{word}
		if sourceline[34] != "" {
			if sourceline[35] != "" {
				sourcesde2 = []string{strings.TrimSpace(sourceline[35])}
			}
			words = append(words, data.Word{Written: sourceline[34], Sources: sourcesde2})
		}
		german := data.Translation{"de", words}
		entry.Translations = append(entry.Translations, german)
	}

	return
}

func konvertifluon(fluo *bufio.Reader) (terminaro data.Terminaro, err error) {
	var result data.Entry
	for i, er := fluo.ReadString('\n'); er == nil; i, er = fluo.ReadString('\n') {
		result, err = konverti(i)
		if err != nil {
			return
		}
		terminaro = append(terminaro, result)
	}
	return
}

func main() {
	flag.Parse()
	var fileFrom = flag.Arg(0)

	filefrom, err := os.Open(fileFrom)
	defer filefrom.Close()
	if err != nil {
		panic(err.Error())
	}

	terminaro, err := konvertifluon(bufio.NewReader(filefrom))
	if err != nil {
		panic(err.Error())
	}

	jsonencb, err := json.MarshalIndent(terminaro, "", " ")
	if err != nil {
		panic(err.Error())
	}

	jsonenc := bytes.NewBuffer(jsonencb)
	fmt.Println(jsonenc)
}
