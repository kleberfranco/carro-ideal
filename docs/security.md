# Segurança

- Senhas são validadas e armazenadas com hash bcrypt.
- Sessões usam token aleatório, hash no banco, expiração e cookie `HttpOnly`.
- Em produção, cookies recebem `Secure`; `SameSite=Strict` reduz envio cruzado.
- Requisições mutáveis exigem token CSRF em cookie e cabeçalho.
- CORS aceita somente origens configuradas e o rate limit usa token bucket por IP.
- Respostas incluem CSP, proteção contra framing e MIME sniffing.
- `html/template` escapa saída HTML; conteúdo inserido via JavaScript passa por
  `text()` antes de compor marcação. Dados são preservados no banco e escapados
  no ponto de saída, evitando dupla codificação.
- SQL usa parâmetros posicionais para todos os valores externos.
- Logs não registram senhas, cookies, tokens nem corpos de requisição.

Nunca use `SESSION_SECRET=change-me` em produção. Restrinja o acesso ao painel
administrativo, mantenha dependências e imagens atualizadas e trate HTTPS,
backups, rotação de segredos e revisão de privilégios como requisitos de release.
