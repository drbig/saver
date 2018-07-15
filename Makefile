OSES    := linux darwin windows
ARCHS   := amd64 386
SRCS    := $(wildcard *.go)
VER     := $(shell grep -Eo 'VERSION = `(.*)`' main.go | cut -d'`' -f2)
TGTS    := $(foreach os,$(OSES),$(foreach arch,$(ARCHS),bin/saver-$(os)-$(arch)))
BUILD   := $(shell echo `whoami`@`hostname -s` on `date`)
LDFLAGS := -ldflags='-X "main.build=$(BUILD)"'

.PHONY: clean dev

all: $(TGTS) bin/checksums.md5

clean:
	@rm -f bin/*

dev: $(SRCS)
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o saver-$@-$(VER) .

$(TGTS): $(SRCS)
	GOOS=$(word 2,$(subst -, ,$@)) GOARCH=$(word 3,$(subst -, ,$@)) go build $(LDFLAGS) -o $@-$(VER) .

$(SRCS):

bin/checksums.md5:
	cd bin && md5sum * > checksums.md5
