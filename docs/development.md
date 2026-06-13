# Desenvolvimento

Pré-requisitos: Docker e Docker Compose. O Go deve ser executado somente em
contêineres neste projeto.

```sh
docker compose up --build
docker run --rm -v "$PWD":/app -w /app golang:1.22-alpine go test ./...
docker run --rm -v "$PWD":/app -w /app golang:1.22-alpine gofmt -w app config
```

A aplicação fica em `http://localhost:8080`; o PostgreSQL de desenvolvimento usa
o volume `postgres_data`. Mudanças de esquema devem ser novas migrations `up` e
`down`, sem alterar migrations já aplicadas. Para falhas locais, consulte
`docker compose logs app db` e confirme se a porta 8080 está livre.
