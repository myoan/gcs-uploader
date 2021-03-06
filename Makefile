GOFLAGS :=
GOVERSION=$(shell go version)
GOOS=$(word 1,$(subst /, ,$(lastword $(GOVERSION))))
GOARCH=$(word 2,$(subst /, ,$(lastword $(GOVERSION))))
BUILD_DIR=build/$(GOOS)-$(GOARCH)

.PHONY: all build clean deps package package-zip package-targz

all: build

build: deps
	mkdir -p $(BUILD_DIR)
	cd ${BUILD_DIR} && GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(GOFLAGS) -o gcs-uploader ../../*.go

clean:
	rm -rf build package

package:
	$(MAKE) package-targz GOOS=linux GOARCH=amd64
	$(MAKE) package-targz GOOS=linux GOARCH=386
	$(MAKE) package-targz GOOS=linux GOARCH=arm64
	$(MAKE) package-targz GOOS=linux GOARCH=arm
	$(MAKE) package-zip GOOS=darwin GOARCH=amd64
	$(MAKE) package-zip GOOS=windows GOARCH=amd64
	$(MAKE) package-zip GOOS=windows GOARCH=386

package-zip: build
	mkdir -p package
	cd $(BUILD_DIR) && zip ../../package/gcs-uploader_$(GOOS)_$(GOARCH).zip gcs-uploader

package-targz: build
	mkdir -p package
	cd $(BUILD_DIR) && tar zcvf ../../package/gcs-uploader_$(GOOS)_$(GOARCH).tar.gz gcs-uploader