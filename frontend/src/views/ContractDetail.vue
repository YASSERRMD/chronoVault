<template>
  <div class="contract-detail">
    <aside class="sidebar">
      <div class="logo"><h2>ChronoVault</h2></div>
      <nav>
        <router-link to="/">Dashboard</router-link>
        <router-link to="/contracts">Contracts</router-link>
        <router-link to="/obligations">Obligations</router-link>
        <router-link to="/financial">Financial</router-link>
        <router-link to="/audit">Audit</router-link>
      </nav>
      <div class="user-info">
        <button @click="handleLogout" class="btn-logout">Logout</button>
      </div>
    </aside>
    
    <main class="main-content">
      <header>
        <div>
          <button @click="$router.back()" class="btn-back">← Back</button>
          <h1>{{ contract.title }}</h1>
        </div>
        <button @click="showClauseModal = true" class="btn-primary">Add Clause</button>
      </header>
      
      <div class="contract-info">
        <div class="info-card">
          <h3>Counterparty</h3>
          <p>{{ contract.counterparty }}</p>
        </div>
        <div class="info-card">
          <h3>Start Date</h3>
          <p>{{ contract.start_date }}</p>
        </div>
        <div class="info-card">
          <h3>End Date</h3>
          <p>{{ contract.end_date }}</p>
        </div>
        <div class="info-card">
          <h3>Value</h3>
          <p>${{ formatNumber(contract.value) }}</p>
        </div>
        <div class="info-card">
          <h3>Status</h3>
          <span class="status-badge" :class="contract.status">{{ contract.status }}</span>
        </div>
      </div>
      
      <div class="clauses-section">
        <h2>Clauses</h2>
        <div v-if="clauses.length === 0" class="empty-state">No clauses yet</div>
        <div v-else class="clause-list">
          <div v-for="clause in clauses" :key="clause.id" class="clause-card">
            <div class="clause-header">
              <h3>{{ clause.title }}</h3>
              <div class="clause-actions">
                <button @click="addObligation(clause.id)" class="btn-small">+ Obligation</button>
                <button @click="deleteClause(clause.id)" class="btn-small danger">Delete</button>
              </div>
            </div>
            <p v-if="clause.description">{{ clause.description }}</p>
            
            <div v-if="clause.obligations && clause.obligations.length" class="obligations-list">
              <h4>Obligations</h4>
              <div v-for="ob in clause.obligations" :key="ob.id" class="obligation-item">
                <div class="ob-info">
                  <span class="ob-desc">{{ ob.description }}</span>
                  <span class="status-badge" :class="ob.status">{{ ob.status }}</span>
                </div>
                <div class="ob-details">
                  <span v-if="ob.due_date">Due: {{ ob.due_date }}</span>
                  <span v-if="ob.penalty_amount">Penalty: ${{ ob.penalty_amount }}</span>
                  <span v-if="ob.responsible_party">Responsible: {{ ob.responsible_party }}</span>
                </div>
                <button v-if="ob.status !== 'fulfilled'" @click="fulfillObligation(ob.id)" class="btn-small success">Fulfill</button>
              </div>
            </div>
          </div>
        </div>
      </div>
      
      <div v-if="showClauseModal" class="modal">
        <div class="modal-content">
          <h2>Add Clause</h2>
          <form @submit.prevent="createClause">
            <div class="form-group">
              <label>Title</label>
              <input v-model="newClause.title" required />
            </div>
            <div class="form-group">
              <label>Description</label>
              <textarea v-model="newClause.description"></textarea>
            </div>
            <div class="modal-actions">
              <button type="button" @click="showClauseModal = false" class="btn-secondary">Cancel</button>
              <button type="submit" class="btn-primary">Create</button>
            </div>
          </form>
        </div>
      </div>
      
      <div v-if="showObligationModal" class="modal">
        <div class="modal-content">
          <h2>Add Obligation</h2>
          <form @submit.prevent="createObligation">
            <div class="form-group">
              <label>Description</label>
              <input v-model="newObligation.description" required />
            </div>
            <div class="form-group">
              <label>Due Date</label>
              <input v-model="newObligation.due_date" type="date" />
            </div>
            <div class="form-group">
              <label>Penalty Amount</label>
              <input v-model.number="newObligation.penalty_amount" type="number" />
            </div>
            <div class="form-group">
              <label>Penalty Type</label>
              <select v-model="newObligation.penalty_type">
                <option value="fixed">Fixed</option>
                <option value="daily">Daily</option>
              </select>
            </div>
            <div class="form-group">
              <label>Responsible Party</label>
              <input v-model="newObligation.responsible_party" />
            </div>
            <div class="modal-actions">
              <button type="button" @click="showObligationModal = false" class="btn-secondary">Cancel</button>
              <button type="submit" class="btn-primary">Create</button>
            </div>
          </form>
        </div>
      </div>
    </main>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import api from '../services/api'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

