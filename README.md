# Carro Ideal

Aplicação web em Go que recomenda veículos a partir das preferências do usuário. O projeto é desenvolvido como TCC com arquitetura em camadas, PostgreSQL, API REST, templates server-side e execução via Docker.

## Estado atual — MVP Completo ✅

Todas as 6 fases do planejamento foram concluídas:

| Fase | Escopo | Status |
|------|--------|--------|
| 1 — Fundação | Docker, PostgreSQL, roteamento, health check | ✅ Concluído |
| 2 — Autenticação | Registro, login, sessões, middleware | ✅ Concluído |
| 3 — Motor de recomendação | Questionário ponderado, algoritmo de scoring, histórico | ✅ Concluído |
| 4 — Painel Admin | CRUD de veículos, categorias, perguntas e opções | ✅ Concluído |
| 5 — Polimento | Segurança, cache, observabilidade, dark mode | ✅ Concluído |
| 6 — QA e Docs | Testes unitários (87% cobertura), API contract, E2E, documentação final | ✅ Concluído |

**Funcionalidades entregues**:
- Registro e login com sessões seguras (bcrypt + SHA-256 + HttpOnly cookie)
- Questionário ponderado com cálculo de perfil de preferências
- Algoritmo de recomendação por score de compatibilidade, top-10 ranqueado
- Histórico de recomendações por usuário
- Painel administrativo com dashboard e CRUD completo
- Logs JSON estruturados, request ID, recuperação de panics, CORS, rate limit, CSRF
- Comparação de veículos lado a lado, dark mode persistente, skeleton loading
- Cache TTL para catálogo (veículos e perguntas), paginação otimizada

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
docs/           guias de API, banco, algoritmo, segurança e implantação
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

Executar suite de testes unitários e de contrato:

```sh
docker run --rm \
  -e GOCACHE=/tmp/go-cache \
  -e GOMODCACHE=/tmp/go-mod \
  -v "$PWD":/src \
  -w /src \
  golang:1.25-alpine \
  go test ./...
```

**Cobertura atual**:

| Pacote | Cobertura |
|--------|-----------|
| `app/service` | 87.3% |
| `app/internal/platform` | 51.4% |
| `app/internal/api` | 54.0% |

Os testes cobrem: ciclo de vida de sessão, algoritmo de scoring, validações de admin, middleware de autenticação/CSRF/rate-limit, handlers HTTP (contrato de API) e fluxo E2E completo (registro → recomendação → histórico → logout).

Para resultados de QA manual e testes de segurança, veja [docs/qa-manual.md](docs/qa-manual.md).

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
| GET | `/api/admin/dashboard` | Métricas administrativas |
| CRUD | `/api/admin/vehicles` | Gestão do catálogo |
| CRUD | `/api/admin/categories` | Gestão de categorias |
| CRUD | `/api/admin/questions` | Gestão de perguntas e opções |

## Acesso administrativo

O cadastro público sempre cria usuários com papel `user`. Para promover uma conta no ambiente local:

```sh
docker compose exec db \
  psql -U postgres -d carro_ideal \
  -c "UPDATE users SET role='admin' WHERE email='seu-email@exemplo.com';"
```

Depois do login, o painel fica disponível em `http://localhost:8080/admin`.

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
- `ALLOWED_ORIGINS`
- `RATE_LIMIT_REQUESTS`, `RATE_LIMIT_WINDOW_SECONDS`
- `CACHE_TTL_SECONDS`
- `TLS_CERT_FILE`, `TLS_KEY_FILE`

Em produção, cookies de sessão recebem também a flag `Secure`.

## Documentação

- [API](docs/api.md)
- [Banco de dados](docs/database.md)
- [Motor de recomendação](docs/recommendation.md)
- [Desenvolvimento](docs/development.md)
- [Implantação](docs/deployment.md)
- [Segurança](docs/security.md)
