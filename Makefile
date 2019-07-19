build:
	@make maps
	@make assets

maps:
	@make wwwdirs
	@make nextzenjs
	@make tangram
	@make styles

assets:	self
	if test ! -d bin; then mkdir bin; fi
	@GOPATH=$(GOPATH) go build -o bin/go-bindata ./vendor/github.com/zendesk/go-bindata/go-bindata/
	@GOPATH=$(GOPATH) go build -o bin/go-bindata-assetfs vendor/github.com/elazarl/go-bindata-assetfs/go-bindata-assetfs/main.go
	rm -f static/*~ static/css/*~ static/javascript/*~ static/tangram/*~
	@PATH=$(PATH):$(CWD)/bin bin/go-bindata-assetfs -pkg mapzenjs www static/javascript static/css static/tangram
	mv bindata.go assets.go

wwwdirs:
	if test ! -d static/javascript; then mkdir static/javascript; fi
	if test ! -d static/css; then mkdir static/css; fi
	if test ! -d static/tangram; then mkdir static/tangram; fi

tangram: 
	curl -s -o static/javascript/tangram.js https://www.nextzen.org/tangram/tangram.debug.js
	curl -s -o static/javascript/tangram.min.js https://www.nextzen.org/tangram/tangram.min.js

styles: refill walkabout

refill:
	curl -s -o static/tangram/refill-style.zip https://www.nextzen.org/carto/refill-style/refill-style.zip
	curl -s -o static/tangram/refill-style-themes-label.zip https://www.nextzen.org/carto/refill-style/themes/label-10.zip

walkabout:
	curl -s -o static/tangram/walkabout-style.zip https://www.nextzen.org/carto/refill-style/walkabout-style.zip

nextzenjs:
	curl -s -o static/css/nextzen.js.css https://www.nextzen.org/js/nextzen.css
	curl -s -o static/javascript/nextzen.js https://www.nextzen.org/js/nextzen.js
	curl -s -o static/javascript/nextzen.min.js https://www.nextzen.org/js/nextzen.min.js

# http://leafletjs.com/download.html

leaflet:
	curl -s -o static/css/leaflet.css https://unpkg.com/leaflet@1.5.1/dist/leaflet.css
	curl -s -o static/javascript/leaflet.js https://unpkg.com/leaflet@1.5.1/dist/leaflet.js 
	@make sritest FILE=static/css/leaflet.css EXPECTED=xwE/Az9zrjBIphAcBb3F6JVqxf46+CDLwfLMHloNu6KEQCAWi6HcDUbeOfBIptF7tcCzusKFjFw2yuvEpDL9wQ==
	@make sritest FILE=static/javascript/leaflet.js EXPECTED=GffPMF3RvMeYyc1LWMHtK8EbPv0iNZ8/oTtHPx9/cc2ILxQ+u905qIwdpULaqDkyBKgOaB57QTMg7ztg8Jm2Og==

# https://developer.mozilla.org/en-US/docs/Web/Security/Subresource_Integrity
# https://github.com/chrisgreg/sri-gen/blob/master/sri-gen.go

sritest:
	$(eval SRIHASH := $(shell shasum -b -a 512 $(FILE)  | xxd -r -p | base64))
	if test "$(SRIHASH)" != "$(EXPECTED)"; then echo "$(FILE) has unexpected SRI hash"; exit 1; fi
	@echo "$(FILE) has expected SRI hash"