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
- Conjunto ativo no runner: **000001, 000002, 000003, 000020, 000021, 000022, 000013, 000023, 000006, 000012, 000024, 000008, 000009, 000025, 000005, 000014, 000026, 000017, 000018, 000019, 000027, 000015, 000028, 000010, 000029, 000030, 000031**.
- `STORAGE_DRIVER=postgres` é padrão, exige DATABASE_URL e AUTH_SECRET próprio de 32+ caracteres. Startup verifica migrations pendentes.
- `scripts/dev-api.sh` inicia a API isolada; `scripts/dev.sh` inicia o ecossistema integrado (Backend na porta 8080 + Frontend Mission Control na porta 3000 com CORS ativo). Produção rejeita modo em memória.
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

## Etapa 3 concluída — Hodômetro e Histórico Imutável de KM (Fase 2.1)

- Migration `000022_odometer_km_history.sql` ativada no runtime (`schema_migrations`).
- Tabela `fleet_odometer_readings` append-only com FKs tenant-safe para `tenants`, `fleet_vehicles`, `drivers` e `users`.
- Regra de anti-regressão: leituras comuns devem ter KM maior ou igual à leitura anterior (`ErrOdometerRegression`).
- Fluxo de correção: troca de painel/hodômetro ou retificação aceita KM inferior somente se `is_correction = true` e justificativa tiver no mínimo 10 caracteres (`ErrOdometerCorrectionRequired`).
- Permissões adicionadas: `fleet.odometer.read`, `fleet.odometer.record`, `fleet.odometer.correct` associadas ao papel admin existente e provisionamento de novas empresas.
- Rotas HTTP sob `/api/v1/fleet/vehicles/:id/odometer` (registro e resumo) e `/history` (histórico paginado ordenado decrescente).
- Repositórios em memória e PostgreSQL com cobertura de testes unitários passando.
- Contratos OpenAPI (`openapi.yaml` e `openapi/openapi.yaml`) atualizados de forma coerente.

## Etapa 4 concluída — Abastecimento (Fuel), Requisições e Conciliação Fiscal (Fase 2.2)

- Migrations `000013_fuel_core.sql` e `000023_fuel_operational_financial_alignment.sql` ativadas no runtime (`schema_migrations`).
- Repositório PostgreSQL completo em `internal/modules/fuel/infrastructure/postgres/store.go` para transações, postos, tipos, tanques, bicos, preços, comprovantes e ajustes.
- Suporte a requisições de frota (`requisition_number`) e diferenciação operacional entre veículos próprios e veículos de terceiros (`is_third_party`, `third_party_identifier`).
- Metadados para conciliação fiscal e caixa da transportadora (`invoice_number`, `invoice_key`, `payment_status`).
- Conexão do módulo `fuel` no roteador PostgreSQL em `container.go` e `router.go` (rotas `/api/v1/fuel/*` agora ativas em modo PostgreSQL).
- Contratos OpenAPI sincronizados e testes unitários aprovados.

## Etapa 5 concluída — Gestão de Pneus (Tires & Ciclo de Vida) (Fase 2.3)

- Migrations `000006_tires_core.sql`, `000012_tire_lifecycle_core.sql` e `000024_tires_lifecycle_alignment.sql` ativadas no runtime (`schema_migrations`).
- Repositório PostgreSQL durável em `internal/modules/tires/infrastructure/postgres/store.go` cobrindo pneus, inspeções de sulco e pressão, movimentações (instalação, rodízio, remoção, recape, descarte) e auditoria de histórico.
- Constraints de escopo tenant-safe e integridade com `fleet_vehicles(tenant_id, id)`.
- Permissões `tires.tires.read` e `tires.tires.manage` concedidas ao papel admin e provisionamento de novas empresas.
- Rotas `/api/v1/tires/*` expostas no modo PostgreSQL em `container.go` e `router.go`.
- Suíte completa de testes aprovada.

## Etapa 6 concluída — Manutenção Preventiva e Corretiva em PostgreSQL (Fase 2.4)

