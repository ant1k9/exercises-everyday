.PHONY: all

commands = server migrate

all: $(commands)

$(commands): %: cmd/%/main.go
	go build -o exercises-$@ $<

clean:
	@rm -f exercises-*
