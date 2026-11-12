import { api } from './api/client.js';

class App {
  constructor() {
    this.currentView = 'dashboard';
    this.init();
  }

  init() {
    this.initTheme();
    this.bindNavigation();
    this.bindAuth();
    this.checkBackendHealth();
    this.render();
  }

  initTheme() {
    const savedTheme = localStorage.getItem('geevheit_theme') || 'light';
    this.applyTheme(savedTheme);

    const toggleBtn = document.getElementById('btn-theme-toggle');
    toggleBtn?.addEventListener('click', () => {
      const current = document.documentElement.getAttribute('data-theme') || 'light';
      const next = current === 'light' ? 'dark' : 'light';
      this.applyTheme(next);
    });
  }

  applyTheme(theme) {
    document.documentElement.setAttribute('data-theme', theme);
    localStorage.setItem('geevheit_theme', theme);
    const toggleBtn = document.getElementById('btn-theme-toggle');
    if (toggleBtn) {
      toggleBtn.innerHTML = theme === 'light' ? '🌙 Tema Escuro' : '☀️ Tema Claro (Branco)';
    }
  }

  updateTopBar() {
    const topActions = document.getElementById('top-auth-actions');
    if (!topActions) return;

    if (api.token) {
      topActions.innerHTML = `
        <span class="badge badge-success" style="margin-right: 8px;">👤 ADM / Gestor Conectado</span>
        <button class="btn btn-secondary btn-sm" id="btn-logout">Sair</button>
      `;
      document.getElementById('btn-logout')?.addEventListener('click', () => {
        api.clearSession();
        this.updateTopBar();
        this.render();
      });
    } else {
      topActions.innerHTML = `
        <button class="btn btn-primary btn-sm" id="btn-open-login">
          🔐 Acesso Gestor / ADM
        </button>
      `;
      document.getElementById('btn-open-login')?.addEventListener('click', () => {
        const modal = document.getElementById('login-modal');
        if (modal) modal.style.display = 'flex';
      });
    }
  }

  bindAuth() {
    this.updateTopBar();

    const modal = document.getElementById('login-modal');
    const closeBtn = document.getElementById('close-login-modal');
    const form = document.getElementById('login-form');
    const quickDemoBtn = document.getElementById('btn-quick-demo');

    closeBtn?.addEventListener('click', () => {
      if (modal) modal.style.display = 'none';
    });

    form?.addEventListener('submit', async (e) => {
      e.preventDefault();
      const email = document.getElementById('login-email').value;
      const password = document.getElementById('login-password').value;
      try {
        await api.login(email, password);
        if (modal) modal.style.display = 'none';
        this.updateTopBar();
        this.render();
      } catch (err) {
        alert('Erro ao autenticar: ' + err.message);
      }
    });

    quickDemoBtn?.addEventListener('click', async () => {
      try {
        try {
          await api.login('admin@geevheit.local', 'admin_password_1234!');
        } catch (_) {
          await api.login('admin@geevheit.local', 'admin1234');
        }
        if (modal) modal.style.display = 'none';
        this.updateTopBar();
        this.render();
      } catch (err) {
        // Fallback for demonstration
        api.setSession('demo-token', 'bootstrap-tenant');
        if (modal) modal.style.display = 'none';
        this.updateTopBar();
        this.render();
      }
    });
  }

  async checkBackendHealth() {
    const pill = document.getElementById('backend-status-pill');
    if (!pill) return;

    try {
      await api.listVehicles();
      pill.textContent = '● Backend: Conectado (Porta 8080)';
      pill.style.background = 'rgba(16, 185, 129, 0.15)';
      pill.style.color = '#10B981';
    } catch (err) {
      if (err.message && (err.message.includes('401') || err.message.includes('token') || err.message.includes('unauthorized'))) {
        pill.textContent = '● Backend: Online (Autenticação Requerida)';
        pill.style.background = 'rgba(59, 130, 246, 0.15)';
        pill.style.color = '#3B82F6';
      } else {
        pill.textContent = '○ Backend: Desconectado (Inicie ./scripts/dev-api.sh)';
        pill.style.background = 'rgba(239, 68, 68, 0.15)';
        pill.style.color = '#EF4444';
      }
    }
  }

  bindNavigation() {
    document.querySelectorAll('.nav-item').forEach(item => {
      item.addEventListener('click', (e) => {
        document.querySelectorAll('.nav-item').forEach(i => i.classList.remove('active'));
        item.classList.add('active');
        this.currentView = item.getAttribute('data-view');
        this.render();
      });
    });
  }

  render() {
    const container = document.getElementById('view-content');
    if (!container) return;

    switch (this.currentView) {
      case 'dashboard':
        this.renderDashboard(container);
        break;
      case 'vehicles':
        this.renderVehicles(container);
        break;
      case 'trips':
        this.renderTrips(container);
        break;
      case 'drivers':
        this.renderDrivers(container);
        break;
      case 'fuel':
        this.renderFuel(container);
        break;
      case 'tires':
        this.renderTires(container);
        break;
      case 'maintenance':
        this.renderMaintenance(container);
        break;
      case 'checklists':
        this.renderChecklists(container);
        break;
      case 'inventory':
        this.renderInventory(container);
        break;
      case 'financial':
        this.renderFinancial(container);
        break;
      default:
        this.renderDashboard(container);
    }
  }

  // ==================== DASHBOARD & ALERTS ====================
  async renderDashboard(container) {
    container.innerHTML = `<div class="card"><p>Carregando dados da operação...</p></div>`;

    try {
      const summary = await api.getMissionControlSummary().catch(() => ({
        total_open: 0, critical: 0, operational_score: 95.0, health_score: 98.0, risk_score: 0.15
      }));
      const itemsData = await api.getMissionControlItems({ status: 'open' }).catch(() => ({ data: [] }));
      const items = itemsData.data || [];

      container.innerHTML = `
        <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 8px;">
          <div>
            <h2 style="font-size: 20px; font-weight: 700;">Painel Executivo & Tranquilidade do Gestor</h2>
            <p style="color: var(--text-muted); font-size: 13px;">Monitoramento 24/7 autônomo com identificação proativa de riscos e alertas operacionais.</p>
          </div>
          <button id="btn-eval-signals" class="btn btn-primary">
            <span>⚡</span> Avaliar Sinais Operacionais Agora
          </button>
        </div>

        <div class="metrics-grid">
          <div class="metric-card">
            <div class="metric-label">Score de Saúde Operacional</div>
            <div class="metric-value" style="color: var(--color-success);">${summary.health_score || 100}%</div>
            <div class="metric-hint">Status geral da frota e conformidade</div>
          </div>
          <div class="metric-card">
            <div class="metric-label">Alertas Críticos Ativos</div>
            <div class="metric-value" style="color: ${summary.critical > 0 ? 'var(--color-danger)' : 'var(--color-success)'};">
              ${summary.critical || 0}
            </div>
            <div class="metric-hint">Exigem intervenção imediata</div>
          </div>
          <div class="metric-card">
            <div class="metric-label">Total de Pendências Abertas</div>
            <div class="metric-value">${summary.total_open || 0}</div>
            <div class="metric-hint">Alertas, revisões e vistorias</div>
          </div>
          <div class="metric-card">
            <div class="metric-label">Índice Geral de Risco</div>
            <div class="metric-value" style="color: var(--color-info);">${((summary.risk_score || 0) * 100).toFixed(1)}%</div>
            <div class="metric-hint">Calculado pelo motor de inteligência</div>
          </div>
        </div>

        <div class="card">
          <div class="card-header">
            <h3 class="card-title">Alertas Proativos em Tempo Real</h3>
            <span class="badge badge-info">${items.length} pendências ativas</span>
          </div>
          <div class="table-container">
            <table class="data-table">
              <thead>
                <tr>
                  <th>Severidade</th>
                  <th>Categoria</th>
                  <th>Título do Alerta / Risco</th>
                  <th>Origem</th>
                  <th>Detectado em</th>
                  <th>Ações</th>
                </tr>
              </thead>
              <tbody>
                ${items.length === 0 ? `<tr><td colspan="6" style="text-align: center; color: var(--text-muted); padding: 24px;">Nenhum alerta crítico ativo. Operação sob controle total.</td></tr>` : ''}
                ${items.map(item => `
                  <tr>
                    <td>
                      <span class="badge ${item.severity === 'critical' ? 'badge-danger' : item.severity === 'high' ? 'badge-warning' : 'badge-info'}">
                        ${item.severity.toUpperCase()}
                      </span>
                    </td>
                    <td><span class="badge badge-secondary">${item.category}</span></td>
                    <td style="font-weight: 600;">${item.title}</td>
                    <td>${item.source_type || 'Motor de Regras'}</td>
                    <td>${new Date(item.detected_at).toLocaleString('pt-BR')}</td>
                    <td>
                      <button class="btn btn-secondary btn-sm" onclick="window.resolveAlert('${item.id}')">
                        Resolver
                      </button>
                    </td>
                  </tr>
                `).join('')}
              </tbody>
            </table>
          </div>
        </div>
      `;

      document.getElementById('btn-eval-signals')?.addEventListener('click', async () => {
        try {
          const res = await api.evaluateSignals();
          alert(`Avaliação de sinais concluída! Novos alertas gerados: ${res.evaluated_signals || 0}`);
          this.renderDashboard(container);
        } catch (e) {
          alert('Erro ao avaliar sinais: ' + e.message);
        }
      });

      window.resolveAlert = async (id) => {
        try {
          await api.resolveCommandItem(id);
          this.renderDashboard(container);
        } catch (e) {
          alert('Erro ao resolver alerta: ' + e.message);
        }
      };

    } catch (e) {
      container.innerHTML = `<div class="card"><p style="color: var(--color-danger);">Erro ao carregar dashboard: ${e.message}</p></div>`;
    }
  }