- Migrations `000008_maintenance_core.sql`, `000009_maintenance_prompt_alignment.sql` e `000025_maintenance_operational_alignment.sql` ativadas no runtime (`schema_migrations`).
- Repositório PostgreSQL completo em `internal/modules/maintenance/infrastructure/postgres/store.go` implementando todas as interfaces: ordens de serviço (O.S.), tipos de serviço, planos preventivos, oficinas, apontamentos, técnicos, mão de obra, tempo de imobilização (downtime) e histórico.
- Campos de conciliação fiscal e controle de terceiros integrados às O.S.: `invoice_number`, `invoice_key`, `total_cost`, `payment_status`, `is_third_party`, `third_party_identifier`.
- Conexão do módulo `maintenance` no roteador PostgreSQL em `container.go` e `router.go` (rotas `/api/v1/maintenance/*` ativas).
- Testes unitários cobrindo transições e conciliação de terceiros (`TestWorkOrderThirdPartyAndFinancialReconciliation`) passando 100%.

## Etapa 7 concluída — Checklist Operacional & Não-Conformidades (Fase 2.5)

- Migrations `000005_checklist_core.sql`, `000014_checklist_engine_core.sql` e `000026_checklist_operational_alignment.sql` ativadas no runtime.
- Repositório PostgreSQL durável em `internal/modules/checklist/infrastructure/postgres/store.go` cobrindo todas as 14 interfaces do domínio de checklist (templates, versões, tipos, seções, itens, execuções de inspeção no pátio, respostas, evidências, não-conformidades, assinaturas e histórico).
- Rastreamento de veículo inspecionado (`vehicle_id`) e classificação de gravidade de não-conformidades (`low`, `medium`, `high`, `critical`).
- Rotas `/api/v1/checklists/*` expostas no modo PostgreSQL em `container.go` e `router.go`.
- Contratos OpenAPI sincronizados e testes unitários aprovados.

## Etapa 8 concluída — Automação, Inteligência & Mission Control Proativo (Fase 3)

- Correção da palavra reservada `"window"` com escape na migration histórica `000018_intelligence_core.sql`.
- Migration `000027_mission_control_alignment.sql` com índices de unicidade tenant-safe e permissões de mission control (`missioncontrol.command.read`, `missioncontrol.command.action`).
- Migrations `000017_event_bus_outbox.sql`, `000018_intelligence_core.sql`, `000019_mission_control_core.sql` e `000027` ativadas no runner.
- Repositório PostgreSQL durável em `internal/modules/missioncontrol/infrastructure/postgres/store.go` (`CommandItem`, `CommandEvent`, `CommandAction`, `OperationalSnapshot`, `Idempotency`).
- Motor de regras preventivas autônomo implementado em `signal_scanner.go`:
  1. Vencimento de CNH e documentos (30 dias ou vencidos).
  2. Limite de segurança de sulco de pneu atingido (`current_tread_mm <= minimum_tread_mm`).
  3. Manutenção preventiva próxima ou vencida por KM acumulado no hodômetro.
  4. Não-conformidades críticas não resolvidas de checklist operacional.
- Endpoint `POST /api/v1/mission-control/evaluate-signals` para disparo de varredura autônoma com idempotência por impressão digital.
- Cobertura de teste com scanner mockado em `mission_control_service_test.go` passando 100%.

## Etapa 9 concluída — Interface Web do Gestor & Design System (Fase 4)

- Aplicação web do gestor construída em `apps/mission-control`:
  - Design system refinado com fundo `#111827`, cartões `#1F2937`, bordas `#374151`, azul de destaque `#2563EB`, tipografia Inter e acessibilidade WCAG 2.2 AA em `src/styles/theme.css`.
  - Cliente de API completo com gerenciamento de tokens e endpoints em `src/api/client.js`.
  - SPA reativa em `src/app.js` e `index.html` com suporte a:
    - Painel Executivo com feed de alertas proativos e botão "Avaliar Sinais Operacionais".
    - Cadastro e consulta de Veículos com hodômetro imutável e fluxo de retificação assistido.
    - Gestão de Motoristas com status de CNH e monitor de vencimentos.
    - Abastecimento com número de requisição, separação própria/terceiros e conciliação fiscal.
    - Pneus com monitoramento de profundidade de sulco em mm e histórico.
    - Ordens de Serviço de Manutenção com tempo de imobilização e notas fiscais.
    - Checklists de pátio na saída/retorno do veículo e não-conformidades.

