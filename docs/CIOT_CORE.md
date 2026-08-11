# CIOT Core

Sprint 014 cria o dominio `internal/modules/ciot` sem integracoes externas e sem dependencias com Fleet, Drivers, Financial, Suppliers, Maintenance, Tires, Checklist ou Intelligence.

## Estado

Transicoes permitidas:

- `draft -> pending`
- `pending -> generated`
- `generated -> active`
- `active -> suspended`
- `suspended -> active`
- `active -> closed`
- `pending -> canceled`
- `generated -> canceled`
- `active -> canceled`
- `error -> pending`

Qualquer outra mudanca de status retorna `invalid_transition`.

## TAC Agregado

O tipo `tac_agregado` suporta operacoes longas. Um CIOT pode permanecer `active` durante o periodo operacional sem preencher `actual_end_date`.

Campos relevantes:

- `start_date`
- `expected_end_date`
- `actual_end_date`
- `operational_period`
- `contract_reference`
- `notes`

## API

Base protegida por `ciot.ciot.manage`:

- `GET /api/v1/ciot`
- `POST /api/v1/ciot`
- `GET /api/v1/ciot/{id}`
- `PATCH /api/v1/ciot/{id}`
- `DELETE /api/v1/ciot/{id}`
- `GET /api/v1/ciot/{id}/history`
- `POST /api/v1/ciot/{id}/submit`
- `POST /api/v1/ciot/{id}/generated`
- `POST /api/v1/ciot/{id}/activate`
- `POST /api/v1/ciot/{id}/suspend`
- `POST /api/v1/ciot/{id}/reactivate`
- `POST /api/v1/ciot/{id}/close`
- `POST /api/v1/ciot/{id}/cancel`
- `POST /api/v1/ciot/{id}/payments`
- `POST /api/v1/ciot/{id}/provider-attempts`
- `POST /api/v1/ciot/{id}/external-reference`
- `GET /api/v1/ciot/types`

## Limites

Nao ha cliente RoadCard, Pamcard, e-Frete ou qualquer provider. `CIOTProviderAttempt` e `CIOTExternalReference` sao estruturas provider-agnostic para uma sprint futura.
