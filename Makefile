TARGET = ddns

# Dependencies will be installed here
#GOPATH = ${CURDIR}/env
# This is where go will be installed.
# Change this if you want to use a preexisting install.
GOROOT = ${CURDIR}/go
# The version of go to install
GOTAG = release

#GO = $(GOROOT)/bin/go

all: build

build: build_darwin_amd64 build_darwin_386 \
	build_linux_amd64 build_linux_386 build_linux_arm \
	build_freebsd_amd64 build_freebsd_386 build_freebsd_arm \
	build_windows_amd64 build_windows_386

build_darwin_%: GOOS = darwin
build_linux_%: GOOS = linux
build_freebsd_%: GOOS = freebsd
build_windows_%: GOOS = windows
build_windows_%: EXT = .exe

build_%_amd64: GOARCH = amd64
build_%_386: GOARCH = 386
build_%_arm: GOARCH = arm

build_%:
	export GOOS=$(GOOS); export GOARCH=$(GOARCH); \
	go build -o build/$(TARGET).$(GOOS)_$(GOARCH)$(EXT) .

clean_all: clean
	rm -rf $(GOROOT)
	rm -rf $(GOPATH)

clean:
	rm -rf build
