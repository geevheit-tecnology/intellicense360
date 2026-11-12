# Preparação para entrega ao cliente — Geevheit Intelligence 360

Diagnóstico e critérios detalhados: [docs/DIAGNOSTICO_ENTREGA.md](docs/DIAGNOSTICO_ENTREGA.md) | Continuidade: [docs/CONTINUIDADE.md](docs/CONTINUIDADE.md) | Diretrizes de Negócio: [services/split01.md](services/split01.md).

---

## 📌 Visão de Entrega (Transportadora Autônoma & Segurança Operacional)
O sistema deve ser o coração da gestão de frota e motoristas, operando com máxima autonomia para empresas enxutas (poucos funcionários), prevenindo falhas, garantindo conformidade regulatória e fornecendo ao gestor a tranquilidade de que os ativos estão seguros, controlados e rentáveis.

---

## Fase 1 — Fundação, Identidade e Cadastros Base em PostgreSQL (Concluída)
- [x] Revisar estrutura, módulos, composição da API e arquitetura monolítica modular.
- [x] Identificar bloqueios técnicos e organizar critérios objetivos de entrega.
- [x] Implementar persistência PostgreSQL de identidade (empresas, usuários, perfis, sessões, refresh token, auditoria).
- [x] Implementar cadastro de frota em PostgreSQL (veículos, marcas, modelos, categorias, tipos e ativos vinculados).
- [x] Validar migrations de identidade/frota (000001, 000002, 000003, 000020) e transações em PostgreSQL descartável local.
- [x] Implementar módulo de motoristas em PostgreSQL (status, contato, documento fiscal, vínculo atual com veículo, soft delete).
- [x] Implementar metadados de documentos em PostgreSQL (CNH, CRLV, seguros, cronograma UTC, regras de vencimento em 30 dias).
- [x] Validar permissões (`drivers.*`, `documents.*`), migração 000021 e auditoria transacional com suíte de testes passando.

---

## Fase 2 — Núcleo Operacional da Transportadora (Persistência & Regras de Negócio)
*Prioridade imediata para ligar a operação da transportadora a dados persistentes em PostgreSQL.*

### 2.1. Hodômetro & Histórico Imutável de KM (Bloqueador do fluxo operacional)
- [x] Criar migration `000022_odometer_km_history.sql` com tabela de leituras imutáveis de KM por veículo/empresa.
- [x] Implementar regra anti-regressão de hodômetro (KM novo >= KM anterior, com fluxo auditado explícito para substituição de painel).
- [x] Implementar repositório PostgreSQL, camada de aplicação e rotas HTTP em `internal/modules/fleet`.
- [x] Validar concorrência, idempotência e isolamento multi-tenant com testes de unidade e regras de integridade.

### 2.2. Gestão de Abastecimento (Combustível, Requisições & Conciliação Fiscal)
- [x] Ativar e adaptar migrations de combustível (`000013_fuel_core.sql` e `000023_fuel_operational_financial_alignment.sql`) para o runner transacional.
- [x] Implementar repositório PostgreSQL durável para transações de combustível, postos, tanques, bombas, preços e comprovantes.
- [x] Implementar suporte operacional completo para requisições da frota (`requisition_number`) e lançamento manual de terceiros (`is_third_party`, `third_party_identifier`).
- [x] Integrar metadados fiscais para conciliação financeira de notas e caixa (`invoice_number`, `invoice_key`, `payment_status`).
- [x] Validar auditoria transacional, detecção de divergência financeira e testes de ciclo de vida com suíte passando.

### 2.3. Gestão de Pneus (Tires & Ciclo de Vida)
- [x] Ativar e adaptar migrations de pneus (`000006_tires_core.sql`, `000012_tire_lifecycle_core.sql` e `000024_tires_lifecycle_alignment.sql`).
- [x] Implementar repositório PostgreSQL durável para pneus, inspeções de sulco/pressão, movimentações (instalação, rodízio, recape, descarte) e auditoria de histórico.
- [x] Implementar histórico de sulco (medidas de profundidade) com validação de limites de segurança.
- [x] Conectar movimentação de pneus ao hodômetro do veículo (`vehicle_id`, posição no eixo/estepe e KM atual do pneu).

