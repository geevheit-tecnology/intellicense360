# Sprint 017 — Mission Control Backend & Operational Command Center

## 1. Objetivo

Implementar a **Sprint 017 — Mission Control Backend & Operational Command Center** no backend do projeto **Geevheit Intelligence 360°**, criando a camada de comando operacional do sistema.

O Mission Control deve funcionar como um **Command Center determinístico**, responsável por consolidar estado operacional, prioridades, riscos, alertas, pendências, indicadores e ações recomendadas em **read models próprios**, sem assumir responsabilidade pelos módulos operacionais.

### Regra fundamental desta Sprint

> **Mission Control NÃO deve integrar operacionalmente os módulos nesta Sprint.**

Ele deve preparar a arquitetura para futuras integrações, mas permanecer desacoplado.

Não implementar:

* integração direta com Fleet;
* integração direta com Drivers;
* integração direta com Maintenance;
* integração direta com Tires;
* integração direta com Inventory;
* integração direta com Fuel;
* integração direta com Suppliers;
* integração direta com Financial;
* integração direta com CIOT;
* integração direta com Checklist;
* integração direta com Intelligence;
* Flutter;
* frontend;
* deploy;
* cloud infrastructure;
* chamadas HTTP entre módulos;
* acesso direto às tabelas de outros bounded contexts.

O objetivo é construir o **núcleo backend do Command Center**, deixando os adapters/integration ports preparados para uma próxima etapa.

---

# 2. Contexto arquitetural

O projeto utiliza arquitetura modular baseada em:

```text
domain
application
infrastructure
transport
ports
```

O Mission Control deve respeitar essa arquitetura.

Estrutura conceitual:

```text
backend/api/internal/modules/missioncontrol/
├── domain/
├── application/
├── infrastructure/
├── transport/
└── ports/
```

Não criar dependência estrutural do Mission Control sobre módulos operacionais.

### Regra de dependência

```text
Mission Control
      │
      ├── Domain
      ├── Application
      ├── Infrastructure
      ├── Transport
      └── Ports
            │
            └── contratos abstratos
```

Não permitir:

```text
Mission Control → Fleet implementation
Mission Control → Tire repository
Mission Control → Maintenance repository
Mission Control → Fuel repository
Mission Control → Financial repository
```

A comunicação futura deverá ocorrer por contratos, eventos ou read models.

---

# 3. Resultado esperado

Ao final da Sprint 017, o backend deverá possuir um **Mission Control funcional**, capaz de representar:

* estado operacional consolidado;
* indicadores;
* prioridades;
* alertas;
* riscos;
* oportunidades;
* incidentes;
* pendências;
* ações recomendadas;
* command items;
* severidade;
* criticidade;
* confiança;
* origem;
* status;
* SLA;
* timestamps;
* auditoria;
* deduplicação;
* idempotência.

Tudo deverá existir de forma determinística e testável.

---

# 4. Domínio

Criar as entidades/value objects necessários para o Mission Control.

## 4.1 CommandCenter

Representa uma visão lógica do centro de comando.

Campos sugeridos:

```text
ID
TenantID
Name
Status
LastCalculatedAt
CreatedAt
UpdatedAt
```

Status:

```text
active
paused
archived
```

---

# 5. CommandItem

Criar a entidade principal do Mission Control:

```text
CommandItem
```

Representa qualquer item operacional que precise de atenção.

Campos:

```text
ID
TenantID
Type
Category
Severity
Priority
Status
Title
Description
Source
SourceType
SourceID
Confidence
ImpactScore
RiskScore
UrgencyScore
DueAt
DetectedAt
AcknowledgedAt
ResolvedAt
AssignedTo
Fingerprint
Metadata
CreatedAt
UpdatedAt
```

---

# 6. Tipos de CommandItem

Implementar tipos determinísticos:

```text
alert
risk
incident
opportunity
recommendation
task
anomaly
warning
insight
```

Não permitir tipos arbitrários.

---

# 7. Categorias

Criar categorias extensíveis dentro de enum/controlado:

```text
operational
maintenance
fleet
tire
fuel
inventory
financial
compliance
driver
document
ciot
safety
performance
cost
```

Nesta Sprint essas categorias são apenas classificações.

**Não criar integração operacional com elas.**

---

# 8. Severidade

Implementar:

```text
info
low
medium
high
critical
```

Criar regras determinísticas de ordenação:

