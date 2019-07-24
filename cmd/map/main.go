package main

import (
	"flag"
	"fmt"
	"github.com/whosonfirst/go-http-nextzenjs"
	"log"
	"net/http"
)

func MapHandler() http.Handler {

	index := `
<!doctype html>
<html lang="en-us">
  <head>
    <meta charset="utf-8">
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=1.0, user-scalable=no">
    <title>Map</title>
    <style>
        body {
            margin: 0px;
            border: 0px;
            padding: 0px;
        }

        #map {
            height: 100%;
            width: 100%;
            position: absolute;
        }

    </style>
  </head>

  <body>
    <div id="map"></div>

    <script>

    	var api_key = document.body.getAttribute("data-nextzen-api-key");

        var map = L.Nextzen.map('map', {apiKey: api_key,attribution: '<a href="https://github.com/tangrams" target="_blank">Tangram</a> | <a href="http://www.openstreetmap.org/copyright" target="_blank">&copy; OpenStreetMap contributors</a> | <a href="https://www.nextzen.org/" target="_blank">Nextzen</a>',
            tangramOptions: {
                scene: {
                    import: [
                        '/tangram/refill-style.zip',
                    ],
                    sources: {
                        mapzen: {
                            url: 'https://{s}.tile.nextzen.org/tilezen/vector/v1/512/all/{z}/{x}/{y}.mvt',
                            url_subdomains: ['a', 'b', 'c', 'd'],
                            url_params: {api_key: api_key},
                            tile_size: 512,
                            max_zoom: 16
                        }
                    }
                }
            }
        });
        map.setView([33.0, -12.3], 2);
        L.Nextzen.hash({map: map});

    </script>

  </body>
</html>`

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		rsp.Write([]byte(index))
	}

	return http.HandlerFunc(fn)
}

func main() {

	api_key := flag.String("api-key", "", "...")
	host := flag.String("host", "localhost", "...")
	port := flag.Int("port", 8080, "...")

	flag.Parse()

	opts := nextzenjs.DefaultNextzenJSOptions()
	opts.APIKey = *api_key

	mux := http.NewServeMux()

	map_handler := MapHandler()

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
