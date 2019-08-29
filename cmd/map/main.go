package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/whosonfirst/go-http-nextzenjs"
	"github.com/whosonfirst/go-http-nextzenjs/assets/templates"
	"html/template"
	"log"
	"net/http"
)

func MapHandler(templates *template.Template) (http.Handler, error) {

	t := templates.Lookup("map")

	if t == nil {
		return nil, errors.New("Missing 'map' template")
	}

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		err := t.Execute(rsp, nil)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusInternalServerError)
		}

		return
	}

	return http.HandlerFunc(fn), nil
}

func main() {

	api_key := flag.String("api-key", "", "...")
	host := flag.String("host", "localhost", "...")
	port := flag.Int("port", 8080, "...")

	path_templates := flag.String("templates", "", "An optional string for local templates. This is anything that can be read by the 'templates.ParseGlob' method.")

	flag.Parse()

	t := template.New("example")

	var err error

	if *path_templates != "" {

		t, err = t.ParseGlob(*path_templates)

		if err != nil {
			log.Fatal(err)
		}

	} else {

		for _, name := range templates.AssetNames() {

			body, err := templates.Asset(name)

			if err != nil {
				log.Fatal(err)
			}

			t, err = t.Parse(string(body))

			if err != nil {
				log.Fatal(err)
			}
		}
	}

	opts := nextzenjs.DefaultNextzenJSOptions()
	opts.APIKey = *api_key

	mux := http.NewServeMux()

	map_handler, err := MapHandler(t)

	if err != nil {
		log.Fatal(err)
	}

	nextzenjs_handler, err := nextzenjs.NextzenJSHandler(map_handler, opts)

	if err != nil {
		log.Fatal(err)
	}

	mux.Handle("/", nextzenjs_handler)

	err = nextzenjs.AppendAssetHandlers(mux)

	if err != nil {
		log.Fatal(err)
	}

	endpoint := fmt.Sprintf("%s:%d", *host, *port)
	log.Printf("Listening for requests on %s\n", endpoint)

	err = http.ListenAndServe(endpoint, mux)

	if err != nil {
		log.Fatal(err)
	}
}
