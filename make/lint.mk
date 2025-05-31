CURDIR=$(shell pwd)
BINDIR=${CURDIR}/bin
GOVER=$(shell go version | perl -nle '/(go\d\S+)/; print $$1;')
LINTVER=v1.60.3
LINTBIN=bin/golangci-lint

GREEN=\033[32m
RED=\033[31m
RESET=\033[0m

LINTER_CYCLO_THRESHOLD=10

bindir:
	mkdir -p ${BINDIR}


install-lint: bindir
	test -f ${LINTBIN} || \
		(GOBIN=${BINDIR} go install github.com/golangci/golangci-lint/cmd/golangci-lint@${LINTVER} && \
		mv ${BINDIR}/golangci-lint ${LINTBIN})
	test -f ${BINDIR}/gocognit || \
		(GOBIN=${BINDIR} go install github.com/uudashr/gocognit/cmd/gocognit@latest)
	test -f ${BINDIR}/gocyclo || \
		(GOBIN=${BINDIR} go install github.com/fzipp/gocyclo/cmd/gocyclo@latest)

define lint
	@if [ -f "$(1)/go.mod" ]; then \
		output=$$(${LINTBIN} --config=.golangci.yaml run $(1)/... 2>&1); \
		exit_code=$$?; \
		echo "$$output"; \
		if [ $$exit_code -ne 0 ]; then \
			if echo "$$output" | grep -q "no go files to analyze"; then \
				exit 0; \
			else \
				exit $$exit_code; \
			fi \
		fi \
	fi
endef


run-linters:
	@LINT_OUT=$($call lint,cart > &1 || true); \
	echo "$$LINT_OUT\n$$GOCG_OUT\n$$GOCY_OUT"


cart-lint:
	$(call lint,cart)
	@OUTPUT=$$(make --no-print-directory run-linters 2>&1); \
    	if [ -z "$$OUTPUT" ]; then \
    		echo "$(GREEN)✔ Linter passed$(RESET)"; \
    	else \
    		echo "$(RED)✘ Linter failed$(RESET)"; \
    		echo "$$OUTPUT"; \
    	fi

loms-lint:
	$(call lint,loms)

notifier-lint:
	$(call lint,notifier)

comments-lint:
	$(call lint,comments)
