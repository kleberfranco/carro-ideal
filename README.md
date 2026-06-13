# Carro Ideal

Aplicação web em Go que recomenda veículos a partir das preferências do usuário. O projeto é desenvolvido como TCC com arquitetura em camadas, PostgreSQL, API REST, templates server-side e execução via Docker.

## Estado atual

- Fase 1: fundação e infraestrutura concluídas
- Fase 2: autenticação principal concluída
- Fase 3: motor de recomendação funcional, com expansão de seeds e testes de integração pendente
- Cadastro, login, consulta da sessão e logout funcionais
- Sessões aleatórias de 32 bytes persistidas no PostgreSQL apenas como hash
- Questionário, ranking, detalhes de veículos e histórico funcionais
- Painel administrativo ainda não implementado

O planejamento detalhado está em [tasks.md](tasks.md) e os artefatos do Spec Kit estão em `specs/001-carro-ideal-mvp/`.

## Arquitetura

```text
HTTP handlers
      |
      v
Services (regras de negócio e autenticação)
      |
      v
Repositories (database/sql)
      |
      v
PostgreSQL
```

Diretórios principais:

```text
app/internal/   handlers, middleware e rotas
app/service/    regras de negócio
app/repository/ acesso ao PostgreSQL
app/models/     estruturas de domínio
migrations/     migrations versionadas
web/            templates e arquivos estáticos
specs/          especificação e plano do Spec Kit
```

## Execução

O Go deve ser executado somente via Docker.

```sh
docker compose up --build
```

A aplicação fica disponível em `http://localhost:8080`.

Health check:

```sh
curl http://localhost:8080/health
```

Parar os containers:

```sh
docker compose down
```

O volume do PostgreSQL é preservado por padrão. Não use `docker compose down -v` se quiser manter os dados.

## Testes

```sh
docker run --rm \
  -e GOCACHE=/tmp/go-cache \
  -e GOMODCACHE=/tmp/go-mod \
  -v "$PWD":/src \
  -w /src \
  golang:1.22-alpine \
  /usr/local/go/bin/go test ./...
```

## Endpoints

| Método | Caminho | Descrição |
|---|---|---|
| GET | `/health` | Verifica aplicação e banco |
| POST | `/api/auth/register` | Cria usuário e inicia sessão |
| POST | `/api/auth/login` | Autentica e inicia sessão |
| POST | `/api/auth/logout` | Invalida a sessão |
| GET | `/api/auth/me` | Retorna o usuário autenticado |
| GET | `/api/user/` | Endpoint protegido ainda não implementado |
| GET | `/api/questions` | Retorna perguntas e opções ativas |
| POST | `/api/recommendations/generate` | Gera e persiste um ranking |
| GET | `/api/recommendations` | Lista o histórico do usuário |
| GET | `/api/recommendations/{id}` | Abre um resultado salvo |
| GET | `/api/vehicles` | Lista o catálogo ativo |
| GET | `/api/vehicles/{id}` | Retorna detalhes do veículo |

As respostas da API seguem:

```json
{
  "success": true,
  "data": {}
}
```

Erros incluem `error` e `code`.

## Banco de dados

As migrations são aplicadas automaticamente na inicialização da aplicação. O estado atual inclui:

- `users`: usuários, papéis, senha com bcrypt e status ativo
- `sessions`: hash do token, usuário e expiração
- `questions` e `answer_options`: questionário ponderado
- `vehicle_categories` e `vehicles`: catálogo e perfis de compatibilidade
- `recommendations`, `recommendation_items` e `user_answers`: histórico reproduzível

As migrations `0006` e `0007` atualizam de forma não destrutiva bancos criados pelas versões iniciais do projeto.

As migrations `0008` e `0009` criam e populam o domínio de recomendação.

## Configuração

Variáveis disponíveis em `.env.example`:

- `ENVIRONMENT`
- `PORT`
- `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`
- `DB_MAX_OPEN_CONNS`, `DB_MAX_IDLE_CONNS`
- `SESSION_SECRET`
- `LOG_LEVEL`

Em produção, cookies de sessão recebem também a flag `Secure`.
