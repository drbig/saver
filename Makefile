OSES  := linux darwin windows
ARCHS := amd64 386
SRCS  := $(wildcard *.go)
TGTS  := $(foreach os,$(OSES),$(foreach arch,$(ARCHS),bin/saver-$(os)-$(arch)))

all: $(TGTS) bin/checksums.md5

$(TGTS): $(SRCS)
	GOOS=$(word 2,$(subst -, ,$@)) GOARCH=$(word 3,$(subst -, ,$@)) go build -o $@ .

$(SRCS):

bin/checksums.md5:
	cd bin && md5sum * > checksums.md5
