ifndef LINT_OUTPUT
LINT_OUTPUT=false
endif

include $(GFW)/make/common.mk

all: build

.PHONY: build
build: lint $(SUBDIRS)

.PHONY: lint
lint:
	@SCRIPTSDIR=$(GFW)/lint/repo-scripts $(GFW)/lint/run.sh

python:
	make -C sdk/python setup
	$(MAKE) L=python build
	make -C sdk/python sdk

.PHONY: publish
publish: build docs python
	make -C publish


.PHONY: docs
docs:
	$(MAKE) -C website

.PHONY: serve
serve:
	cd docs ; python3 -m http.server

.PHONY: clean
clean:
	$(MAKE) -C website clean
