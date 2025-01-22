run:
	go run ./cmd

swag:
	swag init -g ./internal/app/app.go
	@echo "Done!"