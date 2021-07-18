FROM golang:latest
WORKDIR /app
COPY . .

RUN apt-get update
RUN apt-get -y install postgresql-client

RUN chmod +x ./scripts/wait-for-postgres.sh

RUN go mod download
RUN go build -o cost-management-api ./cmp/cost-management-api
EXPOSE 8080
CMD ["./cost-management-api"]