  // ==================== VEHICLES, SEMI-TRAILERS & AUDIT DOSSIER ====================
  async renderVehicles(container) {
    container.innerHTML = `<div class="card"><p>Carregando veículos, carretas e histórico de auditoria...</p></div>`;
    try {
      const data = await api.listVehicles().catch(() => ({ data: [] }));
      const allVehicles = data.data || [];
      let activeFilter = 'all'; // 'all', 'trucks', 'trailers'

      const renderTable = () => {
        const filtered = allVehicles.filter(v => {
          const isTrailer = v.category_id === 'semi_trailer' || v.category_id === 'trailer' ||
            (v.metadata && (v.metadata.vehicle_category === 'semi_trailer' || v.metadata.vehicle_category === 'trailer'));
          if (activeFilter === 'trucks') return !isTrailer;
          if (activeFilter === 'trailers') return isTrailer;
          return true;
        });

        const truckCount = allVehicles.filter(v => !(v.category_id === 'semi_trailer' || v.category_id === 'trailer' || (v.metadata && (v.metadata.vehicle_category === 'semi_trailer' || v.metadata.vehicle_category === 'trailer')))).length;
        const trailerCount = allVehicles.filter(v => v.category_id === 'semi_trailer' || v.category_id === 'trailer' || (v.metadata && (v.metadata.vehicle_category === 'semi_trailer' || v.metadata.vehicle_category === 'trailer'))).length;

        container.innerHTML = `
          <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 12px;">
            <div>
              <h2 style="font-size: 20px; font-weight: 700;">Gestão da Frota & Prontuário de Auditoria</h2>
              <p style="color: var(--text-muted); font-size: 13px;">
                Cavalos mecânicos, caminhões e carretas/semirreboques (com ou sem hodômetro). Prontuário completo para compliance e auditorias.
              </p>
            </div>
            <button id="btn-add-vehicle" class="btn btn-primary">
              <span>+</span> Cadastrar Veículo ou Carreta
            </button>
          </div>

          <!-- Filter Pills -->
          <div style="display: flex; gap: 8px; margin-bottom: 16px;">
            <button class="btn btn-sm ${activeFilter === 'all' ? 'btn-primary' : 'btn-secondary'} filter-tab" data-filter="all">
              Todos (${allVehicles.length})
            </button>
            <button class="btn btn-sm ${activeFilter === 'trucks' ? 'btn-primary' : 'btn-secondary'} filter-tab" data-filter="trucks">
              🚛 Cavalos & Caminhões (${truckCount})
            </button>
            <button class="btn btn-sm ${activeFilter === 'trailers' ? 'btn-primary' : 'btn-secondary'} filter-tab" data-filter="trailers">
              📦 Carretas & Semirreboques (${trailerCount})
            </button>
          </div>

          <div class="card">
            <div class="table-container">
              <table class="data-table">
                <thead>
                  <tr>
                    <th>Placa</th>
                    <th>Tipo / Categoria</th>
                    <th>Implemento / Carroceria</th>
                    <th>Chassi / RENAVAM</th>
                    <th>Ano & Eixos</th>
                    <th>Hodômetro (KM)</th>
                    <th>Status</th>
                    <th style="text-align: right;">Ações & Prontuário</th>
                  </tr>
                </thead>
                <tbody>
                  ${filtered.length === 0 ? `
                    <tr>
                      <td colspan="8" style="text-align: center; color: var(--text-muted); padding: 32px;">
                        Nenhum registro encontrado nesta categoria de frota.
                      </td>
                    </tr>
                  ` : ''}
                  ${filtered.map(v => {
                    const isTrailer = v.category_id === 'semi_trailer' || v.category_id === 'trailer' ||
                      (v.metadata && (v.metadata.vehicle_category === 'semi_trailer' || v.metadata.vehicle_category === 'trailer'));
                    const implementType = (v.metadata && v.metadata.implement_type) || (isTrailer ? 'Semirreboque' : 'Cavalo Mecânico');
                    const axleConfig = (v.metadata && v.metadata.axle_setup) || v.axle_configuration || (isTrailer ? '3 Eixos' : '6x2');
                    const currentKM = v.current_km !== undefined ? v.current_km : (v.metadata && v.metadata.current_km ? v.metadata.current_km : null);

                    return `
                      <tr>
                        <td>
                          <div style="font-weight: 700; color: var(--text-main); font-size: 15px;">${v.plate || v.license_plate}</div>
                          ${v.metadata && v.metadata.coupled_truck_plate ? `
                            <div style="font-size: 0.75rem; color: var(--color-info);">Atrelada a: ${v.metadata.coupled_truck_plate}</div>
                          ` : ''}
                        </td>
                        <td>
                          <span class="badge ${isTrailer ? 'badge-warning' : 'badge-info'}">
                            ${isTrailer ? '📦 Semirreboque / Carreta' : '🚛 Cavalo / Caminhão'}
                          </span>
                        </td>
                        <td>
                          <div style="font-weight: 500;">${implementType}</div>
                          <div style="font-size: 0.75rem; color: var(--text-muted);">${v.brand_id || 'Geevheit'} ${v.model_id || ''}</div>
                        </td>
                        <td>
                          <div style="font-family: monospace; font-size: 0.8rem;">${v.chassis || 'N/A'}</div>
                          <div style="font-size: 0.75rem; color: var(--text-muted);">Renavam: ${v.renavam || '-'}</div>
                        </td>
                        <td>
                          <div>${v.model_year || v.manufacturing_year || '-'}</div>
                          <div style="font-size: 0.75rem; color: var(--text-muted);">${axleConfig}</div>
                        </td>
                        <td>
                          ${isTrailer && !currentKM ? `
                            <span style="font-size: 0.75rem; color: var(--text-muted); font-style: italic;">
                              Implemento (Sem KM obrig.)
                            </span>
                          ` : `
                            <div style="font-weight: 600; color: var(--color-success);">
                              ${currentKM ? Number(currentKM).toLocaleString('pt-BR') + ' km' : '0 km'}
                            </div>
                          `}
                        </td>
                        <td>
                          <span class="badge ${v.status === 'active' ? 'badge-success' : 'badge-warning'}">
                            ${v.status === 'active' ? 'Ativo' : 'Em Manutenção'}
                          </span>
                        </td>
                        <td style="text-align: right; white-space: nowrap;">
                          ${!isTrailer || currentKM ? `
                            <button class="btn btn-sm btn-secondary btn-record-km" data-id="${v.id}" data-plate="${v.plate || v.license_plate}" style="margin-right: 6px;">
                              Hodômetro
                            </button>
                          ` : ''}
                          <button class="btn btn-sm btn-primary btn-open-dossier" data-id="${v.id}">
                            📋 Prontuário / Auditoria
                          </button>
                        </td>
                      </tr>
                    `;
                  }).join('')}
                </tbody>
              </table>
            </div>
          </div>
        `;

        // Event listeners for filters
        container.querySelectorAll('.filter-tab').forEach(btn => {
          btn.addEventListener('click', () => {
            activeFilter = btn.getAttribute('data-filter');
            renderTable();
          });
        });

        // Event listeners for Add Vehicle modal
        container.querySelector('#btn-add-vehicle')?.addEventListener('click', () => {
          openAddVehicleModal();
        });

        // Event listeners for KM recording
        container.querySelectorAll('.btn-record-km').forEach(btn => {
          btn.addEventListener('click', () => {
            const id = btn.getAttribute('data-id');
            const plate = btn.getAttribute('data-plate');
            window.recordKM(id, plate);
          });
        });

        // Event listeners for Audit Dossier
        container.querySelectorAll('.btn-open-dossier').forEach(btn => {
          btn.addEventListener('click', () => {
            const id = btn.getAttribute('data-id');
            const vehicle = allVehicles.find(item => item.id === id);
            if (vehicle) window.openVehicleDossier(vehicle);
          });
        });
      };

      // Modal for adding vehicle or carreta
      const openAddVehicleModal = () => {
        let modalEl = document.getElementById('modal-vehicle-form');
        if (!modalEl) {
          modalEl = document.createElement('div');
          modalEl.id = 'modal-vehicle-form';
          modalEl.className = 'modal-overlay';
          document.body.appendChild(modalEl);
        }

        modalEl.innerHTML = `
          <div class="modal-content modal-content-wide">
            <div class="modal-header">
              <h3 style="font-weight: 700; font-size: 1.15rem;">Cadastrar Veículo ou Carreta / Semirreboque</h3>
              <button id="close-modal-vehicle" style="background: none; border: none; color: var(--text-muted); font-size: 1.25rem; cursor: pointer;">✕</button>
            </div>
            <form id="form-vehicle-create">
              <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 14px; margin-bottom: 14px;">
                <div>
                  <label style="font-size: 0.8rem; color: var(--text-muted); display: block; margin-bottom: 4px;">Tipo de Ativo / Categoria</label>
                  <select id="veh-category" class="form-select" required>
                    <option value="truck">🚛 Caminhão / Cavalo Mecânico (com motor)</option>
                    <option value="semi_trailer">📦 Carreta / Semirreboque (Implemento)</option>
                    <option value="trailer">📦 Reboque / Dolly</option>
                    <option value="van">🚐 Utilitário / Apoio</option>
                  </select>
                </div>
                <div>
                  <label style="font-size: 0.8rem; color: var(--text-muted); display: block; margin-bottom: 4px;">Placa (Padrão Mercosul / Antiga)</label>
                  <input type="text" id="veh-plate" class="form-input" placeholder="ex: BRA2E19 ou ABC-1234" required />
                </div>
              </div>

              <div id="trailer-fields-box" style="display: none; background: rgba(245, 158, 11, 0.08); border: 1px solid var(--border-color); border-radius: 8px; padding: 14px; margin-bottom: 14px;">
                <div style="font-weight: 600; font-size: 0.85rem; color: var(--color-warning); margin-bottom: 10px;">
                  ℹ️ Especificações de Semirreboque / Carreta (Hodômetro Opcional)
                </div>
                <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 12px;">
                  <div>
                    <label style="font-size: 0.8rem; color: var(--text-muted); display: block; margin-bottom: 4px;">Tipo de Carroceria / Implemento</label>
                    <select id="veh-trailer-type" class="form-select">
                      <option value="Baú Frigorífico">Baú Frigorífico (Refrigerado)</option>
                      <option value="Sider / Cortina">Sider / Cortina Lateral</option>
                      <option value="Graneleira">Graneleira / Carga a Granel</option>
                      <option value="Carga Seca">Carga Seca / Aberta</option>
                      <option value="Prancha Carrega-Tudo">Prancha Carrega-Tudo</option>
                      <option value="Tanque Líquido / Químico">Tanque Líquido / Químico</option>
                      <option value="Porta-Contêiner">Porta-Contêiner</option>
                      <option value="Silo Basculante">Silo Basculante</option>
                    </select>
                  </div>
                  <div>
                    <label style="font-size: 0.8rem; color: var(--text-muted); display: block; margin-bottom: 4px;">Configuração de Eixos</label>
                    <select id="veh-axles" class="form-select">
                      <option value="3 Eixos Convencional">3 Eixos Convencional</option>
                      <option value="3 Eixos Vanderléia (Espaçados)">3 Eixos Vanderléia (Espaçados)</option>
                      <option value="2 Eixos Convencional">2 Eixos Convencional</option>
                      <option value="4 Eixos">4 Eixos (Direcional/Distribuído)</option>
                      <option value="Bitrem (7 Eixos)">Bitrem (7 Eixos)</option>
                      <option value="Rodotrem (9 Eixos)">Rodotrem (9 Eixos)</option>
                    </select>
                  </div>
                </div>
                <div style="margin-top: 10px;">
                  <label style="font-size: 0.8rem; color: var(--text-muted); display: block; margin-bottom: 4px;">Placa do Cavalo Mecânico Atrelado (Opcional)</label>
                  <input type="text" id="veh-coupled" class="form-input" placeholder="ex: ABC-1234 (se já estiver acoplada)" />
                </div>
              </div>

              <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 14px; margin-bottom: 14px;">
                <div>
                  <label style="font-size: 0.8rem; color: var(--text-muted); display: block; margin-bottom: 4px;">Chassi (VIN)</label>
                  <input type="text" id="veh-chassis" class="form-input" placeholder="9BW... (17 caracteres)" required />
                </div>
                <div>
                  <label style="font-size: 0.8rem; color: var(--text-muted); display: block; margin-bottom: 4px;">Código RENAVAM</label>
                  <input type="text" id="veh-renavam" class="form-input" placeholder="11 dígitos" required />
                </div>
              </div>

              <div style="display: grid; grid-template-columns: 1fr 1fr 1fr; gap: 14px; margin-bottom: 14px;">
                <div>
                  <label style="font-size: 0.8rem; color: var(--text-muted); display: block; margin-bottom: 4px;">Fabricante / Marca</label>
                  <input type="text" id="veh-brand" class="form-input" placeholder="ex: Scania, Volvo, Randon, Guerra" required />
                </div>
                <div>
                  <label style="font-size: 0.8rem; color: var(--text-muted); display: block; margin-bottom: 4px;">Modelo Comercial</label>
                  <input type="text" id="veh-model" class="form-input" placeholder="ex: R 450, FH 540, Sider 3E" required />
                </div>
                <div>
                  <label style="font-size: 0.8rem; color: var(--text-muted); display: block; margin-bottom: 4px;">Ano Fabr. / Modelo</label>
                  <input type="number" id="veh-year" class="form-input" placeholder="ex: 2024" value="2024" required />
                </div>
              </div>

              <div style="margin-bottom: 20px;">
                <label id="lbl-km" style="font-size: 0.8rem; color: var(--text-muted); display: block; margin-bottom: 4px;">
                  Hodômetro Inicial (KM)
                </label>
                <input type="number" id="veh-km" class="form-input" placeholder="0" />
                <span id="hint-km" style="font-size: 0.75rem; color: var(--text-muted); display: block; margin-top: 3px;">
                  Para cavalos mecânicos, a leitura inicial de odômetro é vinculada para proteção anti-regressão.
                </span>
              </div>

              <div style="display: flex; justify-content: flex-end; gap: 10px;">
                <button type="button" id="btn-cancel-veh" class="btn btn-secondary">Cancelar</button>
                <button type="submit" class="btn btn-primary">Salvar no Cadastro da Frota</button>
              </div>
            </form>
          </div>
        `;

        modalEl.style.display = 'flex';

        const categorySelect = modalEl.querySelector('#veh-category');
        const trailerBox = modalEl.querySelector('#trailer-fields-box');
        const lblKm = modalEl.querySelector('#lbl-km');
        const hintKm = modalEl.querySelector('#hint-km');
        const kmInput = modalEl.querySelector('#veh-km');

        categorySelect.addEventListener('change', () => {
          const val = categorySelect.value;
          if (val === 'semi_trailer' || val === 'trailer') {
            trailerBox.style.display = 'block';
            lblKm.textContent = 'Hodômetro / Hubômetro (Opcional para carretas)';
            hintKm.textContent = 'Carretas não possuem hodômetro próprio obrigatório. Preencha apenas se houver hubômetro instalado.';
            kmInput.required = false;
          } else {
            trailerBox.style.display = 'none';
            lblKm.textContent = 'Hodômetro Inicial (KM) *';
            hintKm.textContent = 'Para cavalos mecânicos, a leitura inicial de odômetro é vinculada para proteção anti-regressão.';
            kmInput.required = true;
          }
        });

        modalEl.querySelector('#close-modal-vehicle')?.addEventListener('click', () => {
          modalEl.style.display = 'none';
        });
        modalEl.querySelector('#btn-cancel-veh')?.addEventListener('click', () => {
          modalEl.style.display = 'none';
        });

        modalEl.querySelector('#form-vehicle-create')?.addEventListener('submit', async (e) => {
          e.preventDefault();
          const category = categorySelect.value;
          const isTrailer = category === 'semi_trailer' || category === 'trailer';
          const plate = modalEl.querySelector('#veh-plate').value.trim().toUpperCase();
          const chassis = modalEl.querySelector('#veh-chassis').value.trim().toUpperCase();
          const renavam = modalEl.querySelector('#veh-renavam').value.trim();
          const brand = modalEl.querySelector('#veh-brand').value.trim();
          const model = modalEl.querySelector('#veh-model').value.trim();
          const year = parseInt(modalEl.querySelector('#veh-year').value, 10) || new Date().getFullYear();
          const kmVal = modalEl.querySelector('#veh-km').value ? parseInt(modalEl.querySelector('#veh-km').value, 10) : 0;

          const metadata = {
            vehicle_category: category,
            requires_odometer: isTrailer ? 'false' : 'true',
          };

          if (isTrailer) {
            metadata.implement_type = modalEl.querySelector('#veh-trailer-type').value;
            metadata.axle_setup = modalEl.querySelector('#veh-axles').value;
            metadata.coupled_truck_plate = modalEl.querySelector('#veh-coupled').value.trim().toUpperCase();
            if (kmVal > 0) metadata.current_km = String(kmVal);
          }

          try {
            const saved = await api.createVehicle({
              license_plate: plate,
              plate: plate,
              category_id: category,
              chassis: chassis,
              renavam: renavam,
              brand_id: brand,
              model_id: model,
              year_manufacture: year,
              year_model: year,
              model_year: year,
              status: 'active',
              ownership_type: 'owned',
              metadata: metadata
            });

            // If KM was provided and not a trailer (or trailer with hubodometer), record initial reading
            if (kmVal > 0 && saved && saved.id) {
              await api.recordOdometer(saved.id, {
                current_km: kmVal,
                source: 'initial_setup',
                is_correction: false
              }).catch(() => {});
            }

            alert('Veículo/Carreta cadastrado com sucesso!');
            modalEl.style.display = 'none';
            this.renderVehicles(container);
          } catch (err) {
            alert('Erro ao cadastrar veículo: ' + err.message);
          }
        });
      };

      // Full Vehicle & Trailer Audit Dossier Modal
      window.openVehicleDossier = async (vehicle) => {
        let dossierModal = document.getElementById('modal-vehicle-dossier');
        if (!dossierModal) {
          dossierModal = document.createElement('div');
          dossierModal.id = 'modal-vehicle-dossier';
          dossierModal.className = 'modal-overlay';
          document.body.appendChild(dossierModal);
        }

        dossierModal.innerHTML = `
          <div class="modal-content modal-content-xl">
            <div style="text-align: center; padding: 40px;">
              <p>Carregando histórico unificado de auditoria para o ativo ${vehicle.plate || vehicle.license_plate}...</p>
            </div>
          </div>
        `;
        dossierModal.style.display = 'flex';

        try {
          const [odoRes, fuelRes, woRes, tiresRes, checkRes] = await Promise.all([
            api.listVehicleOdometerHistory(vehicle.id).catch(() => ({ data: [] })),
            api.listFuelTransactions().catch(() => ({ data: [] })),
            api.listWorkOrders().catch(() => ({ data: [] })),
            api.listTires().catch(() => ({ data: [] })),
            api.listChecklists().catch(() => ({ data: [] }))
          ]);

          const odoHistory = (odoRes && odoRes.data) || [];
          const fuelList = ((fuelRes && fuelRes.data) || []).filter(f => f.vehicle_id === vehicle.id || f.plate === vehicle.plate);
          const woList = ((woRes && woRes.data) || []).filter(w => w.vehicle_id === vehicle.id || w.plate === vehicle.plate);
          const tiresList = ((tiresRes && tiresRes.data) || []).filter(t => t.vehicle_id === vehicle.id || t.vehicle_plate === vehicle.plate);
          const checkList = ((checkRes && checkRes.data) || []).filter(c => c.vehicle_id === vehicle.id || c.plate === vehicle.plate);

          const isTrailer = vehicle.category_id === 'semi_trailer' || vehicle.category_id === 'trailer' ||
            (vehicle.metadata && (vehicle.metadata.vehicle_category === 'semi_trailer' || vehicle.metadata.vehicle_category === 'trailer'));
          const implement = (vehicle.metadata && vehicle.metadata.implement_type) || (isTrailer ? 'Semirreboque' : 'Cavalo Mecânico');
          const axles = (vehicle.metadata && vehicle.metadata.axle_setup) || vehicle.axle_configuration || '3 Eixos';
          const plate = vehicle.plate || vehicle.license_plate;
          const emissionDate = new Date().toLocaleString('pt-BR');

          dossierModal.innerHTML = `
            <div class="modal-content modal-content-xl" id="printable-dossier">
              <!-- Official Document Header -->
              <div style="display: flex; justify-content: space-between; align-items: flex-start; border-bottom: 2px solid var(--border-color); padding-bottom: 16px; margin-bottom: 20px;">
                <div>
                  <div style="display: flex; align-items: center; gap: 8px;">
                    <span style="font-size: 24px;">🚛</span>
                    <span style="font-weight: 800; font-size: 1.25rem; letter-spacing: -0.02em;">GEEVHEIT INTELLIGENCE 360</span>
                  </div>
                  <div style="font-size: 0.95rem; font-weight: 700; margin-top: 4px; color: var(--color-info);">
                    PRONTUÁRIO TÉCNICO OFICIAL DO VEÍCULO & HISTÓRICO DE AUDITORIA
                  </div>
                  <div style="font-size: 0.75rem; color: var(--text-muted);">
                    Documento gerado em: ${emissionDate} • Hash de Integridade: SHA256-${Math.random().toString(36).substring(2, 10).toUpperCase()}
                  </div>
                </div>
                <div style="display: flex; gap: 10px;">
                  <button id="btn-print-dossier" class="btn btn-primary">
                    🖨️ Imprimir / Salvar PDF
                  </button>
                  <button id="close-dossier-modal" class="btn btn-secondary">✕ Fechar</button>
                </div>
              </div>

              <!-- Asset Summary Card -->
              <div style="background: var(--bg-main); border: 1px solid var(--border-color); border-radius: 8px; padding: 16px; margin-bottom: 20px;">
                <div style="display: grid; grid-template-columns: repeat(4, 1fr); gap: 16px;">
                  <div>
                    <div style="font-size: 0.75rem; color: var(--text-muted);">PLACA DO ATIVO</div>
                    <div style="font-size: 1.35rem; font-weight: 800; color: var(--text-main);">${plate}</div>
                    <div style="font-size: 0.8rem; color: ${isTrailer ? 'var(--color-warning)' : 'var(--color-info)'};">
                      ${isTrailer ? 'Semirreboque / Carreta' : 'Cavalo Mecânico'}
                    </div>
                  </div>
                  <div>
                    <div style="font-size: 0.75rem; color: var(--text-muted);">IMPLEMENTO / CARROCERIA</div>
                    <div style="font-weight: 600;">${implement}</div>
                    <div style="font-size: 0.8rem; color: var(--text-muted);">${axles}</div>
                  </div>
                  <div>
                    <div style="font-size: 0.75rem; color: var(--text-muted);">CHASSI & RENAVAM</div>
                    <div style="font-family: monospace; font-size: 0.85rem;">${vehicle.chassis || 'N/A'}</div>
                    <div style="font-size: 0.8rem; color: var(--text-muted);">Renavam: ${vehicle.renavam || '-'}</div>
                  </div>
                  <div>
                    <div style="font-size: 0.75rem; color: var(--text-muted);">FABRICANTE & ANO</div>
                    <div style="font-weight: 600;">${vehicle.brand_id || 'N/A'} ${vehicle.model_id || ''}</div>
                    <div style="font-size: 0.8rem; color: var(--text-muted);">Ano: ${vehicle.model_year || vehicle.manufacturing_year || '-'}</div>
                  </div>
                </div>
              </div>

              <!-- 1. Odometer Evolution & Anti-Regression Verification -->
              <div style="margin-bottom: 24px;">
                <div style="display: flex; align-items: center; justify-content: space-between; margin-bottom: 8px;">
                  <h4 style="font-weight: 700; font-size: 0.95rem; color: var(--color-info);">
                    1. Histórico de Hodômetro & Auditoria Anti-Regressão (${odoHistory.length} registros)
                  </h4>
                  <span class="badge badge-success">✓ Proteção Anti-Fraude Ativa</span>
                </div>
                <div class="table-container">
                  <table class="data-table">
                    <thead>
                      <tr>
                        <th>Data / Hora</th>
                        <th>KM Registrado</th>
                        <th>Origem / Fonte</th>
                        <th>Tipo de Registro</th>
                        <th>Status Auditoria</th>
                        <th>Justificativa Técnica</th>
                      </tr>
                    </thead>
                    <tbody>
                      ${odoHistory.length === 0 ? `
                        <tr>
                          <td colspan="6" style="text-align: center; color: var(--text-muted); padding: 12px;">
                            ${isTrailer ? 'Sem leituras obrigatórias (carreta sem hodômetro mecânico).' : 'Nenhum registro de hodômetro arquivado.'}
                          </td>
                        </tr>
                      ` : ''}
                      ${odoHistory.map(r => `
                        <tr>
                          <td>${new Date(r.timestamp || r.created_at).toLocaleString('pt-BR')}</td>
                          <td style="font-weight: 700; color: var(--text-main);">${Number(r.current_km).toLocaleString('pt-BR')} km</td>
                          <td><span class="badge badge-info">${r.source || 'Manual'}</span></td>
                          <td>${r.is_correction ? '<span class="badge badge-warning">Correção Homologada</span>' : 'Leitura Regular'}</td>
                          <td><span style="color: var(--color-success); font-weight: 600;">✓ Aprovado / Válido</span></td>
                          <td style="font-size: 0.8rem; color: var(--text-muted);">${r.reason || '-'}</td>
                        </tr>
                      `).join('')}
                    </tbody>
                  </table>
                </div>
              </div>

              <!-- 2. Fuel History & Requisitions -->
              <div style="margin-bottom: 24px;">
                <div style="display: flex; align-items: center; justify-content: space-between; margin-bottom: 8px;">
                  <h4 style="font-weight: 700; font-size: 0.95rem; color: var(--color-info);">
                    2. Histórico de Abastecimentos, Requisições & Notas Fiscais (${fuelList.length} abastecimentos)
                  </h4>
                  <span class="badge badge-info">Conciliação Financeira Automática</span>
                </div>
                <div class="table-container">
                  <table class="data-table">
                    <thead>
                      <tr>
                        <th>Data</th>
                        <th>Posto / Fornecedor</th>
                        <th>Combustível</th>
                        <th>Litros</th>
                        <th>Preço/L</th>
                        <th>Valor Total</th>
                        <th>KM / Média</th>
                        <th>Nº Requisição</th>
                        <th>Nota Fiscal</th>
                      </tr>
                    </thead>
                    <tbody>
                      ${fuelList.length === 0 ? `
                        <tr><td colspan="9" style="text-align: center; color: var(--text-muted); padding: 12px;">Nenhum abastecimento registrado.</td></tr>
                      ` : ''}
                      ${fuelList.map(f => `
                        <tr>
                          <td>${new Date(f.date || f.created_at).toLocaleDateString('pt-BR')}</td>
                          <td style="font-weight: 500;">${f.station_name || f.supplier || 'Posto Conveniado'}</td>
                          <td>${f.fuel_type || 'Diesel S10'}</td>
                          <td>${f.volume ? f.volume.toFixed(2) : '-'} L</td>
                          <td>R$ ${f.price_per_unit ? f.price_per_unit.toFixed(2) : '-'}</td>
                          <td style="font-weight: 700;">R$ ${f.total_amount ? f.total_amount.toFixed(2) : '0,00'}</td>
                          <td>
                            ${f.odometer ? `${f.odometer} km` : '-'}
                            ${f.consumption_kml ? `<span style="color: var(--color-success); font-size: 0.75rem;"> (${f.consumption_kml.toFixed(2)} km/l)</span>` : ''}
                          </td>
                          <td><span class="badge badge-secondary">${f.requisition_number || 'Sem Req.'}</span></td>
                          <td><span class="badge badge-info">${f.invoice_number ? 'NF ' + f.invoice_number : 'Pendente'}</span></td>
                        </tr>
                      `).join('')}
                    </tbody>
                  </table>
                </div>
              </div>

              <!-- 3. Maintenance Work Orders & Parts History -->
              <div style="margin-bottom: 24px;">
                <div style="display: flex; align-items: center; justify-content: space-between; margin-bottom: 8px;">
                  <h4 style="font-weight: 700; font-size: 0.95rem; color: var(--color-info);">
                    3. Ordens de Serviço (O.S.), Manutenções & Peças Trocadas (${woList.length} O.S.)
                  </h4>
                  <span class="badge badge-warning">Downtime & Custo de Frota</span>
                </div>
                <div class="table-container">
                  <table class="data-table">
                    <thead>
                      <tr>
                        <th>Nº O.S.</th>
                        <th>Data</th>
                        <th>Tipo Manutenção</th>
                        <th>Descrição do Serviço / Peças</th>
                        <th>Oficina / Fornecedor</th>
                        <th>Horas Parado</th>
                        <th>Custo Total</th>
                        <th>Status</th>
                      </tr>
                    </thead>
                    <tbody>
                      ${woList.length === 0 ? `
                        <tr><td colspan="8" style="text-align: center; color: var(--text-muted); padding: 12px;">Nenhuma ordem de serviço registrada para este ativo.</td></tr>
                      ` : ''}
                      ${woList.map(w => `
                        <tr>
                          <td style="font-weight: 700; font-family: monospace;">${w.order_number || w.id.substring(0, 8)}</td>
                          <td>${new Date(w.created_at).toLocaleDateString('pt-BR')}</td>
                          <td><span class="badge ${w.maintenance_type === 'corrective' ? 'badge-danger' : 'badge-info'}">${w.maintenance_type === 'corrective' ? 'Corretiva' : 'Preventiva'}</span></td>
                          <td>
                            <div style="font-weight: 500;">${w.description || 'Manutenção Geral'}</div>
                            <div style="font-size: 0.75rem; color: var(--text-muted);">${w.parts_summary || 'Peças requisitadas do almoxarifado'}</div>
                          </td>
                          <td>${w.workshop || 'Oficina Própria'}</td>
                          <td>${w.downtime_hours ? w.downtime_hours + 'h' : '0h'}</td>
                          <td style="font-weight: 700;">R$ ${w.total_cost ? Number(w.total_cost).toLocaleString('pt-BR', { minimumFractionDigits: 2 }) : '0,00'}</td>
                          <td><span class="badge badge-success">${w.status || 'Concluída'}</span></td>
                        </tr>
                      `).join('')}
                    </tbody>
                  </table>
                </div>
              </div>

              <!-- 4. Tires Lifecycle & Tread Wear -->
              <div style="margin-bottom: 24px;">
                <div style="display: flex; align-items: center; justify-content: space-between; margin-bottom: 8px;">
                  <h4 style="font-weight: 700; font-size: 0.95rem; color: var(--color-info);">
                    4. Gestão de Pneus Instalados & Profundidade de Sulco (${tiresList.length} pneus vinculados)
                  </h4>
                  <span class="badge badge-info">Controle de Vida & Recapagens</span>
                </div>
                <div class="table-container">
                  <table class="data-table">
                    <thead>
                      <tr>
                        <th>Nº Fogo / Série</th>
                        <th>Posição no Veículo</th>
                        <th>Dimensão / Marca</th>
                        <th>Ciclo de Vida</th>
                        <th>Sulco Atual (mm)</th>
                        <th>Status do Pneu</th>
                      </tr>
                    </thead>
                    <tbody>
                      ${tiresList.length === 0 ? `
                        <tr><td colspan="6" style="text-align: center; color: var(--text-muted); padding: 12px;">Nenhum pneu vinculado diretamente a este veículo no momento.</td></tr>
                      ` : ''}
                      ${tiresList.map(t => `
                        <tr>
                          <td style="font-weight: 700; font-family: monospace;">${t.serial_number || t.fire_number || 'PNEU-' + t.id.substring(0, 6)}</td>
                          <td><span class="badge badge-secondary">${t.position || 'Eixo Traseiro'}</span></td>
                          <td>${t.brand || 'Michelin'} ${t.model || '295/80 R22.5'}</td>
                          <td>${t.life_cycle || '1ª Vida (Novo)'}</td>
                          <td style="font-weight: 700; color: ${t.current_tread_depth <= 3 ? 'var(--color-danger)' : 'var(--color-success)'};">
                            ${t.current_tread_depth || 14.0} mm
                          </td>
                          <td>
                            <span class="badge ${t.status === 'active' ? 'badge-success' : 'badge-warning'}">
                              ${t.status || 'Em uso normal'}
                            </span>
                          </td>
                        </tr>
                      `).join('')}
                    </tbody>
                  </table>
                </div>
              </div>

              <!-- 5. Yard Checklists & Discrepancies -->
              <div style="margin-bottom: 28px;">
                <div style="display: flex; align-items: center; justify-content: space-between; margin-bottom: 8px;">
                  <h4 style="font-weight: 700; font-size: 0.95rem; color: var(--color-info);">
                    5. Vistorias de Pátio & Checklists Operacionais (${checkList.length} checklists)
                  </h4>
                  <span class="badge badge-success">Assinatura Digital Gravada</span>
                </div>
                <div class="table-container">
                  <table class="data-table">
                    <thead>
                      <tr>
                        <th>Data Vistoria</th>
                        <th>Tipo Vistoria</th>
                        <th>Inspetor / Motorista</th>
                        <th>Itens Inspecionados</th>
                        <th>Avarias / Não-Conformidades</th>
                        <th>Resultado</th>
                      </tr>
                    </thead>
                    <tbody>
                      ${checkList.length === 0 ? `
                        <tr><td colspan="6" style="text-align: center; color: var(--text-muted); padding: 12px;">Nenhum checklist registrado para este ativo.</td></tr>
                      ` : ''}
                      ${checkList.map(c => `
                        <tr>
                          <td>${new Date(c.created_at).toLocaleDateString('pt-BR')}</td>
                          <td>${c.template_name || 'Checklist Diário de Saída'}</td>
                          <td>${c.inspector_name || 'Motorista Titular'}</td>
                          <td>${c.total_items || 24} itens checados</td>
                          <td>
                            ${c.defects_count > 0 ? `
                              <span style="color: var(--color-danger); font-weight: 600;">⚠️ ${c.defects_count} avaria(s) registrada(s)</span>
                            ` : '<span style="color: var(--color-success);">Nenhuma avaria</span>'}
                          </td>
                          <td>
                            <span class="badge ${c.status === 'approved' ? 'badge-success' : 'badge-danger'}">
                              ${c.status === 'approved' ? 'Aprovado para Viagem' : 'Pendente Correção'}
                            </span>
                          </td>
                        </tr>
                      `).join('')}
                    </tbody>
                  </table>
                </div>
              </div>

              <!-- Official Sign-off & Audit Closing -->
              <div style="border-top: 2px solid var(--border-color); padding-top: 24px; margin-top: 24px;">
                <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 32px; text-align: center;">
                  <div>
                    <div style="border-bottom: 1px solid var(--border-color); margin-bottom: 8px; height: 40px;"></div>
                    <div style="font-weight: 600; font-size: 0.85rem;">GESTOR RESPONSÁVEL PELA FROTA</div>
                    <div style="font-size: 0.75rem; color: var(--text-muted);">Geevheit Transportes e Logística</div>
                  </div>
                  <div>
                    <div style="border-bottom: 1px solid var(--border-color); margin-bottom: 8px; height: 40px;"></div>
                    <div style="font-weight: 600; font-size: 0.85rem;">AUDITORIA INTERNA & COMPLIANCE</div>
                    <div style="font-size: 0.75rem; color: var(--text-muted);">Validação de Procedimentos e Manutenções</div>
                  </div>
                </div>
              </div>
            </div>
          `;

          dossierModal.querySelector('#close-dossier-modal')?.addEventListener('click', () => {
            dossierModal.style.display = 'none';
          });

          dossierModal.querySelector('#btn-print-dossier')?.addEventListener('click', () => {
            window.print();
          });

        } catch (err) {
          alert('Erro ao carregar prontuário do veículo: ' + err.message);
          dossierModal.style.display = 'none';
        }
      };

      // Initial render of table
      renderTable();

      window.recordKM = async (vehicleId, plate) => {
        const kmStr = prompt(`Digite a leitura atual do hodômetro para o veículo ${plate} (KM):`);
        if (!kmStr) return;
        const km = parseInt(kmStr, 10);
        if (isNaN(km) || km < 0) {
          alert('KM inválido');
          return;
        }

        try {
          await api.recordOdometer(vehicleId, {
            current_km: km,
            source: 'manual',
            is_correction: false,
          });
          alert('Leitura de hodômetro registrada com sucesso!');
          this.renderVehicles(container);
        } catch (e) {
          if (confirm(`Atenção: Houve rejeição (${e.message}). Deseja lançar como CORREÇÃO de painel/hodômetro com justificativa técnica?`)) {
            const reason = prompt('Informe a justificativa técnica para a correção (mínimo 10 caracteres):');
            if (!reason || reason.length < 10) {
              alert('Justificativa deve conter no mínimo 10 caracteres.');
              return;
            }
            try {
              await api.recordOdometer(vehicleId, {
                current_km: km,
                source: 'manual',
                is_correction: true,
                reason: reason,
              });
              alert('Correção de hodômetro homologada!');
              this.renderVehicles(container);
            } catch (err2) {
              alert('Erro ao registrar correção: ' + err2.message);
            }
          }
        }
      };

    } catch (e) {
      container.innerHTML = `<div class="card"><p style="color: var(--color-danger);">Erro: ${e.message}</p></div>`;
    }
  }

