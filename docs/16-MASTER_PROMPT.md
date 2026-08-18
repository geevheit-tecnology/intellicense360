E mantenha exatamente a regra que estamos seguindo: esta sprint ainda não cria automações que alterem outros módulos. O Intelligence apenas observa, calcula, detecta e recomenda.
SPRINT 016 — GEEVHEIT INTELLIGENCE ENGINE

Read first:

docs/00_PRODUCT_MANIFESTO.md
docs/01_MASTER_PROMPT.md
docs/ARCHITECTURE.md
docs/BACKEND_ARCHITECTURE.md
docs/EVENT_BUS.md

Also inspect:

backend/api/internal/events/
backend/api/internal/contracts/events/
backend/api/internal/modules/
backend/api/internal/core/

IMPORTANT

This sprint creates ONLY the Geevheit Intelligence Engine.

DO NOT build the final dashboard.

DO NOT build Flutter.

DO NOT deploy.

DO NOT create external AI integrations.

DO NOT use OpenAI API.

DO NOT use external LLMs.

DO NOT automatically modify Fleet.

DO NOT automatically create Maintenance Orders.

DO NOT automatically modify Tires.

DO NOT automatically modify Fuel.

DO NOT automatically modify Financial.

DO NOT automatically modify Checklist.

DO NOT automatically modify CIOT.

DO NOT create direct business-module dependencies.

The Intelligence Engine is READ/ANALYSIS oriented.

It observes operational data and produces intelligence.

It does NOT execute operational actions.

--------------------------------------------------
OBJECTIVE
--------------------------------------------------

Create the Geevheit Intelligence Engine.

The purpose is to transform operational facts into:

KPIs

Metrics

Trends

Anomalies

Risks

Opportunities

Recommendations

Cost indicators

Efficiency indicators

Confidence scores

Estimated impact

The engine must answer questions such as:

"What is happening?"

"Why is it happening?"

"What is abnormal?"

"How much does it cost?"

"What deserves attention?"

"What could save money?"

"What should the operator investigate?"

--------------------------------------------------
CORE PRINCIPLE
--------------------------------------------------

The Intelligence Engine must distinguish between:

FACT

METRIC

ANOMALY

RISK

OPPORTUNITY

RECOMMENDATION

IMPACT

Do not mix these concepts.

Example:

FACT:
Fuel transaction completed.

METRIC:
Average consumption = 2.45 km/L.

ANOMALY:
Consumption is 14% below historical baseline.

RISK:
Potential abnormal fuel consumption.

OPPORTUNITY:
Potential monthly saving.

RECOMMENDATION:
Investigate the asset's recent fuel pattern.

--------------------------------------------------
MODULE
--------------------------------------------------

Create:

backend/api/internal/modules/intelligence

Architecture:

domain/
application/
infrastructure/
transport/
ports/

Suggested:

domain/
  entities.go
  value_objects.go
  rules.go
  thresholds.go
  insights.go

application/
  services.go
  metric_service.go
  anomaly_service.go
  recommendation_service.go

infrastructure/
  memory_store.go

transport/
  http/

ports/
  repository.go
  event_source.go
  metric_provider.go

--------------------------------------------------
CORE ENTITIES
--------------------------------------------------

IntelligenceMetric

MetricSnapshot

MetricDefinition

Anomaly

Risk

Opportunity

Recommendation

Insight

InsightEvidence

InsightImpact

InsightStatus

IntelligenceRule

IntelligenceThreshold

IntelligencePeriod

--------------------------------------------------
METRIC
--------------------------------------------------

A Metric must contain:

ID

TenantID

MetricType

Name

Value

Unit

PeriodStart

PeriodEnd

Dimension

DimensionValue

Source

CalculatedAt

Confidence

Metadata

Examples:

Fuel Consumption

Fuel Cost

Maintenance Cost

Tire Cost

Checklist Failure Rate

Maintenance Frequency

Expense Total

Revenue Total

Operational Cost

Cost Per Kilometer

