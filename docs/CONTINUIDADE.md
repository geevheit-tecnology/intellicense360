# Continuidade do projeto

Ponto de entrada para outra IA. Estado atualizado em **07/11/2026**, após concluir a segunda etapa local. Leia este arquivo, `AGENTS.md` e `TASKS.md` antes de alterar o projeto.

## Autorização e estado do workspace

O usuário autorizou continuar a implementação local e pediu explicitamente manter o caminho de continuidade aberto. As duas etapas abaixo estão implementadas e validadas. O usuário autorizou salvar essas alterações em um commit local. Push, deploy, serviços pagos, banco remoto e envio de mensagens não foram autorizados.

As alterações das duas etapas e este guia fazem parte do commit local de persistência, motoristas e documentos. Consulte `git log -1 --oneline` e `git status --short` para identificar o commit e eventuais mudanças posteriores. Preserve o histórico e as alterações locais; não faça reset, limpeza de arquivos não rastreados, checkout sobre arquivos modificados nem recrie o projeto. Os geradores temporários usados durante a execução não são necessários: os arquivos Go e SQL são fontes normais, completos.

## Arquitetura e comandos

- Go/Gin, monólito modular em `backend/api/internal/modules`.
- Camadas `domain`, `application`, `ports`, `infrastructure` e `transport`.
- Composição em `backend/api/internal/core/di/container.go` e `operational_adapters.go`.
- Executar comandos Go em `backend/api`, não na raiz do monorepo.
- Não existia `.codegraph` nesta revisão; se existir futuramente, usar CodeGraph antes de procurar código.
- Migrations ficam em `database/migrations`, na raiz do repositório.
- API não faz DDL/seed no startup. `go run ./cmd/migrate` aplica somente a lista explícita em `internal/platform/postgres/migrate.go`, com checksum, lock e transação.
- Conjunto ativo: **000001, 000002, 000003, 000020 e 000021**. Os demais SQL não estão automaticamente homologados ou ativados.
- `STORAGE_DRIVER=postgres` é padrão, exige DATABASE_URL e AUTH_SECRET próprio de 32+ caracteres. Startup verifica migrations pendentes.
- `scripts/dev-api.sh` seleciona a demonstração anterior em memória; motoristas/documentos novos exigem PostgreSQL. Produção rejeita modo em memória.
- Provisionamento explícito: `go run ./cmd/provision -name 'Empresa' -email 'gestor@example.test'`, com PROVISION_ADMIN_PASSWORD definido de maneira privada no ambiente. Cada chamada cria outra empresa; não repetir para atualizar/restaurar uma empresa existente.

## Etapa 1 concluída — identidade e frota

Persistência de empresas, usuários, papéis, permissões, vínculos, sessões/refresh, repositório de auditoria, veículos e seus catálogos. UUID, bcrypt, provisionamento transacional, isolamento por empresa e controle otimista nos cadastros de frota. Logout e usuário/empresa inativos invalidam acesso. Refresh é transacional e de uso único.

Guia: `docs/PERSISTENCIA_IDENTIDADE_FROTA.md`.

## Etapa 2 concluída — motoristas e documentos

- Motoristas: cadastro, status, contato, identificação fiscal opcional, vínculo atual opcional com veículo, paginação/filtros, atualização com versão e exclusão lógica.
- Documentos: metadados de CNH/CRLV/seguro/outros para exatamente um motorista ou veículo, datas opcionais, consulta de vencimentos, renovação/arquivamento com versão e exclusão lógica.
- Proprietário de documento é imutável. Documentos históricos permanecem consultáveis e editáveis/arquiváveis após exclusão lógica do proprietário. Nova criação exige proprietário existente e acessível.
- Datas são calendários UTC: vencido antes de hoje; atenção de hoje a 30 dias inclusive; sem data significa `no_expiry`, sem comprovar regularidade externa.
- Rotas `GET/POST /api/v1/drivers`, `GET/PUT/DELETE /api/v1/drivers/:id` e equivalentes de `/documents`. DELETE exige `?version=N`.
- Permissões separadas `drivers.read`, `drivers.manage`, `documents.read`, `documents.manage`.
- Migration 000021 cria tabelas/índices/FKs e atualiza somente papéis admin existentes que já gerenciavam frota. Novos provisionamentos recebem as quatro permissões. Nenhum usuário/senha é recriado para atualizar uma empresa.
- Create/update/delete e auditoria estão na mesma transação. Auditoria registra ator, vínculos e datas/situações anteriores e atuais, sem duplicar identificadores fiscais/números/documentos/contatos/observações nos metadados.
- Não foi implementado upload/download de arquivos, consulta externa, notificação, viagem ou reserva exclusiva de veículo.

