all: build
build: 
	CGO_ENABLED=0 go build 
deps:
	./install_deps.sh
clean:
	rm -rf venv
	rm *.pyc
	rm -rf vendor 
	rm metrics_poller
test:
	./metrics_poller -inputFile=servers.txt
webserver:
	./run_webserver.sh