### 2.4. Manutenção Preventiva e Corretiva
- [x] Ativar e adaptar migrations de manutenção (`000008_maintenance_core.sql`, `000009_maintenance_prompt_alignment.sql` e `000025_maintenance_operational_alignment.sql`).
- [x] Implementar repositório PostgreSQL durável para ordens de serviço (O.S.), tipos de serviço, oficinas, apontamentos, técnicos e histórico em `internal/modules/maintenance/infrastructure/postgres`.
- [x] Implementar planos preventivos automáticos acionados por KM ou tempo calendário e rastreio de tempo de imobilização (downtime).
- [x] Integrar metadados financeiros e suporte a terceiros (`invoice_number`, `invoice_key`, `total_cost`, `payment_status`, `is_third_party`, `third_party_identifier`) para conciliação contábil/caixa.
- [x] Conectar rotas `/api/v1/maintenance/*` no roteador PostgreSQL em `container.go` e `router.go` com suíte de testes passando.

### 2.5. Checklist Operacional & Não-Conformidades
- [x] Ativar e adaptar migrations de checklist (`000005_checklist_core.sql`, `000014_checklist_engine_core.sql` e `000026_checklist_operational_alignment.sql`).
- [x] Implementar repositório PostgreSQL durável com todas as 14 interfaces de repositório (templates, versões, tipos, seções, itens, execuções de inspeção no pátio, respostas, evidências, não-conformidades, assinaturas e histórico) em `internal/modules/checklist/infrastructure/postgres`.
- [x] Configurar integridade de veículo (`vehicle_id`) e classificação de gravidade de não-conformidades com auditoria de conformidade.
- [x] Conectar rotas `/api/v1/checklists/*` no roteador em `container.go` e `router.go` com suíte de testes passando.

---

## Fase 3 — Automação, Inteligência & "Tranquilidade do Gestor"
*Foco na redução de trabalho manual para empresas com poucos funcionários e monitoramento 24/7.*

- [x] Corrigir bloqueio sintático da migration histórica `000018_intelligence_core.sql` (palavra reservada `"window"` com escape no PostgreSQL) e criar migration de integridade `000027_mission_control_alignment.sql`.
- [x] Ativar migrations de mensageria e controle de missão (`000017_event_bus_outbox.sql`, `000018_intelligence_core.sql`, `000019_mission_control_core.sql`, `000027_mission_control_alignment.sql`).
- [x] Implementar repositório PostgreSQL para Mission Control (`CommandItem`, `CommandEvent`, `CommandAction`, `OperationalSnapshot`, `Idempotency`) em `internal/modules/missioncontrol/infrastructure/postgres`.
- [x] Implementar motor proativo de avaliação de sinais (`SignalScanner`):
  - [x] Alerta de CNH e documentos vencidos ou a vencer em até 30 dias (`drivers_documents`).
  - [x] Alerta de limite de segurança de sulco de pneu atingido (`tires` onde `current_tread_mm <= minimum_tread_mm`).
  - [x] Alerta de manutenção preventiva próxima ou vencida por KM acumulado (`maintenance_preventive_plans` vs `fleet_odometer_readings`).
  - [x] Alerta de não-conformidades críticas não resolvidas de checklist operacional.
- [x] Expor serviço e rota `POST /api/v1/mission-control/evaluate-signals` com testes unitários 100% passando.

---

## Fase 4 — Interface Web do Gestor & Operação Simplificada
*O painel onde o gestor visualiza e opera toda a transportadora seguindo o Design System.*