Guia completo e evidências: `docs/MOTORISTAS_DOCUMENTOS.md`.

Arquivos principais desta etapa:

- `backend/api/internal/modules/drivers/` e `documents/`.
- `backend/api/internal/core/di/operational_adapters.go` e composição/roteador.
- `database/migrations/000021_drivers_documents_core.sql` e lista do runner.
- `backend/api/internal/modules/identity/infrastructure/postgres/provision.go`.
- `backend/api/internal/integration/operational_test.go` e fixture em `postgres_test.go`.
- `backend/api/openapi.yaml` e `backend/api/openapi/openapi.yaml`.

## Validação efetivamente executada

PostgreSQL 16 Docker local, dados sintéticos, schema isolado por teste. Os testes de integração são **ignorados** sem TEST_DATABASE_URL; não interpretar um `ok` sem essa variável como validação do banco.

1. Suíte completa: `TEST_DATABASE_URL=... go test -race ./... -count=1 -timeout=120s` — **exit code 0**, integração em **46.858s**.
2. Após o ajuste final dos documentos históricos/auditoria: `TEST_DATABASE_URL=... go test -race ./internal/modules/documents/... ./internal/integration -count=1 -timeout=120s` — **exit code 0**, integração em **48.273s**.
3. Os oito testes de integração cobrem migrations/replay, atualização de empresa já provisionada, persistência após recriar pool/roteador, CRUD, autorização, isolamento/FKs, vencimentos, controle de versão, duplicatas concorrentes, refresh/logout e rollback da alteração quando a auditoria falha.
4. Parser YAML e resolução das referências locais dos dois contratos passaram (165 e 237 referências). Uma referência antiga inexistente `PathID` foi corrigida para `ID` já definido. Validação completa por ferramenta OpenAPI ainda não foi feita.
5. `git diff --check` passou.

Não houve migração remota, validação de interface, backup/restauração ou produção. O container `gi360-operational-test` usado nesta etapa foi removido após os testes. Não há API iniciada por esta execução; não assumir banco/aplicação local rodando.

Para repetir do zero, com Docker disponível, na raiz do repositório:

```sh
docker run --name gi360-operational-test \
  -e POSTGRES_HOST_AUTH_METHOD=trust -e POSTGRES_DB=gi360_validation \
  -p 127.0.0.1:55439:5432 -d postgres:16-alpine
cd backend/api
TEST_DATABASE_URL='postgres://postgres@127.0.0.1:55439/gi360_validation?sslmode=disable' \
  go test -race ./... -count=1 -timeout=120s
```

Esse container é exclusivamente descartável e local. Aguarde o banco ficar pronto antes dos testes. Remova apenas o container que criou ao terminar. Não reutilize esses parâmetros para dados de cliente ou produção.

## Próximo passo sugerido

**Histórico de KM/hodômetro por veículo**, antes de ligar abastecimentos e manutenção ao mesmo histórico:

1. Revisar Fleet/Fuel/Maintenance e os contratos existentes, sem redesenhar os módulos.
2. Definir registro imutável de leitura, veículo/empresa/ator/data/origem/unidade e motivo de correção. Não misturar KM com horímetro.
3. Criar migration nova (a partir de 000022), repository/application/API e auditoria transacional. Manter as migrations atuais intactas.
4. Tratar concorrência, idempotência, leituras anteriores e correções explicitamente para evitar regressão silenciosa do hodômetro.
5. Testar persistência, isolamento, duplicidade e concorrência em PostgreSQL descartável; atualizar TASKS, OpenAPI e este arquivo.

Isso é a orientação de continuidade, não uma declaração de que KM já foi implementado. Interface e persistência dos fluxos de abastecimento/manutenção/checklist ainda são tarefas posteriores.

## Pendências que não podem ser esquecidas

- Migration histórica `000018_intelligence_core.sql` falha por usar `window` sem tratar a palavra reservada. Não faz parte do conjunto ativo; não afirmar que todo o schema histórico funciona.
- Abastecimento, manutenção, checklist, financeiro, inteligência e outros módulos ainda têm adapters em memória. Suas rotas são desativadas no modo PostgreSQL até ganharem persistência.
- Falta fechar o escopo inicial com o cliente. Autorizações incrementais de implementação não substituem aceite do cliente.
- Falta painel/interface, importação, viagens/escala, upload de documentos, notificações, gestão de papéis completa, recuperação de senha, auditoria dos demais módulos, homologação, monitoramento e backup/restauração.
- Dois contratos OpenAPI existem; manter alterações coerentes nos dois até consolidar um contrato oficial numa etapa própria.
- Não expor credenciais, tokens ou dados reais em logs/relatórios. TEST_DATABASE_URL deve apontar somente para base descartável.
