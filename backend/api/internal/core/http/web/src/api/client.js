// Geevheit Intelligence 360 - API Client
class ApiClient {
  constructor(baseUrl = '') {
    if (!baseUrl) {
      const customUrl = typeof window !== 'undefined' && (window.GI360_API_URL || localStorage.getItem('gi360_api_url'));
      if (customUrl) {
        this.baseUrl = customUrl.replace(/\/+$/, '');
      } else if (typeof window !== 'undefined' && window.location.port === '3000') {
        this.baseUrl = 'http://localhost:8080';
      } else if (typeof window !== 'undefined') {
        this.baseUrl = window.location.origin;
      } else {
        this.baseUrl = 'http://localhost:8080';
      }
    } else {
      this.baseUrl = baseUrl.replace(/\/+$/, '');
    }
    this.token = localStorage.getItem('gi360_token') || '';
    this.tenantId = localStorage.getItem('gi360_tenant_id') || '';
  }

  setSession(token, tenantId) {
    this.token = token;
    this.tenantId = tenantId;
    localStorage.setItem('gi360_token', token);
    localStorage.setItem('gi360_tenant_id', tenantId);
  }

  clearSession() {
    this.token = '';
    this.tenantId = '';
    localStorage.removeItem('gi360_token');
    localStorage.removeItem('gi360_tenant_id');
  }

  async request(path, options = {}) {
    const url = `${this.baseUrl}/api/v1${path}`;
    const headers = {
      'Content-Type': 'application/json',
      ...(this.token ? { 'Authorization': `Bearer ${this.token}` } : {}),
      ...(this.tenantId ? { 'X-Tenant-ID': this.tenantId } : {}),
      ...options.headers,
    };

    try {
      const response = await fetch(url, { ...options, headers });
      if (response.status === 204) return null;
      const data = await response.json();
      if (!response.ok) {
        throw new Error(data.message || data.error || `HTTP ${response.status}`);
      }
      return data;
    } catch (err) {
      console.error(`API Error on ${path}:`, err);
      throw err;
    }
  }

  // Auth
  async login(email, password) {
    const data = await this.request('/identity/auth/login', {
      method: 'POST',
      body: JSON.stringify({ email, password })
    });
    if (data.token) {
      this.setSession(data.token, data.tenant_id || data.tenant?.id || '');
    }
    return data;
  }

  // Mission Control & Proactive Signals
  async getMissionControlSummary() {
    return this.request('/mission-control/summary');
  }

  async getMissionControlItems(filters = {}) {
    const params = new URLSearchParams(filters).toString();
    return this.request(`/mission-control/items${params ? '?' + params : ''}`);
  }

  async evaluateSignals() {
    return this.request('/mission-control/evaluate-signals', { method: 'POST' });
  }

  async acknowledgeCommandItem(id) {
    return this.request(`/mission-control/items/${id}/acknowledge`, { method: 'POST' });
  }

  async resolveCommandItem(id) {
    return this.request(`/mission-control/items/${id}/resolve`, { method: 'POST' });
  }

  // Fleet & Odometer
  async listVehicles() {
    return this.request('/fleet/vehicles');
  }

  async createVehicle(payload) {
    return this.request('/fleet/vehicles', { method: 'POST', body: JSON.stringify(payload) });
  }

  async getVehicleOdometerSummary(vehicleId) {
    return this.request(`/fleet/vehicles/${vehicleId}/odometer`);
  }

  async listVehicleOdometerHistory(vehicleId) {
    return this.request(`/fleet/vehicles/${vehicleId}/odometer/history`);
  }

  async recordOdometer(vehicleId, payload) {
    return this.request(`/fleet/vehicles/${vehicleId}/odometer`, {
      method: 'POST',
      body: JSON.stringify(payload)
    });
  }

  // Drivers & Documents
  async listDrivers() {
    return this.request('/drivers');
  }

  async createDriver(payload) {
    return this.request('/drivers', { method: 'POST', body: JSON.stringify(payload) });
  }

  // Driver Work Shifts & Journey Control (Lei 13.103/2015)
  async listDriverJourney(driverId, limit = 50) {
    return this.request(`/drivers/${driverId}/journey?limit=${limit}`);
  }

  async recordDriverJourney(driverId, payload) {
    return this.request(`/drivers/${driverId}/journey`, {
      method: 'POST',
      body: JSON.stringify(payload)
    });
  }

  async getDriverJourneySummary(driverId) {
    return this.request(`/drivers/${driverId}/journey/summary`);
  }

