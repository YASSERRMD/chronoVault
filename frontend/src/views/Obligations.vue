<template>
  <div class="obligations-page">
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
        <h1>Obligations</h1>
      </header>
      
      <div class="filters">
        <select v-model="statusFilter" @change="fetchObligations">
          <option value="">All Status</option>
          <option value="pending">Pending</option>
          <option value="active">Active</option>
          <option value="fulfilled">Fulfilled</option>
          <option value="breached">Breached</option>
          <option value="expired">Expired</option>
        </select>
      </div>
      
      <div class="obligations-timeline">
        <div v-if="obligations.length === 0" class="empty-state">No obligations found</div>
        <div v-else class="timeline">
          <div v-for="ob in obligations" :key="ob.id" class="timeline-item" :class="ob.status">
            <div class="timeline-marker"></div>
            <div class="timeline-content">
              <div class="ob-header">
                <span class="status-badge" :class="ob.status">{{ ob.status }}</span>
                <span v-if="ob.due_date" class="due-date">Due: {{ ob.due_date }}</span>
              </div>
              <p class="ob-description">{{ ob.description }}</p>
              <div class="ob-meta">
                <span v-if="ob.penalty_amount">Penalty: ${{ ob.penalty_amount }}</span>
                <span v-if="ob.responsible_party">Responsible: {{ ob.responsible_party }}</span>
              </div>
              <button v-if="ob.status !== 'fulfilled'" @click="fulfillObligation(ob.id)" class="btn-small success">Mark Fulfilled</button>
            </div>
          </div>
        </div>
      </div>
    </main>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import wsService from '../services/websocket'
import api from '../services/api'

const router = useRouter()
const authStore = useAuthStore()

const obligations = ref([])
const statusFilter = ref('')

const fetchObligations = async () => {
  try {
    const params = new URLSearchParams()
    if (statusFilter.value) params.append('status', statusFilter.value)
    const response = await api.get(`/obligations?${params}`)
    obligations.value = response.data.data || []
  } catch (err) {
    console.error('Failed to fetch obligations:', err)
  }
}

const fulfillObligation = async (id) => {
  try {
    await api.post(`/obligations/${id}/fulfill`)
    fetchObligations()
  } catch (err) {
    console.error('Failed to fulfill obligation:', err)
  }
}

const handleWebSocketUpdate = () => {
  fetchObligations()
}

const handleLogout = () => {
  authStore.logout()
  router.push('/login')
}

onMounted(() => {
  fetchObligations()
  wsService.on('obligation_activated', handleWebSocketUpdate)
  wsService.on('obligation_breached', handleWebSocketUpdate)
  wsService.on('obligation_fulfilled', handleWebSocketUpdate)
})

onUnmounted(() => {
  wsService.off('obligation_activated', handleWebSocketUpdate)
  wsService.off('obligation_breached', handleWebSocketUpdate)
  wsService.off('obligation_fulfilled', handleWebSocketUpdate)
})
</script>

<style scoped>
.obligations-page { display: flex; min-height: 100vh; }
.sidebar { width: 250px; background: #2c3e50; color: white; padding: 1rem; display: flex; flex-direction: column; }
.logo h2 { padding: 1rem 0; border-bottom: 1px solid #34495e; margin-bottom: 1rem; }
.sidebar nav { flex: 1; }
.sidebar nav a { display: block; padding: 0.75rem 1rem; color: #bdc3c7; text-decoration: none; border-radius: 4px; margin-bottom: 0.25rem; }
.sidebar nav a:hover, .sidebar nav a.router-link-active { background: #34495e; color: white; }
.btn-logout { width: 100%; padding: 0.5rem; background: #e74c3c; color: white; border: none; border-radius: 4px; cursor: pointer; }
.main-content { flex: 1; padding: 2rem; background: #f5f6fa; }
.filters { margin-bottom: 1rem; }
.filters select { padding: 0.5rem; border: 1px solid #ddd; border-radius: 4px; }
.obligations-timeline { background: white; border-radius: 8px; box-shadow: 0 2px 8px rgba(0,0,0,0.1); padding: 1.5rem; }
.timeline { position: relative; padding-left: 2rem; }
.timeline::before { content: ''; position: absolute; left: 0.5rem; top: 0; bottom: 0; width: 2px; background: #e0e0e0; }
.timeline-item { position: relative; padding-bottom: 1.5rem; }
.timeline-marker { position: absolute; left: -1.75rem; top: 0; width: 12px; height: 12px; border-radius: 50%; background: #ccc; }
.timeline-item.pending .timeline-marker { background: #8e44ad; }
.timeline-item.active .timeline-marker { background: #3498db; }
.timeline-item.fulfilled .timeline-marker { background: #27ae60; }
.timeline-item.breached .timeline-marker { background: #e74c3c; }
.timeline-item.expired .timeline-marker { background: #95a5a6; }
.timeline-content { background: #f8f9fa; padding: 1rem; border-radius: 8px; }
.ob-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 0.5rem; }
.status-badge { padding: 0.25rem 0.75rem; border-radius: 12px; font-size: 0.75rem; font-weight: 600; }
.status-badge.pending { background: #ebdef0; color: #8e44ad; }
.status-badge.active { background: #d6eaf8; color: #3498db; }
.status-badge.fulfilled { background: #d5f5e3; color: #27ae60; }
.status-badge.breached { background: #fadbd8; color: #e74c3c; }
.status-badge.expired { background: #f0f0f0; color: #95a5a6; }
.due-date { font-size: 0.875rem; color: #666; }
.ob-description { margin: 0.5rem 0; }
.ob-meta { font-size: 0.75rem; color: #666; display: flex; gap: 1rem; margin-bottom: 0.5rem; }
.btn-small { padding: 0.375rem 0.75rem; font-size: 0.75rem; background: #27ae60; color: white; border: none; border-radius: 4px; cursor: pointer; }
.empty-state { padding: 2rem; text-align: center; color: #999; }
</style>
