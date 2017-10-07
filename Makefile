env = GOOS=linux GOARCH=amd64

all:
	 $(env) go build -o prog-image .

.PHONY: all
