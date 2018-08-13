package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
)

const (
	envVanityVCS    string = "VANITY_VCS"
	envVanityVCSURL string = "VANITY_VCS_URL"
	protoHeader     string = "X-Forwarded-Proto"
)

var tmpl = template.Must(template.New("html").Parse(`<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>
<meta name="go-import" content="{{.Host}} {{.VCS}} {{.VCSURL}}">
</head>
</html>
`))

type tmplData struct {
	Host   string
	VCS    string
	VCSURL string
}

func main() {
	vcs := os.Getenv(envVanityVCS)
	if vcs == "" {
		vcs = "git"
	}
	vcsURL := os.Getenv(envVanityVCSURL)
	if vcsURL == "" {
		log.Fatalf("%s must be set, e.g. https://github.com/username", envVanityVCSURL)
	}
	u, err := url.Parse(vcsURL)
	if err != nil {
		log.Fatalf("error building VCS URL: %v", err)
	}
	if u.Scheme != "https" {
		log.Fatalf("%s scheme must be HTTPS", envVanityVCSURL)
	}

	http.HandleFunc("/", handler(vcs, u))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handler(vcs string, vcsURL *url.URL) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Scheme != "https" {
			proto := r.Header.Get(protoHeader)
			if proto == "" || proto != "https" {
				http.Redirect(w, r, fmt.Sprintf("https://%s%s", r.Host, r.URL.RequestURI()), http.StatusMovedPermanently)
				return
			}
		}

		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		u, err := url.Parse(fmt.Sprintf("https://%s%s", vcsURL.Host, path.Join(vcsURL.Path, r.URL.Path)))
		if err != nil {
			http.Error(w, fmt.Sprintf("error building VCS URL: %v", err), http.StatusInternalServerError)
			return
		}

		if r.URL.Query().Get("go-get") != "1" || len(r.URL.Path) < 2 {
			http.Redirect(w, r, u.String(), http.StatusTemporaryRedirect)
			return
		}

		data := &tmplData{
			Host:   path.Join(r.Host, r.URL.Path),
			VCS:    vcs,
			VCSURL: u.String(),
		}

		var buf bytes.Buffer
		err = tmpl.Execute(&buf, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Cache-Control", "no-store")
		w.Write(buf.Bytes())
	}
}
