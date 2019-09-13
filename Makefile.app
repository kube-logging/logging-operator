.PHONY: check
check: test license

.PHONY: license
license:
	./scripts/check-header.sh