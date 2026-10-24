DATA_SOURCE ?= japan-mobile-career-ip-address

.PHONY: embed sync-data check-data

embed:
	git submodule update --init --recursive --remote
	$(MAKE) sync-data

sync-data:
	cp "$(DATA_SOURCE)"/data/*.json data/

check-data:
	@for carrier in docomo au softbank rakuten; do cmp "$(DATA_SOURCE)"/data/$$carrier.json data/$$carrier.json || exit 1; done