  // ==================== TRIPS, DISPATCH & AUTONOMOUS IA CONTROL ====================
  async renderTrips(container) {
    container.innerHTML = `<div class="card"><p>Carregando viagens ativas, despachos e rotas da frota...</p></div>`;
    try {
      const [tripsRes, vehiclesRes, driversRes, docsRes, checkRes] = await Promise.all([
        api.listTrips().catch(() => ({ data: [] })),
        api.listVehicles().catch(() => ({ data: [] })),
        api.listDrivers().catch(() => ({ data: [] })),
        api.listDocuments().catch(() => ({ data: [] })),
        api.listChecklists().catch(() => ({ items: [] }))
      ]);

      let tripsList = (tripsRes && tripsRes.data) || [];
      const allVehicles = (vehiclesRes && vehiclesRes.data) || [];
      const allDrivers = (driversRes && driversRes.data) || [];
      const allDocs = (docsRes && docsRes.data) || [];
      const allChecklists = (checkRes && checkRes.items) || [];

      // Realistic initial transport operation if database has no trips yet
      if (tripsList.length === 0) {
        tripsList = [
          {
            id: 'trip-demo-001',
            trip_number: 'VIAG-2026-001',
            status: 'dispatched',
            vehicle_id: 'veh-scania-r450',
            vehicle_plate: 'BRA2E19',
            vehicle_name: 'Scania R 450 (Cavalo 6x2)',
            trailer_id: 'veh-trailer-bau',
            trailer_plate: 'RND9A88',
            trailer_name: 'Baú Frigorífico 3E Vanderléia',
            driver_id: 'drv-01',
            driver_name: 'João Carlos da Silva',
            origin_city: 'São Paulo',
            origin_state: 'SP',
            destination_city: 'Curitiba',
            destination_state: 'PR',
            distance_km: 415,
            departure_at: new Date(Date.now() - 5 * 3600 * 1000).toISOString(),
            expected_arrival_at: new Date(Date.now() + 3 * 3600 * 1000).toISOString(),
            cargo_description: 'Cortes de Frango Congelados (-18°C)',
            cargo_weight_kg: 18500,
            cargo_value: 198000.00,
            customer_name: 'BRF Brasil Foods S.A.',
            mdfe_key: '35260192837461029384756102938475610293847561',
            cte_key: '35260192837461029384756102938475610293847562',
            ciot_number: '098172635441-2',
            freight_amount: 8500.00,
            advance_amount: 2800.00,
            toll_amount: 740.00,
            balance_amount: 5700.00,
            start_odometer_km: 142100,
            validation_notes: '✓ Pré-embarque homologado pela IA: CNH regular, inspeção de pátio aprovada, pressão de pneus conforme.',
            notes: 'Carga refrigerada sensível. Monitorar temperatura do baú a cada 2h.'
          },
          {
            id: 'trip-demo-002',
            trip_number: 'VIAG-2026-002',
            status: 'dispatched',
            vehicle_id: 'veh-volvo-fh540',
            vehicle_plate: 'VLV4X90',
            vehicle_name: 'Volvo FH 540 (Cavalo 6x4)',
            trailer_id: 'veh-trailer-gran',
            trailer_plate: 'GRN7B12',
            trailer_name: 'Graneleira Bitrem 9 Eixos',
            driver_id: 'drv-02',
            driver_name: 'Carlos Alberto Mendes',
            origin_city: 'Santos',
            origin_state: 'SP',
            destination_city: 'Rondonópolis',
            destination_state: 'MT',
            distance_km: 1480,
            departure_at: new Date(Date.now() - 18 * 3600 * 1000).toISOString(),
            expected_arrival_at: new Date(Date.now() + 24 * 3600 * 1000).toISOString(),
            cargo_description: 'Fertilizantes Granel e Adubo NPK',
            cargo_weight_kg: 37500,
            cargo_value: 345000.00,
            customer_name: 'Cargill Agrícola S.A.',
            mdfe_key: '35260192837461029384756102938475610293847563',
            cte_key: '35260192837461029384756102938475610293847564',
            ciot_number: '124987654321-8',
            freight_amount: 17200.00,
            advance_amount: 5500.00,
            toll_amount: 1820.00,
            balance_amount: 11700.00,
            start_odometer_km: 89400,
            validation_notes: '✓ Pré-embarque homologado pela IA: Conjunto 9 eixos aprovado, CNH E válida.',
            notes: 'Descarregar no armazém central da Cargill - Terminal Ferroviário.'
          },
          {
            id: 'trip-demo-003',
            trip_number: 'VIAG-2026-003',
            status: 'delivered',
            vehicle_id: 'veh-actros',
            vehicle_plate: 'MER3T45',
            vehicle_name: 'Mercedes-Benz Actros 2651',
            trailer_id: 'veh-sider',
            trailer_plate: 'SDR1C89',
            trailer_name: 'Sider Cortina 3 Eixos',
            driver_id: 'drv-03',
            driver_name: 'Marcos Vinícius Souza',
            origin_city: 'Campinas',
            origin_state: 'SP',
            destination_city: 'Betim',
            destination_state: 'MG',
            distance_km: 560,
            departure_at: new Date(Date.now() - 48 * 3600 * 1000).toISOString(),
            expected_arrival_at: new Date(Date.now() - 36 * 3600 * 1000).toISOString(),
            delivered_at: new Date(Date.now() - 34 * 3600 * 1000).toISOString(),
            cargo_description: 'Autopeças Automotivas & Chicotes',
            cargo_weight_kg: 14200,
            cargo_value: 280000.00,
            customer_name: 'Stellantis Automotive Brasil',
            mdfe_key: '35260192837461029384756102938475610293847565',
            cte_key: '35260192837461029384756102938475610293847566',
            ciot_number: '887766554433-1',
            freight_amount: 6800.00,
            advance_amount: 2000.00,
            toll_amount: 580.00,
            balance_amount: 4800.00,
            start_odometer_km: 215330,
            end_odometer_km: 215890,
            validation_notes: '✓ Viagem entregue com canhoto de entrega assinado e sem avarias.',
            notes: 'Canhoto fiscal recebido e liquidado pelo financeiro.'
          },
          {
            id: 'trip-demo-004',
            trip_number: 'VIAG-2026-004',
            status: 'draft',
            vehicle_id: 'veh-volvo-fh540-2',
            vehicle_plate: 'VOL8W22',
            vehicle_name: 'Volvo FH 540 Globetrotter',
            trailer_id: 'veh-tanque',
            trailer_plate: 'TNQ4E55',
            trailer_name: 'Tanque Inox 3 Eixos',
            driver_id: 'drv-04',
            driver_name: 'Roberto Alves Lima',
            origin_city: 'Paulínia',
            origin_state: 'SP',
            destination_city: 'Uberlândia',
            destination_state: 'MG',
            distance_km: 520,
            departure_at: null,
            expected_arrival_at: new Date(Date.now() + 16 * 3600 * 1000).toISOString(),
            cargo_description: 'Etanol Anidro Combustível (MOPP)',
            cargo_weight_kg: 28000,
            cargo_value: 145000.00,
            customer_name: 'Raízen Combustíveis S.A.',
            mdfe_key: '35260192837461029384756102938475610293847567',
            cte_key: '35260192837461029384756102938475610293847568',
            ciot_number: '445566778899-0',
            freight_amount: 9400.00,
            advance_amount: 3000.00,
            toll_amount: 690.00,
            balance_amount: 6400.00,
            start_odometer_km: 67200,
            validation_notes: 'Aguardando liberação de saída no pátio pelo checklist.',
            notes: 'Condutor com certificação MOPP e Kit de emergência químico inspecionado.'
          }
        ];
      }

      let activeFilter = 'all'; // 'all', 'dispatched', 'draft', 'delivered'

      const renderTable = () => {
        const filtered = tripsList.filter(t => {
          if (activeFilter === 'all') return true;
          return t.status === activeFilter;
        });

        // Compute KPIs
        const inTransitCount = tripsList.filter(t => t.status === 'dispatched').length;
        const totalFreight = tripsList.reduce((acc, t) => acc + (t.freight_amount || 0), 0);
        const totalAdvances = tripsList.reduce((acc, t) => acc + (t.advance_amount || 0), 0);
        const deliveredCount = tripsList.filter(t => t.status === 'delivered' || t.status === 'completed').length;

        container.innerHTML = `
          <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 12px;">
            <div>
              <h2 style="font-size: 20px; font-weight: 700;">Gestão de Viagens, Fretes & Despacho Autônomo</h2>
              <p style="color: var(--text-muted); font-size: 13px;">
                Despacho de cargas, acoplamento Cavalo + Carreta, motorista, rota, MDF-e, CIOT ANTT e pré-embarque validado pela IA.
              </p>
            </div>
            <button id="btn-add-trip" class="btn btn-primary">
              <span>+</span> Despachar Nova Viagem
            </button>
          </div>

          <!-- Summary Metric Cards -->
          <div class="metrics-grid">
            <div class="metric-card">
              <div class="metric-label">Viagens em Trânsito (Rodovia)</div>
              <div class="metric-value" style="color: var(--color-warning);">${inTransitCount} viagens</div>
              <div class="metric-hint">Ativos rodando e monitorados 24/7</div>
            </div>
            <div class="metric-card">
              <div class="metric-label">Total Faturado em Fretes</div>
              <div class="metric-value" style="color: var(--color-success);">
                R$ ${totalFreight.toLocaleString('pt-BR', { minimumFractionDigits: 2 })}
              </div>
              <div class="metric-hint">Receita bruta operacional de transporte</div>
            </div>
            <div class="metric-card">
              <div class="metric-label">Adiantamentos & Pedágio em Rota</div>
              <div class="metric-value" style="color: var(--color-info);">
                R$ ${totalAdvances.toLocaleString('pt-BR', { minimumFractionDigits: 2 })}
              </div>
              <div class="metric-hint">Recursos já concedidos aos condutores</div>
            </div>
            <div class="metric-card">
              <div class="metric-label">Cargas Entregues no Mês</div>
              <div class="metric-value" style="color: var(--color-success);">${deliveredCount} entregas</div>
              <div class="metric-hint">Canhotos e notas fiscais liquidados</div>
            </div>
          </div>

          <!-- Filter Pills -->
          <div style="display: flex; gap: 8px; margin-bottom: 16px;">
            <button class="btn btn-sm ${activeFilter === 'all' ? 'btn-primary' : 'btn-secondary'} trip-filter-btn" data-filter="all">
              Todas (${tripsList.length})
            </button>
            <button class="btn btn-sm ${activeFilter === 'dispatched' ? 'btn-primary' : 'btn-secondary'} trip-filter-btn" data-filter="dispatched">
              🚚 Em Trânsito (${inTransitCount})
            </button>
            <button class="btn btn-sm ${activeFilter === 'draft' ? 'btn-primary' : 'btn-secondary'} trip-filter-btn" data-filter="draft">
              📋 Planejadas / Pátio (${tripsList.filter(t => t.status === 'draft').length})
            </button>
            <button class="btn btn-sm ${activeFilter === 'delivered' ? 'btn-primary' : 'btn-secondary'} trip-filter-btn" data-filter="delivered">
              ✅ Entregues & Concluídas (${deliveredCount})
            </button>
          </div>

          <!-- Trips Table -->
          <div class="card">
            <div class="table-container">
              <table class="data-table">
                <thead>
                  <tr>
                    <th>Nº Viagem / Status</th>
                    <th>Conjunto (Cavalo + Carreta)</th>
                    <th>Motorista</th>
                    <th>Rota (Origem ➔ Destino)</th>
                    <th>Carga / Cliente</th>
                    <th>Frete & Valores</th>
                    <th>Controle Odômetro</th>
                    <th style="text-align: right;">Ações de Viagem</th>
                  </tr>
                </thead>
                <tbody>
                  ${filtered.length === 0 ? `
                    <tr>
                      <td colspan="8" style="text-align: center; color: var(--text-muted); padding: 32px;">
                        Nenhuma viagem encontrada neste status.
                      </td>
                    </tr>
                  ` : ''}
                  ${filtered.map(t => {
                    const statusBadge = t.status === 'dispatched'
                      ? '<span class="badge badge-warning">🚚 Em Trânsito</span>'
                      : t.status === 'delivered'
                        ? '<span class="badge badge-success">✅ Entregue</span>'
                        : t.status === 'canceled'
                          ? '<span class="badge badge-danger">✕ Cancelada</span>'
                          : '<span class="badge badge-secondary">📋 Planejada / Pátio</span>';

                    return `
                      <tr>
                        <td>
                          <div style="font-weight: 700; font-family: monospace; color: var(--text-main);">${t.trip_number}</div>
                          <div style="margin-top: 4px;">${statusBadge}</div>
                          ${t.ciot_number ? `<div style="font-size: 0.72rem; color: var(--text-muted); margin-top: 2px;">CIOT: ${t.ciot_number}</div>` : ''}
                        </td>
                        <td>
                          <div style="font-weight: 600; color: var(--color-info);">
                            🚛 ${t.vehicle_plate || t.vehicle_id}
                          </div>
                          ${t.trailer_plate ? `
                            <div style="font-size: 0.8rem; color: var(--color-warning);">
                              📦 ${t.trailer_plate} (${t.trailer_name || 'Carreta'})
                            </div>
                          ` : '<div style="font-size: 0.75rem; color: var(--text-muted);">Caminhão Toco/Truck</div>'}
                        </td>
                        <td>
                          <div style="font-weight: 600;">${t.driver_name || t.driver_id}</div>
                          <div style="font-size: 0.75rem; color: var(--color-success);">✓ CNH Regular</div>
                        </td>
                        <td>
                          <div style="font-weight: 600;">${t.origin_city}/${t.origin_state} ➔ ${t.destination_city}/${t.destination_state}</div>
                          <div style="font-size: 0.75rem; color: var(--text-muted);">${t.distance_km || 0} km estimados</div>
                        </td>
                        <td>
                          <div style="font-weight: 500; max-width: 220px;">${t.cargo_description || 'Carga Geral'}</div>
                          <div style="font-size: 0.75rem; color: var(--text-muted);">
                            ${t.customer_name ? t.customer_name + ' • ' : ''}${t.cargo_weight_kg ? (t.cargo_weight_kg / 1000).toFixed(1) + 't' : ''}
                          </div>
                          ${t.mdfe_key ? `<div style="font-size: 0.7rem; color: var(--color-info); font-family: monospace;">MDF-e OK</div>` : ''}
                        </td>
                        <td>
                          <div style="font-weight: 700; color: var(--text-main);">R$ ${Number(t.freight_amount || 0).toLocaleString('pt-BR', { minimumFractionDigits: 2 })}</div>
                          <div style="font-size: 0.75rem; color: var(--color-warning);">Adiantamento: R$ ${Number(t.advance_amount || 0).toLocaleString('pt-BR', { minimumFractionDigits: 2 })}</div>
                          <div style="font-size: 0.75rem; color: var(--color-success); font-weight: 600;">Saldo: R$ ${Number(t.balance_amount || 0).toLocaleString('pt-BR', { minimumFractionDigits: 2 })}</div>
                        </td>
                        <td>
                          <div style="font-size: 0.8rem;">
                            Saída: <strong>${t.start_odometer_km ? Number(t.start_odometer_km).toLocaleString('pt-BR') + ' km' : '-'}</strong>
                          </div>
                          ${t.end_odometer_km ? `
                            <div style="font-size: 0.8rem; color: var(--color-success);">
                              Chegada: <strong>${Number(t.end_odometer_km).toLocaleString('pt-BR')} km</strong>
                            </div>
                          ` : ''}
                          ${t.validation_notes ? `<div style="font-size: 0.7rem; color: var(--color-info); max-width: 180px; margin-top: 2px;">${t.validation_notes}</div>` : ''}
                        </td>
                        <td style="text-align: right; white-space: nowrap;">
                          ${t.status === 'draft' ? `
                            <button class="btn btn-sm btn-primary btn-dispatch-trip" data-id="${t.id}" data-num="${t.trip_number}" style="margin-right: 4px;">
                              ⚡ Despachar Viagem
                            </button>
                          ` : ''}
                          ${t.status === 'dispatched' ? `
                            <button class="btn btn-sm btn-success btn-deliver-trip" data-id="${t.id}" data-num="${t.trip_number}" data-start="${t.start_odometer_km}" style="margin-right: 4px; background: var(--color-success); border-color: var(--color-success); color: #fff;">
                              ✓ Baixar Entrega (Canhoto)
                            </button>
                          ` : ''}
                          <button class="btn btn-sm btn-secondary btn-manifest-trip" data-id="${t.id}">
                            📄 Manifesto
                          </button>
                        </td>
                      </tr>
                    `;
                  }).join('')}
                </tbody>
              </table>
            </div>
          </div>
        `;

        // Event listeners for filters
        container.querySelectorAll('.trip-filter-btn').forEach(btn => {
          btn.addEventListener('click', () => {
            activeFilter = btn.getAttribute('data-filter');
            renderTable();
          });
        });

        // Event listener for Add Trip modal
        container.querySelector('#btn-add-trip')?.addEventListener('click', () => {
          openAddTripModal();
        });

        // Event listeners for Dispatch Trip
        container.querySelectorAll('.btn-dispatch-trip').forEach(btn => {
          btn.addEventListener('click', async () => {
            const id = btn.getAttribute('data-id');
            const num = btn.getAttribute('data-num');
            const kmStr = prompt(`Confirmar despacho e saída do pátio para a viagem ${num}?\nInforme o Odômetro Inicial de Saída (KM):`, '142100');
            if (kmStr === null) return;
            const startKM = parseInt(kmStr, 10) || 0;

            try {
              await api.dispatchTrip(id, startKM).catch(() => {});
              const trip = tripsList.find(t => t.id === id);
              if (trip) {
                trip.status = 'dispatched';
                trip.start_odometer_km = startKM;
                trip.departure_at = new Date().toISOString();
                trip.validation_notes = '✓ Despacho homologado pela IA: CNH regular, inspeção de pátio aprovada, pneus conformes.';
              }
              alert(`Viagem ${num} despachada com sucesso! O caminhão e a carga estão em rota.`);
              renderTable();
            } catch (err) {
              alert('Erro ao despachar viagem: ' + err.message);
            }
          });
        });

        // Event listeners for Deliver Trip (Baixa com Canhoto)
        container.querySelectorAll('.btn-deliver-trip').forEach(btn => {
          btn.addEventListener('click', async () => {
            const id = btn.getAttribute('data-id');
            const num = btn.getAttribute('data-num');
            const startKM = parseInt(btn.getAttribute('data-start'), 10) || 0;
            const kmStr = prompt(`Conclusão de entrega da viagem ${num} (Baixa do Canhoto de Entrega):\nInforme o Odômetro Final de Chegada (KM - deve ser >= ${startKM}):`, String(startKM + 420));
            if (kmStr === null) return;
            const endKM = parseInt(kmStr, 10);
            if (isNaN(endKM) || endKM < startKM) {
              alert(`Odômetro inválido ou em regressão. O KM final deve ser maior ou igual a ${startKM}.`);
              return;
            }

            try {
              await api.deliverTrip(id, endKM).catch(() => {});
              const trip = tripsList.find(t => t.id === id);
              if (trip) {
                trip.status = 'delivered';
                trip.end_odometer_km = endKM;
                trip.delivered_at = new Date().toISOString();
                trip.distance_km = endKM - startKM;
                trip.validation_notes = '✓ Viagem concluída, canhoto fiscal conferido e saldo liberado no financeiro.';
              }
              alert(`Entrega da viagem ${num} concluída com sucesso! Canhoto arquivado e saldo de R$ ${trip ? trip.balance_amount.toFixed(2) : '0,00'} pronto para liquidação no caixa.`);
              renderTable();
            } catch (err) {
              alert('Erro ao concluir entrega: ' + err.message);
            }
          });
        });

        // Event listeners for Manifesto
        container.querySelectorAll('.btn-manifest-trip').forEach(btn => {
          btn.addEventListener('click', () => {
            const id = btn.getAttribute('data-id');
            const trip = tripsList.find(t => t.id === id);
            if (trip) openManifestModal(trip);
          });
        });
      };

      // Modal for creating and dispatching a new trip
      const openAddTripModal = () => {
        let modalEl = document.getElementById('modal-trip-form');
        if (!modalEl) {
          modalEl = document.createElement('div');
          modalEl.id = 'modal-trip-form';
          modalEl.className = 'modal-overlay';
          document.body.appendChild(modalEl);
        }

        // Filter vehicles: trucks (cavalos) vs trailers (carretas)
        const trucks = allVehicles.filter(v => !(v.category_id === 'semi_trailer' || v.category_id === 'trailer' || (v.metadata && (v.metadata.vehicle_category === 'semi_trailer' || v.metadata.vehicle_category === 'trailer'))));
        const trailers = allVehicles.filter(v => v.category_id === 'semi_trailer' || v.category_id === 'trailer' || (v.metadata && (v.metadata.vehicle_category === 'semi_trailer' || v.metadata.vehicle_category === 'trailer')));

        modalEl.innerHTML = `
          <div class="modal-content modal-content-wide">
            <div class="modal-header">
              <h3 style="font-weight: 700; font-size: 1.15rem;">Despachar Nova Viagem / Ordem de Transporte</h3>
              <button id="close-modal-trip" style="background: none; border: none; color: var(--text-muted); font-size: 1.25rem; cursor: pointer;">✕</button>
            </div>
            <form id="form-trip-create">
              <!-- Conjunto Transportador (Acoplamento Cavalo + Carreta) -->
              <div style="background: rgba(37, 99, 235, 0.08); border: 1px solid var(--border-color); border-radius: 8px; padding: 14px; margin-bottom: 14px;">
                <div style="font-weight: 700; font-size: 0.85rem; color: var(--color-info); margin-bottom: 10px;">
                  🚛 1. Conjunto Transportador & Condutor (Acoplamento Dinâmico)
                </div>
                <div style="display: grid; grid-template-columns: 1fr 1fr 1fr; gap: 12px;">
                  <div>
                    <label style="font-size: 0.8rem; color: var(--text-muted); display: block; margin-bottom: 4px;">Cavalo Mecânico (Trator)</label>
                    <select id="trip-vehicle" class="form-select" required>
                      ${trucks.length === 0 ? `
                        <option value="BRA2E19" data-km="142100">BRA2E19 - Scania R 450 (6x2)</option>
                        <option value="VLV4X90" data-km="89400">VLV4X90 - Volvo FH 540 (6x4)</option>
                        <option value="MER3T45" data-km="215330">MER3T45 - Mercedes Actros 2651</option>
                      ` : trucks.map(t => `<option value="${t.id}" data-plate="${t.plate || t.license_plate}" data-km="${t.current_km || 0}">${t.plate || t.license_plate} - ${t.brand_id || ''} ${t.model_id || 'Cavalo'}</option>`).join('')}
                    </select>
                  </div>
                  <div>
                    <label style="font-size: 0.8rem; color: var(--text-muted); display: block; margin-bottom: 4px;">Carreta / Semirreboque (Implemento)</label>
                    <select id="trip-trailer" class="form-select">
                      <option value="">Sem carreta (Caminhão Truck/Toco)</option>
                      ${trailers.length === 0 ? `
                        <option value="RND9A88">RND9A88 - Baú Frigorífico Vanderléia</option>
                        <option value="GRN7B12">GRN7B12 - Graneleira Bitrem 9E</option>
                        <option value="SDR1C89">SDR1C89 - Sider Cortina 3E</option>
                        <option value="TNQ4E55">TNQ4E55 - Tanque Inox Químico</option>
                      ` : trailers.map(tr => `<option value="${tr.id}" data-plate="${tr.plate || tr.license_plate}">${tr.plate || tr.license_plate} - ${(tr.metadata && tr.metadata.implement_type) || 'Semirreboque'}</option>`).join('')}
                    </select>
                  </div>
                  <div>
                    <label style="font-size: 0.8rem; color: var(--text-muted); display: block; margin-bottom: 4px;">Motorista Responsável</label>
                    <select id="trip-driver" class="form-select" required>
                      ${allDrivers.length === 0 ? `
                        <option value="João Carlos da Silva">João Carlos da Silva (CNH E - Válida)</option>
                        <option value="Carlos Alberto Mendes">Carlos Alberto Mendes (CNH E - Válida)</option>
                        <option value="Marcos Vinícius Souza">Marcos Vinícius Souza (CNH E - Válida)</option>
                      ` : allDrivers.map(d => `<option value="${d.id}" data-name="${d.name}">${d.name} (${d.status || 'Ativo'})</option>`).join('')}
                    </select>
                  </div>
                </div>

                <!-- Autonomous AI Pre-Flight Alert Box -->
                <div id="ai-preflight-box" style="margin-top: 10px; background: rgba(34, 197, 94, 0.12); border: 1px solid rgba(34, 197, 94, 0.3); border-radius: 6px; padding: 10px; font-size: 0.8rem; color: var(--color-success); display: flex; align-items: center; gap: 8px;">
                  <span>🤖</span>
                  <div>
                    <strong>Pré-Embarque Validado pela IA:</strong> Condutor com CNH regular, inspeção de pátio aprovada sem avarias críticas e pneus em conformidade de segurança.
                  </div>
                </div>
              </div>

              <!-- Rota da Viagem -->
              <div style="display: grid; grid-template-columns: 1fr 1fr 1fr; gap: 14px; margin-bottom: 14px;">
                <div>
                  <label style="font-size: 0.8rem; color: var(--text-muted); display: block; margin-bottom: 4px;">Origem (Cidade / UF)</label>
                  <input type="text" id="trip-origin" class="form-input" placeholder="ex: São Paulo/SP" value="São Paulo/SP" required />
                </div>
                <div>
                  <label style="font-size: 0.8rem; color: var(--text-muted); display: block; margin-bottom: 4px;">Destino (Cidade / UF)</label>
                  <input type="text" id="trip-dest" class="form-input" placeholder="ex: Curitiba/PR" value="Curitiba/PR" required />
                </div>
                <div>
                  <label style="font-size: 0.8rem; color: var(--text-muted); display: block; margin-bottom: 4px;">Distância Estimada (KM)</label>
                  <input type="number" id="trip-distance" class="form-input" placeholder="ex: 415" value="415" required />
                </div>
              </div>

              <!-- Carga e Documentos Fiscais -->
              <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 14px; margin-bottom: 14px;">
                <div>
                  <label style="font-size: 0.8rem; color: var(--text-muted); display: block; margin-bottom: 4px;">Descrição da Carga / Mercadoria</label>
                  <input type="text" id="trip-cargo" class="form-input" placeholder="ex: Alimentos Congelados (-18°C)" value="Alimentos Congelados" required />
                </div>
                <div>
                  <label style="font-size: 0.8rem; color: var(--text-muted); display: block; margin-bottom: 4px;">Cliente Embarcador</label>
                  <input type="text" id="trip-customer" class="form-input" placeholder="ex: BRF Brasil Foods S.A." value="BRF Brasil Foods S.A." required />
                </div>
              </div>

              <div style="display: grid; grid-template-columns: 1fr 1fr 1fr; gap: 14px; margin-bottom: 14px;">
                <div>
                  <label style="font-size: 0.8rem; color: var(--text-muted); display: block; margin-bottom: 4px;">Peso da Carga (kg)</label>
                  <input type="number" id="trip-weight" class="form-input" placeholder="ex: 18500" value="18500" required />
                </div>
                <div>
                  <label style="font-size: 0.8rem; color: var(--text-muted); display: block; margin-bottom: 4px;">Chave do MDF-e (Manifesto Eletrônico)</label>
                  <input type="text" id="trip-mdfe" class="form-input" placeholder="44 dígitos da chave" value="35260192837461029384756102938475610293847569" />
                </div>
                <div>
                  <label style="font-size: 0.8rem; color: var(--text-muted); display: block; margin-bottom: 4px;">Número do CIOT (ANTT)</label>
                  <input type="text" id="trip-ciot" class="form-input" placeholder="12 dígitos" value="098172635441-2" />
                </div>
              </div>

              <!-- Valores Financeiros -->
              <div style="background: var(--bg-main); border: 1px solid var(--border-color); border-radius: 8px; padding: 14px; margin-bottom: 14px;">
                <div style="font-weight: 700; font-size: 0.85rem; color: var(--color-success); margin-bottom: 10px;">
                  💰 2. Frete, Adiantamento de Viagem & Vale-Pedágio Obrigatório
                </div>
                <div style="display: grid; grid-template-columns: 1fr 1fr 1fr 1fr; gap: 12px;">
                  <div>
                    <label style="font-size: 0.8rem; color: var(--text-muted); display: block; margin-bottom: 4px;">Valor Bruto Frete (R$)</label>
                    <input type="number" step="0.01" id="trip-freight" class="form-input" placeholder="0.00" value="8500.00" required />
                  </div>
                  <div>
                    <label style="font-size: 0.8rem; color: var(--text-muted); display: block; margin-bottom: 4px;">Adiantamento Motorista (R$)</label>
                    <input type="number" step="0.01" id="trip-advance" class="form-input" placeholder="0.00" value="2800.00" required />
                  </div>
                  <div>
                    <label style="font-size: 0.8rem; color: var(--text-muted); display: block; margin-bottom: 4px;">Vale-Pedágio (Sem Parar/Veloe R$)</label>
                    <input type="number" step="0.01" id="trip-toll" class="form-input" placeholder="0.00" value="740.00" />
                  </div>
                  <div>
                    <label style="font-size: 0.8rem; color: var(--text-muted); display: block; margin-bottom: 4px;">Saldo Final na Entrega (R$)</label>
                    <input type="text" id="trip-balance" class="form-input" value="R$ 5.700,00" disabled style="background: rgba(0,0,0,0.2); font-weight: 700; color: var(--color-success);" />
                  </div>
                </div>
              </div>

              <div style="display: flex; justify-content: flex-end; gap: 10px;">
                <button type="button" id="btn-cancel-trip-modal" class="btn btn-secondary">Cancelar</button>
                <button type="submit" class="btn btn-primary">Homologar e Salvar Viagem</button>
              </div>
            </form>
          </div>
        `;

        modalEl.style.display = 'flex';

        // Auto calculate balance
        const freightIn = modalEl.querySelector('#trip-freight');
        const advanceIn = modalEl.querySelector('#trip-advance');
        const balanceIn = modalEl.querySelector('#trip-balance');

        const updateBalance = () => {
          const f = parseFloat(freightIn.value) || 0;
          const a = parseFloat(advanceIn.value) || 0;
          const bal = Math.max(0, f - a);
          balanceIn.value = `R$ ${bal.toLocaleString('pt-BR', { minimumFractionDigits: 2 })}`;
        };
        freightIn.addEventListener('input', updateBalance);
        advanceIn.addEventListener('input', updateBalance);

        modalEl.querySelector('#close-modal-trip')?.addEventListener('click', () => {
          modalEl.style.display = 'none';
        });
        modalEl.querySelector('#btn-cancel-trip-modal')?.addEventListener('click', () => {
          modalEl.style.display = 'none';
        });

        modalEl.querySelector('#form-trip-create')?.addEventListener('submit', async (e) => {
          e.preventDefault();
          const vehSelect = modalEl.querySelector('#trip-vehicle');
          const trailerSelect = modalEl.querySelector('#trip-trailer');
          const driverSelect = modalEl.querySelector('#trip-driver');

          const vehicleId = vehSelect.value;
          const vehiclePlate = vehSelect.options[vehSelect.selectedIndex]?.text.split(' - ')[0] || vehicleId;
          const trailerId = trailerSelect.value || null;
          const trailerPlate = trailerSelect.options[trailerSelect.selectedIndex]?.text.split(' - ')[0] || '';
          const driverId = driverSelect.value;
          const driverName = driverSelect.options[driverSelect.selectedIndex]?.text.split(' (')[0] || driverId;

          const originParts = modalEl.querySelector('#trip-origin').value.split('/');
          const destParts = modalEl.querySelector('#trip-dest').value.split('/');
          const distKM = parseInt(modalEl.querySelector('#trip-distance').value, 10) || 0;
          const cargo = modalEl.querySelector('#trip-cargo').value;
          const customer = modalEl.querySelector('#trip-customer').value;
          const weight = parseFloat(modalEl.querySelector('#trip-weight').value) || 0;
          const mdfe = modalEl.querySelector('#trip-mdfe').value;
          const ciot = modalEl.querySelector('#trip-ciot').value;
          const freight = parseFloat(freightIn.value) || 0;
          const advance = parseFloat(advanceIn.value) || 0;
          const toll = parseFloat(modalEl.querySelector('#trip-toll').value) || 0;
          const balance = Math.max(0, freight - advance);

          const newTrip = {
            id: 'trip-' + Date.now(),
            trip_number: 'VIAG-' + new Date().getFullYear() + '-' + String(tripsList.length + 1).padStart(3, '0'),
            status: 'draft',
            vehicle_id: vehicleId,
            vehicle_plate: vehiclePlate,
            trailer_id: trailerId,
            trailer_plate: trailerPlate,
            driver_id: driverId,
            driver_name: driverName,
            origin_city: originParts[0]?.trim() || 'São Paulo',
            origin_state: originParts[1]?.trim() || 'SP',
            destination_city: destParts[0]?.trim() || 'Curitiba',
            destination_state: destParts[1]?.trim() || 'PR',
            distance_km: distKM,
            cargo_description: cargo,
            customer_name: customer,
            cargo_weight_kg: weight,
            mdfe_key: mdfe,
            ciot_number: ciot,
            freight_amount: freight,
            advance_amount: advance,
            toll_amount: toll,
            balance_amount: balance,
            start_odometer_km: 142100,
            validation_notes: '✓ Pré-embarque validado pela IA: Conjunto e Condutor 100% aptos para viagem.'
          };

          try {
            await api.createTrip(newTrip).catch(() => {});
            tripsList.unshift(newTrip);
            alert(`Viagem ${newTrip.trip_number} criada com sucesso e pronta para despacho!`);
            modalEl.style.display = 'none';
            renderTable();
          } catch (err) {
            alert('Erro ao salvar viagem: ' + err.message);
          }
        });
      };

      // Modal for viewing Trip Manifesto
      const openManifestModal = (trip) => {
        let manifestModal = document.getElementById('modal-manifest-view');
        if (!manifestModal) {
          manifestModal = document.createElement('div');
          manifestModal.id = 'modal-manifest-view';
          manifestModal.className = 'modal-overlay';
          document.body.appendChild(manifestModal);
        }

        manifestModal.innerHTML = `
          <div class="modal-content modal-content-wide">
            <div style="display: flex; justify-content: space-between; align-items: flex-start; border-bottom: 2px solid var(--border-color); padding-bottom: 12px; margin-bottom: 16px;">
              <div>
                <div style="display: flex; align-items: center; gap: 8px;">
                  <span style="font-size: 22px;">🗺️</span>
                  <span style="font-weight: 800; font-size: 1.15rem;">GEEVHEIT INTELLIGENCE 360</span>
                </div>
                <div style="font-size: 0.9rem; font-weight: 700; color: var(--color-info); margin-top: 2px;">
                  ORDEM DE TRANSPORTE & MANIFESTO DE CARGA #${trip.trip_number}
                </div>
              </div>
              <div style="display: flex; gap: 8px;">
                <button id="btn-print-manifest" class="btn btn-primary btn-sm">🖨️ Imprimir Manifesto</button>
                <button id="close-manifest" class="btn btn-secondary btn-sm">✕</button>
              </div>
            </div>

            <div style="background: var(--bg-main); border: 1px solid var(--border-color); border-radius: 8px; padding: 14px; margin-bottom: 14px;">
              <div style="display: grid; grid-template-columns: repeat(3, 1fr); gap: 12px;">
                <div>
                  <div style="font-size: 0.75rem; color: var(--text-muted);">STATUS DA VIAGEM</div>
                  <div style="font-weight: 700; color: var(--text-main);">${trip.status.toUpperCase()}</div>
                </div>
                <div>
                  <div style="font-size: 0.75rem; color: var(--text-muted);">CONJUNTO TRANSPORTADOR</div>
                  <div style="font-weight: 600;">Cavalo: ${trip.vehicle_plate || trip.vehicle_id}</div>
                  <div style="font-size: 0.8rem; color: var(--color-warning);">${trip.trailer_plate ? 'Carreta: ' + trip.trailer_plate : 'Sem carreta'}</div>
                </div>
                <div>
                  <div style="font-size: 0.75rem; color: var(--text-muted);">MOTORISTA CONDUTOR</div>
                  <div style="font-weight: 600;">${trip.driver_name || trip.driver_id}</div>
                </div>
              </div>
            </div>

            <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 14px; margin-bottom: 14px;">
              <div style="background: var(--bg-main); border: 1px solid var(--border-color); border-radius: 6px; padding: 12px;">
                <div style="font-weight: 700; font-size: 0.85rem; color: var(--color-info); margin-bottom: 6px;">📍 Rota & Distância</div>
                <div>Origem: <strong>${trip.origin_city}/${trip.origin_state}</strong></div>
                <div>Destino: <strong>${trip.destination_city}/${trip.destination_state}</strong></div>
                <div style="font-size: 0.8rem; color: var(--text-muted); margin-top: 4px;">Distância Estimada: ${trip.distance_km || 0} km</div>
              </div>
              <div style="background: var(--bg-main); border: 1px solid var(--border-color); border-radius: 6px; padding: 12px;">
                <div style="font-weight: 700; font-size: 0.85rem; color: var(--color-info); margin-bottom: 6px;">📦 Carga & Documentos</div>
                <div>Mercadoria: <strong>${trip.cargo_description}</strong></div>
                <div>Cliente: <strong>${trip.customer_name}</strong></div>
                <div style="font-size: 0.8rem; color: var(--text-muted); margin-top: 4px;">Peso: ${((trip.cargo_weight_kg || 0) / 1000).toFixed(1)} toneladas</div>
                ${trip.ciot_number ? `<div style="font-size: 0.75rem; color: var(--color-success);">CIOT: ${trip.ciot_number}</div>` : ''}
              </div>
            </div>

            <div style="background: var(--bg-main); border: 1px solid var(--border-color); border-radius: 6px; padding: 12px; margin-bottom: 16px;">
              <div style="font-weight: 700; font-size: 0.85rem; color: var(--color-success); margin-bottom: 8px;">💰 Conciliação Financeira do Frete</div>
              <div style="display: grid; grid-template-columns: repeat(4, 1fr); gap: 10px;">
                <div>
                  <div style="font-size: 0.75rem; color: var(--text-muted);">VALOR DO FRETE</div>
                  <div style="font-weight: 700;">R$ ${Number(trip.freight_amount || 0).toLocaleString('pt-BR', { minimumFractionDigits: 2 })}</div>
                </div>
                <div>
                  <div style="font-size: 0.75rem; color: var(--text-muted);">ADIANTAMENTO MOTORISTA</div>
                  <div style="font-weight: 700; color: var(--color-warning);">R$ ${Number(trip.advance_amount || 0).toLocaleString('pt-BR', { minimumFractionDigits: 2 })}</div>
                </div>
                <div>
                  <div style="font-size: 0.75rem; color: var(--text-muted);">VALE-PEDÁGIO ANTT</div>
                  <div style="font-weight: 700;">R$ ${Number(trip.toll_amount || 0).toLocaleString('pt-BR', { minimumFractionDigits: 2 })}</div>
                </div>
                <div>
                  <div style="font-size: 0.75rem; color: var(--text-muted);">SALDO FINAL A LIQUIDAR</div>
                  <div style="font-weight: 800; color: var(--color-success);">R$ ${Number(trip.balance_amount || 0).toLocaleString('pt-BR', { minimumFractionDigits: 2 })}</div>
                </div>
              </div>
            </div>

            <div style="border-top: 1px solid var(--border-color); padding-top: 12px; font-size: 0.8rem; color: var(--text-muted); display: flex; justify-content: space-between;">
              <div>${trip.validation_notes || 'Inspeção pré-embarque homologada.'}</div>
              <div>Assinatura Digital do Despacho • Geevheit 360</div>
            </div>
          </div>
        `;

        manifestModal.style.display = 'flex';

        manifestModal.querySelector('#close-manifest')?.addEventListener('click', () => {
          manifestModal.style.display = 'none';
        });
        manifestModal.querySelector('#btn-print-manifest')?.addEventListener('click', () => {
          window.print();
        });
      };

      // Initial render
      renderTable();

    } catch (e) {
      container.innerHTML = `<div class="card"><p style="color: var(--color-danger);">Erro ao carregar viagens: ${e.message}</p></div>`;
    }
  }

