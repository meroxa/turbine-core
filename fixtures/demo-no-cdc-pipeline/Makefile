.PHONY: all
all: clean run

.PHONY: run
run: clean
	conduit -pipelines.path ./ -pipelines.exit-on-error

.PHONY: clean
clean:
	rm -rf conduit.db