```text
critical > high > medium > low > info
```

---

# 9. Prioridade

Implementar:

```text
low
normal
high
urgent
critical
```

A prioridade não deve depender de texto livre.

Criar value object ou enum apropriado.

---

# 10. Status

Implementar lifecycle:

```text
open
acknowledged
in_progress
resolved
dismissed
expired
```

Criar state machine.

Não permitir transições arbitrárias.

Exemplo:

```text
open
 ├── acknowledged
 ├── dismissed
 └── resolved

acknowledged
 ├── in_progress
 ├── dismissed
 └── resolved

in_progress
 ├── resolved
 └── dismissed
```

Transições inválidas devem retornar erro de domínio.

---

# 11. Risk Score

Criar modelo determinístico para representar risco.

Faixa:

```text
0.0 → 1.0
```

ou equivalente inteiro normalizado.

Regras:

```text
0.00–0.20 = very_low
0.21–0.40 = low
0.41–0.60 = medium
0.61–0.80 = high
0.81–1.00 = critical
```

Não usar IA generativa.

Não usar LLM.

Não criar machine learning.

---

# 12. Impact Score

Representar impacto potencial:

```text
0.0 → 1.0
```

Deve ser determinístico.

O domínio deve impedir:

```text
< 0
> 1
```

---

# 13. Confidence

Representar confiança na informação:

```text
0.0 → 1.0
```

Também deve possuir validação de domínio.

---

# 14. Fingerprint / Deduplicação

Criar mecanismo determinístico de deduplicação.

O fingerprint deverá considerar, quando aplicável:

```text
TenantID
Type
Category
SourceType
SourceID
Title/semantic key
```

O objetivo é impedir criação duplicada do mesmo CommandItem.

Criar índice único apropriado.

---

# 15. SLA

Criar estrutura para representar prazo operacional.

Campos:

```text
DueAt
SLAStatus
SLAHours
```

Estados:

```text
within_sla
at_risk
breached
not_applicable
```

O cálculo deve ser determinístico.

---

# 16. CommandAction

Criar entidade para representar uma ação associada ao CommandItem.

Campos:

```text
ID
TenantID
CommandItemID
Type
Label
Status
Priority
Payload
CreatedAt
UpdatedAt
```

Tipos:

```text
acknowledge
assign
start
resolve
dismiss
escalate
review
```

Nesta Sprint as ações devem ser **ações internas de domínio**.

Não executar operações nos módulos externos.

Exemplo:

```text
resolve command item
```

é permitido.

Mas:

```text
resolve maintenance order
```

não é permitido nesta Sprint.

---

# 17. CommandEvent

Criar histórico/auditoria do Mission Control.

Campos:

```text
ID
TenantID
CommandItemID
EventType
PreviousStatus
NewStatus
ActorID
Payload
OccurredAt
```

Eventos:

```text
created
acknowledged
assigned
started
resolved
dismissed
expired
priority_changed
severity_changed
```

O histórico deve ser append-only.

---

# 18. Operational Snapshot

Criar entidade/read model:

```text
OperationalSnapshot
```

Representa uma fotografia do estado do Command Center.

Campos sugeridos:

```text
ID
TenantID
SnapshotAt
OpenItems
CriticalItems
HighPriorityItems
ActiveRisks
ActiveAlerts
OpenIncidents
Opportunities
BreachedSLAs
AverageResolutionTime
OperationalScore
RiskScore
HealthScore
```

Importante:

> Nesta Sprint o snapshot não deve buscar dados diretamente dos módulos operacionais.

Ele será alimentado exclusivamente pelo próprio contexto do Mission Control e por contratos futuros.

---

# 19. Command Center Summary

Criar read model de resumo:

```text
MissionControlSummary
```

Deve permitir obter:

```text
total open
critical
high
medium
risks
alerts
incidents
opportunities
recommendations
breached SLA
operational score
risk score
health score
last updated
```

Todos os valores devem ser derivados de dados persistidos no próprio contexto.

---

# 20. Ranking

Criar serviço de domínio/application responsável por ordenar CommandItems.

Critérios:

1. severity;
2. priority;
3. risk score;
4. impact score;
5. SLA;
6. created/detected time.

O ranking deve ser determinístico.

Para dois itens equivalentes, utilizar:

```text
DetectedAt
ID
```

como desempate estável.

---

# 21. Operational Score

