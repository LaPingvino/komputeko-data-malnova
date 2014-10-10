package main

import (
	"encoding/json"
	"flag"
	"fmt"
	data "github.com/komputeko/komputeko-data"
	"html/template"
	"os"
	"strings"
)

type langWord [2]string

var tmpl *template.Template

func writeHTML(ctx data.Entry, dirname string, references map[langWord]string) error {
	filename := template.URLQueryEscaper(ctx.Translations[0].Language + "_" + ctx.Translations[0].Words[0].Written + ".html")

	file, err := os.Create(filename)
	defer file.Close()
	if err != nil {
		return err
	}

	showme := struct {
		Title string
		Body  template.HTML
	}{
		ctx.Translations[0].Words[0].Written,
		template.HTML(""),
	}

	for _, translation := range ctx.Translations {
		showme.Body += template.HTML("<table class='search' style='border-collapse: collapse; border-width: 2px; display: inline-block;' cellpadding=2><tr style=''><th style='width: 12em;'>" + translation.Language + "</th></tr><tr><td><span style='background: #ffffbb;'>\n")
		for _, word := range translation.Words {
			showme.Body += template.HTML("<p>" +
				template.HTMLEscapeString(word.Written+" ("+
					strings.Join(word.Sources, ",")+")") +
				"</p>")
			references[langWord{translation.Language, word.Written}] = filename
		}
		showme.Body += template.HTML("</span></td></tr></table>")
	}

	err = tmpl.Execute(file, showme)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	references := make(map[langWord]string)
	var terminaro data.Terminaro
	flag.Parse()
	filename := flag.Arg(0)
	dirname := flag.Arg(1)

	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		panic(err.Error())
	}

	err = json.NewDecoder(file).Decode(&terminaro)
	if err != nil {
		panic(err.Error())
	}

	tmpl, err = template.ParseFiles("page.tmpl")
	if err != nil {
		panic(err.Error())
	}

	for _, entry := range terminaro {
		err := writeHTML(entry, dirname, references)
		if err != nil {
			fmt.Println(err, "\nEntry with error:\n", entry)
		}
	}

}
