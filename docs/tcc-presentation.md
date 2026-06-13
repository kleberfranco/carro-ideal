# Roteiro de Apresentação TCC — Carro Ideal

**Tempo estimado**: 20 minutos de apresentação + 10 minutos de perguntas  
**Ferramentas**: slides + demonstração ao vivo (Docker Compose)

---

## Slide 1 — Problema e Motivação (2 min)

**Título**: "Como escolher o carro certo sem ser especialista?"

- O mercado brasileiro tem +500 modelos de veículos disponíveis
- Usuários se perdem em comparativos técnicos (consumo, câmbio, categoria)
- Decisão errada custa R$50k-200k e gera insatisfação por anos
- Solução: sistema de recomendação personalizado baseado em preferências declaradas

**Ponto-chave para a banca**: O projeto resolve um problema real com tecnologia acessível, sem depender de IA cara ou dados históricos.

---

## Slide 2 — Solução: Carro Ideal (2 min)

**Título**: "Questionário → Perfil → Recomendação"

Fluxo do usuário:
1. Registra-se (30 segundos)
2. Responde 7-8 perguntas sobre uso, orçamento e preferências (5 minutos)
3. Recebe ranking de veículos com score e justificativa (imediato)
4. Consulta histórico de recomendações

**Demo rápida**: mostrar a tela de questionário e o resultado de recomendações.

---

## Slide 3 — Arquitetura Técnica (3 min)

**Título**: "Arquitetura em Camadas com Go e PostgreSQL"

```
Browser ──► Chi Router ──► Handlers
                               │
                          Services
                               │
                        Repositories
                               │
                          PostgreSQL
```

**Tecnologias**:
- Go 1.25 — binário único, performance nativa, tipagem estática
- PostgreSQL 15 — JSONB para score profiles, transações ACID
- Docker Compose — desenvolvimento e deployment em um comando
- go-chi/chi — roteamento HTTP leve, compatível com net/http

**Destaques de design**:
- Interfaces em todos os repositórios → testabilidade sem banco real
- `html/template` → proteção XSS automática
- Sessões server-side com SHA-256 hash → token seguro mesmo com leak de banco

---

## Slide 4 — Algoritmo de Recomendação (3 min)

**Título**: "Score de Compatibilidade: Simples, Explicável, Extensível"

**Como funciona**:
1. Cada resposta tem um `score_profile` (ex: `{"urban": 0.9, "efficiency": 0.8}`)
2. Questão tem peso configurável (ex: orçamento = 0.3, combustível = 0.2)
3. Perfil do usuário = soma ponderada dos profiles das respostas escolhidas
4. Score do veículo = interseção normalizada entre perfil do usuário e perfil do veículo

**Exemplo**:
- Usuário: `{"urban": 2.0, "efficiency": 1.5, "automatic": 1.0}`
- Toyota Corolla: `{"urban": 0.9, "efficiency": 0.8, "automatic": 1.0}`
- Score = 87%

**Vantagens**:
- Zero dados históricos necessários
- Admin ajusta pesos sem deploy (JSONB + painel admin)
- Cada recomendação tem razão textual gerada automaticamente
- Arquitetura pronta para substituir por embeddings/ChatGPT

---

## Slide 5 — Segurança Implementada (2 min)

**Título**: "Segurança por Design"

| Camada | Mecanismo |
|--------|-----------|
| Autenticação | Sessão 32 bytes (crypto/rand) + hash SHA-256 no banco |
| XSS | html/template escapa automaticamente |
| SQL Injection | 100% prepared statements ($1, $2...) |
| CSRF | Token duplo cookie+header (X-CSRF-Token) |
| Senhas | bcrypt custo 10 |
| Cookies | HttpOnly, Secure (prod), SameSite=Lax |
| Rate Limit | 100 req/min por IP (token bucket) |

---

## Slide 6 — Qualidade de Código (2 min)

**Título**: "Testes e Cobertura"

**Cobertura atual**:
- Camada de serviço: **87.3%** (47 testes)
- Middleware de segurança: 51.4% (CSRF, rate limit, panic recovery)
- Handlers HTTP: 54.0% (registro, login, recomendações, E2E)

**Tipos de teste**:
- Unitários com fakes (sem banco): AuthService, UserService, RecommendationService, AdminService
- Contrato de API: todos os endpoints críticos via httptest
- Fluxo E2E: registro → questionário → recomendação → histórico → logout
- Benchmark: algoritmo de scoring (93ns/veículo, zero alocações)

---

## Slide 7 — Painel Administrativo (1 min)

**Título**: "Controle Total pelo Admin"

- Dashboard com métricas em tempo real
- CRUD de veículos com validação e soft delete
- CRUD de categorias, perguntas e opções de resposta
- Ajuste de pesos do algoritmo via painel (sem deploy)
- Autorização: apenas role=admin acessa /admin/*

---

## Slide 8 — Demonstração ao Vivo (5 min)

**Roteiro**:
1. `docker compose up --build` → mostrar que sobe em ~5 segundos
2. `curl localhost:8080/health` → resposta JSON healthy
3. Registrar novo usuário no browser
4. Completar questionário (7 perguntas)
5. Ver resultado de recomendações com scores
6. Mostrar painel admin → dashboard de métricas
7. Editar score_profile de uma opção → regenerar recomendação → ver mudança no ranking

---

## Slide 9 — Conclusão e Próximos Passos (1 min)

**Entregues**:
- MVP funcional com 40+ endpoints REST
- Algoritmo de recomendação explicável e configurável
- Suite de testes com 87% de cobertura no core
- Documentação completa: ADRs, API, banco, algoritmo, segurança, QA

**Fase 2 planejada**:
- Integração ChatGPT para recomendações em linguagem natural
- Feedback de usuário para aprimoramento iterativo do algoritmo
- Pipeline CI/CD com GitHub Actions

---

## Perguntas Comuns da Banca

**"Por que não JWT?"**
JWT não permite invalidação imediata (logout não funciona de verdade). Sessões server-side com hash SHA-256 oferecem segurança equivalente com invalidação confiável.

**"Por que não usar ML?"**
ML requer dados históricos (preferências + escolhas reais). Com zero usuários no lançamento, o algoritmo ponderado é mais adequado e igualmente eficaz para o MVP.

**"Como escala?"**
Atualmente single-server. Para multi-instância: trocar cache in-memory por Redis e sessões por Redis também. A arquitetura já prevê isso (interface `SessionRepository`).

**"Como o score é calculado?"**
Score = produto interno normalizado entre o perfil do usuário e o perfil do veículo, multiplicado por 100. Detalhes em `docs/recommendation.md`.
