cur = $(shell pwd)

code-check:
	@make ast-check
	@make test
	@echo "success"


test:
	@echo "testing"
	@go test -v mock/*.go | sed ''s/PASS/$$(printf "\033[32mPASS\033[0m")/'' | sed ''s/FAIL/$$(printf "\033[31mFAIL\033[0m")/'' | sed ''s/RUN/$$(printf "\033[33mRUN\033[0m")/''

ast-check:
	@echo "ast checking"
	@golangci-lint run -c tools/.golangci.yml --path-prefix $(cur) --trace-path tools/.golangci.trace
	@echo "function coverage checking"
	@go run tools/ast.go