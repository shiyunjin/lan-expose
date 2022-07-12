APPEXE :=
#NAMESPACE := $(shell pwd | rev | cut -d'/' -f2 | rev)

##################################################################################################################################
ifeq ($(OS),Windows_NT)
	APPEXE = .exe
endif

ifdef GOBIN
PATH := $(GOBIN):$(PATH)
else
PATH := $(subst :,/bin:,$(shell go env GOPATH))/bin:$(PATH)
endif

COMMIT := $(shell git rev-parse --short HEAD)
COMMIT_BRANCH := $(shell git branch -l --points-at HEAD --format "%(refname:strip=2)")
COMMIT_TAG := $(shell git tag -l --points-at HEAD --format "%(refname:strip=2)")

LDFLAGS:=-s -w -X github.com/shiyunjin/lan-expose/pkg/version.commit=${COMMIT} -X github.com/shiyunjin/lan-expose/pkg/version.branch=${COMMIT_BRANCH}
ifeq ("${COMMIT_TAG}", "")
	LDFLAGS:=$(LDFLAGS) -X github.com/shiyunjin/lan-expose/pkg/version.tag=${COMMIT_TAG}
endif

##################################################################################################################################
.PHONY: all
all: upgrade proxy

.PHONY: upgrade
upgrade:
	go build -mod vendor -ldflags "$(LDFLAGS)" -o bin/lan_expose_upgrade github.com/shiyunjin/lan-expose/cmd/upgrade


.PHONY: proxy
proxy:
	go build -mod vendor -ldflags "$(LDFLAGS)" -o bin/lan_expose_proxy github.com/shiyunjin/lan-expose/cmd/proxy

##################################################################################################################################

.PHONY: deps
deps:
	# set git config
	git config --global core.ignorecase false
	go mod vendor
