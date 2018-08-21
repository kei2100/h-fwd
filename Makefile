.PHONY: test.nocache test fmt vet lint setup vendor bin version.list version.new release

PACKAGES := $(shell go list ./...)
DIRS := $(shell go list -f '{{.Dir}}' ./...)

setup:
	which dep > /dev/null 2>&1 || go get -u github.com/golang/dep/cmd/dep
	which golint > /dev/null 2>&1 || go get -u github.com/golang/lint/golint
	which goimports > /dev/null 2>&1 || go get -u golang.org/x/tools/cmd/goimports
	which richgo > /dev/null 2>&1 || go get -u github.com/kyoh86/richgo
	which goreleaser > /dev/null 2>&1 || go get -u github.com/goreleaser/goreleaser

vendor: vendor/.timestamp

vendor/.timestamp: $(shell find $(DIRS) -name '*.go')
	dep ensure -v
	touch vendor/.timestamp

vet:
	go vet $(PACKAGES)

lint:
	! find $(DIRS) -name '*.go' | xargs goimports -d | grep '^'
	echo $(PACKAGES) | xargs -n 1 golint -set_exit_status

fmt:
	find $(DIRS) -name '*.go' | xargs goimports -w

test:
	richgo test -v -race $(PACKAGES)

test.nocache:
	richgo test -count=1 -v -race $(PACKAGES)

bin: vendor bin/hfwd

bin/hfwd: $(shell find $(DIRS) -name '*.go')
	go build -o $@ ./hfwd.go

version.list:
	@ echo 'latest releases:'
	@ git tag -l | sort -nr | head -5

version.new: version.list
ifndef NEW_VERSION
	$(eval NEW_VERSION = $(shell \
	  read -p $$'\e[33mPlease enter the version name for the new release\e[0m: ' new_version; \
	  echo $${new_version} \
	 ))
endif

release: version.new
	@ read -p $$'\e[33mRelease $(NEW_VERSION) ? (y/n)\e[0m: ' answer; \
	if [ "$${answer}" != "y" ]; then exit 1; fi
	git tag -a $(NEW_VERSION) -m ""
	git push origin $(NEW_VERSION)
	goreleaser --rm-dist