const contract = ref({})
const clauses = ref([])
const showClauseModal = ref(false)
const showObligationModal = ref(false)
const selectedClauseId = ref(null)
const newClause = ref({ title: '', description: '' })
const newObligation = ref({ description: '', due_date: '', penalty_amount: 0, penalty_type: 'fixed', responsible_party: '' })

const fetchContract = async () => {
  try {
    const response = await api.get(`/contracts/${route.params.id}`)
    contract.value = response.data
  } catch (err) {
    console.error('Failed to fetch contract:', err)
  }
}

const fetchClauses = async () => {
  try {
    const response = await api.get(`/contracts/${route.params.id}/clauses`)
    clauses.value = response.data
    for (const clause of clauses.value) {
      const obRes = await api.get(`/obligations?clause_id=${clause.id}`)
      clause.obligations = obRes.data.data || []
    }
  } catch (err) {
    console.error('Failed to fetch clauses:', err)
  }
}

const formatNumber = (num) => num ? new Intl.NumberFormat().format(num) : '0'

const createClause = async () => {
  try {
    await api.post('/clauses', { ...newClause.value, contract_id: route.params.id })
    showClauseModal.value = false
    newClause.value = { title: '', description: '' }
    fetchClauses()
  } catch (err) {
    console.error('Failed to create clause:', err)
  }
}

const deleteClause = async (id) => {
  if (!confirm('Delete this clause?')) return
  try {
    await api.delete(`/clauses/${id}`)
    fetchClauses()
  } catch (err) {
    console.error('Failed to delete clause:', err)
  }
}

const addObligation = (clauseId) => {
  selectedClauseId.value = clauseId
  showObligationModal.value = true
}

const createObligation = async () => {
  try {
    await api.post('/obligations', { ...newObligation.value, clause_id: selectedClauseId.value })
    showObligationModal.value = false
    newObligation.value = { description: '', due_date: '', penalty_amount: 0, penalty_type: 'fixed', responsible_party: '' }
    fetchClauses()
  } catch (err) {
    console.error('Failed to create obligation:', err)
  }
}

const fulfillObligation = async (id) => {
  try {
    await api.post(`/obligations/${id}/fulfill`)
    fetchClauses()
  } catch (err) {
    console.error('Failed to fulfill obligation:', err)
  }
}

const handleLogout = () => {
  authStore.logout()
  router.push('/login')
}

onMounted(() => { fetchContract(); fetchClauses() })
</script>

