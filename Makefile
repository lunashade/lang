TARGET = lang
BINDIR = bin
TMPDIR = tmp

clean:
	rm -rf $(BINDIR)
	rm -rf $(TMPDIR)

build:
	mkdir -p $(BINDIR)
	go build -o $(BINDIR)/$(TARGET) .

test: clean build
	go test ./...
	./test.sh $(BINDIR)/$(TARGET)
