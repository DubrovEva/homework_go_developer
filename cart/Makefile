BINDIR=${CURDIR}/bin
PACKAGE=route256/cart

GREEN=\033[32m
RED=\033[31m
RESET=\033[0m

bindir:
	mkdir -p ${BINDIR}

build: bindir
	echo "build cart"

e2e:
	@touch ./tests/e2e/container.lock
	@docker compose -f ../docker-compose.yaml up --build -d
	@{ \
		go test ./tests/e2e/... -v; \
		EXIT_CODE=$$?; \
		docker compose down -v; \
		rm ./tests/e2e/container.lock; \
		if [ $$EXIT_CODE -eq 0 ]; then \
			echo "$(GREEN)Tests passed$(RESET)"; \
		else \
			echo "$(RED)Tests failed$(RESET)"; \
			exit $$EXIT_CODE; \
		fi \
	}

e2e-bench:
	@touch ./tests/benchmarks/container.lock
	@docker compose -f ../docker-compose.yaml up --build -d
	@{ \
		go test -bench=. ./tests/benchmarks; \
		EXIT_CODE=$$?; \
		docker compose down -v; \
		rm ./tests/benchmarks/container.lock; \
		if [ $$EXIT_CODE -eq 0 ]; then \
			echo "$(GREEN)Tests passed$(RESET)"; \
		else \
			echo "$(RED)Tests failed$(RESET)"; \
			exit $$EXIT_CODE; \
		fi \
	}


bench:
	go test -bench=. ./internal/repository