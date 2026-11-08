# Preparação para entrega ao cliente

Diagnóstico e critérios detalhados: [docs/DIAGNOSTICO_ENTREGA.md](docs/DIAGNOSTICO_ENTREGA.md).

- [x] Revisar estrutura, módulos, composição da API e aplicações.
- [x] Identificar bloqueios e organizar critérios de entrega.
- [ ] Definir e aprovar o escopo inicial com o cliente.
- [x] Implementar persistência PostgreSQL de identidade e frota (primeira etapa).
- [ ] Estender persistência a abastecimentos, manutenção, checklist e demais módulos do escopo.
- [x] Validar migrations de identidade/frota e transações em PostgreSQL descartável local.
- [ ] Validar migrations dos demais módulos; corrigir bloqueio SQL da migration 000018 antes de ativar inteligência.
- [ ] Preparar identidade, senhas, perfis e auditoria para operação.
- [x] Implementar motoristas, vínculo atual com veículo e metadados de documentos em PostgreSQL.
- [x] Validar vencimentos, permissões, atualização de empresas existentes e auditoria transacional em banco local.
- [ ] Implementar upload/armazenamento de arquivos dos documentos quando definido no escopo.
- [ ] Concluir histórico de KM e vínculos entre fluxos operacionais.
- [ ] Construir painel web conectado à API.
- [ ] Conectar Mission Control aos sinais operacionais do escopo.
- [ ] Consolidar contrato OpenAPI oficial.
- [ ] Validar fluxos HTTP, banco real de teste e interface de ponta a ponta.
- [ ] Preparar homologação, monitoramento e backup com restauração testada.
- [ ] Validar importação dos dados do cliente.
- [ ] Executar piloto, treinar usuários e obter aceite de produção.

Esta lista registra pendências; não autoriza implantação ou contratação de serviços.

Execução desta etapa e limites: [docs/PERSISTENCIA_IDENTIDADE_FROTA.md](docs/PERSISTENCIA_IDENTIDADE_FROTA.md).
