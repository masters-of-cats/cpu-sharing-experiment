all: app
.PHONY: app

app:
	GOOS=linux go build -o bin/app app/main.go
clean:
	rm bin/app

