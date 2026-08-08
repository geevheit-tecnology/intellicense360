Sprint 04 — Checklist Core (Clean Architecture)
Objetivo

Implementar o módulo Checklist totalmente independente utilizando Clean Architecture e Domain Driven Design.

IMPORTANTE

Nesta sprint NÃO integrar com:

Fleet
Tires
Maintenance
Fuel
CIOT
Intelligence

O módulo apenas armazenará um vehicle_id (UUID). Não consultar Fleet, não validar existência do veículo e não realizar chamadas entre módulos.

Arquitetura

Seguir exatamente o padrão do módulo Fleet.

backend/api/internal/modules/checklist/

domain/
application/
ports/
repository/
http/
dto/
mapper/

Separar rigorosamente:

Domain
Use Cases
Ports
Repository
HTTP Handlers
DTOs

Sem regras de negócio dentro dos handlers HTTP.

Entidades
Checklist

Campos:

id UUID

tenant_id UUID

vehicle_id UUID

title

description

type

status

started_at

finished_at

driver_name

driver_document

created_by

updated_by

deleted_at

created_at

updated_at

Status:

draft

in_progress

completed

cancelled
Checklist Item
id

checklist_id

title

description

category

required

order_index

answer_type

expected_value

created_at

updated_at

Tipos de resposta

boolean

text

number

photo

signature

select
Checklist Answer
id

checklist_item_id

answer

notes

photo_url

answered_by

answered_at
Casos de uso

Implementar:

CreateChecklist

UpdateChecklist

DeleteChecklist (Soft Delete)

GetChecklist

ListChecklists

StartChecklist

FinishChecklist

CancelChecklist

Itens

AddItem

UpdateItem

DeleteItem

ListItems

Respostas

AnswerItem

ListAnswers
Regras

Soft Delete.

Tenant Isolation.

Paginação.

Busca.

Filtros.

Ordenação.

Validação.

Repository

Criar implementação em memória semelhante ao Fleet.

Garantir:

isolamento por tenant;
soft delete;
busca;
paginação;
ordenação;
filtros.
HTTP

Criar rotas:

GET    /api/v1/checklists

POST   /api/v1/checklists

GET    /api/v1/checklists/{id}

PUT    /api/v1/checklists/{id}

DELETE /api/v1/checklists/{id}

POST   /api/v1/checklists/{id}/start

POST   /api/v1/checklists/{id}/finish

POST   /api/v1/checklists/{id}/cancel

Itens

GET

POST

PUT

DELETE

Respostas

POST /answers

GET /answers
Banco

Criar migration

database/migrations/000005_checklist_core.sql

Criar tabelas

checklists

checklist_items

checklist_answers

Todas contendo:

tenant_id

created_at

updated_at

deleted_at

Índices para:

tenant_id

vehicle_id

status

created_at
OpenAPI

Atualizar:

backend/api/openapi.yaml

backend/api/openapi/openapi.yaml

Documentar todas as rotas.

Testes

Criar testes para:

CRUD Checklist
CRUD Itens
Respostas
Paginação
Busca
Tenant Isolation
Soft Delete
Status
Casos inválidos

Executar:

go test ./...

Todos os testes devem passar.

Documentação

Criar:

docs/CHECKLIST_CORE.md

Documentar:

arquitetura;
entidades;
fluxos;
casos de uso;
endpoints;
regras de negócio;
exemplos de requisição e resposta.
Critérios de aceite

A sprint somente será considerada concluída quando:

Todos os testes passarem (go test ./...).
As migrations executarem sem erros.
O OpenAPI estiver atualizado.
O módulo Checklist estiver completamente desacoplado dos demais módulos.
As rotas responderem corretamente em testes de smoke.
Não existir nenhuma dependência direta entre Checklist e Fleet, Tires, Maintenance, Fuel, CIOT ou Intelligence.

Resultado esperado: um módulo Checklist Core totalmente funcional, independente e pronto para ser integrado aos demais módulos apenas em uma sprint futura.