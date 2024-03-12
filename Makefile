setup-dev-db:
	RESUME_AI_ENV=development go run scripts/setup_dev_env.go

run-dev-api-server:
	RESUME_AI_ENV=development go run main.go

run-mongo-db:
	docker-compose up -d

restart-mongo-db:
	docker-compose restart