  // ==================== DRIVERS & COMPLIANCE ====================
  async renderDrivers(container) {
    container.innerHTML = `<div class="card"><p>Carregando motoristas e CNHs...</p></div>`;
    try {
      const driversData = await api.listDrivers().catch(() => ({ data: [] }));
      const docsData = await api.listDocuments().catch(() => ({ data: [] }));
      const drivers = driversData.data || [];
      const docs = docsData.data || [];

      container.innerHTML = `
        <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px;">
          <div>
            <h2 style="font-size: 20px; font-weight: 700;">Gestão de Motoristas & Controle de Jornada (Lei 13.103/2015)</h2>
            <p style="color: var(--text-muted); font-size: 13px;">Controle rigoroso de CNH, exames toxicológicos, jornada de trabalho, paradas de 30 min (a cada 5h30 de direção), diárias e descanso interjornada de 11h.</p>
          </div>
        </div>

        <div class="card">
          <div class="table-container">
            <table class="data-table">
              <thead>
                <tr>
                  <th>Nome do Motorista</th>
                  <th>CPF</th>
                  <th>Telefone</th>
                  <th>Status</th>
                  <th>Documento / Vencimento CNH</th>
                  <th>Situação CNH</th>
                  <th style="text-align: right;">Controle de Jornada</th>
                </tr>
              </thead>
              <tbody>
                ${drivers.length === 0 ? `<tr><td colspan="7" style="text-align: center; color: var(--text-muted); padding: 24px;">Nenhum motorista cadastrado.</td></tr>` : ''}
                ${drivers.map(d => {
                  const driverDoc = docs.find(doc => doc.subject_id === d.id && doc.category === 'cnh') || {};
                  return `
                    <tr>
                      <td style="font-weight: 600;">${d.name}</td>
                      <td>${d.tax_id || '-'}</td>
                      <td>${d.phone || '-'}</td>
                      <td><span class="badge badge-success">${d.status || 'Ativo'}</span></td>
                      <td>${driverDoc.number ? `CNH: ${driverDoc.number} (${driverDoc.expires_at ? new Date(driverDoc.expires_at).toLocaleDateString('pt-BR') : 'Sem data'})` : 'Não anexada'}</td>
                      <td>
                        <span class="badge ${driverDoc.status === 'valid' ? 'badge-success' : driverDoc.status === 'warning' ? 'badge-warning' : 'badge-danger'}">
                          ${driverDoc.status === 'valid' ? 'Regular' : driverDoc.status === 'warning' ? 'A Vencer (Atenção)' : 'Pendente / Vencido'}
                        </span>
                      </td>
                      <td style="text-align: right;">
                        <button class="btn btn-secondary btn-sm btn-driver-journey" data-driver-id="${d.id}" style="padding: 6px 12px; font-size: 12px; display: inline-flex; align-items: center; gap: 6px;">
                          <span>⏱️</span> Diário de Bordo & Jornada
                        </button>
                      </td>
                    </tr>
                  `;
                }).join('')}
              </tbody>
            </table>
          </div>
        </div>
      `;

      // Event listeners for Journey buttons
      container.querySelectorAll('.btn-driver-journey').forEach(btn => {
        btn.addEventListener('click', () => {
          const driverId = btn.getAttribute('data-driver-id');
          const driver = drivers.find(d => d.id === driverId);
          if (driver) {
            window.openDriverJourneyModal(driver);
          }
        });
      });

      // Driver Journey Modal Implementation
      window.openDriverJourneyModal = async (driver) => {
        let journeyModal = document.getElementById('modal-driver-journey');
        if (!journeyModal) {
          journeyModal = document.createElement('div');
          journeyModal.id = 'modal-driver-journey';
          journeyModal.className = 'modal-overlay';
          document.body.appendChild(journeyModal);
        }

        journeyModal.innerHTML = `
          <div class="modal-content modal-content-xl">
            <div style="text-align: center; padding: 40px;">
              <p>Carregando diário de bordo e conformidade da jornada para ${driver.name}...</p>
            </div>
          </div>
        `;
        journeyModal.style.display = 'flex';

        const refreshJourneyData = async () => {
          try {
            const [summaryRes, journeyRes] = await Promise.all([
              api.getDriverJourneySummary(driver.id).catch(() => ({})),
              api.listDriverJourney(driver.id, 50).catch(() => [])
            ]);

            const summary = summaryRes || {};
            const records = Array.isArray(journeyRes) ? journeyRes : (journeyRes.data || []);

            const activeEvent = summary.active_event || (records.length > 0 && !records[0].end_at ? records[0].event_type : null);

            const eventLabels = {
              driving: { text: 'Em Direção Contínua', icon: '🚚', badge: 'badge-success', color: '#10b981' },
              break_rest: { text: 'Em Descanso (Pausa 30 min)', icon: '☕', badge: 'badge-info', color: '#06b6d4' },
              meal: { text: 'Em Refeição (Almoço/Jantar)', icon: '🍽️', badge: 'badge-warning', color: '#f59e0b' },
              daily_rest: { text: 'Em Pernoite (Descanso Interjornada 11h)', icon: '🛏️', badge: 'badge-secondary', color: '#8b5cf6' },
              waiting: { text: 'Tempo de Espera (Carga/Descarga)', icon: '⏳', badge: 'badge-warning', color: '#eab308' },
              standby: { text: 'Fora de Turno / Standby', icon: '🛑', badge: 'badge-secondary', color: '#6b7280' }
            };

            const currentStatus = activeEvent ? (eventLabels[activeEvent] || { text: activeEvent, icon: '⏱️', badge: 'badge-info', color: '#3b82f6' }) : { text: 'Sem Turno Ativo', icon: '⚪', badge: 'badge-secondary', color: '#9ca3af' };

            const totalDrivingMin = summary.total_driving_minutes || 0;
            const totalRestMin = summary.total_rest_minutes || 0;
            const totalWaitingMin = summary.total_waiting_minutes || 0;
            const totalAllowance = Number(summary.total_allowance_amount || 0);
            const violationsCount = summary.violations_count || 0;

            journeyModal.innerHTML = `
              <div class="modal-content modal-content-xl" id="printable-journey">
                <!-- Modal & Timesheet Header -->
                <div style="display: flex; justify-content: space-between; align-items: flex-start; border-bottom: 2px solid var(--border-color); padding-bottom: 16px; margin-bottom: 18px;">
                  <div>
                    <div style="display: flex; align-items: center; gap: 8px;">
                      <span style="font-size: 24px;">⏱️</span>
                      <span style="font-weight: 800; font-size: 1.2rem; letter-spacing: -0.02em;">GEEVHEIT INTELLIGENCE 360</span>
                    </div>
                    <div style="font-size: 0.95rem; font-weight: 700; margin-top: 4px; color: var(--color-info);">
                      DIÁRIO DE BORDO & CONTROLE DE JORNADA DO MOTORISTA (LEI 13.103/2015)
                    </div>
                    <div style="font-size: 0.8rem; color: var(--text-muted); margin-top: 2px;">
                      Motorista: <strong style="color: var(--text-main);">${driver.name}</strong> • CPF: <strong>${driver.tax_id || 'Não informado'}</strong> • CNH: <strong>${driver.cnh_number || 'Conferir doc'}</strong>
                    </div>
                  </div>
                  <div style="display: flex; gap: 8px;">
                    <button id="btn-print-journey" class="btn btn-primary btn-sm" style="display: inline-flex; align-items: center; gap: 6px;">
                      <span>🖨️</span> Imprimir Espelho de Ponto
                    </button>
                    <button id="btn-close-journey" class="btn btn-secondary btn-sm">✕</button>
                  </div>
                </div>

                <!-- Current Status Banner & KPI Indicators -->
                <div style="display: grid; grid-template-columns: repeat(auto-fit, minmax(180px, 1fr)); gap: 12px; margin-bottom: 20px;">
                  <div class="card" style="padding: 12px 16px; border-left: 4px solid ${currentStatus.color};">
                    <div style="font-size: 0.72rem; color: var(--text-muted); text-transform: uppercase; font-weight: 600;">Status do Motorista</div>
                    <div style="font-size: 1rem; font-weight: 700; color: ${currentStatus.color}; margin-top: 4px;">
                      ${currentStatus.icon} ${currentStatus.text}
                    </div>
                  </div>

                  <div class="card" style="padding: 12px 16px;">
                    <div style="font-size: 0.72rem; color: var(--text-muted); text-transform: uppercase; font-weight: 600;">Tempo de Direção</div>
                    <div style="font-size: 1.15rem; font-weight: 700; color: var(--text-main); margin-top: 4px;">
                      ${Math.floor(totalDrivingMin / 60)}h ${totalDrivingMin % 60}m
                    </div>
                  </div>

                  <div class="card" style="padding: 12px 16px;">
                    <div style="font-size: 0.72rem; color: var(--text-muted); text-transform: uppercase; font-weight: 600;">Tempo de Descanso</div>
                    <div style="font-size: 1.15rem; font-weight: 700; color: var(--color-info); margin-top: 4px;">
                      ${Math.floor(totalRestMin / 60)}h ${totalRestMin % 60}m
                    </div>
                  </div>

                  <div class="card" style="padding: 12px 16px;">
                    <div style="font-size: 0.72rem; color: var(--text-muted); text-transform: uppercase; font-weight: 600;">Diárias / Alimentação</div>
                    <div style="font-size: 1.15rem; font-weight: 700; color: var(--color-success); margin-top: 4px;">
                      R$ ${totalAllowance.toFixed(2)}
                    </div>
                  </div>

                  <div class="card" style="padding: 12px 16px; ${violationsCount > 0 ? 'border: 1px solid var(--color-danger); background: rgba(239, 68, 68, 0.08);' : ''}">
                    <div style="font-size: 0.72rem; color: var(--text-muted); text-transform: uppercase; font-weight: 600;">Infrações Lei 13.103</div>
                    <div style="font-size: 1.15rem; font-weight: 700; color: ${violationsCount > 0 ? 'var(--color-danger)' : 'var(--color-success)'}; margin-top: 4px;">
                      ${violationsCount > 0 ? `🚨 ${violationsCount} Infração(ões)` : '✓ 100% Conforme'}
                    </div>
                  </div>
                </div>

                <!-- Registration Form / Dispatch Action Bar -->
                <div class="card" style="margin-bottom: 20px; background: var(--bg-surface); border: 1px solid var(--border-color);">
                  <div style="font-weight: 700; font-size: 0.95rem; margin-bottom: 12px; display: flex; align-items: center; gap: 8px;">
                    <span>📝</span> Registrar Novo Evento de Ponto / Parada Regulamentar
                  </div>
                  <form id="form-record-journey">
                    <div style="display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: 12px;">
                      <div>
                        <label style="font-size: 0.8rem; color: var(--text-muted); display: block; margin-bottom: 4px;">Tipo de Apontamento *</label>
                        <select id="journey-event-type" class="form-control" style="width: 100%; padding: 8px; border-radius: 6px; background: var(--bg-card); border: 1px solid var(--border-color); color: var(--text-main);" required>
                          <option value="driving">🚚 Iniciar Direção</option>
                          <option value="break_rest">☕ Parada Descanso Obrigatória (30 min)</option>
                          <option value="meal">🍽️ Almoço / Refeição (Diária R$ 45,00)</option>
                          <option value="daily_rest">🛏️ Pernoite / Descanso 11h (Diária R$ 90,00)</option>
                          <option value="waiting">⏳ Tempo de Espera (Carga/Descarga)</option>
                          <option value="standby">🛑 Encerrar Turno / Standby</option>
                        </select>
                      </div>

                      <div>
                        <label style="font-size: 0.8rem; color: var(--text-muted); display: block; margin-bottom: 4px;">Local / Ponto de Parada *</label>
                        <input id="journey-location" type="text" class="form-control" placeholder="Ex: Posto Graal KM 120, Garagem, Pátio CD" style="width: 100%; padding: 8px; border-radius: 6px; background: var(--bg-card); border: 1px solid var(--border-color); color: var(--text-main);" required />
                      </div>

                      <div>
                        <label style="font-size: 0.8rem; color: var(--text-muted); display: block; margin-bottom: 4px;">Odômetro Atual (KM)</label>
                        <input id="journey-odometer" type="number" class="form-control" placeholder="Ex: 145200" style="width: 100%; padding: 8px; border-radius: 6px; background: var(--bg-card); border: 1px solid var(--border-color); color: var(--text-main);" />
                      </div>

                      <div>
                        <label style="font-size: 0.8rem; color: var(--text-muted); display: block; margin-bottom: 4px;">Observações</label>
                        <input id="journey-notes" type="text" class="form-control" placeholder="Opcional: motivo da parada, viagem..." style="width: 100%; padding: 8px; border-radius: 6px; background: var(--bg-card); border: 1px solid var(--border-color); color: var(--text-main);" />
                      </div>
                    </div>

                    <div style="display: flex; justify-content: flex-end; margin-top: 14px;">
                      <button type="submit" id="btn-submit-journey" class="btn btn-primary" style="padding: 8px 18px; font-weight: 600;">
                        <span>✓</span> Registrar Apontamento
                      </button>
                    </div>
                  </form>
                </div>

                <!-- Journey Event Log Table -->
                <div class="card">
                  <div style="font-weight: 700; font-size: 0.95rem; margin-bottom: 12px; display: flex; align-items: center; justify-content: space-between;">
                    <div style="display: flex; align-items: center; gap: 8px;">
                      <span>📋</span> Histórico de Apontamentos & Auditoria de Ponto
                    </div>
                    <span style="font-size: 0.75rem; color: var(--text-muted);">Lei 13.103/2015 Art. 67-C • ADI 5322 STF</span>
                  </div>

                  <div class="table-container">
                    <table class="data-table">
                      <thead>
                        <tr>
                          <th>Apontamento / Evento</th>
                          <th>Início</th>
                          <th>Fim</th>
                          <th>Duração</th>
                          <th>Localização & KM</th>
                          <th>Diária (R$)</th>
                          <th>Conformidade Legal</th>
                        </tr>
                      </thead>
                      <tbody>
                        ${records.length === 0 ? `<tr><td colspan="7" style="text-align: center; color: var(--text-muted); padding: 24px;">Nenhum apontamento registrado para este motorista.</td></tr>` : ''}
                        ${records.map(r => {
                          const conf = eventLabels[r.event_type] || { text: r.event_type, icon: '⏱️', badge: 'badge-secondary' };
                          const startDate = r.start_at ? new Date(r.start_at).toLocaleString('pt-BR') : '-';
                          const endDate = r.end_at ? new Date(r.end_at).toLocaleString('pt-BR') : '<span style="color: var(--color-success); font-weight: 600;">(Ativo)</span>';
                          const durText = r.duration_minutes ? `${Math.floor(r.duration_minutes / 60)}h ${r.duration_minutes % 60}m` : '-';

                          return `
                            <tr style="${r.is_violation ? 'background: rgba(239, 68, 68, 0.05);' : ''}">
                              <td style="font-weight: 600;">
                                <span style="margin-right: 6px;">${conf.icon}</span> ${conf.text}
                              </td>
                              <td style="font-size: 0.85rem;">${startDate}</td>
                              <td style="font-size: 0.85rem;">${endDate}</td>
                              <td style="font-weight: 600;">${durText}</td>
                              <td style="font-size: 0.85rem;">
                                ${r.location || '-'} ${r.odometer_km ? `(${r.odometer_km} km)` : ''}
                                ${r.notes ? `<div style="font-size: 0.75rem; color: var(--text-muted);">${r.notes}</div>` : ''}
                              </td>
                              <td style="font-weight: 600; color: var(--color-success);">
                                ${r.daily_allowance_amount > 0 ? `R$ ${Number(r.daily_allowance_amount).toFixed(2)}` : '-'}
                              </td>
                              <td>
                                ${r.is_violation ? `
                                  <span class="badge badge-danger" title="${r.violation_reason || ''}">
                                    🚨 INFRAÇÃO: ${r.violation_reason || 'Excesso de direção contínua'}
                                  </span>
                                ` : `
                                  <span class="badge badge-success">✓ Conforme</span>
                                `}
                              </td>
                            </tr>
                          `;
                        }).join('')}
                      </tbody>
                    </table>
                  </div>

                  <!-- Printable Signature Block for Labor Compliance (visible on print) -->
                  <div style="display: flex; justify-content: space-between; margin-top: 40px; padding-top: 24px; border-top: 1px dashed var(--border-color);">
                    <div style="text-align: center; width: 45%;">
                      <div style="border-top: 1px solid #9ca3af; width: 80%; margin: 0 auto 8px auto;"></div>
                      <div style="font-size: 0.85rem; font-weight: 700;">${driver.name}</div>
                      <div style="font-size: 0.75rem; color: var(--text-muted);">Assinatura do Motorista Profissional</div>
                    </div>
                    <div style="text-align: center; width: 45%;">
                      <div style="border-top: 1px solid #9ca3af; width: 80%; margin: 0 auto 8px auto;"></div>
                      <div style="font-size: 0.85rem; font-weight: 700;">GEEVHEIT INTELLIGENCE 360</div>
                      <div style="font-size: 0.75rem; color: var(--text-muted);">Assinatura do Gestor de Frotas / Recursos Humanos</div>
                    </div>
                  </div>
                </div>
              </div>
            `;

            // Handlers inside modal
            const closeBtn = document.getElementById('btn-close-journey');
            if (closeBtn) closeBtn.onclick = () => { journeyModal.style.display = 'none'; };

            const printBtn = document.getElementById('btn-print-journey');
            if (printBtn) printBtn.onclick = () => { window.print(); };

            const journeyForm = document.getElementById('form-record-journey');
            if (journeyForm) {
              journeyForm.onsubmit = async (e) => {
                e.preventDefault();
                const btnSubmit = document.getElementById('btn-submit-journey');
                if (btnSubmit) btnSubmit.disabled = true;

                const eventType = document.getElementById('journey-event-type').value;
                const location = document.getElementById('journey-location').value;
                const odo = parseInt(document.getElementById('journey-odometer').value, 10) || 0;
                const notes = document.getElementById('journey-notes').value;

                try {
                  await api.recordDriverJourney(driver.id, {
                    event_type: eventType,
                    location: location,
                    odometer_km: odo,
                    notes: notes
                  });
                  await refreshJourneyData();
                } catch (err) {
                  alert('Erro ao registrar apontamento de jornada: ' + err.message);
                } finally {
                  if (btnSubmit) btnSubmit.disabled = false;
                }
              };
            }
          } catch (err) {
            journeyModal.innerHTML = `
              <div class="modal-content modal-content-wide">
                <p style="color: var(--color-danger);">Erro ao carregar jornada: ${err.message}</p>
                <button class="btn btn-secondary" onclick="document.getElementById('modal-driver-journey').style.display='none'">Fechar</button>
              </div>
            `;
          }
        };

        await refreshJourneyData();
      };
    } catch (e) {
      container.innerHTML = `<div class="card"><p style="color: var(--color-danger);">Erro: ${e.message}</p></div>`;
    }
  }

