build: ## builds for current OS and architecture
	go build -o gather main.go

install: build
	mv gather ~/go/bin/