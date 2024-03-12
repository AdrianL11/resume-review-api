run-dev-api-server:
	RESUME_AI_ENV=development go run main.go

run-mongo-db:
	docker-compose up -d

restart-mongo-db:
	docker-compose restart