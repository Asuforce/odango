export GO111MODULE=on

PKG = github.com/Asuforce/odango
COMMIT = $$(git describe --tags --always)
DATE = $$(date -u '+%Y-%m-%d_%H:%M:%S')
BUILD_LDFLAGS = -X $(PKG).commit=$(COMMIT) -X $(PKG).date=$(DATE)
RELEASE_BUILD_LDFLAGS = -s -w $(BUILD_LDFLAGS)

.PHONY: build
build:
	go build -ldflags="$(BUILD_LDFLAGS)"

.PHONY: test
test:
	go test -v ./...

.PHONY: devel-deps
devel-deps:
	GO111MODULE=off go get -v \
	github.com/motemen/gobump/cmd/gobump \
	github.com/Songmu/ghch/cmd/ghch \
	github.com/Songmu/goxz/cmd/goxz \
	github.com/tcnksm/ghr

.PHONY: crossbuild
crossbuild: devel-deps
	$(eval ver = $(shell gobump show -r))
	goxz -pv=v$(ver) -os=linux,darwin,windows -arch=amd64 -build-ldflags="$(RELEASE_BUILD_LDFLAGS)" -d=./dist/v$(ver)

.PHONY: release
release: devel-deps
	_tools/release.sh
	_tools/upload_artifacts

.PHONY: lint
lint:
	go vet ./...
	golint -set_exit_status ./...
