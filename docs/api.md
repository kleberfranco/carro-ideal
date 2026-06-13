# API

Base local: `http://localhost:8080`. Respostas JSON usam `{"data": ...}` em sucesso
e `{"error":"...","code":"..."}` em falhas.

## Autenticação e CSRF

A autenticação usa o cookie de sessão `carro_ideal_session`. Operações `POST`, `PUT`,
`PATCH` e `DELETE` também exigem o valor do cookie `carro_csrf` no cabeçalho
`X-CSRF-Token`. Faça um `GET` inicial para receber o token.

| Método | Rota | Autenticação | Descrição |
|---|---|---:|---|
| POST | `/api/auth/register` | Não | Cria usuário com `name`, `email`, `password` |
| POST | `/api/auth/login` | Não | Abre sessão com `email`, `password` |
| POST | `/api/auth/logout` | Sim | Encerra a sessão |
| GET | `/api/auth/me` | Sim | Retorna o usuário da sessão |
| GET | `/api/questions/` | Sim | Lista questionário ativo |
| POST | `/api/recommendations/generate` | Sim | Gera ranking a partir de `answers` |
| GET | `/api/recommendations/?page=1&limit=10` | Sim | Lista histórico |
| GET | `/api/recommendations/{id}` | Sim | Detalha recomendação do usuário |
| GET | `/api/vehicles?category_id={id}` | Não | Lista veículos ativos |
| GET | `/api/vehicles/{id}` | Não | Detalha veículo ativo |
| GET | `/api/admin/dashboard` | Admin | Estatísticas do sistema |
| GET/POST | `/api/admin/vehicles` | Admin | Lista ou cria veículos |
| PUT/DELETE | `/api/admin/vehicles/{id}` | Admin | Atualiza ou desativa veículo |
| GET/POST | `/api/admin/categories` | Admin | Lista ou cria categorias |
| PUT/DELETE | `/api/admin/categories/{id}` | Admin | Atualiza ou desativa categoria |
| GET/POST | `/api/admin/questions` | Admin | Lista ou cria perguntas |
| PUT/DELETE | `/api/admin/questions/{id}` | Admin | Atualiza ou desativa pergunta |
| POST | `/api/admin/questions/{id}/options` | Admin | Cria opção |
| PUT/DELETE | `/api/admin/questions/{id}/options/{optionID}` | Admin | Atualiza ou desativa opção |

Exemplo de geração:

```json
{
  "answers": [
    {"question_id": 1, "answer_option_id": 2}
  ]
}
```

Erros comuns: `VALIDATION_ERROR` (400), `UNAUTHORIZED` (401), `FORBIDDEN` (403),
`CSRF_INVALID` (403), `NOT_FOUND` (404), `RATE_LIMITED` (429) e
`INTERNAL_ERROR` (500).
