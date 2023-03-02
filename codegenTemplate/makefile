files = $(shell find ./*/  -name "makefile" -maxdepth 1 | xargs -n1 -I {} dirname {} | grep -v "./tools")

.PHONY: help

help: ## This help.
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help

code_check: ## check your code
	go get -d github.com/golangci/golangci-lint/cmd/golangci-lint@v1.45.2
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.45.2
	@for file in $(files); \
	do \
		echo "code check "$$file; \
		make code-check -C $$file; \
	done
	@echo "\033[32mSuccess\033[0m"

init: ## init code check environment and add auto code check before git push for all projects
	@echo "\033[33mNotice: the old file will be replaced \033[0m"
	@for file in $(files); \
	do \
		echo "init "$$file; \
		cp -af tools/* $$file; \
	done
	@make init_git
	@echo "finished"

init_proj: ## init code check environment for target project
	@cp -ai tools/* $(shell read -p "Enter target project:" name; echo $$name)

init_proj_force: ## init code check environment for target project without covering ask
	@echo "\033[33mNotice: the old file will be replaced\033[0m"
	@cp -af tools/* $(shell read -p "Enter target project:" name; echo $$name)

init_git: ## add auto code check before git push
	@echo "\033[33mNotice: file '.git/hooks/pre-push' will be cleaned\033[0m"
	@rm -f .git/hooks/pre-push
	@echo "make -i code_check" >> .git/hooks/pre-push
	@chmod 700 .git/hooks/pre-push
	@echo "finished"