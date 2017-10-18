ifeq ($(shell uname), Darwin)
	URL := https://github.com/walter-cd/walter/releases/download/v2.0.0/walter_v2.0.0_darwin_amd64.zip
endif
ifeq ($(shell uname), Linux)
	URL := https://github.com/walter-cd/walter/releases/download/v2.0.0/walter_v2.0.0_linux_amd64.zip
endif

pipeline: build/bin/walter
	 if test -z $(VERSION) ; then echo need VERSION environment variable ; exit 1 ; fi
	 build/bin/walter -build -config walter.yaml

build/bin/walter:
	mkdir -p build/bin
	curl -sLo build/bin/walter.zip $(URL)
	unzip build/bin/walter.zip -d build/bin