- [x] Estruturar aplicação web do gestor em `apps/mission-control` (HTML5, Vanilla JS SPA, tema escuro nativo `#111827`, WCAG 2.2 AA).
- [x] Desenvolver cliente de API (`src/api/client.js`) com gerenciamento de sessão, autenticação e integração a todos os módulos da API.
- [x] Construir telas completas de administração e operação (`src/app.js`):
  - [x] Dashboard Executivo com painel de alertas proativos de missão ("Tranquilidade do Gestor") e botão para disparo de avaliação de sinais.
  - [x] Gestão de Veículos e Hodômetro/KM (com validação anti-regressão e fluxo de correção assistido).
  - [x] Gestão de Motoristas e Conformidade (CNH, status operacional, monitor de vencimentos).
  - [x] Abastecimento & Conciliação Fiscal (número de requisição, suporte frota própria vs. terceiros, número de NF, chave e status de pagamento).
  - [x] Painel de Pneus & Ciclo de Vida (medidas de sulco mm, limites mínimos e inspeções).
  - [x] Ordens de Serviço de Manutenção (tipo, oficina, tempo de imobilização, valores totais e dados fiscais).
  - [x] Checklists Operacionais (inspeções de pátio na saída/retorno e acompanhamento de não-conformidades).
- [x] Interface responsiva para uso no desktop e dispositivos móveis no pátio ou posto.

---

## Fase 5 — Consolidação de Contratos, Segurança & Qualidade
- [x] Sincronizar contratos OpenAPI (`backend/api/openapi.yaml` e `backend/api/openapi/openapi.yaml`) com novos endpoints de hodômetro, abastecimento/conciliação, manutenção, mission control e financeiro/fluxo de caixa.
- [x] Validação da suíte de testes Go (100% dos testes unitários passando em todos os pacotes).
- [ ] Executar teste de integração ponta a ponta com Docker e PostgreSQL descartável local quando autorizado.

---

## Fase 6 — Homologação, Dados do Cliente & Deploy em Produção
- [x] Desenvolver ferramenta CLI de importação inicial de dados legados do cliente (`backend/api/cmd/import/main.go`) com suporte a veículos, leitura inicial de hodômetro e motoristas a partir de CSV.
- [x] Configurar empacotamento com Docker multi-stage para backend (`backend/api/Dockerfile`) e frontend (`apps/mission-control/Dockerfile`).
- [ ] Preparar ambiente de homologação e execução de testes com equipe do cliente quando autorizado.

---

## Fase 7 — Módulo Financeiro, Contas a Pagar & Fechamento de Caixa
*Integração contábil e de caixa para que nenhuma nota fiscal, boleto, requisição de combustível ou O.S. de manutenção (frota própria ou terceiros) fique perdida.*

- [x] Ativar migrations de financeiro (`000015_financial_cost_core.sql` e `000028_financial_cashflow_reconciliation.sql`) no runtime runner.
- [x] Implementar repositório PostgreSQL durável em `internal/modules/financial/infrastructure/postgres/store.go` cobrindo transações, categorias, tipos, centros de custo, contas, métodos de pagamento, períodos fiscais, orçamentos, ajustes e fechamentos de caixa.
- [x] Implementar conciliação operacional direta de despesas de abastecimento e manutenção (`ReconcileOperationalExpense`) com suporte a frota própria x terceiros (`is_third_party`, `third_party_identifier`), número de requisição, número de nota, chave NF-e e código de barras de boletos.
- [x] Implementar motor de cálculo de fluxo de caixa operacional (`GetCashFlowSummary`), separando entradas, custos de frota própria, despesas de terceiros e pendências a pagar.
- [x] Implementar rotina auditável de fechamento diário de caixa (`CloseCash`).
- [x] Conectar rotas `/api/v1/financial/*` no modo PostgreSQL em `container.go` e `router.go` com suíte de testes passando 100%.
- [x] Adicionar tela de visualização e operação do Financeiro & Caixa na interface web (`apps/mission-control`).

---

## Fase 8 — Almoxarifado/Estoque de Peças, Carretas & Prontuário de Auditoria (Concluída)
*Atendimento às demandas operacionais de estoque completo (pneus, peças, equipamentos), cadastro de semirreboques/carretas com KM opcional e prontuário histórico auditável com impressão de relatórios.*

