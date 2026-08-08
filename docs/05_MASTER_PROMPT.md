Sprint 05 — Tires Core (Gestão de Pneus)
Objetivo

Desenvolver o módulo Tires Core totalmente independente, seguindo os mesmos princípios utilizados em Fleet e Checklist.

NÃO integrar nesta sprint com:

Fleet
Checklist
Maintenance
Fuel
CIOT
Intelligence

O módulo armazenará apenas referências como vehicle_id e position_id, sem consultar outros módulos.

Estrutura
backend/api/internal/modules/tires/

domain/
application/
ports/
repository/
http/
dto/
mapper/

Seguir exatamente o padrão dos módulos anteriores.

Entidades
Tire

Campos:

id UUID
tenant_id UUID

serial_number
fire_number
brand
model

size

construction

tire_type
position_type

manufacturing_date

purchase_date

purchase_value

supplier

dot

current_tread_mm

original_tread_mm

minimum_tread_mm

status

vehicle_id

position

current_km

total_km

recap_count

notes

created_by
updated_by

created_at
updated_at
deleted_at
TireInspection
id
tenant_id

tire_id

inspection_date

tread_mm

pressure

temperature

condition

observations

inspector

created_at
updated_at
deleted_at
TireMovement
id

tenant_id

tire_id

movement_type

vehicle_id

position

km

reason

performed_by

movement_date

created_at
updated_at
deleted_at
Status
new

in_stock

installed

removed

recapping

disposed
Tipos de movimentação
installation

rotation

removal

recapping

repair

disposal
Casos de uso

Implementar:

Tire
CreateTire

UpdateTire

DeleteTire

GetTire

ListTires
Inspeções
RegisterInspection

UpdateInspection

DeleteInspection

ListInspections
Movimentações
RegisterMovement

ListMovements
Operações
InstallTire

RemoveTire

RotateTires

SendToRecap

ReturnFromRecap

DisposeTire
Regras

Implementar:

Soft Delete
Tenant Isolation
Busca
Paginação
Ordenação
Filtros

Validar:

DOT obrigatório
Número de fogo único por tenant
Número de série único por tenant
Sulco nunca negativo
Sulco mínimo não pode ser maior que o original
Não permitir movimentação em pneu descartado
Repository

Criar repositório em memória contendo:

paginação
filtros
busca
ordenação
tenant isolation
soft delete
HTTP

Rotas

GET    /api/v1/tires

POST   /api/v1/tires

GET    /api/v1/tires/{id}

PUT    /api/v1/tires/{id}

DELETE /api/v1/tires/{id}

Inspeções

GET

POST

PUT

DELETE

Movimentações

GET

POST

Operações

POST /install

POST /remove

POST /rotate

POST /recap

POST /return

POST /dispose
Banco

Criar migration

database/migrations/000006_tires_core.sql

Criar tabelas

tires

tire_inspections

tire_movements

Todas contendo:

tenant_id

created_at

updated_at

deleted_at

Criar índices para:

tenant_id
vehicle_id
serial_number
fire_number
status
OpenAPI

Atualizar:

backend/api/openapi.yaml

backend/api/openapi/openapi.yaml

Documentar todos os endpoints.

Testes

Cobrir:

CRUD Tire
CRUD Inspection
Movimentações
Instalação
Rotação
Recapagem
Descarte
Tenant Isolation
Soft Delete
Paginação
Busca
Casos inválidos

Executar:

go test ./...

Todos os testes devem passar.

Documentação

Criar:

docs/TIRES_CORE.md

Documentar:

arquitetura;
entidades;
regras de negócio;
endpoints;
fluxos de movimentação;
exemplos de requisição e resposta.
Critérios de aceite

A Sprint 05 somente será considerada concluída quando:

✅ go test ./... passar sem falhas.
✅ Smoke HTTP validar o fluxo básico (login, criação de pneu, inspeção, movimentação e consulta).
✅ OpenAPI estiver atualizado.
✅ Migration executar sem erros.
✅ Documentação estiver completa.
✅ O módulo permanecer totalmente desacoplado de Fleet, Checklist, Maintenance, Fuel, CIOT e Intelligence, armazenando apenas referências (vehicle_id, etc.) sem dependências diretas.

Ao concluir essa sprint, você terá três módulos centrais (Fleet, Checklist e Tires) independentes e prontos para integração em uma fase posterior.