ifdef SUBDIRS
.PHONY: $(SUBDIRS)
$(SUBDIRS):
	$(MAKE) -C $@
endif

ifdef PROTO_FILES
PWD=$(shell pwd)

# Doing this will make sure to use the full path as the names
# of the protobufs to avoid conflicts.
PROTO_PATH:=$(PWD:$(TOP)/%=%)
BUILD_PROTO_FILES:=$(addprefix $(PROTO_PATH)/, $(PROTO_FILES))

.PHONY: $(PROTO_FILES)
$(PROTO_FILES): $(BUILD_PROTO_FILES)

$(BUILD_PROTO_FILES):
	@cd $(TOP) && grpcfw-go $@
	@cd $(TOP) && grpcfw-rest $@
	@cd $(TOP) && grpcfw-doc $@
	@cd $(TOP) && grpcfw-lint \
		--config=$(TOP)/rules.yml \
		--set-exit-status \
		$@
endif