# Persistência local — identidade e frota

Etapa implementada em 07/11/2026. Este documento descreve o escopo persistente atual; o diagnóstico inicial permanece como registro da situação anterior.

> Etapa seguinte implementada: [motoristas e documentos](MOTORISTAS_DOCUMENTOS.md). O conjunto ativo passou a incluir 000021; a configuração e os comandos abaixo continuam válidos.

## Escopo disponível

- Identidade: empresas, usuários, papéis, permissões, vínculos, sessões, refresh tokens e repositório de auditoria.
- Frota: veículos, marcas, modelos, categorias, tipos e ativos de frota (`fleet_assets`).
- UUIDs compatíveis com as tabelas existentes; senhas de novos usuários com bcrypt.
- Isolamento por empresa nas consultas e chaves estrangeiras compostas nos vínculos de identidade/frota.
- Exclusão lógica e controle otimista de versão nos cadastros de frota.
- Login, refresh e logout persistidos. Rotação de refresh transacional com bloqueio da linha: o token anterior só pode ser usado uma vez. Logout e desativação do usuário invalidam o acesso.
- Provisionamento explícito cria empresa, administrador, papel, permissões e evento de auditoria na mesma transação. O administrador da empresa não recebe permissões globais de administração de empresas/permissões.

## Modos da API

`STORAGE_DRIVER=postgres` é o padrão. Exige `DATABASE_URL` e `AUTH_SECRET` com pelo menos 32 caracteres. Expõe identidade, frota, motoristas e documentos. As rotas dos demais módulos retornam 404 porque seus adapters ainda são voláteis.

`STORAGE_DRIVER=memory` preserva a demonstração anterior com todos os módulos. É permitido apenas fora de `APP_ENV=production`. Dados não sobrevivem ao reinício. `scripts/dev-api.sh` seleciona esse modo explicitamente quando nenhum modo é informado.

A API não executa DDL ou seeds no startup. No modo PostgreSQL, verifica a presença das migrations necessárias e encerra se estiverem pendentes. Não cria conta nem senha de demonstração.

## Execução local com PostgreSQL

Na raiz do repositório, disponibilize o PostgreSQL local:

```sh
docker compose -f infra/docker/docker-compose.yml up -d postgres
cd backend/api
export DATABASE_URL='postgres://geevheit:geevheit_dev@localhost:5432/geevheit_intelligence?sslmode=disable'
export STORAGE_DRIVER=postgres
```

A URL acima usa as credenciais de desenvolvimento do Compose e deve ser usada somente localmente. Defina `AUTH_SECRET` no ambiente com um valor aleatório próprio, sem registrá-lo no Git.

Execute o conjunto de migrations desta etapa:

```sh
go run ./cmd/migrate
```

O comando aplica **somente** `000001`, `000002`, `000003`, `000020` e `000021`. Registra checksums em `schema_migrations`, usa lock transacional e permite reexecução sem repetir migrations. O lote inteiro é revertido se alguma etapa falhar. A seed `000004_dev_seed.sql` não é executada.

Antes do provisionamento, defina `PROVISION_ADMIN_PASSWORD` no ambiente por entrada privada ou gerenciador de segredos. O comando não aceita senha como argumento e não imprime a senha:

```sh
go run ./cmd/provision -name 'Empresa de teste' -email 'gestor@example.test'
unset PROVISION_ADMIN_PASSWORD
go run ./cmd/api
```

Guarde o `tenant_id` retornado pelo provisionamento. O login precisa desse ID e das credenciais escolhidas. Cada execução de provisionamento cria uma nova empresa; não é comando de atualização nem de redefinição de senha. Não repetir para recuperar acesso a uma empresa existente.

## Contrato de atualização de frota

Ao fazer PUT, envie os campos do cadastro e a `version` recebida no GET/listagem ou na resposta de criação. A API incrementa a versão. Uma gravação com versão anterior retorna **409**. Referências inválidas retornam **400**; ID de outra empresa retorna **404**; falta de autenticação retorna **401** e falta de permissão retorna **403**.

## Validação

Os testes em `backend/api/internal/integration/postgres_test.go` exigem `TEST_DATABASE_URL` de um **banco descartável local**. Criam schemas aleatórios, instalam as extensões necessárias no schema público dessa base de teste e removem apenas os schemas que criaram. Não apontar para produção ou dados do cliente.

```sh
cd backend/api
TEST_DATABASE_URL='postgres://postgres@127.0.0.1:55439/gi360_validation?sslmode=disable' \
  go test -race ./... -count=1 -timeout=120s
```

Sem `TEST_DATABASE_URL`, os testes de integração são explicitamente ignorados; o sucesso da suíte nesse caso não comprova PostgreSQL.

Resultado executado nesta revisão, com PostgreSQL 16 em Docker local e dados sintéticos:

```text
$ TEST_DATABASE_URL=... go test -race ./... -count=1 -timeout=120s
ok  github.com/geevheit/intelligence360/backend/api/internal/integration 18.312s
ok  github.com/geevheit/intelligence360/backend/api/internal/modules/identity/application 6.691s
ok  github.com/geevheit/intelligence360/backend/api/internal/modules/fleet/application 6.514s
```

A suíte completa terminou com **exit code 0**, incluindo os três testes PostgreSQL (`TestPostgresRuntime`, `TestPostgresTransactionRollback`, `TestPostgresIdentityRelations`). O detector de concorrência não reportou falhas. Houve ensaios anteriores que falharam na migration histórica 000018 e no isolamento inicial do fixture de teste; a validação acima foi executada após delimitar as migrations do runtime e corrigir o fixture.

Cobertura de integração: migrations e reexecução, provisionamento, bcrypt, cadastros, filtros/paginação, retenção ao fechar e recriar pool/roteador, isolamento de empresas, chaves estrangeiras, controle de versão, exclusão lógica, refresh concorrente, logout, usuário/empresa inativos, remoção de permissão com papel excluído e rollback transacional.

## Pendências preservadas

- Nenhuma interface foi construída nesta etapa.
- Os demais módulos ainda precisam de persistência e ativação das suas migrations, além da conexão dos fluxos operacionais.
- Ao ensaiar o conjunto histórico completo em banco descartável, `000018_intelligence_core.sql` falhou com `syntax error at or near "window"`. Esse conjunto não está homologado; nenhuma migration histórica foi alterada. Será necessário definir a correção e o procedimento de aplicação ao ativar inteligência.
- O repositório de auditoria e o evento de provisionamento são duráveis; o middleware HTTP ainda registra logs técnicos. Auditoria de todas as mutações operacionais continua pendente.
- Gestão completa de atribuições de papéis pela interface/API, recuperação de senha, rate limiting e importação de dados ficam para as etapas de identidade operacional e interface.
- Integração validada em PostgreSQL Docker local com dados sintéticos. Não houve migração de banco remoto, publicação, teste de backup ou homologação com o cliente.