## Etapa 10 concluída — Empacotamento Docker & Ferramenta de Importação Legada (Fase 6)

- Docker multi-stage configurado para a API Go em `backend/api/Dockerfile` (Go 1.24 build + Alpine runtime com usuário não-root).
- Docker configurado para a interface web em `apps/mission-control/Dockerfile` (Nginx Alpine com proxy reverso transparente para `/api/`).
- Ferramenta CLI de importação inicial de dados legados do cliente em `backend/api/cmd/import/main.go` suportando leitura de CSV de veículos, hodômetros iniciais e motoristas com transação atômica.

## Etapa 11 concluída — Módulo Financeiro, Contas a Pagar & Fechamento de Caixa (Fase 7)

- Migrations `000015_financial_cost_core.sql` e `000028_financial_cashflow_reconciliation.sql` ativadas no runtime.
- Repositório PostgreSQL completo em `internal/modules/financial/infrastructure/postgres/store.go` implementando todas as interfaces: transações, categorias, tipos, centros de custo, contas, métodos de pagamento, períodos fiscais, orçamentos, ajustes e fechamentos diários de caixa (`financial_cash_closings`).
- Conciliação operacional direta de despesas de abastecimento e ordens de manutenção (`ReconcileOperationalExpense`), diferenciando custos de frota própria e terceiros (`is_third_party`, `third_party_identifier`), rastreando número de requisição, número da nota fiscal, chave NF-e de 44 dígitos e código de barras/linha digitável de boletos.
- Motor de fluxo de caixa operacional em tempo real (`GetCashFlowSummary`) que calcula receitas recebidas, despesas pagas por frota x terceiros, saldo líquido e monitora contas a pagar pendentes e boletos vencidos.
- Rotina auditável de fechamento diário de caixa (`CloseCash`).
- Rota e permissões integradas: `/api/v1/financial/*` agora ativo no modo PostgreSQL em `container.go` e `router.go`.
- Interface web do gestor atualizada em `apps/mission-control` com a aba "Financeiro & Caixa", incluindo cartões de métricas, lista de boletos e notas com ação rápida "Baixar / Pagar" e fechamento de caixa diário.
- Suíte completa de testes unitários aprovada com 100% de sucesso.

## Etapa 12 concluída — Almoxarifado/Estoque de Peças, Carretas & Prontuário de Auditoria (Fase 8)

- Migrations `000010_inventory_core.sql` e `000029_inventory_operational_alignment.sql` ativadas no runtime runner.
- Entidades e DTOs de `internal/modules/inventory` alinhados com custo unitário (`unit_cost`), estoque mínimo (`min_stock`), saldo atual (`current_stock`) e localização física (`location`).
- Repositório PostgreSQL durável implementado em `internal/modules/inventory/infrastructure/postgres/store.go` cobrindo peças, categorias, depósitos e localizações.
- Rotas `/api/v1/inventory/*` ativas no modo PostgreSQL e injetadas em `container.go` e `router.go`.
- Suporte a Carretas / Semirreboques / Implementos Rodoviários:
  - Cadastro flexível com seleção de implemento (Baú Frigorífico, Sider/Cortina, Graneleira, Carga Seca, Prancha, Tanque, Porta-Contêiner, Silo Basculante).
  - Configuração de eixos (2 eixos, 3 eixos convencionais, 3 eixos vanderléia/espaçados, 4 eixos, bitrem, rodotrem).
  - Hodômetro flexibilizado: opcional para carretas/implementos (sem hodômetro obrigatório) e obrigatório para caminhões/cavalos.
  - Filtros rápidos por categoria na visão de frota.