Criar serviço:

```text
OperationalScoreService
```

Responsável por calcular score de forma determinística.

Não utilizar IA.

Exemplo conceitual:

```text
OperationalScore =
    HealthComponent
    - RiskPenalty
    - CriticalPenalty
    - SLAOverduePenalty
```

A fórmula deve ser encapsulada em domínio e possuir testes.

Não espalhar cálculo em handlers ou repositories.

---

# 22. Risk Aggregation

Criar:

```text
RiskAggregationService
```

Responsável por consolidar riscos existentes no Mission Control.

Deve considerar:

```text
risk score
severity
impact
confidence
SLA
```

Gerar:

```text
overall risk score
risk level
top risks
```

Tudo determinístico.

---

# 23. Health Score

Criar:

```text
HealthScoreService
```

Resultado:

```text
0 → 100
```

Faixas:

```text
90–100 = excellent
75–89  = good
60–74  = attention
40–59  = warning
0–39   = critical
```

Criar testes de boundary.

---

# 24. Recommendation Engine

Criar uma estrutura de recomendações internas.

Não implementar IA generativa.

Criar regras determinísticas.

Exemplos:

```text
IF critical risk exists
THEN create recommendation

IF SLA breached
THEN create escalation recommendation

IF repeated anomaly exists
THEN create review recommendation

IF high impact + high risk
THEN create urgent recommendation
```

As regras devem ser isoladas em serviços/policies.

---

# 25. Rule Engine

Criar estrutura:

```text
MissionControlRule
```

Com:

```text
ID
Code
Name
Description
Enabled
Priority
Version
```

Não criar DSL complexa.

Nesta Sprint utilizar regras tipadas em código.

Preparar arquitetura para futuras regras configuráveis.

---

# 26. Idempotência

Todos os comandos críticos devem ser idempotentes.

Especialmente:

```text
create command item
acknowledge
assign
start
resolve
dismiss
create recommendation
create snapshot
```

Não permitir duplicação por retry.

Criar:

```text
idempotency_key
```

quando aplicável.

---

# 27. Multi-tenancy

Todas as entidades devem possuir:

```text
TenantID
```

Nenhuma query pode retornar dados de outro tenant.

Criar testes explícitos:

```text
tenant A cannot read tenant B
tenant A cannot modify tenant B
tenant A cannot resolve tenant B command
tenant A cannot access tenant B snapshot
```

---

# 28. Repository Ports

Criar interfaces no domínio/application apropriado.

Exemplo:

```go
type CommandItemRepository interface {
    Create(...)
    GetByID(...)
    Update(...)
    List(...)
    Count(...)
}
```

Também criar ports para:

```text
CommandEventRepository
CommandActionRepository
OperationalSnapshotRepository
MissionControlSummaryRepository
IdempotencyRepository
```

Os repositories devem ser abstrações.

---

# 29. Infrastructure

Implementar persistência PostgreSQL.

Criar migrations específicas da Sprint 017.

Sugestão:

```text
database/migrations/
```

Criar tabelas necessárias para:

```text
mission_control_command_items
mission_control_command_actions
mission_control_command_events
mission_control_snapshots
mission_control_idempotency
```

Adicionar:

* PK;
* FK internas;
* tenant indexes;
* status indexes;
* severity indexes;
* priority indexes;
* created_at indexes;
* due_at indexes;
* fingerprint unique index;
* constraints;
* check constraints quando apropriado.

---

# 30. Transport / HTTP

Criar endpoints REST versionados:

```text
/api/v1/mission-control
```

## Summary

```http
GET /api/v1/mission-control/summary
```

## Command Items

```http
GET    /api/v1/mission-control/items
GET    /api/v1/mission-control/items/:id
POST   /api/v1/mission-control/items
PATCH  /api/v1/mission-control/items/:id
```

## Lifecycle

```http
POST /api/v1/mission-control/items/:id/acknowledge
POST /api/v1/mission-control/items/:id/start
POST /api/v1/mission-control/items/:id/resolve
POST /api/v1/mission-control/items/:id/dismiss
```

## Actions

```http
GET  /api/v1/mission-control/items/:id/actions
POST /api/v1/mission-control/items/:id/actions
```

## History

```http
GET /api/v1/mission-control/items/:id/history
```

## Snapshot

```http
GET /api/v1/mission-control/snapshot
POST /api/v1/mission-control/snapshot/rebuild
```

