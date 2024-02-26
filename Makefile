.PHONY: build
build:
	dgen
	cd dgen-calendar && yarn && yarn build && cd -
	cp -r dgen-calendar/dist docs/calendar 