Cost Per Asset

Cost Per Operation

Do not calculate cross-module metrics unless the required data is already available through an abstraction.

--------------------------------------------------
METRIC TYPES
--------------------------------------------------

Operational

Financial

Fuel

Tire

Maintenance

Checklist

Inventory

CIOT

Asset

Efficiency

Risk

Cost

Do not hardcode tenant-specific thresholds.

--------------------------------------------------
ANOMALY
--------------------------------------------------

Represent:

ID

TenantID

Type

Severity

Metric

ObservedValue

ExpectedValue

Deviation

DeviationPercentage

DetectedAt

Period

Evidence

Confidence

Status

--------------------------------------------------
ANOMALY SEVERITY
--------------------------------------------------

Info

Low

Medium

High

Critical

--------------------------------------------------
RISK
--------------------------------------------------

Represent:

Risk ID

Category

Severity

Probability

Impact

Confidence

Evidence

DetectedAt

Status

Examples:

Abnormal Fuel Consumption

Repeated Maintenance Failures

Critical Tire Condition

Checklist Failure Pattern

Unusual Cost Increase

CIOT Operational Risk

Inventory Stock Risk

Do not automatically execute remediation.

--------------------------------------------------
OPPORTUNITY
--------------------------------------------------

Represent:

Opportunity ID

Category

Estimated Impact

Estimated Saving

Confidence

Evidence

DetectedAt

ValidUntil

Status

Examples:

Fuel Efficiency

Maintenance Optimization

Tire Optimization

Inventory Optimization

Operational Efficiency

Cost Reduction

--------------------------------------------------
RECOMMENDATION
--------------------------------------------------

A recommendation must explain:

What happened

Why it matters

Evidence

Suggested action

Expected impact

Confidence

Priority

Do not execute the action automatically.

Example:

"Fuel consumption for asset X is 13.4% above its historical baseline. Investigate recent fueling events and operating conditions. Estimated avoidable cost: R$ 2,840/month."

--------------------------------------------------
INSIGHT
--------------------------------------------------

Insight is the user-facing intelligence object.

It should combine:

Title

Summary

Category

Severity

Evidence

Metric

Anomaly

Risk

Opportunity

Recommendation

Estimated Impact

Confidence

CreatedAt

Status

--------------------------------------------------
INSIGHT STATUS
--------------------------------------------------

New

Acknowledged

In Progress

Resolved

Dismissed

Expired

The Intelligence Engine may update its own insight status.

It must NOT modify operational domains.

--------------------------------------------------
EVIDENCE
--------------------------------------------------

Every important insight should be explainable.

Evidence must contain:

Source

Source Type

Reference ID

Metric

Observed Value

Expected Value

Period

Timestamp

Explanation

Do not generate unexplained recommendations.

--------------------------------------------------
CONFIDENCE
--------------------------------------------------

Create a confidence model.

Possible values:

0.0 to 1.0

Represent:

Data quality

Sample size

Historical consistency

Rule confidence

Do not pretend confidence is AI probability.

Document how confidence is calculated.

--------------------------------------------------
RULE ENGINE
--------------------------------------------------

Create a deterministic rule engine.

Rules must be:

Explicit

Versioned

Testable

Tenant-aware

Explainable

Do not use machine learning in this sprint.

Do not use an LLM.

Do not use external AI services.

--------------------------------------------------
RULE TYPES
--------------------------------------------------

ThresholdRule

DeviationRule

TrendRule

FrequencyRule

PatternRule

CostRule

--------------------------------------------------
EXAMPLE RULES
--------------------------------------------------

Fuel consumption deviation:

IF observed consumption deviates significantly from historical baseline

THEN create anomaly.

Fuel cost increase:

IF fuel cost increases above configured threshold

THEN create anomaly.

Repeated checklist failure:

IF the same checklist issue occurs repeatedly

THEN create risk.

Repeated tire problem:

IF abnormal tire events repeat

THEN create risk.

