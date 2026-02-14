<template>
  <div class="contracts-page">
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
        <p>{{ authStore.currentUser?.full_name }}</p>
        <button @click="handleLogout" class="btn-logout">Logout</button>
      </div>
    </aside>
    
    <main class="main-content">
      <header>
        <h1>Contracts</h1>
        <button @click="showCreateModal = true" class="btn-primary">New Contract</button>
      </header>
      
      <div class="filters">
        <select v-model="statusFilter" @change="fetchContracts">
          <option value="">All Status</option>
          <option value="draft">Draft</option>
          <option value="active">Active</option>
          <option value="expired">Expired</option>
        </select>
      </div>
      
      <div class="contracts-table">
        <table>
          <thead>
            <tr>
              <th>Title</th>
              <th>Counterparty</th>
              <th>Start Date</th>
              <th>End Date</th>
              <th>Value</th>
              <th>Status</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="contract in contracts" :key="contract.id">
              <td>{{ contract.title }}</td>
              <td>{{ contract.counterparty }}</td>
              <td>{{ contract.start_date }}</td>
              <td>{{ contract.end_date }}</td>
              <td>${{ formatNumber(contract.value) }}</td>
              <td><span class="status-badge" :class="contract.status">{{ contract.status }}</span></td>
              <td>
                <button @click="viewContract(contract.id)" class="btn-small">View</button>
              </td>
            </tr>
          </tbody>
        </table>
        <div v-if="contracts.length === 0" class="empty-state">No contracts found</div>
      </div>
      
      <div v-if="showCreateModal" class="modal">
        <div class="modal-content">
          <h2>Create Contract</h2>
          <form @submit.prevent="createContract">
            <div class="form-group">
              <label>Title</label>
              <input v-model="newContract.title" required />
            </div>
            <div class="form-group">
              <label>Counterparty</label>
              <input v-model="newContract.counterparty" required />
            </div>
            <div class="form-group">
              <label>Start Date</label>
              <input v-model="newContract.start_date" type="date" required />
            </div>
            <div class="form-group">
              <label>End Date</label>
              <input v-model="newContract.end_date" type="date" required />
            </div>
            <div class="form-group">
              <label>Value</label>
              <input v-model.number="newContract.value" type="number" />
            </div>
            <div class="form-group">
              <label>Status</label>
              <select v-model="newContract.status">
                <option value="draft">Draft</option>
                <option value="active">Active</option>
              </select>
            </div>
            <div class="modal-actions">
              <button type="button" @click="showCreateModal = false" class="btn-secondary">Cancel</button>
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
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import api from '../services/api'

const router = useRouter()
const authStore = useAuthStore()

const contracts = ref([])
const statusFilter = ref('')
const showCreateModal = ref(false)
const newContract = ref({
  title: '',
  counterparty: '',
  start_date: '',
  end_date: '',
  value: 0,
  status: 'draft'
})

const fetchContracts = async () => {
  try {
    const params = new URLSearchParams()
    if (statusFilter.value) params.append('status', statusFilter.value)
    const response = await api.get(`/contracts?${params}`)
    contracts.value = response.data.data || []
  } catch (err) {
    console.error('Failed to fetch contracts:', err)
  }
}

const formatNumber = (num) => {
  if (!num) return '0'
  return new Intl.NumberFormat().format(num)
}

const viewContract = (id) => {
  router.push(`/contracts/${id}`)
}

const createContract = async () => {
  try {
    await api.post('/contracts', newContract.value)
    showCreateModal.value = false
    newContract.value = { title: '', counterparty: '', start_date: '', end_date: '', value: 0, status: 'draft' }
    fetchContracts()
  } catch (err) {
    console.error('Failed to create contract:', err)
  }
}

const handleLogout = () => {
  authStore.logout()
  router.push('/login')
}

onMounted(() => fetchContracts())
</script>

<style scoped>
.contracts-page { display: flex; min-height: 100vh; }
.sidebar { width: 250px; background: #2c3e50; color: white; padding: 1rem; display: flex; flex-direction: column; }
.logo h2 { padding: 1rem 0; border-bottom: 1px solid #34495e; margin-bottom: 1rem; }
.sidebar nav { flex: 1; }
.sidebar nav a { display: block; padding: 0.75rem 1rem; color: #bdc3c7; text-decoration: none; border-radius: 4px; margin-bottom: 0.25rem; }
.sidebar nav a:hover, .sidebar nav a.router-link-active { background: #34495e; color: white; }
.user-info { border-top: 1px solid #34495e; padding-top: 1rem; }
.btn-logout { margin-top: 1rem; width: 100%; padding: 0.5rem; background: #e74c3c; color: white; border: none; border-radius: 4px; cursor: pointer; }
.main-content { flex: 1; padding: 2rem; background: #f5f6fa; }
header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 2rem; }
.btn-primary { padding: 0.75rem 1.5rem; background: #667eea; color: white; border: none; border-radius: 6px; cursor: pointer; }
.filters { margin-bottom: 1rem; }
.filters select { padding: 0.5rem; border: 1px solid #ddd; border-radius: 4px; }
.contracts-table { background: white; border-radius: 8px; box-shadow: 0 2px 8px rgba(0,0,0,0.1); overflow: hidden; }
table { width: 100%; border-collapse: collapse; }
th, td { padding: 1rem; text-align: left; border-bottom: 1px solid #eee; }
th { background: #f8f9fa; font-weight: 600; }
.status-badge { padding: 0.25rem 0.75rem; border-radius: 12px; font-size: 0.75rem; font-weight: 600; }
.status-badge.draft { background: #ffeaa7; color: #d68910; }
.status-badge.active { background: #d5f5e3; color: #27ae60; }
.status-badge.expired { background: #fadbd8; color: #e74c3c; }
.btn-small { padding: 0.375rem 0.75rem; background: #3498db; color: white; border: none; border-radius: 4px; cursor: pointer; }
.empty-state { padding: 2rem; text-align: center; color: #999; }
.modal { position: fixed; top: 0; left: 0; right: 0; bottom: 0; background: rgba(0,0,0,0.5); display: flex; align-items: center; justify-content: center; }
.modal-content { background: white; padding: 2rem; border-radius: 8px; width: 100%; max-width: 500px; }
.modal-content h2 { margin: 0 0 1.5rem; }
.form-group { margin-bottom: 1rem; }
.form-group label { display: block; margin-bottom: 0.5rem; font-weight: 500; }
.form-group input, .form-group select { width: 100%; padding: 0.75rem; border: 1px solid #ddd; border-radius: 4px; box-sizing: border-box; }
.modal-actions { display: flex; gap: 1rem; justify-content: flex-end; margin-top: 1.5rem; }
.btn-secondary { padding: 0.75rem 1.5rem; background: #95a5a6; color: white; border: none; border-radius: 6px; cursor: pointer; }
</style>
