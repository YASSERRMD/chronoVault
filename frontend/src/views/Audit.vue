<template>
  <div class="audit-page">
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
        <h1>Audit History</h1>
      </header>
      
      <div class="audit-table">
        <table>
          <thead>
            <tr>
              <th>Timestamp</th>
              <th>Action</th>
              <th>Entity Type</th>
              <th>Entity ID</th>
              <th>User ID</th>
              <th>Details</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="log in logs" :key="log.id">
              <td>{{ formatDate(log.created_at) }}</td>
              <td><span class="action-badge" :class="log.action">{{ log.action }}</span></td>
              <td>{{ log.entity_type }}</td>
              <td class="entity-id">{{ log.entity_id.substring(0, 8) }}...</td>
              <td>{{ log.user_id ? log.user_id.substring(0, 8) + '...' : '-' }}</td>
              <td>
                <button @click="showDetails(log)" class="btn-small">View</button>
              </td>
            </tr>
          </tbody>
        </table>
        <div v-if="logs.length === 0" class="empty-state">No audit logs found</div>
      </div>
      
      <div v-if="selectedLog" class="modal" @click.self="selectedLog = null">
        <div class="modal-content">
          <h2>Audit Details</h2>
          <div class="detail-row">
            <strong>Action:</strong>
            <span>{{ selectedLog.action }}</span>
          </div>
          <div class="detail-row">
            <strong>Entity Type:</strong>
            <span>{{ selectedLog.entity_type }}</span>
          </div>
          <div class="detail-row">
            <strong>Entity ID:</strong>
            <span>{{ selectedLog.entity_id }}</span>
          </div>
          <div class="detail-row">
            <strong>Timestamp:</strong>
            <span>{{ formatDate(selectedLog.created_at) }}</span>
          </div>
          <div v-if="selectedLog.old_values" class="detail-section">
            <strong>Old Values:</strong>
            <pre>{{ formatJSON(selectedLog.old_values) }}</pre>
          </div>
          <div v-if="selectedLog.new_values" class="detail-section">
            <strong>New Values:</strong>
            <pre>{{ formatJSON(selectedLog.new_values) }}</pre>
          </div>
          <button @click="selectedLog = null" class="btn-primary" style="margin-top: 1rem;">Close</button>
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

const logs = ref([])
const selectedLog = ref(null)

const fetchLogs = async () => {
  try {
    const response = await api.get('/audit')
    logs.value = response.data.data || []
  } catch (err) {
    console.error('Failed to fetch audit logs:', err)
  }
}

const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleString()
}

const formatJSON = (str) => {
  try {
    return JSON.stringify(JSON.parse(str), null, 2)
  } catch {
    return str
  }
}

const showDetails = (log) => {
  selectedLog.value = log
}

const handleLogout = () => {
  authStore.logout()
  router.push('/login')
}

onMounted(() => fetchLogs())
</script>

<style scoped>
.audit-page { display: flex; min-height: 100vh; }
.sidebar { width: 250px; background: #2c3e50; color: white; padding: 1rem; display: flex; flex-direction: column; }
.logo h2 { padding: 1rem 0; border-bottom: 1px solid #34495e; margin-bottom: 1rem; }
.sidebar nav { flex: 1; }
.sidebar nav a { display: block; padding: 0.75rem 1rem; color: #bdc3c7; text-decoration: none; border-radius: 4px; margin-bottom: 0.25rem; }
.sidebar nav a:hover, .sidebar nav a.router-link-active { background: #34495e; color: white; }
.btn-logout { width: 100%; padding: 0.5rem; background: #e74c3c; color: white; border: none; border-radius: 4px; cursor: pointer; }
.main-content { flex: 1; padding: 2rem; background: #f5f6fa; }
.audit-table { background: white; border-radius: 8px; box-shadow: 0 2px 8px rgba(0,0,0,0.1); overflow: hidden; }
table { width: 100%; border-collapse: collapse; }
th, td { padding: 1rem; text-align: left; border-bottom: 1px solid #eee; }
th { background: #f8f9fa; font-weight: 600; }
.entity-id { font-family: monospace; color: #666; }
.action-badge { padding: 0.25rem 0.5rem; border-radius: 4px; font-size: 0.75rem; font-weight: 600; }
.action-badge.create { background: #d5f5e3; color: #27ae60; }
.action-badge.update { background: #d6eaf8; color: #3498db; }
.action-badge.delete { background: #fadbd8; color: #e74c3c; }
.btn-small { padding: 0.375rem 0.75rem; background: #3498db; color: white; border: none; border-radius: 4px; cursor: pointer; }
.empty-state { padding: 2rem; text-align: center; color: #999; }
.modal { position: fixed; top: 0; left: 0; right: 0; bottom: 0; background: rgba(0,0,0,0.5); display: flex; align-items: center; justify-content: center; }
.modal-content { background: white; padding: 2rem; border-radius: 8px; width: 100%; max-width: 600px; max-height: 80vh; overflow-y: auto; }
.modal-content h2 { margin: 0 0 1rem; }
.detail-row { display: flex; gap: 1rem; padding: 0.5rem 0; border-bottom: 1px solid #eee; }
.detail-row strong { min-width: 120px; }
.detail-section { margin-top: 1rem; }
.detail-section strong { display: block; margin-bottom: 0.5rem; }
.detail-section pre { background: #f8f9fa; padding: 1rem; border-radius: 4px; overflow-x: auto; font-size: 0.875rem; }
.btn-primary { padding: 0.75rem 1.5rem; background: #667eea; color: white; border: none; border-radius: 6px; cursor: pointer; }
</style>