- [x] Ativar migrations de estoque (`000010_inventory_core.sql` e `000029_inventory_operational_alignment.sql`) com suporte a custo unitário, estoque mínimo, saldo atual e localização física.
- [x] Implementar repositório PostgreSQL durável para peças, categorias, depósitos e localizações em `internal/modules/inventory/infrastructure/postgres/store.go`.
- [x] Conectar rotas `/api/v1/inventory/*` no modo PostgreSQL no roteador e injeção de dependências (`container.go` e `router.go`).
- [x] Suporte completo a Carretas / Semirreboques / Implementos Rodoviários:
  - [x] Cadastro com identificação de implemento (Baú Frigorífico, Sider, Graneleira, Carga Seca, Prancha, Tanque, Porta-Contêiner, Silo Basculante).
  - [x] Configuração de eixos (2 eixos, 3 eixos convencional, 3 eixos vanderleia, bitrem, rodotrem).
  - [x] Flexibilização de hodômetro (opcional para carretas, obrigatório para cavalos/caminhões).
  - [x] Filtros por categoria (Todos, Cavalos/Caminhões, Carretas/Semirreboques).
- [x] Dossiê / Prontuário Completo de Auditoria do Veículo e Carreta:
  - [x] Histórico de hodômetro auditado com validação anti-regressão.
  - [x] Histórico de abastecimentos com requisições, notas fiscais e médias km/l.
  - [x] Histórico de ordens de serviço (O.S.) com peças aplicadas, downtime e custos.
  - [x] Histórico de pneus instalados e profundidade de sulco (mm).
  - [x] Histórico de checklists de pátio com registro de avarias/não-conformidades e assinaturas digitais.
  - [x] Relatório oficial de auditoria pronto para impressão (`@media print`) e exportação em PDF com termo de encerramento e assinaturas.
- [x] Módulo Web de Estoque & Almoxarifado (`apps/mission-control`):
  - [x] Gestão de todas as categorias operacionais: Pneus (novos, recapados, sucatas), Peças Mecânicas & Freios, Filtros & Óleos/Arla, Elétrica & Baterias, Equipamentos de Bordo & Amarração (catracas, cintas, tacógrafos, extintores), EPIs e Acessórios de Carreta (lonas, engates de ar).
  - [x] Métricas de valor total do estoque, saldo físico e alertas de estoque mínimo para reposição.
  - [x] Modal para cadastro de novos itens e movimentações de entrada (compra/NF) e saída (aplicação em manutenção/veículo).

---

## Fase 9 — Viagens, Fretes & Despacho Autônomo com IA (Concluída)
*O coração operacional da transportadora: despacho de cargas, acoplamento dinâmico Cavalo + Carreta, rotas, adiantamentos, vale-pedágio e validação pré-embarque por IA.*

- [x] Criar migration `000030_trips_dispatch_core.sql` com tabela `operational_trips`, índices e permissões `trips.read` e `trips.manage`.
- [x] Registrar `000030_trips_dispatch_core.sql` no runtime runner (`internal/platform/postgres/migrate.go`).
- [x] Implementar módulo `trips` completo em Go:
  - [x] Entidades e domínio em `internal/modules/trips/domain/entities.go`.
  - [x] Contratos e portas em `internal/modules/trips/ports/repository.go` e `service.go`.
  - [x] Camada de aplicação com pré-embarque por IA e validação anti-regressão de hodômetro (`internal/modules/trips/application/services.go`).
  - [x] Repositório PostgreSQL durável em `internal/modules/trips/infrastructure/postgres/store.go`.
  - [x] Camada de transporte HTTP, DTOs e mappers (`internal/modules/trips/transport/http/handler.go`).
- [x] Injeção de dependências e roteador em `container.go` e `router.go` (rotas `/api/v1/trips/*`).
- [x] Sincronização dos contratos OpenAPI (`openapi.yaml` e `openapi/openapi.yaml`).
- [x] Suíte de testes unitários passando 100%.
- [x] Interface Web de Viagens & Despacho em `apps/mission-control`:
  - [x] Menu de navegação "🗺️ Viagens & Despacho".
  - [x] Métricas em tempo real: Viagens em Trânsito, Faturamento em Fretes, Adiantamentos/Pedágio em Rota e Cargas Entregues.
  - [x] Filtros por status: Todas, Em Trânsito, Planejadas/Pátio, Entregues/Concluídas.
  - [x] Modal de Despacho de Nova Viagem com acoplamento dinâmico (Cavalo + Carreta), condutor, rota, carga, peso, MDF-e, CIOT ANTT e valores financeiros.
  - [x] Autonomous Pre-Flight Check da IA: validação automática de CNH regular, checklist de pátio sem avaria crítica e pneus operacionais.
  - [x] Conclusão e Baixa de Entrega (Canhoto Fiscal) com registro de KM final e cálculo de saldo para o financeiro.
  - [x] Emissão de Manifesto de Viagem com suporte para impressão.