<style scoped>
.contract-detail { display: flex; min-height: 100vh; }
.sidebar { width: 250px; background: #2c3e50; color: white; padding: 1rem; display: flex; flex-direction: column; }
.logo h2 { padding: 1rem 0; border-bottom: 1px solid #34495e; margin-bottom: 1rem; }
.sidebar nav { flex: 1; }
.sidebar nav a { display: block; padding: 0.75rem 1rem; color: #bdc3c7; text-decoration: none; border-radius: 4px; margin-bottom: 0.25rem; }
.sidebar nav a:hover, .sidebar nav a.router-link-active { background: #34495e; color: white; }
.btn-logout { width: 100%; padding: 0.5rem; background: #e74c3c; color: white; border: none; border-radius: 4px; cursor: pointer; }
.main-content { flex: 1; padding: 2rem; background: #f5f6fa; }
header { display: flex; justify-content: space-between; align-items: flex-start; margin-bottom: 2rem; }
.btn-back { background: none; border: none; color: #3498db; cursor: pointer-bottom: 0; margin.5rem; }
.contract-info { display: grid; grid-template-columns: repeat(auto-fit, minmax(150px, 1fr)); gap: 1rem; margin-bottom: 2rem; }
.info-card { background: white; padding: 1rem; border-radius: 8px; box-shadow: 0 2px 8px rgba(0,0,0,0.1); }
.info-card h3 { margin: 0 0 0.5rem; font-size: 0.75rem; color: #666; }
.info-card p { margin: 0; font-size: 1.125rem; font-weight: 600; }
.status-badge { padding: 0.25rem 0.75rem; border-radius: 12px; font-size: 0.75rem; font-weight: 600; }
.status-badge.draft { background: #ffeaa7; color: #d68910; }
.status-badge.active { background: #d5f5e3; color: #27ae60; }
.status-badge.pending { background: #ebdef0; color: #8e44ad; }
.status-badge.breached { background: #fadbd8; color: #e74c3c; }
.status-badge.fulfilled { background: #d5f5e3; color: #27ae60; }
.clauses-section { background: white; padding: 1.5rem; border-radius: 8px; box-shadow: 0 2px 8px rgba(0,0,0,0.1); }
.clauses-section h2 { margin: 0 0 1rem; }
.clause-card { border: 1px solid #eee; border-radius: 8px; padding: 1rem; margin-bottom: 1rem; }
.clause-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 0.5rem; }
.clause-header h3 { margin: 0; }
.clause-actions { display: flex; gap: 0.5rem; }
.btn-small { padding: 0.375rem 0.75rem; font-size: 0.75rem; background: #3498db; color: white; border: none; border-radius: 4px; cursor: pointer; }
.btn-small.danger { background: #e74c3c; }
.btn-small.success { background: #27ae60; }
.obligations-list { margin-top: 1rem; padding-top: 1rem; border-top: 1px solid #eee; }
.obligations-list h4 { margin: 0 0 0.5rem; font-size: 0.875rem; color: #666; }
.obligation-item { background: #f8f9fa; padding: 0.75rem; border-radius: 4px; margin-bottom: 0.5rem; display: flex; justify-content: space-between; align-items: center; flex-wrap: wrap; gap: 0.5rem; }
.ob-info { display: flex; align-items: center; gap: 0.5rem; }
.ob-details { font-size: 0.75rem; color: #666; display: flex; gap: 1rem; }
.modal { position: fixed; top: 0; left: 0; right: 0; bottom: 0; background: rgba(0,0,0,0.5); display: flex; align-items: center; justify-content: center; }
.modal-content { background: white; padding: 2rem; border-radius: 8px; width: 100%; max-width: 500px; }
.modal-content h2 { margin: 0 0 1.5rem; }
.form-group { margin-bottom: 1rem; }
.form-group label { display: block; margin-bottom: 0.5rem; font-weight: 500; }
.form-group input, .form-group select, .form-group textarea { width: 100%; padding: 0.75rem; border: 1px solid #ddd; border-radius: 4px; box-sizing: border-box; }
.modal-actions { display: flex; gap: 1rem; justify-content: flex-end; margin-top: 1.5rem; }
.btn-primary { padding: 0.75rem 1.5rem; background: #667eea; color: white; border: none; border-radius: 6px; cursor: pointer; }
.btn-secondary { padding: 0.75rem 1.5rem; background: #95a5a6; color: white; border: none; border-radius: 6px; cursor: pointer; }
.empty-state { padding: 2rem; text-align: center; color: #999; }
</style>
