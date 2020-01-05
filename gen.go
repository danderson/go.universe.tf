//go:generate go run gen.go

package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	Package("netboot", "danderson/netboot")
	Package("tcpproxy", "google/tcpproxy")
	Package("ppp", "danderson/goppp")
	Package("virtuakube", "danderson/virtuakube")
	Package("natlab", "danderson/natlab")
	Program("conduits", "danderson/conduits")
	Program("metallb", "danderson/metallb")
	must(ioutil.WriteFile("out/_redirects", []byte(strings.Join(redirects, "\n")+"\n"), 0644))
}

func Package(pkg, gh string) {
	must(os.MkdirAll(fmt.Sprintf("out/%s", pkg), 0755))
	writeTemplate(fmt.Sprintf("out/%s/index.html", pkg), packageIndex, map[string]string{
		"pkg": pkg,
		"gh":  gh,
	})
	redirects = append(redirects, fmt.Sprintf("/%s/* /%s/index.html 200", pkg, pkg))
}

func Program(pkg, gh string) {
	must(os.MkdirAll(fmt.Sprintf("out/%s", pkg), 0755))
	writeTemplate(fmt.Sprintf("out/%s/index.html", pkg), packageIndex, map[string]string{
		"pkg": pkg,
		"gh":  gh,
	})
	redirects = append(redirects, fmt.Sprintf("/%s/* /%s/index.html 200", pkg, pkg))
}

var (
	packageIndex = template.Must(template.New("pkg").Parse(`<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>
<meta name="go-import" content="go.universe.tf/{{ .pkg }} git https://github.com/{{ .gh }}">
<meta name="go-source" content="go.universe.tf/{{ .pkg }} https://github.com/{{ .gh }} https://github.com/{{ .gh }}/tree/master{/dir} https://github.com/{{ .gh }}/blob/master{/dir}/{file}#L{line}">
<meta http-equiv="refresh" content="0; url=https://godoc.org/go.universe.tf/{{ .pkg }}">
</head>
<body>
Nothing to see here; <a href="https://godoc.org/go.universe.tf/{{ .pkg }}">move along</a>.
</body>
</html>`))
	redirects []string
)

func init() {
	must(os.RemoveAll("out"))
}

func writeTemplate(path string, tpl *template.Template, vars map[string]string) {
	var b bytes.Buffer
	must(tpl.Execute(&b, vars))
	must(ioutil.WriteFile(path, b.Bytes(), 0644))
}

func must(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func fatal(msg string, args ...interface{}) {
	fmt.Printf(msg+"\n", args...)
	os.Exit(1)
}
