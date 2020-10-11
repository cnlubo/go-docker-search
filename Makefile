BUILD_VERSION   := $(shell cat build_version)
BUILD_DATE      := $(shell date '+%Y-%m-%d %H:%M:%S')
COMMIT_SHA1     := $(shell git rev-parse --short HEAD)

VERSION_PKG     := github.com/cnlubo/go-docker-search
DEST_DIR        := dist
APP             := go-docker-search

all:
	gox -osarch="windows/amd64 darwin/amd64 linux/amd64" \
        -output='${DEST_DIR}/${APP}_{{.OS}}_{{.Arch}}' \
    	-ldflags   "-X '${VERSION_PKG}/version.Version=${BUILD_VERSION}' \
                            -X '${VERSION_PKG}/version.BuildTime=${BUILD_DATE}' \
                            -X '${VERSION_PKG}/version.GitCommit=${COMMIT_SHA1}'" \
                            ./cmd
release: all
	ghr -u cnlubo -t $(GITHUB_RELEASE_TOKEN) -replace -recreate --debug ${BUILD_VERSION} dist

pre-release: all
	ghr -u cnlubo -t $(GITHUB_RELEASE_TOKEN) -replace -recreate -prerelease --debug ${BUILD_VERSION} dist

clean:
	rm -rf ${DEST_DIR}

.PHONY : all release clean install

.EXPORT_ALL_VARIABLES:

GO111MODULE = on