Maintenance recurrence:

IF the same maintenance issue repeatedly occurs

THEN create risk.

Cost opportunity:

IF historical data indicates potential reduction

THEN create opportunity.

Do not activate cross-module rules unless required data is available.

--------------------------------------------------
BASELINE
--------------------------------------------------

Create baseline abstractions.

Support:

Historical Average

Moving Average

Rolling Window

Previous Period

Comparable Period

Do not hardcode one algorithm.

--------------------------------------------------
TREND
--------------------------------------------------

Support:

Increasing

Decreasing

Stable

Volatile

Insufficient Data

A trend must include:

Direction

Magnitude

Period

Confidence

--------------------------------------------------
TIME WINDOWS
--------------------------------------------------

Support:

Daily

Weekly

Monthly

Quarterly

Custom

Do not hardcode timezone assumptions.

Use tenant/application timezone abstractions where appropriate.

--------------------------------------------------
COST INTELLIGENCE
--------------------------------------------------

Prepare metrics for:

Total Cost

Variable Cost

Fixed Cost

Cost Per Kilometer

Cost Per Asset

Cost Per Operation

Fuel Cost

Tire Cost

Maintenance Cost

Inventory Cost

Do NOT implement automatic financial posting.

--------------------------------------------------
DATA SOURCES
--------------------------------------------------

Create provider abstractions for intelligence data.

Possible sources:

Financial

Fuel

Tires

Maintenance

Checklist

Inventory

CIOT

Assets

Fleet

Do not import business modules directly.

Use:

interfaces

repositories

event projections

read models

or other decoupled abstractions.

The preferred architecture is event/read-model based.

--------------------------------------------------
EVENT CONSUMPTION
--------------------------------------------------

The Intelligence Engine may subscribe to Domain Events through the Event Bus.

Examples:

fuel.transaction.completed.v1

fuel.transaction.adjusted.v1

tire.removed.v1

tire.inspected.v1

tire.retreaded.v1

maintenance.order.created.v1

maintenance.order.completed.v1

checklist.execution.completed.v1

checklist.non_conformity.created.v1

financial.expense.created.v1

financial.expense.paid.v1

ciot.activated.v1

ciot.closed.v1

IMPORTANT:

Consuming events must NOT modify the originating business module.

The Intelligence Engine creates/updates its own read models and intelligence data.

--------------------------------------------------
READ MODELS
--------------------------------------------------

Create intelligence-oriented read models.

Examples:

FuelUsageSnapshot

MaintenanceCostSnapshot

TireLifecycleSnapshot

ChecklistRiskSnapshot

FinancialCostSnapshot

CIOTOperationalSnapshot

AssetPerformanceSnapshot

These are projections.

They are NOT replacements for source-of-truth business entities.

--------------------------------------------------
EVENT PROCESSING
--------------------------------------------------

Every consumed event must be idempotent.

Use EventID.

Do not process the same event twice.

Use the Event Bus idempotency infrastructure where appropriate.

--------------------------------------------------
INSIGHT DEDUPLICATION
--------------------------------------------------

Do not create thousands of identical insights.

Create a deduplication strategy based on:

Tenant

Insight Type

Dimension

Time Window

Rule Version

Example:

Do not generate the same fuel anomaly every minute.

--------------------------------------------------
PRIORITY
--------------------------------------------------

Create priority calculation:

Low

Medium

High

Critical

Priority may consider:

Severity

Financial Impact

Operational Impact

Confidence

Recurrence

--------------------------------------------------
IMPACT
--------------------------------------------------

Impact must support:

Estimated Cost

Estimated Saving

Operational Impact

Risk Impact

Confidence

Currency

Do not claim precise savings without sufficient data.

Use:

Estimated

Range

Confidence

when appropriate.

--------------------------------------------------
API
--------------------------------------------------

Create:

/api/v1/intelligence/metrics

/api/v1/intelligence/metrics/{id}

/api/v1/intelligence/anomalies

/api/v1/intelligence/risks

