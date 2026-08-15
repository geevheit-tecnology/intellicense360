# Event Bus

Sprint 015 cria a infraestrutura interna de Domain Events. Ela permite que modulos publiquem eventos de negocio sem importar outros dominios diretamente.

## Proposito

O Event Bus e infraestrutura. Ele nao contem regra de negocio, nao executa integracoes externas e nao substitui transacoes de dominio.

Fluxo alvo:

```text
Module -> Domain Event -> Event Bus -> Subscribers -> futuros Workflows/Intelligence
```

## Domain Event

Envelope padrao:

```json
{
  "event_id": "...",
  "event_type": "tire.removed.v1",
  "event_version": 1,
  "occurred_at": "...",
  "tenant_id": "...",
  "aggregate_id": "...",
  "aggregate_type": "tire",
  "correlation_id": "...",
  "causation_id": "...",
  "actor_id": "...",
  "actor_type": "...",
  "metadata": {},
  "payload": {}
}
```

Payloads devem ser JSON-compatible e nao devem expor detalhes internos de banco.

## Nomenclatura

Use nomes estaveis, semanticamente ligados ao negocio:

- `tire.removed.v1`
- `fuel.transaction.completed.v1`
- `checklist.execution.completed.v1`
- `maintenance.order.created.v1`
- `financial.expense.paid.v1`
- `ciot.closed.v1`

Nao use nomes de pacotes Go como identificador de evento.

## Versionamento

O sufixo `v1`, `v2` etc. faz parte do contrato. O significado de uma versao existente nao deve ser alterado silenciosamente. Mudancas incompativeis exigem novo tipo versionado.

## Event Bus

Interface em `backend/api/internal/events`:

- `Publish`
- `PublishBatch`
- `Subscribe`
- `Unsubscribe`

O adapter atual e `events/inmemory`, deterministico para desenvolvimento e testes. Ele isola falhas de handlers: um handler com erro nao impede a chamada dos demais, mas o erro e retornado ao publicador.

## Handlers

Handlers declaram explicitamente quais `event_type` processam. Nesta sprint eles podem observar, logar ou testar eventos, mas nao alteram outros dominios.

## Outbox

O outbox prepara publicacao confiavel futura:

- `pending`
- `processing`
- `published`
- `failed`
- `dead_letter`

O outbox em memoria cobre testes e desenvolvimento. A migration `000017_event_bus_outbox.sql` prepara tabelas PostgreSQL para persistencia futura.

## Idempotencia

`EventConsumerStore` registra `consumer_name + event_id`. O mesmo consumidor logico nao deve processar o mesmo evento duas vezes.

## Retry e Dead Letter

Retries possuem limite configuravel. Quando o limite e excedido, o evento vai para `dead_letter`; ele nao e descartado silenciosamente.

## Correlation e Causation

`correlation_id` agrupa um fluxo de negocio. `causation_id` aponta o evento/comando que causou o evento atual.

## Seguranca

Payloads nunca devem conter senhas, JWTs, tokens, API keys, credenciais de pagamento, credenciais de providers ou informacoes sensiveis de autenticacao. `tenant_id` deve estar sempre explicito.

## Observabilidade

O bus registra logs estruturados para:

- evento publicado;
- evento tratado;
- falha de handler;
- retry e dead letter no outbox.

Campos relevantes: `event_id`, `event_type`, `tenant_id`, `correlation_id`, `handler`, `duration`, `result`.

## Limite Transacional

A transacao de dominio, a persistencia no outbox e a publicacao no bus sao fronteiras diferentes. O adapter em memoria nao fornece garantias distribuidas.

## Estrategia Futura

Adapters futuros podem implementar a mesma interface:

- `KafkaEventBus`
- `RabbitMQEventBus`
- `NATSEventBus`
- `CloudEventBus`

Nenhum broker externo foi instalado nesta sprint.