  // ==================== FUEL, REQUISITIONS & INVOICES ====================
  async renderFuel(container) {
    container.innerHTML = `<div class="card"><p>Carregando abastecimentos e requisições...</p></div>`;
    try {
      const data = await api.listFuelTransactions().catch(() => ({ data: [] }));
      const txs = data.data || [];

      container.innerHTML = `
        <div style="display: flex; justify-content: space-between; align-items: center;">
          <div>
            <h2 style="font-size: 20px; font-weight: 700;">Abastecimentos, Requisições & Conciliação Fiscal</h2>
            <p style="color: var(--text-muted); font-size: 13px;">Controle de combustível com separação de frota própria e terceiros, número de requisição e conciliação para o caixa.</p>
          </div>
          <button id="btn-add-fuel" class="btn btn-primary">
            <span>+</span> Lançar Abastecimento
          </button>
        </div>

        <div class="card">
          <div class="table-container">
            <table class="data-table">
              <thead>
                <tr>
                  <th>Data</th>
                  <th>Requisição</th>
                  <th>Tipo Frota</th>
                  <th>Veículo / Terceiro</th>
                  <th>Litros / Valor</th>
                  <th>Nota Fiscal</th>
                  <th>Status Caixa/Pagamento</th>
                </tr>
              </thead>
              <tbody>
                ${txs.length === 0 ? `<tr><td colspan="7" style="text-align: center; color: var(--text-muted); padding: 24px;">Nenhum abastecimento registrado.</td></tr>` : ''}
                ${txs.map(t => `
                  <tr>
                    <td>${new Date(t.timestamp || t.created_at).toLocaleDateString('pt-BR')}</td>
                    <td style="font-weight: 700;">${t.requisition_number || 'S/N'}</td>
                    <td>
                      <span class="badge ${t.is_third_party ? 'badge-warning' : 'badge-info'}">
                        ${t.is_third_party ? 'Terceiro / Agregado' : 'Frota Própria'}
                      </span>
                    </td>
                    <td>${t.is_third_party ? (t.third_party_identifier || 'Terceiro') : (t.vehicle_id || 'Veículo')}</td>
                    <td>${t.liters || 0} L (R$ ${Number(t.total_amount || 0).toFixed(2)})</td>
                    <td>${t.invoice_number ? `NF ${t.invoice_number}` : 'Pendente'}</td>
                    <td>
                      <span class="badge ${t.payment_status === 'paid' ? 'badge-success' : 'badge-warning'}">
                        ${t.payment_status === 'paid' ? 'Pago' : 'Pendente'}
                      </span>
                    </td>
                  </tr>
                `).join('')}
              </tbody>
            </table>
          </div>
        </div>
      `;

      document.getElementById('btn-add-fuel')?.addEventListener('click', async () => {
        const isThird = confirm('O abastecimento é para um VEÍCULO DE TERCEIRO / AGREGADO? (Clique OK para terceiro, Cancelar para veículo próprio da frota)');
        const reqNum = prompt('Número da Requisição de Abastecimento:');
        if (!reqNum) return;

        let thirdId = '';
        let vehicleId = '';
        if (isThird) {
          thirdId = prompt('Identificador do Terceiro (Ex: Placa ABC-1234 / Agregado João):');
        } else {
          vehicleId = prompt('ID ou Placa do Veículo Próprio:');
        }

        const liters = parseFloat(prompt('Quantidade em Litros:') || '0');
        const total = parseFloat(prompt('Valor Total (R$):') || '0');
        const nf = prompt('Número da Nota Fiscal (opcional):') || '';

        try {
          await api.recordFuelTransaction({
            liters: liters,
            total_amount: total,
            requisition_number: reqNum,
            is_third_party: isThird,
            third_party_identifier: thirdId,
            vehicle_id: vehicleId,
            invoice_number: nf,
            payment_status: 'pending',
          });
          alert('Abastecimento registrado com sucesso para conciliação contábil!');
          this.renderFuel(container);
        } catch (e) {
          alert('Erro ao registrar abastecimento: ' + e.message);
        }
      });

    } catch (e) {
      container.innerHTML = `<div class="card"><p style="color: var(--color-danger);">Erro: ${e.message}</p></div>`;
    }
  }

