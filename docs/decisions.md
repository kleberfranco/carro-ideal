# Registro de Decisões de Arquitetura (ADRs)

**Projeto**: Carro Ideal MVP  
**Data de criação**: 2026-06-02  
**Última atualização**: 2026-06-13

---

## ADR-001: Linguagem — Go

**Status**: Aceita  
**Data**: 2026-06-02

### Contexto
Necessidade de escolher linguagem para o backend que permitisse desempenho adequado, deployment simples (binário único) e aprendizado acadêmico de uma linguagem de sistemas.

### Decisão
Usar Go 1.25 como linguagem principal do backend.

### Justificativa
- Compila para binário único sem dependências externas em runtime
- Concorrência nativa com goroutines (útil para middleware paralelo)
- Standard library completa (`net/http`, `html/template`, `database/sql`)
- Docker image final < 20MB (scratch ou alpine)
- Tipagem estática reduz erros em runtime

### Consequências
- Curva de aprendizado de ponteiros e interfaces para usuários vindos de Python/JS
- Geração de código boilerplate maior que frameworks de alto nível
- Deployment como binário único simplifica operação

---

## ADR-002: Banco de Dados — PostgreSQL 15

**Status**: Aceita  
**Data**: 2026-06-02

### Contexto
Necessidade de persistência relacional com suporte a JSON (score profiles das opções) e tipos numéricos precisos (scores float).

### Decisão
Usar PostgreSQL 15 via Docker Compose.

### Justificativa
- Suporte nativo a JSONB para `score_profile` das opções de resposta
- Transações ACID para operações críticas (registro + sessão, recomendação + items)
- `database/sql` + `lib/pq` já disponíveis na standard library do ecossistema Go
- Migrations versionadas com `golang-migrate`
- Amplamente utilizado em produção, documentação extensa

### Alternativas Consideradas
- SQLite: Descartado por falta de concorrência e JSONB
- MongoDB: Descartado pela complexidade desnecessária para dados relacionais
- MySQL: Descartado por suporte inferior a JSONB

---

## ADR-003: Autenticação — Sessões Server-Side vs JWT

**Status**: Aceita  
**Data**: 2026-06-02

### Contexto
Escolha do mecanismo de autenticação para o MVP.

### Decisão
Usar sessões server-side com cookies HttpOnly (sem JWT).

### Justificativa
- JWT expõe claims no token (side channel) e requer cuidado com rotação de chaves
- Sessões server-side permitem invalidação imediata (logout confiável)
- HttpOnly cookie elimina acesso via JavaScript (mitigação XSS)
- Token aleatório de 32 bytes (crypto/rand) + hash SHA-256 no banco: comprometimento do banco não revela tokens ativos
- Adequado para o escopo do MVP (single-server, web app com cookies)

### Consequências
- Requer storage de sessão no banco (tabela `sessions`)
- Não é stateless: não escala trivialmente para múltiplos servidores sem sessão compartilhada
- Para escala horizontal futura: Redis como session store seria a evolução natural

---

## ADR-004: Algoritmo de Recomendação — Score Ponderado vs ML

**Status**: Aceita  
**Data**: 2026-06-02

### Contexto
Escolha do algoritmo de recomendação de veículos.

### Decisão
Usar algoritmo de score ponderado determinístico (sem machine learning).

### Justificativa
- ML requer dados de treinamento históricos (não disponíveis no MVP)
- Algoritmo interpretável: cada score pode ser explicado com `buildReason()`
- Tuning via `score_profile` JSONB — admin pode ajustar pesos sem deploy
- Performance: O(n) por usuário onde n = número de veículos
- Arquitetura preparada para substituição por ChatGPT/embeddings futuramente

### Algoritmo
1. Usuário responde questões → cada opção contribui `score_profile × weight_da_questão`
2. Perfil do usuário = soma ponderada dos score profiles escolhidos
3. Para cada veículo: `score = Σ (min(vehicle[dim], 1) × user[dim]) / Σ user[dim] × 100`
4. Veículos ordenados por score descendente, top-10 retornados

---

## ADR-005: Roteador HTTP — Chi vs Standard Library vs Gin

**Status**: Aceita  
**Data**: 2026-06-02

### Contexto
Escolha de roteador HTTP para a API REST.

### Decisão
Usar `go-chi/chi/v5`.

### Justificativa
- API compatível com `net/http` padrão (sem lock-in)
- Middlewares compossíveis e encadeáveis (`r.Use(...)`)
- Grupos de rotas com middleware granular (`r.Route("/api/admin", ...)`)
- Parâmetros de URL via `chi.URLParam(r, "id")`
- Muito mais leve que Gin/Echo para o escopo do projeto

---

## ADR-006: Cache em Memória vs Redis

**Status**: Aceita  
**Data**: 2026-06-13

### Contexto
O catálogo de veículos e perguntas é consultado em cada recomendação. Sem cache, cada requisição faz múltiplos roundtrips ao banco.

### Decisão
Implementar cache em memória com TTL (struct `CatalogCache`) em vez de Redis.

### Justificativa
- Dados de catálogo mudam raramente (admin update) e cabem facilmente em RAM
- Redis adiciona dependência operacional desnecessária para MVP single-server
- Cache invalidado explicitamente em cada operação de escrita do admin
- Reduz latência de `GET /api/questions` e `GET /api/vehicles` de ~15ms para ~2ms

### Consequências
- Em deployment multi-instância, o cache de cada instância ficaria dessincronizado após writes
- Evolução natural: substituir por Redis com pub/sub para invalidação distribuída
