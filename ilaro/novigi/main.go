package main

import (
	"flag"
	"bufio"
	"bytes"
	"strings"
	"regexp"
	"fmt"
	"os"
	"encoding/json"
)

type Terminaro	[]Entry
type Entry	struct {
	Wordtype	string
	Translations	[]Translation
}
type Translation struct {
	Language	string
	Words	[]Word
}
type Word struct {
	Written	string
	Sources	[]string
//	Frequency float32
}

func konverti(de string) (entry Entry, err error) {
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
		word := Word{
			Written: regexp.MustCompile("^[^( ]+").FindString(sourceline[3]),
		}
		words := []Word{word}
		if sourceline[4] != "" {
			words = append(words, Word{ Written: sourceline[4] })
		} else {
			match := regexp.MustCompile(`\([^)]+\)`).FindString(sourceline[3])
			if match != "" {
				matchsplit := strings.Split(match, ",")
				for i := 0; i < len(matchsplit); i++ {
					words = append(words, Word{ Written: strings.Trim(matchsplit[i], "() ") })
				}
			}
		}
		english := Translation{ "en", words}
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
		word := Word{
			Written: sourceline[5],
			Sources: sources,
		}
		words := []Word{word}
		if sourceline[10] != "" {
			for _, source := range sourceline[11:13] {
				if source != "" {
					sources2 = append(sources2, source)
				}
			}
			words = append(words, Word{ Written: sourceline[10], Sources: sources2, })
		}
		if sourceline[13] != "" {
			if sourceline[14] != "" {
				sources3 = []string{ sourceline[14] }
			}
			words = append(words, Word{ Written: sourceline[13], Sources: sources3, })
		}
		esperanto := Translation{ "eo", words}
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
		if sourceline[15] != "" { extrainfo = ", " + sourceline[15]}
		word := Word{
			Written: sourceline[16] + extrainfo,
			Sources: sourcesnl,
		}
		words := []Word{word}
		if sourceline[21] != "" {
			if sourceline[22] != "" {
				sourcesnl2 = []string{ sourceline[22] }
			}
			words = append(words, Word{ Written: sourceline[21], Sources: sourcesnl2, })
		}
		dutch := Translation{ "nl", words}
		entry.Translations = append(entry.Translations, dutch)
	}

	return
}

func konvertifluon(fluo *bufio.Reader) (terminaro Terminaro, err error) {
	var result Entry
	for i, er := fluo.ReadString('\n'); er == nil; i, er = fluo.ReadString('\n') {
		result, err = konverti(i)
		if err != nil { return }
		terminaro = append(terminaro, result)
	}
	return
}

func main() {
	var fileTo string
	flag.StringVar(&fileTo, "o", "/tmp/output.json", "Write resulting file to this file.")
	flag.Parse()
	var fileFrom = flag.Arg(0)
	
	filefrom, err := os.Open(fileFrom)
	defer filefrom.Close()
	if err != nil { panic(err.Error()) }

	fileto, err := os.Create(fileTo)
	defer fileto.Close()
	if err != nil { panic(err.Error()) }

	terminaro, err := konvertifluon(bufio.NewReader(filefrom))
	if err != nil { panic(err.Error()) }

	jsonencb, err := json.MarshalIndent(terminaro, "", " ")
	if err != nil { panic(err.Error()) }
	
	jsonenc := bytes.NewBuffer(jsonencb)
	fmt.Println(jsonenc)
}