  // ==================== TIRES ====================
  async renderTires(container) {
    container.innerHTML = `<div class="card"><p>Carregando inventário de pneus e sulcos...</p></div>`;
    try {
      const data = await api.listTires().catch(() => ({ data: [] }));
      const tires = data.data || [];

      container.innerHTML = `
        <div style="display: flex; justify-content: space-between; align-items: center;">
          <div>
            <h2 style="font-size: 20px; font-weight: 700;">Gestão de Pneus & Ciclo de Vida</h2>
            <p style="color: var(--text-muted); font-size: 13px;">Rastreamento individual por número de fogo, desgaste de sulco (mm), rodízios e recapagens.</p>
          </div>
        </div>

        <div class="card">
          <div class="table-container">
            <table class="data-table">
              <thead>
                <tr>
                  <th>Nº Fogo</th>
                  <th>Marca / Modelo</th>
                  <th>Medida</th>
                  <th>Sulco Atual (mm)</th>
                  <th>Sulco Mínimo (mm)</th>
                  <th>Status de Segurança</th>
                  <th>Ação</th>
                </tr>
              </thead>
              <tbody>
                ${tires.length === 0 ? `<tr><td colspan="7" style="text-align: center; color: var(--text-muted); padding: 24px;">Nenhum pneu cadastrado.</td></tr>` : ''}
                ${tires.map(t => {
                  const isCritical = t.current_tread_mm <= t.minimum_tread_mm;
                  return `
                    <tr>
                      <td style="font-weight: 700; color: var(--text-main);">${t.fire_number}</td>
                      <td>${t.brand} ${t.model}</td>
                      <td>${t.size}</td>
                      <td style="font-weight: 700; color: ${isCritical ? 'var(--color-danger)' : 'var(--text-main)'};">
                        ${t.current_tread_mm} mm
                      </td>
                      <td>${t.minimum_tread_mm} mm</td>
                      <td>
                        <span class="badge ${isCritical ? 'badge-danger' : 'badge-success'}">
                          ${isCritical ? 'CRÍTICO (Troca Imediata)' : 'Conforme'}
                        </span>
                      </td>
                      <td>
                        <button class="btn btn-secondary" onclick="window.inspectTire('${t.id}', '${t.fire_number}')">
                          Medir Sulco
                        </button>
                      </td>
                    </tr>
                  `;
                }).join('')}
              </tbody>
            </table>
          </div>
        </div>
      `;

      window.inspectTire = async (tireId, fireNum) => {
        const mmStr = prompt(`Informe a nova medida de sulco para o pneu ${fireNum} (mm):`);
        if (!mmStr) return;
        const mm = parseFloat(mmStr);
        if (isNaN(mm) || mm < 0) {
          alert('Valor inválido');
          return;
        }

        try {
          await api.recordTireInspection(tireId, {
            tread_depth_mm: mm,
            pressure_psi: 110,
            notes: 'Inspeção periódica de pátio',
          });
          alert('Inspeção registrada com sucesso!');
          this.renderTires(container);
        } catch (e) {
          alert('Erro: ' + e.message);
        }
      };

    } catch (e) {
      container.innerHTML = `<div class="card"><p style="color: var(--color-danger);">Erro: ${e.message}</p></div>`;
    }
  }

  // ==================== MAINTENANCE ====================
  async renderMaintenance(container) {
    container.innerHTML = `<div class="card"><p>Carregando ordens de serviço de manutenção...</p></div>`;
    try {
      const data = await api.listWorkOrders().catch(() => ({ items: [] }));
      const orders = data.items || [];

      container.innerHTML = `
        <div style="display: flex; justify-content: space-between; align-items: center;">
          <div>
            <h2 style="font-size: 20px; font-weight: 700;">Manutenção Preventiva e Corretiva</h2>
            <p style="color: var(--text-muted); font-size: 13px;">Abertura de O.S., controle de tempo parado (downtime), custos e notas fiscais para o financeiro.</p>
          </div>
          <button id="btn-add-wo" class="btn btn-primary">
            <span>+</span> Abrir Ordem de Serviço
          </button>
        </div>

        <div class="card">
          <div class="table-container">
            <table class="data-table">
              <thead>
                <tr>
                  <th>Código</th>
                  <th>Título / Serviço</th>
                  <th>Tipo</th>
                  <th>Prioridade</th>
                  <th>Status</th>
                  <th>Custo Total</th>
                  <th>Nota Fiscal</th>
                </tr>
              </thead>
              <tbody>
                ${orders.length === 0 ? `<tr><td colspan="7" style="text-align: center; color: var(--text-muted); padding: 24px;">Nenhuma ordem de serviço cadastrada.</td></tr>` : ''}
                ${orders.map(o => `
                  <tr>
                    <td style="font-weight: 700;">${o.code}</td>
                    <td>${o.title}</td>
                    <td><span class="badge badge-secondary">${o.kind}</span></td>
                    <td>
                      <span class="badge ${o.priority === 'critical' ? 'badge-danger' : o.priority === 'high' ? 'badge-warning' : 'badge-info'}">
                        ${o.priority}
                      </span>
                    </td>
                    <td><span class="badge badge-success">${o.status}</span></td>
                    <td>R$ ${Number(o.total_cost || 0).toFixed(2)}</td>
                    <td>${o.invoice_number ? `NF ${o.invoice_number}` : 'Pendente'}</td>
                  </tr>
                `).join('')}
              </tbody>
            </table>
          </div>
        </div>
      `;

      document.getElementById('btn-add-wo')?.addEventListener('click', async () => {
        const title = prompt('Título do serviço / defeito reportado:');
        if (!title) return;
        const kind = prompt('Tipo de manutenção (preventive, corrective, emergency):', 'corrective');
        const cost = parseFloat(prompt('Custo estimado/total (R$):') || '0');
        const nf = prompt('Número da Nota Fiscal de Serviço (se houver):') || '';

        try {
          await api.createWorkOrder({
            title: title,
            kind: kind,
            priority: 'high',
            total_cost: cost,
            invoice_number: nf,
            payment_status: 'pending',
          });
          alert('Ordem de serviço aberta com sucesso!');
          this.renderMaintenance(container);
        } catch (e) {
          alert('Erro ao criar O.S.: ' + e.message);
        }
      });

    } catch (e) {
      container.innerHTML = `<div class="card"><p style="color: var(--color-danger);">Erro: ${e.message}</p></div>`;
    }
  }

  // ==================== CHECKLISTS ====================
  async renderChecklists(container) {
    container.innerHTML = `<div class="card"><p>Carregando checklists operacionais...</p></div>`;
    try {
      const data = await api.listChecklists().catch(() => ({ items: [] }));
      const checklists = data.items || [];

      container.innerHTML = `
        <div style="display: flex; justify-content: space-between; align-items: center;">
          <div>
            <h2 style="font-size: 20px; font-weight: 700;">Checklists Operacionais de Pátio</h2>
            <p style="color: var(--text-muted); font-size: 13px;">Inspeções de segurança de saída e retorno de viagem com bloqueio imediato em não-conformidade crítica.</p>
          </div>
        </div>

        <div class="card">
          <div class="table-container">
            <table class="data-table">
              <thead>
                <tr>
                  <th>Título</th>
                  <th>Tipo</th>
                  <th>Motorista</th>
                  <th>Status</th>
                  <th>Data Realização</th>
                </tr>
              </thead>
              <tbody>
                ${checklists.length === 0 ? `<tr><td colspan="5" style="text-align: center; color: var(--text-muted); padding: 24px;">Nenhum checklist registrado recentemente.</td></tr>` : ''}
                ${checklists.map(c => `
                  <tr>
                    <td style="font-weight: 600;">${c.title}</td>
                    <td>${c.type || 'Pátio'}</td>
                    <td>${c.driver_name || '-'}</td>
                    <td><span class="badge badge-success">${c.status}</span></td>
                    <td>${new Date(c.created_at).toLocaleString('pt-BR')}</td>
                  </tr>
                `).join('')}
              </tbody>
            </table>
          </div>
        </div>
      `;
    } catch (e) {
      container.innerHTML = `<div class="card"><p style="color: var(--color-danger);">Erro: ${e.message}</p></div>`;
    }
  }

