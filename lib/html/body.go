package html

// To learn more visit generator.go
//go:generate go run generator.go

import (
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
	"tldr-translation-progress/lib/tldr"
)

const DefaultFileMask = 0740

// If you create new files, which contain html templates, don't forget to add them to the file
// ../../resources/tailwind.config.js, otherwise used css classes may not be available in the purged css file.
const styleFilename = "style.css"
const indexFilename = "index.html"
const htmlSite string = `
{{- define "site" -}}
<!DOCTYPE HTML>
<html lang="en">
<head>
	<title>tldr translation progress</title>
	<meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>
	<link href="{{ style_name }}" rel="stylesheet">
</head>
<body>
<div class="container mx-auto text-center">
	<h1 class="font-bold text-4xl m-10">tldr translation progress</h1>
	<div class="mb-10 mt-2">
		<h3 class="text-2xl p-5">Legend</h3>
		<table class="border-collapse mx-auto">
		  <tbody>
			<tr>
			  <td class="border-b border-gray-400 px-4 py-2">✔</td>
			  <td class="border-b border-gray-400 px-4 py-2">translated & same number of entries as the english version </td>
			</tr>
			<tr>
			  <td class="border-b border-gray-400 px-4 py-2">⚠</td>
			  <td class="border-b border-gray-400 px-4 py-2">not up-to-date (different number of entries than the english version)</td>
			</tr>
			<tr>
			  <td class="px-4 py-2">✖</td>
			  <td class="px-4 py-2">not translated</td>
			</tr>
		  </tbody>
		</table>
	</div>
	{{- template "table" . -}}
	<div class="my-6 text-center text-gray-700">
		Thanks for using this site • 
		Generated by <a href="https://github.com/LukWebsForge/TldrProgress">tldr-translation-progress</a> • 
		Last updated {{ current_date_time }}
	</div>
</div>
</body>
</html>
{{- end -}}
`

// Generates a html file path/index.html, which shows the progress of translating the tldr pages.
// This information nis provided by the index.
// A css file path/style.css used for styling the website also will be copied.
func GenerateHtml(index *tldr.Index, path string) error {
	// Adding custom function
	funcs := make(template.FuncMap)
	funcs["status2html"] = statusToHtml
	funcs["print_percentage"] = printPercentage
	funcs["style_name"] = func() string {
		return styleFilename
	}
	funcs["current_date_time"] = func() string {
		return time.Now().Format(time.RFC850)
	}

	tmpl := template.New("page")
	tmpl.Funcs(funcs)

	// Parsing the templates
	_, err := tmpl.Parse(htmlTable)
	if err != nil {
		return err
	}

	_, err = tmpl.Parse(htmlSite)
	if err != nil {
		return err
	}

	// Opening the file
	err = os.MkdirAll(path, DefaultFileMask)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(filepath.Join(path, indexFilename), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, DefaultFileMask)
	if err != nil {
		return err
	}
	defer file.Close()

	// Executing the template
	err = tmpl.ExecuteTemplate(file, "site", index)
	if err != nil {
		return err
	}

	err = copyStyle(filepath.Join(path, styleFilename))
	if err != nil {
		return err
	}

	return nil
}

// Copies the style file, which is stored internally, to an external file at the given path
func copyStyle(path string) error {
	return ioutil.WriteFile(path, []byte(styleFromAssets()), DefaultFileMask)
}
