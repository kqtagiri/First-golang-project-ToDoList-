include .env
export

start_program:
	@go run main.go

deploy_program:
	@docker compose up -d

undeploy_program:
	@docker compose down