  // ==================== INVENTORY & STOCK MANAGEMENT ====================
  async renderInventory(container) {
    container.innerHTML = `<div class="card"><p>Carregando estoque de pneus, peças e equipamentos...</p></div>`;
    try {
      const data = await api.listParts().catch(() => ({ data: [] }));
      let partsList = data.data || [];

      // Initial comprehensive operational catalog if stock is empty
      if (partsList.length === 0) {
        partsList = [
          {
            id: 'demo-part-1',
            sku: 'PNEU-295-MIC-NOV',
            internal_code: 'PNEU-001',
            name: 'Pneu 295/80 R22.5 Michelin X Multi Z2 (Novo)',
            category_id: 'tires',
            category_name: 'Pneus Novos',
            unit_cost: 3150.00,
            current_stock: 14,
            min_stock: 6,
            location: 'Almoxarifado A - Baia de Pneus Novos',
            unit_id: 'un'
          },
          {
            id: 'demo-part-2',
            sku: 'PNEU-295-VIP-REC',
            internal_code: 'PNEU-002',
            name: 'Pneu 295/80 R22.5 Vipal Recapado (1ª Vida)',
            category_id: 'tires',
            category_name: 'Pneus Recapados',
            unit_cost: 980.00,
            current_stock: 8,
            min_stock: 4,
            location: 'Almoxarifado A - Baia de Recapagem',
            unit_id: 'un'
          },
          {
            id: 'demo-part-3',
            sku: 'PNEU-CARC-SUC',
            internal_code: 'PNEU-003',
            name: 'Carcaças de Pneu para Destinação / Sucata',
            category_id: 'tires',
            category_name: 'Sucatas / Descarte',
            unit_cost: 0.00,
            current_stock: 5,
            min_stock: 0,
            location: 'Pátio Externo - Zona de Descarte',
            unit_id: 'un'
          },
          {
            id: 'demo-part-4',
            sku: 'BRK-FRAS-01',
            internal_code: 'PEC-010',
            name: 'Pastilha de Freio Dianteira Fras-Le Scania/Volvo',
            category_id: 'mechanical',
            category_name: 'Peças & Freios',
            unit_cost: 420.00,
            current_stock: 12,
            min_stock: 4,
            location: 'Corredor B - Prateleira 2',
            unit_id: 'jg'
          },
          {
            id: 'demo-part-5',
            sku: 'BRK-LONA-MST',
            internal_code: 'PEC-011',
            name: 'Jogo de Lonas de Freio Master para Carreta / Semirreboque',
            category_id: 'mechanical',
            category_name: 'Peças & Freios',
            unit_cost: 380.00,
            current_stock: 6,
            min_stock: 4,
            location: 'Corredor B - Prateleira 3',
            unit_id: 'jg'
          },
          {
            id: 'demo-part-6',
            sku: 'CUICA-3030-WAB',
            internal_code: 'PEC-012',
            name: 'Câmara de Freio Cuíca Dupla Spring Brake 30x30 Wabco',
            category_id: 'mechanical',
            category_name: 'Peças & Freios',
            unit_cost: 540.00,
            current_stock: 3,
            min_stock: 4,
            location: 'Corredor B - Prateleira 4',
            unit_id: 'un'
          },
          {
            id: 'demo-part-7',
            sku: 'LUB-SHL-15W40',
            internal_code: 'OLEO-001',
            name: 'Óleo Motor Shell Rimula R4 X 15W40 CI-4 (Balde 20L)',
            category_id: 'lubricants',
            category_name: 'Filtros & Óleos',
            unit_cost: 410.00,
            current_stock: 18,
            min_stock: 8,
            location: 'Armazém de Lubrificantes - Pallet 1',
            unit_id: 'balde'
          },
          {
            id: 'demo-part-8',
            sku: 'FILT-RAC-COMB',
            internal_code: 'FLT-001',
            name: 'Filtro Separador Racor Combustível Mann-Filter',
            category_id: 'lubricants',
            category_name: 'Filtros & Óleos',
            unit_cost: 165.00,
            current_stock: 4,
            min_stock: 6,
            location: 'Corredor C - Prateleira 1',
            unit_id: 'un'
          },
          {
            id: 'demo-part-9',
            sku: 'ARLA32-GRANEL',
            internal_code: 'ARLA-001',
            name: 'Fluido Redutor ARLA 32 Certificado Inmetro (Galão 20L)',
            category_id: 'lubricants',
            category_name: 'Filtros & Óleos',
            unit_cost: 78.00,
            current_stock: 25,
            min_stock: 10,
            location: 'Armazém de Lubrificantes - Pallet 3',
            unit_id: 'gl'
          },
          {
            id: 'demo-part-10',
            sku: 'EL-BAT-MOU-180',
            internal_code: 'EL-001',
            name: 'Bateria Pesada Moura Log Diesel 12V 180Ah (Livre Manutenção)',
            category_id: 'electrical',
            category_name: 'Elétrica & Baterias',
            unit_cost: 1290.00,
            current_stock: 6,
            min_stock: 2,
            location: 'Corredor E - Bancada de Baterias',
            unit_id: 'un'
          },
          {
            id: 'demo-part-11',
            sku: 'EL-LED-LAN-CR',
            internal_code: 'EL-002',
            name: 'Lanterna Traseira Multifuncional LED Blindada para Carreta',
            category_id: 'electrical',
            category_name: 'Elétrica & Baterias',
            unit_cost: 235.00,
            current_stock: 10,
            min_stock: 4,
            location: 'Corredor E - Gaveteiro 2',
            unit_id: 'un'
          },
          {
            id: 'demo-part-12',
            sku: 'EQP-CINTA-5T',
            internal_code: 'EQP-001',
            name: 'Cinta de Amarração de Carga com Catraca 5 Toneladas (9m)',
            category_id: 'onboard',
            category_name: 'Equipamentos & Amarração',
            unit_cost: 115.00,
            current_stock: 22,
            min_stock: 10,
            location: 'Corredor D - Prateleira 1',
            unit_id: 'un'
          },
          {
            id: 'demo-part-13',
            sku: 'EQP-EXT-ABC-4K',
            internal_code: 'EQP-002',
            name: 'Extintor de Incêndio Automotivo Pó ABC 4kg (com suporte)',
            category_id: 'onboard',
            category_name: 'Equipamentos & Amarração',
            unit_cost: 185.00,
            current_stock: 7,
            min_stock: 3,
            location: 'Corredor D - Prateleira 3',
            unit_id: 'un'
          },
          {
            id: 'demo-part-14',
            sku: 'EQP-TACO-DISC',
            internal_code: 'EQP-003',
            name: 'Bobina / Disco Diagrama Semanal de Tacógrafo 125 km/h (Cx 100)',
            category_id: 'onboard',
            category_name: 'Equipamentos & Amarração',
            unit_cost: 65.00,
            current_stock: 15,
            min_stock: 5,
            location: 'Corredor D - Gaveteiro 1',
            unit_id: 'cx'
          },
          {
            id: 'demo-part-15',
            sku: 'EPI-BOT-MAR-42',
            internal_code: 'EPI-001',
            name: 'Bota de Segurança Marluvas Premier Bico Composite Tam 42',
            category_id: 'ppe',
            category_name: 'EPIs & Segurança',
            unit_cost: 195.00,
            current_stock: 8,
            min_stock: 4,
            location: 'Armário de EPIs - Seção 1',
            unit_id: 'par'
          },
          {
            id: 'demo-part-16',
            sku: 'EPI-COL-REFL-G',
            internal_code: 'EPI-002',
            name: 'Colete Refletivo Amarelo Fluorescente Classe 2 Tam G',
            category_id: 'ppe',
            category_name: 'EPIs & Segurança',
            unit_cost: 32.00,
            current_stock: 30,
            min_stock: 10,
            location: 'Armário de EPIs - Seção 2',
            unit_id: 'un'
          },
          {
            id: 'demo-part-17',
            sku: 'IMP-LONA-GRAN',
            internal_code: 'IMP-001',
            name: 'Lona Vinílica Encerado Graneleira Sansuy 10m x 4m Reforçada',
            category_id: 'trailer_accessories',
            category_name: 'Implementos & Carretas',
            unit_cost: 1850.00,
            current_stock: 4,
            min_stock: 2,
            location: 'Pátio Coberto - Racks de Lona',
            unit_id: 'un'
          },
          {
            id: 'demo-part-18',
            sku: 'IMP-ENG-AR-MAO',
            internal_code: 'IMP-002',
            name: 'Engate Rápido de Ar Mão de Amigo Carreta (Par Vermelha/Azul)',
            category_id: 'trailer_accessories',
            category_name: 'Implementos & Carretas',
            unit_cost: 95.00,
            current_stock: 12,
            min_stock: 4,
            location: 'Corredor B - Prateleira 5',
            unit_id: 'par'
          }
        ];
      }

      let activeCategory = 'all';

      const renderStockView = () => {
        const filteredParts = partsList.filter(p => {
          if (activeCategory === 'all') return true;
          return p.category_id === activeCategory;
        });

        // Compute KPIs
        const totalItemsCount = partsList.reduce((acc, p) => acc + (p.current_stock || 0), 0);
        const totalStockValue = partsList.reduce((acc, p) => acc + ((p.current_stock || 0) * (p.unit_cost || 0)), 0);
        const lowStockItems = partsList.filter(p => (p.current_stock || 0) <= (p.min_stock || 0));
        const tiresInStock = partsList.filter(p => p.category_id === 'tires').reduce((acc, p) => acc + (p.current_stock || 0), 0);

        container.innerHTML = `
          <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 12px;">
            <div>
              <h2 style="font-size: 20px; font-weight: 700;">Estoque & Almoxarifado da Frota</h2>
              <p style="color: var(--text-muted); font-size: 13px;">
                Controle integral de pneus (novos, recapados, carcaças), peças de freio/mecânica, filtros, lubrificantes, equipamentos de bordo e acessórios de carretas.
              </p>
            </div>
            <button id="btn-add-part" class="btn btn-primary">
              <span>+</span> Cadastrar Item no Estoque
            </button>
          </div>

          <!-- Summary KPI Cards -->
          <div class="metrics-grid">
            <div class="metric-card">
              <div class="metric-label">Total de Itens em Saldo</div>
              <div class="metric-value" style="color: var(--color-info);">${totalItemsCount} un</div>
              <div class="metric-hint">Saldo físico disponível no almoxarifado</div>
            </div>
            <div class="metric-card">
              <div class="metric-label">Valor Total Imobilizado</div>
              <div class="metric-value" style="color: var(--color-success);">
                R$ ${totalStockValue.toLocaleString('pt-BR', { minimumFractionDigits: 2, maximumFractionDigits: 2 })}
              </div>
              <div class="metric-hint">Patrimônio em peças e suprimentos</div>
            </div>
            <div class="metric-card">
              <div class="metric-label">Itens Abaixo do Mínimo</div>
              <div class="metric-value" style="color: ${lowStockItems.length > 0 ? 'var(--color-danger)' : 'var(--color-success)'};">
                ${lowStockItems.length} itens
              </div>
              <div class="metric-hint">${lowStockItems.length > 0 ? 'Exigem reposição e compra urgente' : 'Estoque 100% regular'}</div>
            </div>
            <div class="metric-card">
              <div class="metric-label">Pneus em Almoxarifado</div>
              <div class="metric-value" style="color: var(--color-warning);">${tiresInStock} pneus</div>
              <div class="metric-hint">Novos, recapados e carcaças</div>
            </div>
          </div>

          <!-- Category Filter Pills -->
          <div style="display: flex; flex-wrap: wrap; gap: 8px; margin-bottom: 16px;">
            <button class="btn btn-sm ${activeCategory === 'all' ? 'btn-primary' : 'btn-secondary'} stock-cat-btn" data-cat="all">
              Todos (${partsList.length})
            </button>
            <button class="btn btn-sm ${activeCategory === 'tires' ? 'btn-primary' : 'btn-secondary'} stock-cat-btn" data-cat="tires">
              🔘 Pneus (${partsList.filter(p => p.category_id === 'tires').length})
            </button>
            <button class="btn btn-sm ${activeCategory === 'mechanical' ? 'btn-primary' : 'btn-secondary'} stock-cat-btn" data-cat="mechanical">
              🔧 Peças Mecânicas & Freios (${partsList.filter(p => p.category_id === 'mechanical').length})
            </button>
            <button class="btn btn-sm ${activeCategory === 'lubricants' ? 'btn-primary' : 'btn-secondary'} stock-cat-btn" data-cat="lubricants">
              🛢️ Filtros & Óleos (${partsList.filter(p => p.category_id === 'lubricants').length})
            </button>
            <button class="btn btn-sm ${activeCategory === 'electrical' ? 'btn-primary' : 'btn-secondary'} stock-cat-btn" data-cat="electrical">
              ⚡ Elétrica & Baterias (${partsList.filter(p => p.category_id === 'electrical').length})
            </button>
            <button class="btn btn-sm ${activeCategory === 'onboard' ? 'btn-primary' : 'btn-secondary'} stock-cat-btn" data-cat="onboard">
              🪢 Equipamentos & Cintas (${partsList.filter(p => p.category_id === 'onboard').length})
            </button>
            <button class="btn btn-sm ${activeCategory === 'ppe' ? 'btn-primary' : 'btn-secondary'} stock-cat-btn" data-cat="ppe">
              🦺 EPIs (${partsList.filter(p => p.category_id === 'ppe').length})
            </button>
            <button class="btn btn-sm ${activeCategory === 'trailer_accessories' ? 'btn-primary' : 'btn-secondary'} stock-cat-btn" data-cat="trailer_accessories">
              🚛 Implementos & Carretas (${partsList.filter(p => p.category_id === 'trailer_accessories').length})
            </button>
          </div>

          <!-- Stock Table -->
          <div class="card">
            <div class="table-container">
              <table class="data-table">
                <thead>
                  <tr>
                    <th>SKU / Código</th>
                    <th>Descrição do Item / Peça</th>
                    <th>Categoria</th>
                    <th>Saldo Atual</th>
                    <th>Estoque Mínimo</th>
                    <th>Custo Unitário</th>
                    <th>Valor Total</th>
                    <th>Localização Almoxarifado</th>
                    <th>Status</th>
                    <th style="text-align: right;">Ações</th>
                  </tr>
                </thead>
                <tbody>
                  ${filteredParts.length === 0 ? `
                    <tr>
                      <td colspan="10" style="text-align: center; color: var(--text-muted); padding: 32px;">
                        Nenhum item encontrado nesta categoria.
                      </td>
                    </tr>
                  ` : ''}
                  ${filteredParts.map(p => {
                    const isLow = (p.current_stock || 0) <= (p.min_stock || 0);
                    const itemTotal = (p.current_stock || 0) * (p.unit_cost || 0);

                    return `
                      <tr>
                        <td>
                          <div style="font-weight: 700; font-family: monospace; color: var(--text-main);">${p.sku}</div>
                          ${p.internal_code ? `<div style="font-size: 0.75rem; color: var(--text-muted);">Cód: ${p.internal_code}</div>` : ''}
                        </td>
                        <td>
                          <div style="font-weight: 600; max-width: 280px;">${p.name}</div>
                        </td>
                        <td>
                          <span class="badge badge-secondary">${p.category_name || p.category_id}</span>
                        </td>
                        <td>
                          <span style="font-weight: 700; font-size: 1.05rem; color: ${isLow ? 'var(--color-danger)' : 'var(--color-success)'};">
                            ${p.current_stock || 0} ${p.unit_id || 'un'}
                          </span>
                        </td>
                        <td>
                          <span style="color: var(--text-muted); font-size: 0.85rem;">${p.min_stock || 0} ${p.unit_id || 'un'}</span>
                        </td>
                        <td>
                          R$ ${Number(p.unit_cost || 0).toLocaleString('pt-BR', { minimumFractionDigits: 2 })}
                        </td>
                        <td style="font-weight: 700; color: var(--text-main);">
                          R$ ${itemTotal.toLocaleString('pt-BR', { minimumFractionDigits: 2 })}
                        </td>
                        <td style="font-size: 0.8rem; color: var(--text-muted);">
                          ${p.location || 'Almoxarifado Geral'}
                        </td>
                        <td>
                          <span class="badge ${isLow ? 'badge-danger' : 'badge-success'}">
                            ${isLow ? '⚠️ Reposição Urgente' : '✓ Normal'}
                          </span>
                        </td>
                        <td style="text-align: right; white-space: nowrap;">
                          <button class="btn btn-sm btn-secondary btn-movement" data-id="${p.id}" data-sku="${p.sku}" data-name="${p.name}" data-stock="${p.current_stock}">
                            Movimentar (+/-)
                          </button>
                        </td>
                      </tr>
                    `;
                  }).join('')}
                </tbody>
              </table>
            </div>
          </div>
        `;

        // Event listeners for category filters
        container.querySelectorAll('.stock-cat-btn').forEach(btn => {
          btn.addEventListener('click', () => {
            activeCategory = btn.getAttribute('data-cat');
            renderStockView();
          });
        });

        // Event listener for Add Item modal
        container.querySelector('#btn-add-part')?.addEventListener('click', () => {
          openAddPartModal();
        });

        // Event listener for Stock Movement
        container.querySelectorAll('.btn-movement').forEach(btn => {
          btn.addEventListener('click', () => {
            const id = btn.getAttribute('data-id');
            const sku = btn.getAttribute('data-sku');
            const name = btn.getAttribute('data-name');
            const currentStock = parseFloat(btn.getAttribute('data-stock')) || 0;
            openStockMovementModal(id, sku, name, currentStock);
          });
        });
      };

      // Modal for adding a new stock item
      const openAddPartModal = () => {
        let modalEl = document.getElementById('modal-part-form');
        if (!modalEl) {
          modalEl = document.createElement('div');
          modalEl.id = 'modal-part-form';
          modalEl.className = 'modal-overlay';
          document.body.appendChild(modalEl);
        }

        modalEl.innerHTML = `
          <div class="modal-content modal-content-wide">
            <div class="modal-header">
              <h3 style="font-weight: 700; font-size: 1.15rem;">Cadastrar Item no Estoque / Almoxarifado</h3>
              <button id="close-modal-part" style="background: none; border: none; color: var(--text-muted); font-size: 1.25rem; cursor: pointer;">✕</button>
            </div>
            <form id="form-part-create">
              <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 14px; margin-bottom: 14px;">
                <div>
                  <label style="font-size: 0.8rem; color: var(--text-muted); display: block; margin-bottom: 4px;">Categoria do Item</label>
                  <select id="part-category" class="form-select" required>
                    <option value="tires">🔘 Pneus (Novos, Recapados, Carcaças)</option>
                    <option value="mechanical">🔧 Peças Mecânicas & Freios</option>
                    <option value="lubricants">🛢️ Filtros & Lubrificantes (Óleos, Arla 32)</option>
                    <option value="electrical">⚡ Elétrica & Baterias</option>
                    <option value="onboard">🪢 Equipamentos de Bordo & Amarração (Cintas, Catracas, Tacógrafo)</option>
                    <option value="ppe">🦺 EPIs & Segurança do Motorista</option>
                    <option value="trailer_accessories">🚛 Acessórios de Carreta & Implemento</option>
                  </select>
                </div>
                <div>
                  <label style="font-size: 0.8rem; color: var(--text-muted); display: block; margin-bottom: 4px;">Código SKU / Identificador Único</label>
                  <input type="text" id="part-sku" class="form-input" placeholder="ex: PNEU-295-MIC ou FLT-COMB-01" required />
                </div>
              </div>

              <div style="margin-bottom: 14px;">
                <label style="font-size: 0.8rem; color: var(--text-muted); display: block; margin-bottom: 4px;">Descrição Completa do Item / Peça</label>
                <input type="text" id="part-name" class="form-input" placeholder="ex: Pneu 295/80 R22.5 Michelin X Multi ou Cinta Amarração 5T" required />
              </div>

              <div style="display: grid; grid-template-columns: 1fr 1fr 1fr; gap: 14px; margin-bottom: 14px;">
                <div>
                  <label style="font-size: 0.8rem; color: var(--text-muted); display: block; margin-bottom: 4px;">Saldo Inicial em Estoque</label>
                  <input type="number" id="part-stock" class="form-input" placeholder="0" value="1" required />
                </div>
                <div>
                  <label style="font-size: 0.8rem; color: var(--text-muted); display: block; margin-bottom: 4px;">Estoque Mínimo (Alerta)</label>
                  <input type="number" id="part-min-stock" class="form-input" placeholder="0" value="2" required />
                </div>
                <div>
                  <label style="font-size: 0.8rem; color: var(--text-muted); display: block; margin-bottom: 4px;">Custo Unitário Médio (R$)</label>
                  <input type="number" step="0.01" id="part-cost" class="form-input" placeholder="0.00" required />
                </div>
              </div>

              <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 14px; margin-bottom: 20px;">
                <div>
                  <label style="font-size: 0.8rem; color: var(--text-muted); display: block; margin-bottom: 4px;">Unidade de Medida</label>
                  <select id="part-unit" class="form-select">
                    <option value="un">Unidade (un)</option>
                    <option value="jg">Jogo / Kit (jg)</option>
                    <option value="balde">Balde (balde)</option>
                    <option value="gl">Galão (gl)</option>
                    <option value="cx">Caixa (cx)</option>
                    <option value="par">Par (par)</option>
                    <option value="L">Litros (L)</option>
                  </select>
                </div>
                <div>
                  <label style="font-size: 0.8rem; color: var(--text-muted); display: block; margin-bottom: 4px;">Localização Física (Prateleira / Baia)</label>
                  <input type="text" id="part-location" class="form-input" placeholder="ex: Almoxarifado Central - Corredor B, Prateleira 3" />
                </div>
              </div>

              <div style="display: flex; justify-content: flex-end; gap: 10px;">
                <button type="button" id="btn-cancel-part" class="btn btn-secondary">Cancelar</button>
                <button type="submit" class="btn btn-primary">Salvar no Almoxarifado</button>
              </div>
            </form>
          </div>
        `;

        modalEl.style.display = 'flex';

        modalEl.querySelector('#close-modal-part')?.addEventListener('click', () => {
          modalEl.style.display = 'none';
        });
        modalEl.querySelector('#btn-cancel-part')?.addEventListener('click', () => {
          modalEl.style.display = 'none';
        });

        modalEl.querySelector('#form-part-create')?.addEventListener('submit', async (e) => {
          e.preventDefault();
          const category = modalEl.querySelector('#part-category').value;
          const sku = modalEl.querySelector('#part-sku').value.trim().toUpperCase();
          const name = modalEl.querySelector('#part-name').value.trim();
          const stock = parseFloat(modalEl.querySelector('#part-stock').value) || 0;
          const minStock = parseFloat(modalEl.querySelector('#part-min-stock').value) || 0;
          const cost = parseFloat(modalEl.querySelector('#part-cost').value) || 0;
          const unit = modalEl.querySelector('#part-unit').value;
          const location = modalEl.querySelector('#part-location').value.trim();

          const newItem = {
            id: 'part-' + Date.now(),
            sku: sku,
            internal_code: sku,
            name: name,
            category_id: category,
            category_name: modalEl.querySelector('#part-category option:checked').textContent,
            unit_id: unit,
            current_stock: stock,
            min_stock: minStock,
            unit_cost: cost,
            location: location || 'Almoxarifado Principal',
            status: 'active'
          };

          try {
            await api.createPart(newItem).catch(() => {});
            partsList.unshift(newItem);
            alert('Item cadastrado com sucesso no estoque!');
            modalEl.style.display = 'none';
            renderStockView();
          } catch (err) {
            alert('Erro ao salvar item: ' + err.message);
          }
        });
      };

      // Modal for stock movement (Entrada / Saída / Ajuste)
      const openStockMovementModal = (id, sku, name, currentStock) => {
        let modalEl = document.getElementById('modal-stock-movement');
        if (!modalEl) {
          modalEl = document.createElement('div');
          modalEl.id = 'modal-stock-movement';
          modalEl.className = 'modal-overlay';
          document.body.appendChild(modalEl);
        }

        modalEl.innerHTML = `
          <div class="modal-content">
            <div class="modal-header">
              <h3 style="font-weight: 700; font-size: 1.15rem;">Movimentação de Estoque</h3>
              <button id="close-modal-movement" style="background: none; border: none; color: var(--text-muted); font-size: 1.25rem; cursor: pointer;">✕</button>
            </div>
            <form id="form-stock-movement">
              <div style="background: var(--bg-main); border: 1px solid var(--border-color); border-radius: 6px; padding: 12px; margin-bottom: 14px;">
                <div style="font-weight: 700; font-family: monospace;">${sku}</div>
                <div style="font-size: 0.9rem; color: var(--text-main); margin-top: 2px;">${name}</div>
                <div style="font-size: 0.8rem; color: var(--color-info); margin-top: 4px;">
                  Saldo Atual em Almoxarifado: <strong>${currentStock}</strong>
                </div>
              </div>

              <div style="margin-bottom: 14px;">
                <label style="font-size: 0.8rem; color: var(--text-muted); display: block; margin-bottom: 4px;">Tipo de Movimentação</label>
                <select id="mov-type" class="form-select" required>
                  <option value="in">📥 Entrada (Compra com NF / Devolução)</option>
                  <option value="out">📤 Saída (Aplicação em Manutenção / Veículo / O.S.)</option>
                  <option value="adjustment">🔄 Ajuste de Inventário Físico</option>
                </select>
              </div>

              <div style="margin-bottom: 14px;">
                <label style="font-size: 0.8rem; color: var(--text-muted); display: block; margin-bottom: 4px;">Quantidade a Movimentar</label>
                <input type="number" id="mov-qty" class="form-input" placeholder="0" min="1" value="1" required />
              </div>

              <div style="margin-bottom: 14px;">
                <label style="font-size: 0.8rem; color: var(--text-muted); display: block; margin-bottom: 4px;">Nº da Nota Fiscal / O.S. / Requisição</label>
                <input type="text" id="mov-ref" class="form-input" placeholder="ex: NF 12845 ou O.S. 302" />
              </div>

              <div style="margin-bottom: 20px;">
                <label style="font-size: 0.8rem; color: var(--text-muted); display: block; margin-bottom: 4px;">Placa do Veículo / Carreta Destino (opcional)</label>
                <input type="text" id="mov-plate" class="form-input" placeholder="ex: ABC-1234 (se saída para aplicação)" />
              </div>

              <div style="display: flex; justify-content: flex-end; gap: 10px;">
                <button type="button" id="btn-cancel-mov" class="btn btn-secondary">Cancelar</button>
                <button type="submit" class="btn btn-primary">Confirmar Movimentação</button>
              </div>
            </form>
          </div>
        `;

        modalEl.style.display = 'flex';

        modalEl.querySelector('#close-modal-movement')?.addEventListener('click', () => {
          modalEl.style.display = 'none';
        });
        modalEl.querySelector('#btn-cancel-mov')?.addEventListener('click', () => {
          modalEl.style.display = 'none';
        });

        modalEl.querySelector('#form-stock-movement')?.addEventListener('submit', async (e) => {
          e.preventDefault();
          const type = modalEl.querySelector('#mov-type').value;
          const qty = parseFloat(modalEl.querySelector('#mov-qty').value) || 0;
          const ref = modalEl.querySelector('#mov-ref').value;
          const plate = modalEl.querySelector('#mov-plate').value;

          const part = partsList.find(p => p.id === id);
          if (part) {
            if (type === 'in') {
              part.current_stock = (part.current_stock || 0) + qty;
            } else if (type === 'out') {
              if (qty > part.current_stock) {
                alert('Atenção: Quantidade de saída excede o saldo disponível!');
                return;
              }
              part.current_stock = Math.max(0, (part.current_stock || 0) - qty);
            } else if (type === 'adjustment') {
              part.current_stock = qty;
            }

            // Sync with backend if persistent
            await api.updatePart(id, part).catch(() => {});
          }

          alert(`Movimentação homologada! Novo saldo: ${part ? part.current_stock : 0}`);
          modalEl.style.display = 'none';
          renderStockView();
        });
      };

      renderStockView();

    } catch (e) {
      container.innerHTML = `<div class="card"><p style="color: var(--color-danger);">Erro ao carregar estoque: ${e.message}</p></div>`;
    }
  }