---

## Fase 10 — Controle de Jornada do Motorista & Diárias (Lei 13.103/2015) (Concluída)
*Blindagem trabalhista contra passivos: controle rigoroso de 5h30 de direção contínua, pausas de 30 min, refeição, pernoite de 11h, tempo de espera e diárias.*

- [x] Criar migration `000031_driver_journey_compliance.sql` com tabela `driver_journey_records`, índices e permissões.
- [x] Registrar migration `000031` no runtime runner (`internal/platform/postgres/migrate.go`).
- [x] Implementar domínio e regras de jornada em `internal/modules/drivers`:
  - [x] Entidades e constantes de eventos (`driving`, `break_rest`, `meal`, `daily_rest`, `waiting`, `standby`).
  - [x] Verificação automática de violação da Lei 13.103/2015 (direção contínua superior a 330 min / 5h30 sem pausa de 30 min).
  - [x] Atribuição de diárias de alimentação (R$ 45,00) e pernoite/estadia (R$ 90,00).
  - [x] Resumo operacional da jornada com tempo de volante, descansos, espera e diárias acumuladas.
  - [x] Repositório PostgreSQL com fechamento automático de eventos anteriores ao abrir um novo.
- [x] Expor endpoints REST:
  - [x] `GET /api/v1/drivers/:id/journey` (histórico de apontamentos).
  - [x] `POST /api/v1/drivers/:id/journey` (registro de novo evento).
  - [x] `GET /api/v1/drivers/:id/journey/summary` (estatísticas e conformidade).
- [x] Sincronizar contratos OpenAPI (`openapi.yaml` e `openapi/openapi.yaml`).
- [x] Suíte de testes unitários passando 100% (`journey_service_test.go`).
- [x] Interface Web no Mission Control (`apps/mission-control`):
  - [x] Botão "⏱️ Diário de Bordo & Jornada" na tabela de motoristas.
  - [x] Modal completo de Diário de Bordo com status ativo em tempo real.
  - [x] Painel de métricas (tempo de direção, descansos, espera, diárias e infrações).
  - [x] Formulário de registro rápido de paradas regulamentares.
  - [x] Tabela de auditoria com alerta visual de infrações.
  - [x] Espelho de ponto oficial impresso (`@media print`) com campos de assinatura para MTE e fiscalização.

---

## Fase 11 — Infraestrutura de Produção & DevOps Autônomo (Concluída)
*Ambiente de deploy pronto para produção com máxima segurança e recuperação de desastres.*

- [x] Elaborar documento estratégico de arquitetura em `docs/PLANO_DEPLOY_PRODUCAO_AUTONOMO.md`.
- [x] Configurar stack unificada de produção com `docker-compose.prod.yml`:
  - [x] PostgreSQL 16 Alpine com volume persistente e healthcheck.
  - [x] Redis 7 Alpine com volume persistente e política de appendonly.
  - [x] Container de migração automática no boot (`/app/migrate up`).
  - [x] Backend API Go de alta performance com reinício automático.
  - [x] Frontend Mission Control com Nginx e proxy reverso integrado.
  - [x] Serviço de backup diário automatizado do PostgreSQL com retenção de 30 dias.
- [x] Criar scripts de backup e restauração (`infra/scripts/backup_postgres.sh` e `infra/scripts/restore_postgres.sh`).
- [x] Criar arquivo de variáveis de ambiente de produção documentado (`.env.production.example`).

---

Esta lista registra pendências e orienta o desenvolvimento; não autoriza implantação sem prévio alinhamento de infraestrutura.
Documentos de suporte: [docs/PERSISTENCIA_IDENTIDADE_FROTA.md](docs/PERSISTENCIA_IDENTIDADE_FROTA.md) | [docs/MOTORISTAS_DOCUMENTOS.md](docs/MOTORISTAS_DOCUMENTOS.md).