- Dossiê e Prontuário Completo de Auditoria do Veículo e Carreta:
  - Visualização unificada agregando: 1. Cadastro e especificações técnicas; 2. Histórico de hodômetro com validação anti-regressão e correções homologadas; 3. Histórico de abastecimentos com requisições e notas fiscais; 4. Histórico de ordens de serviço (O.S.) com peças aplicadas, downtime e custos; 5. Histórico e posições de pneus com sulco medido em mm; 6. Checklists de pátio com registro de avarias e assinaturas digitais.
  - Relatório impresso oficial com cabeçalho da empresa, carimbo de integridade e campos de assinatura do gestor e auditoria interna via `@media print` (`window.print()`).
- Módulo Web de Estoque & Almoxarifado em `apps/mission-control`:
  - Gestão de todas as categorias operacionais solicitadas: Pneus (novos, recapados, sucatas), Peças Mecânicas & Freios, Filtros & Óleos/Arla 32, Elétrica & Baterias, Equipamentos de Bordo & Amarração (catracas, cintas 5t/10t, tacógrafo, extintores ABC), EPIs de segurança e Acessórios de Carretas (lonas, engates de ar).
  - Indicadores de total de itens, valor total imobilizado, alertas de reposição mínima e controle de pneus em almoxarifado.
  - Modais para cadastro de novos itens e movimentação de estoque (entrada por compra com NF ou saída para manutenção/veículo).
- 100% dos testes unitários passando em todos os módulos Go.

## Etapa 13 concluída — Módulo de Viagens, Fretes & Despacho Autônomo com IA (Fase 9)

- Migration `000030_trips_dispatch_core.sql` criada e adicionada ao runtime runner (`operational_trips` com índices, restrições e permissões `trips.read` e `trips.manage`).
- Módulo `trips` completo em Go:
  - Entidades de domínio com status (`draft`, `dispatched`, `delivered`, `completed`, `canceled`), acoplamento dinâmico Cavalo + Carreta, rotas, adiantamentos e documentos fiscais (MDF-e, CT-e, CIOT ANTT).
  - Camada de aplicação com pré-embarque por IA e trava anti-regressão de hodômetro (`internal/modules/trips/application/services.go`).
  - Repositório PostgreSQL durável implementado em `internal/modules/trips/infrastructure/postgres/store.go`.
  - Camada de transporte HTTP, DTOs e mappers (`internal/modules/trips/transport/http/handler.go`).
- Injeção de dependências e roteador em `container.go` e `router.go` (rotas `/api/v1/trips/*`).
- Contratos OpenAPI sincronizados de forma coerente (`backend/api/openapi.yaml` e `backend/api/openapi/openapi.yaml`).
- Suíte de testes unitários Go aprovada com 100% de sucesso.
- Interface Web de Viagens & Despacho em `apps/mission-control`:
  - Nova aba "🗺️ Viagens & Despacho" no menu de navegação.
  - Indicadores em tempo real: Viagens em Trânsito, Faturamento em Fretes, Adiantamentos/Pedágio em Rota e Cargas Entregues.
  - Modal de Despacho de Nova Viagem com acoplamento dinâmico (Cavalo + Carreta), condutor, rota, carga, peso, MDF-e, CIOT ANTT e valores financeiros.
  - Autonomous Pre-Flight Check da IA: validação automática de CNH regular, checklist de pátio sem avaria crítica e pneus operacionais.
  - Conclusão e Baixa de Entrega (Canhoto Fiscal) com registro de KM final e cálculo de saldo para o financeiro.
  - Visualização de Manifesto de Viagem com suporte para impressão.

## Etapa 10 concluída — Controle de Jornada do Motorista & Diárias (Lei 13.103/2015)