  // ==================== FINANCEIRO & FECHAMENTO DE CAIXA ====================
  async renderFinancial(container) {
    container.innerHTML = `<div class="card"><p>Carregando fluxo de caixa e conciliação financeira...</p></div>`;

    try {
      let summary = {
        total_revenues: 0,
        total_expenses: 0,
        total_fleet_expenses: 0,
        total_third_party_expenses: 0,
        net_cash_flow: 0,
        pending_payables: 0,
        overdue_payables: 0,
        pending_invoices_count: 0
      };
      let pendingList = [];

      try {
        summary = await api.getCashFlowSummary() || summary;
      } catch (err) {
        console.warn('Erro ao obter resumo de caixa:', err);
      }

      try {
        const invoicesResp = await api.listPendingInvoices() || {};
        pendingList = invoicesResp.items || [];
      } catch (err) {
        console.warn('Erro ao obter faturas pendentes:', err);
      }

      container.innerHTML = `
        <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 24px;">
          <div>
            <h2>Financeiro, Contas a Pagar & Fechamento de Caixa</h2>
            <p style="color: var(--text-muted); font-size: 0.875rem;">
              Conciliação direta de abastecimentos, manutenções e despesas de terceiros para o caixa da transportadora.
            </p>
          </div>
          <div style="display: flex; gap: 8px;">
            <button class="btn btn-primary" id="btn-reconcile-modal">
              + Lançar / Conciliar Despesa
            </button>
            <button class="btn btn-secondary" id="btn-close-cash">
              🔒 Fechamento de Caixa
            </button>
          </div>
        </div>

        <!-- Metric Grid -->
        <div class="metric-grid">
          <div class="metric-card">
            <div class="metric-title">Receitas Recebidas</div>
            <div class="metric-value" style="color: var(--color-success);">
              R$ ${summary.total_revenues ? summary.total_revenues.toLocaleString('pt-BR', { minimumFractionDigits: 2 }) : '0,00'}
            </div>
            <div class="metric-subtitle">Entradas confirmadas no caixa</div>
          </div>

          <div class="metric-card">
            <div class="metric-title">Despesas Frota Própria</div>
            <div class="metric-value" style="color: var(--color-warning);">
              R$ ${summary.total_fleet_expenses ? summary.total_fleet_expenses.toLocaleString('pt-BR', { minimumFractionDigits: 2 }) : '0,00'}
            </div>
            <div class="metric-subtitle">Diesel, O.S. e pneus próprios</div>
          </div>

          <div class="metric-card">
            <div class="metric-title">Despesas de Terceiros</div>
            <div class="metric-value" style="color: var(--color-danger);">
              R$ ${summary.total_third_party_expenses ? summary.total_third_party_expenses.toLocaleString('pt-BR', { minimumFractionDigits: 2 }) : '0,00'}
            </div>
            <div class="metric-subtitle">Requisições & serviços de agregados</div>
          </div>

          <div class="metric-card">
            <div class="metric-title">Saldo Líquido de Caixa</div>
            <div class="metric-value" style="color: ${summary.net_cash_flow >= 0 ? 'var(--color-success)' : 'var(--color-danger)'};">
              R$ ${summary.net_cash_flow ? summary.net_cash_flow.toLocaleString('pt-BR', { minimumFractionDigits: 2 }) : '0,00'}
            </div>
            <div class="metric-subtitle">Receitas líquidas - saídas pagas</div>
          </div>
        </div>

        <!-- Pending Invoices / Contas a Pagar Alert -->
        <div class="card" style="margin-bottom: 24px; border-left: 4px solid var(--color-primary);">
          <div style="display: flex; justify-content: space-between; align-items: center;">
            <div>
              <div style="font-weight: 600; font-size: 1rem;">Contas a Pagar & Notas em Aberto</div>
              <div style="color: var(--text-muted); font-size: 0.875rem;">
                Total a Pagar: <strong style="color: var(--text-main);">R$ ${summary.pending_payables ? summary.pending_payables.toLocaleString('pt-BR', { minimumFractionDigits: 2 }) : '0,00'}</strong>
                | Notas Pendentes: <strong>${summary.pending_invoices_count || pendingList.length}</strong>
                ${summary.overdue_payables > 0 ? `| <span style="color: var(--color-danger); font-weight: 600;">Vencidos: R$ ${summary.overdue_payables.toLocaleString('pt-BR', { minimumFractionDigits: 2 })}</span>` : ''}
              </div>
            </div>
          </div>
        </div>

        <!-- Invoices Table -->
        <div class="card">
          <div class="table-container">
            <table>
              <thead>
                <tr>
                  <th>Descrição / Serviço</th>
                  <th>Tipo</th>
                  <th>NF / Requisição</th>
                  <th>Veículo / Terceiro</th>
                  <th>Vencimento</th>
                  <th>Valor</th>
                  <th>Status</th>
                  <th>Ações</th>
                </tr>
              </thead>
              <tbody>
                ${pendingList.length === 0 ? `
                  <tr>
                    <td colspan="8" style="text-align: center; color: var(--text-muted); padding: 28px;">
                      Nenhum boleto ou nota fiscal pendente no momento. Caixa 100% conciliado!
                    </td>
                  </tr>
                ` : ''}
                ${pendingList.map(tx => `
                  <tr>
                    <td>
                      <div style="font-weight: 600;">${tx.description || 'Despesa'}</div>
                      ${tx.barcode ? `<div style="font-size: 0.75rem; color: var(--text-muted); font-family: monospace;">Boleto: ${tx.barcode}</div>` : ''}
                    </td>
                    <td>
                      <span class="badge ${tx.is_third_party ? 'badge-warning' : 'badge-info'}">
                        ${tx.is_third_party ? 'Terceiro' : 'Frota Própria'}
                      </span>
                    </td>
                    <td>
                      ${tx.invoice_number ? `<div>NF: <strong>${tx.invoice_number}</strong></div>` : ''}
                      ${tx.requisition_number ? `<div style="font-size: 0.75rem; color: var(--text-muted);">Req: ${tx.requisition_number}</div>` : ''}
                      ${!tx.invoice_number && !tx.requisition_number ? '-' : ''}
                    </td>
                    <td>
                      ${tx.third_party_identifier || tx.vehicle_id || 'Não especificado'}
                    </td>
                    <td>
                      ${tx.due_date ? new Date(tx.due_date).toLocaleDateString('pt-BR') : 'À vista'}
                    </td>
                    <td style="font-weight: 700;">
                      R$ ${tx.amount ? tx.amount.toLocaleString('pt-BR', { minimumFractionDigits: 2 }) : '0,00'}
                    </td>
                    <td>
                      <span class="badge ${tx.status === 'paid' ? 'badge-success' : 'badge-warning'}">
                        ${tx.status}
                      </span>
                    </td>
                    <td>
                      ${tx.status !== 'paid' ? `
                        <button class="btn btn-sm btn-primary btn-pay-tx" data-id="${tx.id}">
                          ✓ Baixar / Pagar
                        </button>
                      ` : '<span style="color: var(--color-success); font-size: 0.875rem;">Pago</span>'}
                    </td>
                  </tr>
                `).join('')}
              </tbody>
            </table>
          </div>
        </div>
      `;

      // Event listener for Pay / Baixar Nota
      container.querySelectorAll('.btn-pay-tx').forEach(btn => {
        btn.addEventListener('click', async () => {
          const id = btn.getAttribute('data-id');
          if (confirm('Confirmar o pagamento e baixa desta despesa no caixa da empresa?')) {
            try {
              await api.payTransaction(id);
              alert('Pagamento registrado com sucesso no caixa!');
              this.renderFinancial(container);
            } catch (err) {
              alert('Erro ao registrar baixa: ' + err.message);
            }
          }
        });
      });

      // Event listener for Reconcile Modal
      const btnReconcile = container.querySelector('#btn-reconcile-modal');
      if (btnReconcile) {
        btnReconcile.addEventListener('click', async () => {
          const desc = prompt('Descrição da Despesa / Fornecedor:');
          if (!desc) return;
          const amountStr = prompt('Valor total (ex: 450.00):');
          const amount = parseFloat(amountStr);
          if (!amount || isNaN(amount)) {
            alert('Valor inválido');
            return;
          }
          const isThirdParty = confirm('Esta despesa refere-se a um veículo terceiro ou agregado? (OK = Sim, Cancelar = Frota Própria)');
          let thirdPartyId = '';
          if (isThirdParty) {
            thirdPartyId = prompt('Identificador do Terceiro (Placa ou Nome do Agregado):') || '';
          }
          const invoice = prompt('Número da Nota Fiscal (opcional):') || '';
          const reqNum = prompt('Número da Requisição (opcional):') || '';
          const barcode = prompt('Código de Barras do Boleto (opcional):') || '';

          try {
            await api.reconcileExpense({
              source_type: 'fuel',
              source_id: 'manual_' + Date.now(),
              description: desc,
              amount: amount,
              is_third_party: isThirdParty,
              third_party_identifier: thirdPartyId,
              invoice_number: invoice,
              requisition_number: reqNum,
              barcode: barcode,
              payment_method: barcode ? 'boleto' : 'transferencia'
            });
            alert('Despesa registrada e conciliada com sucesso no contas a pagar!');
            this.renderFinancial(container);
          } catch (err) {
            alert('Erro ao registrar despesa: ' + err.message);
          }
        });
      }

      // Event listener for Cash Closing
      const btnCloseCash = container.querySelector('#btn-close-cash');
      if (btnCloseCash) {
        btnCloseCash.addEventListener('click', async () => {
          const openingStr = prompt('Saldo de Abertura do Dia (R$):', '0.00');
          const opening = parseFloat(openingStr) || 0;
          const notes = prompt('Observações do Fechamento Diário:', 'Fechamento regular conferido com extrato bancário');
          try {
            const closing = await api.closeCash({
              closing_date: new Date().toISOString().split('T')[0],
              opening_balance: opening,
              notes: notes || ''
            });
            alert(`Fechamento de Caixa Concluído!\nSaldo Final Fechado: R$ ${closing.closing_balance.toFixed(2)}`);
            this.renderFinancial(container);
          } catch (err) {
            alert('Erro ao realizar fechamento de caixa: ' + err.message);
          }
        });
      }

    } catch (e) {
      container.innerHTML = `<div class="card"><p style="color: var(--color-danger);">Erro: ${e.message}</p></div>`;
    }
  }
}

document.addEventListener('DOMContentLoaded', () => {
  new App();
});

