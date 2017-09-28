PROTOC=protoc
GO=go

BUILD_PATH=build
BUILD_CODE_PATH=$(BUILD_PATH)/src/cook
PKGS=config connector daemon http io log mserver os pool stats util

.PHONY: test bench clean $(PKGS:%=test.%) $(PKGS:%=bench.%)

.SECONDEXPANSION:

test: $(PKGS:%=$(BUILD_PATH)/test.%)

bench: $(PKGS:%=$(BUILD_PATH)/bench.%)

clean:
	rm -rf $(BUILD_PATH)

$(BUILD_CODE_PATH):
	mkdir -p $(BUILD_CODE_PATH)

$(PKGS:%=test.%) $(PKGS:%=bench.%): %:$(BUILD_PATH)/%

$(PKGS:%=$(BUILD_PATH)/test.%): $(BUILD_PATH)/test.%:$(BUILD_CODE_PATH) $$(wildcard %/*.go)
	@echo "testing ["$(@:$(BUILD_PATH)/test.%=%)"] ..."
	@rm -rf $(BUILD_CODE_PATH)/$(@:$(BUILD_PATH)/test.%=%)
	@cp -rf $(@:$(BUILD_PATH)/test.%=%) $(BUILD_CODE_PATH)/$(@:$(BUILD_PATH)/test.%=%)
	@cd $(BUILD_PATH) ; GOPATH=${GOPATH}:$(abspath $(BUILD_PATH)) $(GO) test cook/$(@:$(BUILD_PATH)/test.%=%)
	@touch $@
	@echo "done."
	@echo

$(PKGS:%=$(BUILD_PATH)/bench.%): $(BUILD_PATH)/bench.%:$(BUILD_CODE_PATH) $$(wildcard %/*.go)
	@echo "testing ["$(@:$(BUILD_PATH)/bench.%=%)"] ..."
	@rm -rf $(BUILD_CODE_PATH)/$(@:$(BUILD_PATH)/bench.%=%)
	@cp -rf $(@:$(BUILD_PATH)/bench.%=%) $(BUILD_CODE_PATH)/$(@:$(BUILD_PATH)/bench.%=%)
	@cd $(BUILD_PATH) ; GOPATH=${GOPATH}:$(abspath $(BUILD_PATH)) $(GO) test -test.bench "." cook/$(@:$(BUILD_PATH)/bench.%=%)
	@touch $@
	@echo "done."
	@echo
