.PHONY: gomod

gomod:
	go mod tidy && go mod vendor