  async listDocuments(filters = {}) {
    const params = new URLSearchParams(filters).toString();
    return this.request(`/documents${params ? '?' + params : ''}`);
  }

  async createDocument(payload) {
    return this.request('/documents', { method: 'POST', body: JSON.stringify(payload) });
  }

  // Fuel & Invoices
  async listFuelTransactions() {
    return this.request('/fuel/transactions');
  }

  async recordFuelTransaction(payload) {
    return this.request('/fuel/transactions', { method: 'POST', body: JSON.stringify(payload) });
  }

  // Tires
  async listTires() {
    return this.request('/tires');
  }

  async recordTireInspection(tireId, payload) {
    return this.request(`/tires/${tireId}/inspections`, { method: 'POST', body: JSON.stringify(payload) });
  }

  // Maintenance
  async listWorkOrders() {
    return this.request('/maintenance/work-orders');
  }

  async createWorkOrder(payload) {
    return this.request('/maintenance/work-orders', { method: 'POST', body: JSON.stringify(payload) });
  }

  // Checklists
  async listChecklists() {
    return this.request('/checklists');
  }

  async createChecklist(payload) {
    return this.request('/checklists', { method: 'POST', body: JSON.stringify(payload) });
  }

  // Financial & Cash Flow
  async getCashFlowSummary(startDate = '', endDate = '') {
    const params = new URLSearchParams();
    if (startDate) params.set('start_date', startDate);
    if (endDate) params.set('end_date', endDate);
    const qs = params.toString();
    return this.request(`/financial/cash-flow/summary${qs ? '?' + qs : ''}`);
  }

  async listPendingInvoices(filters = {}) {
    const params = new URLSearchParams(filters).toString();
    return this.request(`/financial/pending-invoices${params ? '?' + params : ''}`);
  }

  async reconcileExpense(payload) {
    return this.request('/financial/reconcile/operational', {
      method: 'POST',
      body: JSON.stringify(payload)
    });
  }

  async closeCash(payload) {
    return this.request('/financial/cash-closing', {
      method: 'POST',
      body: JSON.stringify(payload)
    });
  }

  async payTransaction(id) {
    return this.request(`/financial/transactions/${id}/pay`, { method: 'POST' });
  }

  async approveTransaction(id) {
    return this.request(`/financial/transactions/${id}/approve`, { method: 'POST' });
  }

  // Inventory & Stock
  async listParts(filters = {}) {
    const params = new URLSearchParams(filters).toString();
    return this.request(`/inventory/parts${params ? '?' + params : ''}`);
  }

  async createPart(payload) {
    return this.request('/inventory/parts', {
      method: 'POST',
      body: JSON.stringify(payload)
    });
  }

  async updatePart(id, payload) {
    return this.request(`/inventory/parts/${id}`, {
      method: 'PUT',
      body: JSON.stringify(payload)
    });
  }

  async deletePart(id) {
    return this.request(`/inventory/parts/${id}`, { method: 'DELETE' });
  }

  async listCategories() {
    return this.request('/inventory/categories');
  }

  async createCategory(payload) {
    return this.request('/inventory/categories', {
      method: 'POST',
      body: JSON.stringify(payload)
    });
  }

  async listWarehouses() {
    return this.request('/inventory/warehouses');
  }

  async createWarehouse(payload) {
    return this.request('/inventory/warehouses', {
      method: 'POST',
      body: JSON.stringify(payload)
    });
  }

  // Trips & Cargo Dispatch
  async listTrips(filters = {}) {
    const params = new URLSearchParams(filters).toString();
    return this.request(`/trips${params ? '?' + params : ''}`);
  }

  async getTrip(id) {
    return this.request(`/trips/${id}`);
  }

  async getTripSummary() {
    return this.request('/trips/summary');
  }

  async createTrip(payload) {
    return this.request('/trips', {
      method: 'POST',
      body: JSON.stringify(payload)
    });
  }

  async dispatchTrip(id, startKM) {
    return this.request(`/trips/${id}/dispatch`, {
      method: 'POST',
      body: JSON.stringify({ start_odometer_km: startKM })
    });
  }

  async deliverTrip(id, endKM) {
    return this.request(`/trips/${id}/deliver`, {
      method: 'POST',
      body: JSON.stringify({ end_odometer_km: endKM })
    });
  }

  async cancelTrip(id, reason) {
    return this.request(`/trips/${id}/cancel`, {
      method: 'POST',
      body: JSON.stringify({ reason: reason })
    });
  }
}

export const api = new ApiClient();

