package nextzenjs

import (
	"github.com/aaronland/go-http-rewrite"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"io"
	_ "log"
	"net/http"
)

type NextzenJSOptions struct {
	AppendAPIKey bool
	AppendJS     bool
	AppendCSS    bool
	APIKey       string
	JS           []string
	CSS          []string
}

func DefaultNextzenJSOptions() *NextzenJSOptions {

	opts := NextzenJSOptions{
		AppendAPIKey: true,
		AppendJS:     true,
		AppendCSS:    true,
		APIKey:       "nextzen-xxxxxx",
		JS:           []string{"/javascript/nextzen.min.js"},
		CSS:          []string{"/css/nextzen.js.css"},
	}

	return &opts
}

func NextzenJSHandler(next http.Handler, opts *NextzenJSOptions) (http.Handler, error) {

	var cb rewrite.RewriteHTMLFunc

	cb = func(n *html.Node, w io.Writer) {

		if n.Type == html.ElementNode && n.Data == "head" {

			if opts.AppendJS {

				for _, js := range opts.JS {

					script_type := html.Attribute{"", "type", "text/javascript"}
					script_src := html.Attribute{"", "src", js}

					script := html.Node{
						Type:      html.ElementNode,
						DataAtom:  atom.Script,
						Data:      "script",
						Namespace: "",
						Attr:      []html.Attribute{script_type, script_src},
					}

					n.AppendChild(&script)
				}

			}

			if opts.AppendCSS {

				for _, css := range opts.CSS {
					link_type := html.Attribute{"", "type", "text/css"}
					link_rel := html.Attribute{"", "rel", "stylesheet"}
					link_href := html.Attribute{"", "href", css}

					link := html.Node{
						Type:      html.ElementNode,
						DataAtom:  atom.Link,
						Data:      "link",
						Namespace: "",
						Attr:      []html.Attribute{link_type, link_rel, link_href},
					}

					n.AppendChild(&link)
				}
			}
		}

		if n.Type == html.ElementNode && n.Data == "body" {

			if opts.AppendAPIKey {
				api_key_ns := ""
				api_key_key := "data-nextzen-api-key"
				api_key_value := opts.APIKey

				api_key_attr := html.Attribute{api_key_ns, api_key_key, api_key_value}
				n.Attr = append(n.Attr, api_key_attr)
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			cb(c, w)
		}
	}

	return rewrite.RewriteHTMLHandler(next, cb), nil
}

func NextzenJSAssetsHandler() (http.Handler, error) {

	fs := assetFS()
	return http.FileServer(fs), nil
}

func AppendAssetHandlers(mux *http.ServeMux) error {

	assets_handler, err := NextzenJSAssetsHandler()

	if err != nil {
		return err
	}

	mux.Handle("/javascript/nextzen.js", assets_handler)
	mux.Handle("/javascript/nextzen.min.js", assets_handler)
	mux.Handle("/javascript/tangram.js", assets_handler)
	mux.Handle("/javascript/tangram.min.js", assets_handler)
	mux.Handle("/css/nextzen.js.css", assets_handler)
	mux.Handle("/tangram/refill-style.zip", assets_handler)

	return nil
}
