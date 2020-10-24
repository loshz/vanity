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
	envVanityVCS    = "VANITY_VCS"
	envVanityVCSURL = "VANITY_VCS_URL"
)

func main() {
	// check for VCS type
	vcs := os.Getenv(envVanityVCS)
	if vcs == "" {
		vcs = "git"
	}

	// check for VCS URL
	vcsURL := os.Getenv(envVanityVCSURL)
	if vcsURL == "" {
		log.Fatalf("%s must be set, e.g. https://github.com/username", envVanityVCSURL)
	}

	// attempt to parse VCS URL
	u, err := url.Parse(vcsURL)
	if err != nil {
		log.Fatalf("invalid VCS URL: %v", err)
	}

	// check for HTTPS
	if u.Scheme != "https" {
		log.Fatalf("%s scheme must be HTTPS", envVanityVCSURL)
	}

	http.HandleFunc("/", handler(vcs, u))
	log.Printf("starting web server on :8080, VCS: %s, VCS_URL: %s", vcs, vcsURL)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

var tmpl = template.Must(template.New("html").Parse(`<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>
<meta name="go-import" content="{{.Host}} {{.VCS}} {{.VCSURL}}">
</head>
</html>
`))

func handler(vcs string, vcsURL *url.URL) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Scheme != "https" {
			proto := r.Header.Get("X-Forwarded-Proto")
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

		data := struct {
			Host   string
			VCS    string
			VCSURL string
		}{
			path.Join(r.Host, r.URL.Path),
			vcs,
			u.String(),
		}

		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, &data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Cache-Control", "no-store")
		w.Write(buf.Bytes())
	}
}