/api/v1/intelligence/opportunities

/api/v1/intelligence/recommendations

/api/v1/intelligence/insights

/api/v1/intelligence/insights/{id}

/api/v1/intelligence/insights/{id}/acknowledge

/api/v1/intelligence/insights/{id}/resolve

/api/v1/intelligence/insights/{id}/dismiss

/api/v1/intelligence/health

Support:

GET

POST where appropriate

Pagination

Filtering

Sorting

Search

Date ranges

Severity

Category

Status

--------------------------------------------------
IMPORTANT API RULE
--------------------------------------------------

The Intelligence API must NOT expose endpoints that mutate operational domains.

It may only manage intelligence data.

--------------------------------------------------
RBAC
--------------------------------------------------

Create permission:

intelligence.intelligence.manage

If read-only permissions already exist in the architecture, prepare:

intelligence.intelligence.read

Do not weaken existing RBAC.

--------------------------------------------------
AUDIT
--------------------------------------------------

Prepare audit events for:

Metric Calculated

Anomaly Detected

Risk Detected

Opportunity Detected

Recommendation Generated

Insight Created

Insight Acknowledged

Insight Resolved

Insight Dismissed

Do NOT publish external events.

Do not create a feedback loop into operational domains.

--------------------------------------------------
DATABASE
--------------------------------------------------

Create versioned migration:

database/migrations/000018_intelligence_core.sql

Tables should support:

intelligence_metrics

intelligence_anomalies

intelligence_risks

intelligence_opportunities

intelligence_recommendations

intelligence_insights

intelligence_evidence

intelligence_rule_definitions

intelligence_read_models

Use:

UUID

TenantID

timestamps

indexes

versioning

status

confidence

rule version

period fields

--------------------------------------------------
SECURITY
--------------------------------------------------

Never expose:

JWT

passwords

API keys

provider credentials

secrets

internal authentication data

inside intelligence records.

Respect tenant isolation.

--------------------------------------------------
PERFORMANCE
--------------------------------------------------

Do not calculate expensive historical analysis synchronously inside HTTP requests if avoidable.

Prepare application services that can later run asynchronously.

Do not introduce a job queue unless already present.

Keep the architecture ready for asynchronous processing.

--------------------------------------------------
OBSERVABILITY
--------------------------------------------------

Log:

Event consumed

Metric calculated

Rule evaluated

Anomaly detected

Insight generated

Insight deduplicated

Processing failure

Include:

EventID

TenantID

CorrelationID

Rule

Duration

Result

--------------------------------------------------
TESTS
--------------------------------------------------

Create comprehensive tests for:

Metric calculation

Baseline

Moving average

Deviation

Trend

Anomaly detection

Risk detection

Opportunity detection

Recommendation generation

Confidence calculation

Impact calculation

Insight deduplication

Priority calculation

Event consumption

Event idempotency

Tenant isolation

Rule versioning

Invalid data

Insufficient data

API

RBAC

--------------------------------------------------
IMPORTANT TEST SCENARIO
--------------------------------------------------

Create deterministic tests such as:

Historical fuel consumption:

2.8 km/L
2.7 km/L
2.9 km/L
2.8 km/L

Current:

2.35 km/L

The engine must detect a significant deviation according to the configured rule.

It must produce:

Anomaly

Evidence

Confidence

Recommendation

Potential impact when enough data exists.

Do not hardcode a specific financial saving.

--------------------------------------------------
SECOND TEST
--------------------------------------------------

Create repeated operational failure data.

Example:

Same checklist non-conformity repeated several times in a defined period.

The engine must detect a recurrence pattern and create a Risk.

--------------------------------------------------
THIRD TEST
--------------------------------------------------

Create repeated maintenance-related events.

Detect recurrence.

Generate a Risk.

Do not create a Maintenance Order.

--------------------------------------------------
FOURTH TEST
--------------------------------------------------

Create insufficient historical data.

The engine must NOT generate a high-confidence recommendation.

