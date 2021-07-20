build:
	go build -v .
.PHONY: build

docker-pg-development:
	docker run -d \
	 --name pg \
	 --network host
	  -e POSTGRES_PASSWORD=root \
	  -e PGDATA=/var/lib/postgresql/data/pgdata \
	  -v C:\Users\Ivank\Documents\cost-management-api\.database:/var/lib/postgresql/data \
	  -v C:\Users\Ivank\Documents\cost-management-api\migration\init_database.sql:/docker-entrypoint-initdb.d/init.sql \
	  -p 5432:5432 \
	  postgres
.PHONY: docker-pg-development

cost-management-api-development:
	docker run -d --cost-management-api -p 8080:8080
.PHONY: cost-management-api-development

migrate-up: 
	migrate -path ./migration -database 'postgres://pord:root@localhost:5432/cost_management_api?sslmode=disable' up
.PHONY: migrate-up

migrate-down:
	migrate -path ./migration -database 'postgres://pord:root@localhost:5432/cost_management_api?sslmode=disable' down
.PHONY: migrate-down

docker-connect-db:
	docker exec -it pg psql -Upord -dcost_management_api
.PHONY: docker-connect-db

.DEFAULT_GOAL := build

