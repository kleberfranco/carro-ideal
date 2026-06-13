# Lições Aprendidas — Carro Ideal TCC

**Data**: 2026-06-13

---

## O que funcionou bem

### Arquitetura em Camadas (Handlers → Services → Repositories)
A separação estrita de camadas provou seu valor durante os testes da Fase 6: foi possível testar 87% da camada de serviço com fakes simples sem um banco de dados real. A interface `AdminRepositoryInterface` foi adicionada tardiamente, mas a refatoração foi mínima (troca de `*repository.AdminRepository` por interface em dois lugares).

**Lição**: Definir interfaces nos repositórios desde o início evita refatorações posteriores.

### Spec Kit como Bússola
Ter o spec (spec.md) e o plan (tasks.md) como referência durante todo o desenvolvimento manteve o escopo controlado. Cada fase tinha um critério claro de conclusão, evitando o "feature creep".

### Cache Simples com Invalidação Explícita
O `CatalogCache` in-memory com TTL e invalidação explícita em writes resolveu o problema de performance sem adicionar complexidade operacional. Redis seria exagero para o escopo MVP.

### `html/template` para Segurança XSS
O escape automático do package `html/template` do Go eliminou toda uma classe de vulnerabilidades sem nenhum esforço adicional. Comparado a frameworks que exigem chamadas explícitas de escape, a segurança por padrão é uma vantagem significativa.

---

## Desafios encontrados

### AdminService e Testabilidade
O AdminService foi implementado com `*repository.AdminRepository` (struct concreto) em vez de interface, o que impediu testes unitários até a Fase 6. A correção foi simples — criar `AdminRepositoryInterface` — mas o custo foi horas de cobertura baixa.

**Lição**: Use interfaces nos limites de dependências desde o início, mesmo quando há apenas uma implementação.

### Score Profile e Normalização
Os primeiros seeds tinham `score_profile` com valores > 1.0 (ex: `"urban": 1.5`), o que causava scores maiores que 100 no algoritmo. A validação `validateOption()` foi adicionada para garantir que todos os valores estejam em [0, 1].

**Lição**: Validar dados na entrada do sistema (admin) é tão importante quanto validar dados de usuário.

### Migrations Acumuladas
O projeto chegou à migration `0009` por causa de ajustes iterativos no schema. Migrações `0006` e `0007` foram necessárias para compatibilizar bancos criados antes. Isso adiciona complexidade operacional.

**Lição**: Planejar o schema do banco antes de iniciar o desenvolvimento reduz o número de migrations corretivas.

### Sessões e Testes de Middleware
Testar `RequireAuth` exigiu montar um `AuthService` completo com sessão fake. Isso revelou que o middleware estava acoplado ao `AuthService` concreto em vez de uma interface.

**Lição**: Middlewares de autenticação se beneficiam de interfaces pequenas para testabilidade.

---

## O que faria diferente

1. **Definir todas as interfaces de repositório no início** — não apenas `VehicleRepository`, `QuestionRepository` etc., mas também `AdminRepository`.

2. **Planejar o schema do banco antes do código** — evitaria as migrations corretivas 0006-0007.

3. **Escrever testes de handlers (T142) junto com os handlers (Fase 3/4)** em vez de acumular para a Fase 6. Tests-as-you-go evita o "debt de testes".

4. **Separar a configuração do frontend** — o JS de `recommend.js` cresceu e se tornou difícil de testar. Um build step simples (ou Vite) com módulos separados seria mais manutenível.

---

## Melhorias para Fase 2+ (Pós-TCC)

| Funcionalidade | Motivação |
|----------------|-----------|
| Integração ChatGPT | Gerar explicações de recomendação em linguagem natural |
| Feedback do usuário | Coletar avaliação das recomendações para melhorar o algoritmo |
| Redis para sessões | Permitir deployment multi-instância |
| Pipeline CI/CD | GitHub Actions para testes + build automático |
| Analytics de uso | Identificar perguntas mais influentes no ranking |
| Internacionalização | Suporte a inglês para mercado fora do Brasil |
| PWA / App Mobile | Acesso offline e notificações push |