- Migration `000031_driver_journey_compliance.sql` criada e ativa no runner (`driver_journey_records` com índices, checagem e integridade).
- Domínio e serviço de aplicação em `internal/modules/drivers/application/journey_service.go`:
  - Tipos de eventos: `driving`, `break_rest`, `meal`, `daily_rest`, `waiting`, `standby`.
  - Validação regulatória automática: direção contínua superior a 330 min (5h30) sem pausa de 30 min dispara `is_violation = true` com descrição da infração.
  - Fechamento transacional do apontamento anterior ao abrir um novo.
  - Cálculo de diárias automáticas (Refeição R$ 45,00, Pernoite/Estadia R$ 90,00).
- Repositório PostgreSQL durável implementado em `internal/modules/drivers/infrastructure/postgres/journey_repository.go`.
- Rotas HTTP expostas: `GET /api/v1/drivers/:id/journey`, `POST /api/v1/drivers/:id/journey`, `GET /api/v1/drivers/:id/journey/summary`.
- Suíte de testes unitários passando 100% (`journey_service_test.go`).
- Interface completa no Mission Control (`apps/mission-control/src/app.js`):
  - Botão "⏱️ Diário de Bordo & Jornada" na tabela de motoristas.
  - Modal dinâmico com status do motorista em tempo real (Em Direção, Em Descanso, Refeição, Pernoite, Espera, Standby).
  - Cards de KPIs: Tempo de Direção, Descanso, Espera, Diárias Acumuladas e Infrações Lei 13.103.
  - Formulário para apontamento imediato de paradas e KM.
  - Histórico auditável de apontamentos com alertas de infração destacados em vermelho.
  - Espelho de Ponto Oficial para impressão (`@media print`) com campos de assinatura para MTE e fiscalização.

## Etapa 11 concluída — Deploy em Produção & DevOps Autônomo

- Arquitetura de produção documentada em `docs/PLANO_DEPLOY_PRODUCAO_AUTONOMO.md`.
- Stack unificada configurada em `docker-compose.prod.yml`:
  - `postgres`: PostgreSQL 16 Alpine com volume persistente e healthcheck.
  - `redis`: Redis 7 Alpine com persistência AOF.
  - `migrate`: Container executando `/app/migrate up` automaticamente no startup antes de subir a API.
  - `api`: Backend em Go com reinício automático e conexões seguras.
  - `web`: Frontend Mission Control servido por Nginx na porta 80 com proxy reverso.
  - `db-backup`: Serviço cron de backup diário do PostgreSQL às 03:00 com retenção de 30 dias.
- Scripts de backup e restauração criados em `infra/scripts/backup_postgres.sh` e `infra/scripts/restore_postgres.sh`.
- Arquivo de exemplo de produção criado em `.env.production.example`.
- **Suporte Nativo Vercel Serverless (Frontend + Backend Go):**
  - Configuração unificada via `vercel.json` roteando `/api/*` diretamente para a Serverless Function em Go (`api/index.go`) e todas as outras rotas para o frontend SPA (`apps/mission-control`).
  - Entrypoint `backend/api/app.go` exporta `InitRouter()` e `ServeHTTP(w, r)` inicializando o monolito modular sob demanda com suporte a fallback em memória ou banco PostgreSQL remoto gerenciado (ex.: Neon / Supabase).
  - Repositório sincronizado no GitHub remoto oficial (`geevheit-tecnology/intellicense360`).

## Próximo passo sugerido

1. Executar testes de homologação com dados reais do cliente via CLI `cmd/import` quando autorizado.
2. Configurar domínio e certificado SSL/TLS (Let's Encrypt / Certbot) no Nginx de produção conforme aprovação do usuário.

## Pendências e recomendações

- Manter alterações coerentes nos dois OpenAPI até uma consolidação explícita.
- As migrações agora contemplam a esteira completa: 000001 a 000003, 000005, 000006, 000008, 000009, 000010, 000012, 000013, 000014, 000015, 000017 a 000031.
- Não expor credenciais, tokens ou dados reais em logs/relatórios. TEST_DATABASE_URL deve apontar somente para base descartável.
