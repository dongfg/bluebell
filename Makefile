##setup: Install wire and packr globally, and then download all dependency
setup:
	@-echo "instal packr/v2"
	@GO111MODULE=off go get -u github.com/gobuffalo/packr/v2
	@-echo "instal wire"
	@GO111MODULE=off go get -u github.com/google/wire
	@-go mod download

##gen: Generate wire and packr2 file
gen:
	@-echo "execute wire"
	@-cd cmd/bluebell && wire 2> /dev/null
	@-echo "execute packr2"
	@-cd cmd/bluebell && packr2 2> /dev/null


##build: Build binary
build: setup gen
	@-echo "build"
	@-go build -a ./cmd/bluebell

##clean: Clean generate files and build files
clean:
	@-rm -f cmd/bluebell/wire_gen.go  2> /dev/null
	@-rm -f cmd/bluebell/main-packr.go  2> /dev/null
	@-rm -rf cmd/bluebell/packrd  2> /dev/null
	@-rm -rf ./bluebell  2> /dev/null
	@-rm -rf dist  2> /dev/null

.PHONY: help
all: help
help: Makefile
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'