It should return:

Insufficient Data

or low confidence.

--------------------------------------------------
OPENAPI
--------------------------------------------------

Update existing OpenAPI.

Do not create another OpenAPI file.

--------------------------------------------------
DOCUMENTATION
--------------------------------------------------

Create:

docs/INTELLIGENCE_ENGINE.md

Document:

Purpose

Architecture

Facts

Metrics

Baselines

Trends

Anomalies

Risks

Opportunities

Recommendations

Insights

Confidence

Impact

Rule Engine

Read Models

Event Consumption

Idempotency

Deduplication

Security

Performance

Testing

Future AI Strategy

--------------------------------------------------
FUTURE AI STRATEGY
--------------------------------------------------

Document that deterministic intelligence comes first.

Future versions may optionally use ML/LLM capabilities for:

Natural language explanation

Advanced pattern detection

Predictive maintenance

Forecasting

Natural language queries

However:

Do NOT implement those capabilities now.

The deterministic engine must remain the source of explainable operational intelligence.

--------------------------------------------------
ARCHITECTURAL RULE
--------------------------------------------------

The Intelligence Engine must NOT become a god service.

It must not own:

Fleet

Tires

Fuel

Maintenance

Financial

Checklist

CIOT

Inventory

Suppliers

Assets

Those remain source-of-truth modules.

Intelligence owns only:

Metrics

Read Models

Rules

Anomalies

Risks

Opportunities

Recommendations

Insights

--------------------------------------------------
QUALITY
--------------------------------------------------

DDD

Clean Architecture

SOLID

Hexagonal Architecture

Dependency Inversion

Event-driven architecture

CQRS-friendly read models

Deterministic rules

Explainable intelligence

Idempotent event consumption

Tenant isolation

Production-oriented implementation

Do not refactor unrelated modules.

Do not break existing endpoints.

--------------------------------------------------
VALIDATION
--------------------------------------------------

Run:

go test ./...

Verify ALL existing tests continue passing.

Start API using configurable PORT.

Run HTTP smoke tests for:

login

GET /api/v1/intelligence/health

create/query metric

create/query anomaly

create/query risk

create/query opportunity

create/query recommendation

create/query insight

acknowledge insight

resolve insight

RBAC

tenant isolation

Verify configurable PORT.

--------------------------------------------------
ARCHITECTURAL VALIDATION
--------------------------------------------------

Verify:

Intelligence does NOT import Fleet.

Intelligence does NOT import Tires.

Intelligence does NOT import Fuel.

Intelligence does NOT import Maintenance.

Intelligence does NOT import Financial.

Intelligence does NOT import Checklist.

Intelligence does NOT import CIOT.

Intelligence does NOT import Inventory.

Intelligence does NOT import Suppliers.

Intelligence does NOT modify operational domains.

Intelligence consumes events/read models only.

--------------------------------------------------
FINAL REPORT
--------------------------------------------------

Report:

1. Files created
2. Files modified
3. Migration created
4. Intelligence entities
5. Metric architecture
6. Baseline architecture
7. Rule engine
8. Event consumption
9. Read models
10. Anomaly detection
11. Risk detection
12. Opportunity detection
13. Recommendation engine
14. Confidence model
15. Impact model
16. Deduplication strategy
17. Tests executed
18. Smoke tests executed
19. Architecture validation
20. Problems found
21. Future recommendations

IMPORTANT:

DO NOT START SPRINT 017.

DO NOT BUILD FINAL DASHBOARD.

DO NOT BUILD FLUTTER.

DO NOT DEPLOY.

DO NOT CREATE AUTOMATIC OPERATIONAL ACTIONS.

DO NOT CONNECT BUSINESS MODULES DIRECTLY.

STOP AFTER SPRINT 016.

WAIT FOR APPROVAL.

E sem IA generativa ainda. Isso é proposital: primeiro construímos uma inteligência determinística, auditável e explicável. Depois podemos colocar IA por cima para explicar os achados em linguagem natural.