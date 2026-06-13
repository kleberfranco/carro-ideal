# Implantação

1. Configure as variáveis descritas em `.env.example`, usando segredos fortes.
2. Execute `docker compose build`.
3. Execute `docker compose up -d`.
4. Confirme `GET /health` e os logs JSON do contêiner.

As migrations e seeds são aplicadas automaticamente antes do servidor iniciar.
Faça backup do PostgreSQL antes de atualizar uma instalação existente.

Em produção, prefira terminar HTTPS em um proxy reverso. Para TLS direto no Go,
monte o certificado e a chave no contêiner e defina `TLS_CERT_FILE` e
`TLS_KEY_FILE`. Configure `ALLOWED_ORIGINS`, cookies seguros via
`ENVIRONMENT=production`, limites de requisição e recursos do contêiner.

Checklist: banco sem porta pública, segredo de sessão rotacionável, HTTPS ativo,
origens restritas, backup testado, health check monitorado e logs coletados.
Se o app não iniciar, verifique saúde do banco, migrations, caminhos TLS e
permissões dos arquivos montados.
