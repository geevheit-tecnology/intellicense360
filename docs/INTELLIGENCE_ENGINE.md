# Intelligence Engine

Sprint 016 cria uma inteligencia deterministica, auditavel e explicavel. Ela observa fatos operacionais, calcula metricas, detecta anomalias/riscos/oportunidades e gera recomendacoes. Ela nao executa acoes operacionais.

## Arquitetura

O modulo fica em `backend/api/internal/modules/intelligence` com `domain`, `application`, `infrastructure`, `transport` e `ports`.

Fontes de dados entram por read models, eventos de dominio ou interfaces. O modulo nao importa Fleet, Tires, Fuel, Maintenance, Financial, Checklist, CIOT, Inventory, Suppliers ou Assets.

## Conceitos

- Fact: evento/fato observado.
- Metric: valor calculado.
- Anomaly: desvio detectado.
- Risk: risco operacional ou financeiro.
- Opportunity: possibilidade de melhoria.
- Recommendation: acao sugerida, sem execucao automatica.
- Impact: custo, economia ou impacto operacional estimado.

## Baselines e Trends

O baseline inicial usa media movel quando ha amostra suficiente. Com menos de 3 pontos, o motor retorna baixa confianca/insufficient data. Trends podem ser `increasing`, `decreasing`, `stable`, `volatile` ou `insufficient_data`.

## Rule Engine

As regras sao deterministicas, versionadas e explicaveis. Tipos previstos: threshold, deviation, trend, frequency, pattern e cost.

## Confidence

Confidence e uma pontuacao entre 0 e 1 baseada em qualidade de dados, tamanho da amostra, consistencia historica e confianca da regra. Nao e probabilidade de IA.

## Deduplicacao

Insights usam `tenant:type:category:dimension:window:rule_version` para evitar repeticao excessiva.

## Event Consumption

O motor pode consumir `DomainEvent` por `ProjectionService`. O processamento e idempotente por `tenant_id + event_id` e atualiza apenas read models internos.

## API

Rotas protegidas por `intelligence.intelligence.manage`:

- `/api/v1/intelligence/health`
- `/api/v1/intelligence/metrics`
- `/api/v1/intelligence/anomalies`
- `/api/v1/intelligence/risks`
- `/api/v1/intelligence/opportunities`
- `/api/v1/intelligence/recommendations`
- `/api/v1/intelligence/insights`
- `/api/v1/intelligence/insights/{id}/acknowledge`
- `/api/v1/intelligence/insights/{id}/resolve`
- `/api/v1/intelligence/insights/{id}/dismiss`

## Future AI Strategy

IA generativa, ML, forecasting e linguagem natural podem ser adicionados depois. A fonte primaria deve continuar sendo a inteligencia deterministica e explicavel.
