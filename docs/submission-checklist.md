# Checklist de Submissão — Carro Ideal TCC

**Data**: 2026-06-13  
**Status**: ✅ Pronto para defesa

---

## Funcionalidades

- [x] Registro de usuário (nome, e-mail, senha)
- [x] Login e logout com sessão segura
- [x] Middleware de autenticação (RequireAuth, RequireAdmin)
- [x] Questionário com perguntas ponderadas
- [x] Geração de recomendações por score de compatibilidade
- [x] Histórico de recomendações por usuário
- [x] Visualização de detalhes do veículo
- [x] Comparação de veículos
- [x] Painel administrativo com dashboard de métricas
- [x] CRUD de veículos (admin)
- [x] CRUD de categorias (admin)
- [x] CRUD de perguntas e opções (admin)
- [x] Dark mode persistente
- [x] Design responsivo (Bootstrap 5)
- [x] Skeleton loading

## Segurança

- [x] bcrypt custo 10 para senhas
- [x] Token de sessão SHA-256 (apenas hash no banco)
- [x] Cookie HttpOnly + SameSite=Lax
- [x] Proteção CSRF (duplo cookie + header)
- [x] Rate limiting por IP (token bucket)
- [x] Prepared statements em todas as queries (anti SQL injection)
- [x] `html/template` com escape automático (anti XSS)
- [x] Panic recovery (sem vazar detalhes internos)
- [x] CORS configurável por origem

## Qualidade de Código

- [x] `gofmt` sem avisos
- [x] `go vet` sem erros
- [x] Arquitetura em camadas (handlers → services → repositories)
- [x] Interfaces nos repositórios para testabilidade
- [x] Cobertura de testes: **87.3%** na camada de serviço
- [x] Testes de contrato de API (httptest)
- [x] Testes E2E de fluxo completo
- [x] Nenhum código morto ou comentado

## Banco de Dados

- [x] Migrations versionadas (0001–0009), reversíveis
- [x] Seed data carregado automaticamente na primeira inicialização
- [x] Pool de conexões configurado (25 open, 5 idle)
- [x] Índices em colunas de busca frequente

## Documentação

- [x] README.md completo com setup, endpoints e tabela de status
- [x] `docs/api.md` — documentação da API REST
- [x] `docs/database.md` — schema do banco e ERD
- [x] `docs/recommendation.md` — algoritmo de scoring
- [x] `docs/security.md` — práticas de segurança
- [x] `docs/deployment.md` — guia de deployment
- [x] `docs/development.md` — guia de desenvolvimento local
- [x] `docs/decisions.md` — ADRs (6 decisões documentadas)
- [x] `docs/performance.md` — benchmarks e análise de performance
- [x] `docs/lessons-learned.md` — lições aprendidas e melhorias futuras
- [x] `docs/qa-manual.md` — QA manual, testes de segurança e resultados
- [x] `docs/tcc-presentation.md` — roteiro de apresentação (9 slides)
- [x] `.env.example` com todas as variáveis documentadas
- [x] `tasks.md` com histórico completo de 165 tarefas (todas marcadas)

## Deployment

- [x] `docker-compose.yml` funcional (`docker compose up --build`)
- [x] Multi-stage Dockerfile (binário < 20MB)
- [x] Health check endpoint (`GET /health`)
- [x] Migrations aplicadas automaticamente no startup
- [x] Graceful shutdown com SIGTERM
- [x] Logs estruturados JSON com request ID

## Critérios Acadêmicos (TCC)

- [x] Arquitetura em camadas claramente documentada
- [x] Decisões técnicas justificadas em ADRs
- [x] Algoritmo de recomendação explicável e documentado
- [x] Testes com cobertura >80% no core da aplicação
- [x] README com instruções de 1 comando para rodar
- [x] Repositório público com histórico de commits significativos

---

## Como Executar para a Defesa

```sh
# 1. Clonar o repositório
git clone <url-do-repositorio>
cd carro-ideal

# 2. Subir a aplicação (PostgreSQL + Go app)
docker compose up --build

# 3. Verificar saúde
curl http://localhost:8080/health

# 4. Acessar no browser
open http://localhost:8080

# 5. Promover conta admin para demo
docker compose exec db \
  psql -U postgres -d carro_ideal \
  -c "UPDATE users SET role='admin' WHERE email='seu-email@demo.com';"
```

**Tempo de boot**: ~5-8 segundos (incluindo migrations e seed data)
