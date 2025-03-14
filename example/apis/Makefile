export TOP=$(shell pwd)
SUBDIRS = hello

all: $(SUBDIRS)

.PHONY: $(SUBDIRS)
$(SUBDIRS):
	$(MAKE) -C $@
