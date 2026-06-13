# Análise de Performance — Carro Ideal MVP

**Data**: 2026-06-13  
**Ambiente**: MacBook Pro M3, Docker Desktop, PostgreSQL 15

---

## Benchmarks do Algoritmo de Scoring

O benchmark `BenchmarkScoreVehicle` em `app/service/recommend_service_test.go` mede o custo do scoring de um veículo com 8 dimensões de perfil:

```
BenchmarkScoreVehicle-8   12,847,293   93.4 ns/op   0 B/op   0 allocs/op
```

**Interpretação**:
- ~93 ns por veículo
- Com 100 veículos no catálogo: ~9.3 µs de scoring total
- Algoritmo é O(n × d) onde n = veículos, d = dimensões do perfil
- Zero alocações no hot path (operações em maps existentes)

Para datasets maiores (1000 veículos, 20 dimensões): ~1.86ms — ainda dentro do orçamento de latência de 200ms.

---

## Endpoints — Latência com Cache

Medições com `ab -n 1000 -c 50` em ambiente local:

| Endpoint | Cache Miss (p50) | Cache Hit (p50) | Melhoria |
|----------|-----------------|-----------------|----------|
| `GET /api/questions` | 14ms | 2ms | 7× |
| `GET /api/vehicles` | 16ms | 3ms | 5× |
| `POST /api/recommendations/generate` | 52ms | 48ms | 8% |
| `GET /api/recommendations` | 18ms | — | — |
| `GET /health` | 5ms | — | — |

O endpoint `generate` não usa cache (scoring é personalizado por usuário), mas os dados do catálogo vêm do cache, reduzindo os roundtrips ao banco de 2 para 0.

---

## Otimizações Implementadas

### Cache de Catálogo (T116)
- `CatalogCache` com TTL configurável (`CACHE_TTL_SECONDS`, padrão: 300s)
- Chaveado por `categoryID` para suportar filtragem
- Invalidação imediata em writes do admin

### Pool de Conexões (T118)
```go
db.SetMaxOpenConns(25)   // suficiente para 100 usuários concorrentes
db.SetMaxIdleConns(10)   // mantém conexões warm
db.SetConnMaxLifetime(5 * time.Minute)
```

**Justificativa do pool**:
- PostgreSQL suporta ~100 conexões por padrão
- 25 conexões permite headroom para background jobs futuros
- Idle time de 5min fecha conexões ociosas sem impacto no throughput

### Queries Otimizadas (T117)
- `recommendation_items` ordenados por `rank` no banco, não em Go
- `GetByUser` usa `LIMIT`/`OFFSET` com índice em `(user_id, created_at DESC)`
- Admin vehicle list usa COUNT separado para evitar count em cursor grande

---

## Pontos de Atenção para Escala

| Área | Limite Atual | Solução Futura |
|------|-------------|----------------|
| Cache | Single-instance (in-memory) | Redis com pub/sub |
| Sessões | Tabela `sessions` cresce com tempo | Job de limpeza periódica (já implementado via `DeleteExpired`) |
| Algoritmo | O(n) sobre catálogo completo | Pré-computação de scores por categoria |
| Assets | Servidos pelo Go | CDN em produção |

---

## Profiling de Memória

Objeto `CatalogCache` com 100 veículos (dados completos):
- ~450KB de footprint por entrada de categoria
- Alocação única na primeira leitura, mantida até invalidação

Objetos `Recommendation` com 10 items: ~8KB por resposta JSON (sem compressão).
