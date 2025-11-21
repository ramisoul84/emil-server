MAIN_FILE=cmd/emil-server/main.go

dev:
	APP_ENV=development go run $(MAIN_FILE)

prod:
	APP_ENV=production go run $(MAIN_FILE)