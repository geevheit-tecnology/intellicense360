# Motoristas e documentos — etapa local

## O que foi implementado

Módulos `drivers` e `documents`, seguindo domain/application/ports/infrastructure/transport do backend Go existente. Os adapters são PostgreSQL e as rotas são habilitadas somente nesse modo. A demonstração antiga em memória continua com os módulos anteriores; não existe simulação em memória destes dois módulos novos.

Motoristas têm código interno, nome, identificação fiscal opcional (`tax_id`, normalizada para 11 dígitos), telefone, e-mail, situação ativa/inativa, observações e vínculo atual opcional com um veículo da mesma empresa. O formato do identificador não comprova identidade em fonte externa. CNH e seus vencimentos são registrados no módulo de documentos, evitando duplicar datas no perfil.

O vínculo do motorista não é reserva exclusiva do veículo, escala nem viagem. Mudanças desse vínculo registram o veículo anterior e o atual na auditoria. Histórico de viagens e de KM continuam pendentes.

Documentos registram título, número, emissor, datas opcionais de emissão/validade, situação (`active`/`archived`) e observações. Cada documento pertence a exatamente um motorista ou um veículo. Tipos de cadastro: `cnh` (motorista), `crlv` (veículo), `insurance` ou `other` (qualquer um dos dois). O proprietário é imutável depois da criação; para corrigir um proprietário, arquive/exclua logicamente o registro e crie outro. Não há upload/download de arquivos ou verificação externa nesta etapa.

## API

Todas as rotas exigem Bearer token. A empresa é obtida da autenticação; `tenant_id` enviado no corpo não altera o contexto.

| Método | Rota | Uso |
| --- | --- | --- |
| POST | `/api/v1/drivers` | Criar motorista |
| GET | `/api/v1/drivers` | Listar com paginação/filtros |
| GET | `/api/v1/drivers/:id` | Consultar motorista |
| PUT | `/api/v1/drivers/:id` | Atualizar cadastro completo com `version` |
| DELETE | `/api/v1/drivers/:id?version=N` | Excluir logicamente com controle de versão |
| POST | `/api/v1/documents` | Registrar documento |
| GET | `/api/v1/documents` | Listar com paginação/filtros/vencimentos |
| GET | `/api/v1/documents/:id` | Consultar documento |
| PUT | `/api/v1/documents/:id` | Atualizar cadastro completo com `version` |
| DELETE | `/api/v1/documents/:id?version=N` | Excluir logicamente com controle de versão |

Permissões: `drivers.read`, `drivers.manage`, `documents.read`, `documents.manage`. Consultar exige `read`; gravar exige `manage`. A migration concede as quatro permissões somente a papéis `admin` já existentes que possuíam `fleet.vehicles.manage`. Papéis personalizados não recebem poderes novos automaticamente. Novos administradores provisionados recebem as quatro permissões, sem administração global de empresas/permissões.

Exemplo de criação de motorista:

```json
{"code":"MOT-001","name":"Motorista de teste","status":"active","vehicle_id":"UUID_DO_VEICULO"}
```

Exemplo de documento:

```json
{"driver_id":"UUID_DO_MOTORISTA","kind":"cnh","title":"CNH cadastrada","number":"NUMERO_INFORMADO","issued_on":"2024-01-01","expires_on":"2027-01-01"}
```

Para documento de veículo, envie `vehicle_id` e omita `driver_id`. Para omitir datas, não envie os campos ou use `null`; string vazia não é data válida. PUT precisa dos campos do cadastro e da última versão retornada. O backend define IDs, empresa, timestamps, exclusão e versão; valores desses campos no corpo não substituem o contexto/controladores do servidor.

Filtros de motoristas: `search` (nome/código), `status`, `vehicle_id`, `page` e `page_size`. Ordenação estável por criação/ID.

Filtros de documentos: `search` (título/número), `kind`, `status`, `driver_id`, `vehicle_id`, `expiry_status`, `page` e `page_size`. Não combinar os dois proprietários. Ordenação por vencimento, datas nulas por último e ID para desempate. Página padrão 1, tamanho padrão 20/máximo 100. Páginas vazias preservam a contagem total. Busca e filtros são limitados à empresa autenticada.

## Vencimentos

`expiry_status` é calculado a cada consulta com a data civil UTC do servidor; não é gravado nem aceito como decisão enviada pelo cliente:

- `expired`: validade anterior à data atual.
- `due_soon`: vence hoje ou até 30 dias depois, inclusive.
- `valid`: validade posterior à janela de 30 dias.
- `no_expiry`: data não informada; não representa confirmação de regularidade.
- `archived`: registro arquivado, independente da data.

Datas são `YYYY-MM-DD`, sem hora. Validade não pode anteceder a emissão. Essas classificações organizam as datas informadas pelo operador; não consultam órgãos externos nem habilitam/bloqueiam viagens automaticamente.

## Integridade e auditoria

A migration `000021_drivers_documents_core.sql` adiciona tabelas nativas, chaves compostas por empresa, unicidade dos códigos/identificadores de motoristas e dos números de documentos ativos por proprietário/tipo. Registros arquivados preservam histórico e permitem registrar outro documento com o mesmo número. Exclusões são lógicas. Documentos históricos continuam armazenados e podem ser editados/arquivados quando o proprietário é excluído logicamente, mantendo o mesmo proprietário.

Create/update/delete dos dois módulos e seus eventos de auditoria são gravados na mesma transação. Se a auditoria falhar, a alteração operacional é revertida. Os eventos registram empresa, ator, ação, recurso, IDs de vínculos e datas/situação anteriores e atuais dos documentos, sem copiar números de documentos, identificação fiscal, e-mail ou observações para os metadados da auditoria.

Respostas: 400 para campos/filtros inválidos; 401 sem autenticação; 403 sem permissão; 404 para recurso/proprietário/veículo ausente ou de outra empresa; 409 para duplicidade/versão antiga; 500 para falha interna, sem expor SQL ou dados de conexão.

## Executar e retomar

Siga `docs/PERSISTENCIA_IDENTIDADE_FROTA.md` para configurar PostgreSQL/segredo. Antes de iniciar a API, execute em `backend/api`:

```sh
go run ./cmd/migrate
```

O conjunto ativo agora é 000001, 000002, 000003, 000020 e 000021. Empresas já provisionadas recebem a atualização de permissões pela migration, sem criar outra empresa ou trocar a senha. Não repita o comando de provisionamento para atualizar uma empresa.

Consulte `docs/CONTINUIDADE.md` para o estado dos testes, limites de autorização e próxima etapa. Frontend, upload de arquivos, alertas externos, escala/viagens e validação com cliente permanecem pendentes.


## Evidência de validação local

Ambiente: PostgreSQL 16 em Docker, base descartável `gi360_validation`, schemas aleatórios por teste e dados sintéticos. Não foi usado banco remoto nem dados do cliente.

Saída real da suíte completa nesta etapa:

```text
$ TEST_DATABASE_URL=... go test -race ./... -count=1 -timeout=120s
ok  github.com/geevheit/intelligence360/backend/api/internal/integration 46.858s
ok  github.com/geevheit/intelligence360/backend/api/internal/modules/documents/application 5.859s
ok  github.com/geevheit/intelligence360/backend/api/internal/modules/documents/domain 5.127s
```

Depois do ajuste final para arquivar documentos de proprietários excluídos logicamente e preservar datas anteriores na auditoria, os testes afetados e toda a integração foram repetidos:

```text
$ TEST_DATABASE_URL=... go test -race ./internal/modules/documents/... ./internal/integration -count=1 -timeout=120s
ok  github.com/geevheit/intelligence360/backend/api/internal/modules/documents/application 3.592s
ok  github.com/geevheit/intelligence360/backend/api/internal/modules/documents/domain 2.152s
ok  github.com/geevheit/intelligence360/backend/api/internal/integration 48.273s
```

Ambos terminaram com **exit code 0**, sem falha reportada pelo detector de concorrência. A integração inclui os três testes da etapa anterior e cinco novos: `TestDriversDocumentsLifecycle`, `TestOperationalValidationAndPermissions`, `TestOperationalAuditRollback`, `TestOperationalConcurrentCreate` e `TestOperationalMigrationUpgrade`.

Os dois arquivos OpenAPI foram lidos por um parser YAML; as 165 referências locais de `backend/api/openapi.yaml` e as 237 de `backend/api/openapi/openapi.yaml` foram resolvidas. Foi corrigida uma referência antiga a `PathID`, que não existia, para o parâmetro `ID` já definido. Isso não substitui uma validação completa por ferramenta OpenAPI nem a futura consolidação dos dois contratos.

`git diff --check` passou. Interface, upload, serviços externos, homologação com cliente e produção não foram validados nesta etapa.