## Recommendations

```http
GET /api/v1/mission-control/recommendations
POST /api/v1/mission-control/recommendations/evaluate
```

---

# 31. Filtros

O endpoint:

```http
GET /api/v1/mission-control/items
```

deve suportar filtros:

```text
type
category
severity
priority
status
source_type
assigned_to
sla_status
created_from
created_to
due_from
due_to
```

Suportar:

```text
pagination
sorting
```

---

# 32. Segurança

Manter o padrão de autenticação já existente.

Utilizar:

```text
JWT
tenant context
permissions
authorization
```

Criar permissões específicas:

```text
mission_control.read
mission_control.create
mission_control.update
mission_control.acknowledge
mission_control.resolve
mission_control.dismiss
mission_control.snapshot
mission_control.admin
```

Não permitir bypass do middleware existente.

---

# 33. Auditoria

Toda mudança relevante deve gerar `CommandEvent`.

Exemplo:

```text
created
acknowledged
started
resolved
dismissed
priority_changed
severity_changed
assigned
```

Não sobrescrever histórico.

---

# 34. Observabilidade

Preparar logs estruturados para:

```text
command_item_created
command_item_updated
command_item_acknowledged
command_item_resolved
command_item_dismissed
recommendation_generated
snapshot_generated
duplicate_command_item_detected
invalid_state_transition
```

Não adicionar dependências desnecessárias.

Utilizar o mecanismo de logging já existente no projeto.

---

# 35. OpenAPI

Atualizar a documentação OpenAPI existente.

Documentar:

* endpoints;
* request;
* response;
* parâmetros;
* filtros;
* erros;
* autenticação;
* permissões;
* enums;
* paginação.

Não criar documentação paralela.

---

# 36. Erros de domínio

Criar erros específicos:

```text
ErrCommandItemNotFound
ErrInvalidCommandItemType
ErrInvalidStatusTransition
ErrInvalidSeverity
ErrInvalidPriority
ErrInvalidRiskScore
ErrInvalidImpactScore
ErrInvalidConfidence
ErrDuplicateCommandItem
ErrTenantMismatch
ErrCommandActionNotAllowed
ErrInvalidSLA
```

Não utilizar strings espalhadas pelo código.

---

# 37. Testes unitários

Cobertura obrigatória para:

### State Machine

Testar:

```text
open → acknowledged
open → dismissed
open → resolved
acknowledged → in_progress
acknowledged → resolved
in_progress → resolved
```

E transições inválidas.

### Scores

Testar:

```text
0
1
boundary values
negative
> maximum
```

### Ranking

Testar todos os critérios.

### Deduplicação

Testar:

```text
same fingerprint → duplicate
different fingerprint → allowed
different tenant → allowed
```

### SLA

Testar:

```text
within
at risk
breached
```

### Health Score

Testar todos os boundaries.

### Risk Score

Testar todos os boundaries.

### Recommendation Rules

Testar cada regra individualmente.

---

# 38. Testes de integração

Criar testes com PostgreSQL/test database para:

```text
create item
read item
update item
list items
filter items
transition status
persist events
persist actions
generate snapshot
generate summary
deduplication
idempotency
tenant isolation
```

---

# 39. Testes HTTP

Criar testes para os endpoints principais:

```text
GET summary
GET items
GET item
POST item
PATCH item
POST acknowledge
POST start
POST resolve
POST dismiss
GET history
GET recommendations
GET snapshot
```

Testar:

```text
200
201
400
401
403
404
409
422
500
```

quando aplicável.

---

# 40. Contratos futuros

Criar ports para futuras integrações, sem implementá-las.

Exemplo conceitual:

```go
type OperationalSignalProvider interface {
    GetSignals(ctx context.Context, tenantID string) ([]Signal, error)
}
```

Ou contratos equivalentes que façam sentido à arquitetura existente.

Esses ports devem existir apenas como **abstrações**.

Não implementar:

```text
FleetSignalProvider
TireSignalProvider
MaintenanceSignalProvider
FuelSignalProvider
FinancialSignalProvider
```

nesta Sprint, salvo se forem interfaces vazias/contratuais absolutamente necessárias.

---

# 41. Proibição de acoplamento

Durante a implementação, verificar que:

```text
missioncontrol
```

não importa:

```text
fleet
tires
maintenance
inventory
fuel
suppliers
financial
ciot
checklist
intelligence
```

para acessar regras, repositories ou entidades internas.

Se houver necessidade de informação externa, utilizar abstração.

---

# 42. Banco de dados

Não reutilizar tabelas de outros módulos.

Não criar FK para tabelas de outros bounded contexts.

Não fazer:

```sql
JOIN fleet_vehicles ...
JOIN tires ...
JOIN maintenance_orders ...
```

O Mission Control deve possuir seu próprio modelo.

---

# 43. Performance

Adicionar índices para os principais padrões de consulta.

Principalmente:

```text
tenant_id
tenant_id + status
tenant_id + severity
tenant_id + priority
tenant_id + due_at
tenant_id + detected_at
tenant_id + fingerprint
```

Paginação obrigatória em listagens.

Evitar N+1 queries.

---

# 44. Consistência

Garantir:

```text
transactional status update
+
event creation
```

na mesma transação quando necessário.

Exemplo:

```text
resolve item
    ↓
update status
    ↓
create event
    ↓
commit
```

Se uma etapa falhar:

```text
rollback
```

---

# 45. Transaction boundaries

A aplicação deve definir claramente os casos transacionais.

Não deixar transações espalhadas pelos handlers HTTP.

A regra deve permanecer na camada application/use case.

---

# 46. Use Cases

Criar casos de uso separados, por exemplo:

```text
CreateCommandItem
GetCommandItem
ListCommandItems
UpdateCommandItem
AcknowledgeCommandItem
StartCommandItem
ResolveCommandItem
DismissCommandItem
GetCommandHistory
CreateCommandAction
GetCommandActions
GenerateMissionControlSummary
GenerateOperationalSnapshot
EvaluateRecommendations
```

Não colocar lógica de domínio nos controllers.

---

# 47. DTOs

Separar:

```text
HTTP request DTO
application command
domain entity
HTTP response DTO
```

Não retornar entidade de domínio diretamente.

---

# 48. Paginação

Padronizar a resposta.

Exemplo conceitual:

```json
{
  "data": [],
  "pagination": {
    "page": 1,
    "page_size": 20,
    "total": 100,
    "total_pages": 5
  }
}
```

Respeitar o padrão já existente no projeto caso exista outro formato oficial.

---

# 49. Health Check

Não alterar o comportamento existente de:

```http
GET /health
```

Se necessário, adicionar apenas checks internos compatíveis com a infraestrutura atual.

---

# 50. Compatibilidade

Não quebrar:

* autenticação;
* autorização;
* tenant;
* módulos existentes;
* migrations anteriores;
* APIs existentes;
* testes existentes.

Executar:

```bash
go test ./...
```

ao final.

---

# 51. Qualidade de código

Seguir:

* Go idiomático;
* SOLID;
* Clean Architecture;
* DDD quando aplicável;
* dependency inversion;
* interfaces pequenas;
* composição;
* errors explícitos;
* context.Context;
* validação no domínio;
* SQL parametrizado;
* transações explícitas;
* testes determinísticos.

Evitar:

* `any` desnecessário;
* global state;
* singleton desnecessário;
* lógica de negócio em handlers;
* SQL espalhado na aplicação;
* dependência circular;
* importação direta de módulos operacionais.

---

# 52. Não implementar nesta Sprint

É obrigatório NÃO implementar:

```text
Flutter
React
Frontend
Dashboard visual
Deploy
Docker de produção
Cloud
CI/CD novo
LLM
OpenAI
Machine Learning
Computer Vision
Google Vision
Google Lens
integrações externas
Pamcard
RoadCard
CIOT operacional
notificações externas
WhatsApp
e-mail
SMS
push notification
```

Também não implementar ações que alterem entidades dos módulos operacionais.

---

# 53. Critério arquitetural principal

O Mission Control deverá conseguir funcionar com dados próprios mesmo que todos os outros módulos estejam desligados.

Ou seja:

```text
Mission Control
      ↓
PostgreSQL próprio
      ↓
Domain/Application
```

deve funcionar isoladamente.

Os módulos externos serão conectados posteriormente através de contratos/eventos/read models.

---

# 54. Migration

Criar migration versionada seguindo o padrão atual do projeto.

Não modificar migrations antigas.

Não apagar dados.

Não realizar alterações destrutivas.

---

# 55. Documentação

Criar ou atualizar:

```text
docs/
```

com documentação da Sprint 017 contendo:

```text
Architecture
Domain
Entities
State Machine
Scoring
Ranking
Recommendation Rules
API
Database
Security
Tenant Isolation
Idempotency
Testing
Future Integration Points
```

---

# 56. Checklist de conclusão

Antes de considerar a Sprint concluída, validar:

* [ ] Mission Control module criado.
* [ ] Arquitetura domain/application/infrastructure/transport/ports respeitada.
* [ ] CommandItem implementado.
* [ ] CommandAction implementado.
* [ ] CommandEvent implementado.
* [ ] OperationalSnapshot implementado.
* [ ] Summary implementado.
* [ ] State machine implementada.
* [ ] Severity implementada.
* [ ] Priority implementada.
* [ ] Risk score implementado.
* [ ] Impact score implementado.
* [ ] Confidence implementada.
* [ ] SLA implementado.
* [ ] Deduplicação implementada.
* [ ] Idempotência implementada.
* [ ] Ranking implementado.
* [ ] Operational Score implementado.
* [ ] Risk Aggregation implementado.
* [ ] Health Score implementado.
* [ ] Recommendation Engine determinístico implementado.
* [ ] Regras testadas.
* [ ] Repository ports criados.
* [ ] PostgreSQL persistence criada.
* [ ] Migration criada.
* [ ] Tenant isolation validada.
* [ ] JWT/permissions integrados ao padrão existente.
* [ ] Auditoria implementada.
* [ ] Endpoints REST implementados.
* [ ] OpenAPI atualizado.
* [ ] Unit tests implementados.
* [ ] Integration tests implementados.
* [ ] HTTP tests implementados.
* [ ] `go test ./...` executado com sucesso.
* [ ] Nenhum módulo operacional importado diretamente.
* [ ] Nenhuma integração operacional criada.
* [ ] Nenhum Flutter criado.
* [ ] Nenhum deploy realizado.

---

# 57. Definition of Done

A Sprint 017 só será considerada concluída quando:

```text
Mission Control existir como bounded context backend independente,
com domínio próprio,
persistência própria,
APIs próprias,
state machine,
scores determinísticos,
ranking,
recomendações,
auditoria,
idempotência,
tenant isolation,
testes unitários,
testes de integração,
testes HTTP
e documentação,
sem dependência operacional dos demais módulos.
```

O resultado deve ser uma fundação sólida para que, em uma Sprint futura, os sinais provenientes de Fleet, Tires, Maintenance, Fuel, Inventory, Financial, Checklist, CIOT e demais módulos possam ser conectados através de **ports, eventos e read models**, sem necessidade de reestruturar o Mission Control.

---

# 58. Comando final para execução

Execute esta Sprint como uma implementação real no repositório existente.

Antes de modificar código:

1. inspecione a arquitetura atual;
2. identifique padrões existentes de modules;
3. identifique padrões de migration;
4. identifique padrões de repository;
5. identifique padrões de authentication/tenant;
6. identifique padrões de HTTP;
7. identifique padrões de testes;
8. identifique o que já existe relacionado a Mission Control/Intelligence.

Depois:

1. implemente o domínio;
2. implemente application/use cases;
3. implemente ports;
4. implemente infrastructure;
5. implemente migrations;
6. implemente transport;
7. implemente autorização;
8. implemente testes;
9. atualize OpenAPI;
10. atualize documentação;
11. execute todos os testes.

### Regra final

> **Não invente uma arquitetura paralela.**
>
> **Não refatore módulos não relacionados.**
>
> **Não conecte operacionalmente os módulos nesta Sprint.**
>
> **Não implemente frontend.**
>
> **Não faça deploy.**
>
> **Não use IA generativa para decisões do domínio.**
>
> **O Mission Control deve ser determinístico, isolado, testável, multi-tenant e preparado para futuras integrações através de contratos.**

Ao final, apresentar:

```text
1. Arquivos criados
2. Arquivos alterados
3. Migrations criadas
4. Endpoints adicionados
5. Entidades implementadas
6. Use cases implementados
7. Regras determinísticas implementadas
8. Testes executados
9. Resultado de go test ./...
10. Possíveis pendências
11. Pontos preparados para a Sprint 018
```

Não declarar a Sprint como concluída se `go test ./...` falhar ou se houver acoplamento direto do Mission Control com módulos operacionais.
