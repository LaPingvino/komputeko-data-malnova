package main

import (
	"encoding/json"
	"flag"
	"fmt"
	data "github.com/komputeko/komputeko-data"
	"html/template"
	"os"
	"sort"
	"strings"
	"unicode/utf8"
)

type reference [][3]string
type runeSlice []rune

var references reference
var alphabet runeSlice
var langs []string

func (r reference) Len() int {
	return len(r)
}

func (r reference) Less(i, j int) bool {
	return strings.ToLower(r[i][1]) < strings.ToLower(r[j][1])
}

func (r reference) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func (r runeSlice) Len() int {
	return len(r)
}

func (r runeSlice) Less(i, j int) bool {
	return r[i] < r[j]
}

func (r runeSlice) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

var tmpl *template.Template

func writeHTML(ctx data.Entry, dirname string) error {
	filename := strings.Replace(template.URLQueryEscaper(ctx.Translations[0].Language+"_"+ctx.Translations[0].Words[0].Written+".html"), "%", "X", -1)

	file, err := os.Create(filename)
	defer file.Close()
	if err != nil {
		return err
	}

	showme := struct {
		Title string
		Body  template.HTML
	}{
		ctx.Translations[0].Words[0].Written + ": " + ctx.Wordtype,
		template.HTML(""),
	}

	for _, translation := range ctx.Translations {
		showme.Body += template.HTML("<div class=\"word\"><div id=\"" + translation.Language + "\" class=\"wordheader\">" + translation.Language + "</div>\n<div class=\"definitions\">")
		for _, word := range translation.Words {
			showme.Body += template.HTML("<div class=\"definition\">" +
				template.HTMLEscapeString(word.Written))
			for _, sw := range word.Sources {
				showme.Body += template.HTML("<span class=\"source\">" +
					template.HTMLEscapeString(sw) + "</span>")
			}

			showme.Body += template.HTML("</div>")
			references = append(references, [3]string{translation.Language, word.Written, filename})
			found := false
			firstletter, _ := utf8.DecodeRuneInString(strings.ToLower(word.Written))
			for _, letter := range alphabet {
				if firstletter == letter {
					found = true
				}
			}
			if !found {
				alphabet = append(alphabet, firstletter)
			}
		}
		showme.Body += template.HTML("</div></div>")
	}

	err = tmpl.Execute(file, showme)
	if err != nil {
		return err
	}

	return nil
}

func writeLanguageIndex(dirname string, lang string, letter rune) error {
	showme := struct {
		Title string
		Body  template.HTML
	}{
		"Index " + lang + ":" + fmt.Sprint(letter),
		template.HTML(""),
	}

	file, err := os.Create("index_" + lang + "_" + fmt.Sprint(letter) + ".html")
	defer file.Close()
	if err != nil {
		return err
	}

	showme.Body += template.HTML("<ul id=\"results\">")
	for _, entry := range references {
		firstletter, _ := utf8.DecodeRuneInString(strings.ToLower(entry[1]))
		if entry[0] == lang && firstletter == letter {
			showme.Body += template.HTML("<li><a href=\"" +
				entry[2] + "#" + entry[0] + "\">" + entry[1] + "</a></li>")
		}
	}
	showme.Body += template.HTML("</ul>")

	err = tmpl.Execute(file, showme)
	if err != nil {
		return err
	}

	return nil
}

func writeIndexHtml() error {
	showme := struct {
		Title string
		Body  template.HTML
	}{
		"Index",
		template.HTML(""),
	}

	file, err := os.Create("index.html")
	defer file.Close()
	if err != nil {
		return err
	}

	showme.Body += template.HTML("<div id=\"languages\">\n")
	for _, lang := range langs {
		showme.Body += template.HTML("<div id=\"" + lang + "\"><span class=\"language\">" +
			lang + ": </span>\n")
		for _, letter := range alphabet {
			showme.Body += template.HTML("<a href=\"index_" + lang + "_" +
				fmt.Sprint(letter) + ".html\">" + string(letter) + "</a>\n")
		}
		showme.Body += template.HTML("</div>")
	}
	showme.Body += template.HTML("</div>")

	err = tmpl.Execute(file, showme)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	var tname string
	var dirname string
	var terminaro data.Terminaro
	flag.StringVar(&tname, "t", "page.tmpl", "The template being used")
	flag.StringVar(&dirname, "d", "", "The output directory")

	flag.Parse()
	filename := flag.Arg(0)

	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		panic(err.Error())
	}

	err = json.NewDecoder(file).Decode(&terminaro)
	if err != nil {
		panic(err.Error())
	}

	tmpl, err = template.ParseFiles(tname)
	if err != nil {
		panic(err.Error())
	}

	for _, entry := range terminaro[1:] {
		err := writeHTML(entry, dirname)
		if err != nil {
			fmt.Println(err, "\nEntry with error:\n", entry)
		}
	}

	sort.Sort(references)

	var oldpart data.Translation
	for _, part := range terminaro[0].Translations {
		if oldpart.Language != part.Language {
			langs = append(langs, part.Language)
		}
		oldpart = part
	}
	sort.Sort(sort.StringSlice(langs))
	sort.Sort(alphabet)

	for _, lang := range langs {
		for _, letter := range alphabet {
			err := writeLanguageIndex(dirname, lang, letter)
			if err != nil {
				fmt.Println(err, "\nLanguage with error:\n", lang)
			}
		}
	}

	err = writeIndexHtml()
	if err != nil {
		fmt.Println(err, " while writing index.html")
	}
}
