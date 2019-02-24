EXAMPLES = $(wildcard *.yaml)

.PHONY : $(EXAMPLES)
.PHONY : run
.PHONY : help
.DEFAULT_GOAL := run

run: $(EXAMPLES) ## Run all the examples

$(EXAMPLES):
	../yq .Items -- $@

help: ## This help.
	@echo EXAMPLES=$(EXAMPLES)
	@echo Targets:
	@awk 'BEGIN {FS = ":.*?## "} / [#][#] / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
