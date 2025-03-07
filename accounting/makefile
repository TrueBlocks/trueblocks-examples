#-------------------------------------------------
bin=../../bin/examples/

#-------------------------------------------------
exec=accounting
dest=$(bin)/$(exec)

#-------------------------------------------------
all:
	@make app

every:
	@cd ../../build ; make ; cd -
	@make app

app:
	@mkdir -p bin
	@go build -o $(dest) *.go

test:
	@make app
	@cd tests && ./run_all.sh && cd -

data:
	@make every
	@cd tests && go run update_data.go && cd -

#-------------------------------------------------
clean:
	-@$(RM) -f $(dest)
