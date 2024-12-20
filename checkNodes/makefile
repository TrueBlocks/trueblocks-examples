#-------------------------------------------------
bin=bin

#-------------------------------------------------
exec=checkNodes.1
dest=$(bin)/$(exec)

#-------------------------------------------------
all:
	@make app

every:
	@cd ../../build ; make ; cd -
	@make app

app:
	@mkdir -p bin
	@cp -fp checkNodes ./bin # copy the watcher script to ./bin
	@chmod uog+x ./bin/checkNodes
	go build -o $(dest) .

test:
	@make
	@echo "No tests"

#-------------------------------------------------
install:
	@make app
	@cp -p ./bin/checkNodes.1 /usr/local/bin
	@cp -p ./bin/checkNodes /usr/local/bin
	@echo "Installed checkNodes to /usr/local/bin"

uninstall:
	@rm -f /usr/local/bin/checkNodes.1
	@rm -f /usr/local/bin/checkNodes
	@make clean
	@echo "Removed checkNodes from /usr/local/bin"

#-------------------------------------------------
clean:
	-@$(RM) -f ./bin/checkNodes
	-@$(RM) -f $(dest)
