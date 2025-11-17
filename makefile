all:
	@echo "To build any of these examples, cd into the directory and type make."

update:
	@echo "To update any of these examples, cd into the directory and type make update."

goMaker:
	@cd ../dev-tools/goMaker && yarn deploy && cd -

generate:
	@goMaker
