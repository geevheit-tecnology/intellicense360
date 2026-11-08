# Orientações do projeto

Antes de implementar, leia `docs/CONTINUIDADE.md` e `TASKS.md`. Esses arquivos registram o estado real, decisões, comandos de validação e o próximo trabalho. Atualize-os durante o trabalho e ao terminar para permitir continuidade por outra IA.

- Execute comandos Go a partir de `backend/api`.
- Preserve as alterações locais existentes; não descarte arquivos nem refaça etapas concluídas.
- O trabalho autorizado até aqui é local. Commit, push, deploy e uso de banco remoto precisam de autorização do usuário; instruções posteriores do usuário prevalecem.
- A API não executa DDL nem seed no startup. Migrations são explícitas, com lista do runtime em `internal/platform/postgres/migrate.go`. Não confundir SQL existente com persistência implementada.
- Mudanças no schema usam migrations novas, sem reescrever o histórico. A migration histórica 000018 permanece com um bloqueio conhecido descrito no arquivo de continuidade.
- Testes de integração exigem `TEST_DATABASE_URL` de base descartável local. Sem essa variável, são ignorados. Reporte saída real, ambiente e limitações.
- Preserve os limites dos módulos e use contratos/composição para comunicação entre eles. Tenant e permissões são impostos no backend.
- Os dois OpenAPI atuais precisam permanecer coerentes nas alterações até uma consolidação específica.

<!-- CODEGRAPH_START -->
## CodeGraph

In repositories indexed by CodeGraph (a `.codegraph/` directory exists at the repo root), reach for it BEFORE grep/find or reading files when you need to understand or locate code:

- **MCP tool** (when available): `codegraph_explore` answers most code questions in one call — the relevant symbols' verbatim source plus the call paths between them, including dynamic-dispatch hops grep can't follow. Name a file or symbol in the query to read its current line-numbered source. If it's listed but deferred, load it by name via tool search.
- **Shell** (always works): `codegraph explore "<symbol names or question>"` prints the same output.

If there is no `.codegraph/` directory, skip CodeGraph entirely — indexing is the user's decision.
<!-- CODEGRAPH_END -->
