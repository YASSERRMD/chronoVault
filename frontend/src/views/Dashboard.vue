<template>
  <div class="dashboard">
    <aside class="sidebar">
      <div class="logo">
        <h2>ChronoVault</h2>
      </div>
      <nav>
        <router-link to="/">Dashboard</router-link>
        <router-link to="/contracts">Contracts</router-link>
        <router-link to="/obligations">Obligations</router-link>
        <router-link to="/financial">Financial</router-link>
        <router-link to="/audit">Audit</router-link>
      </nav>
      <div class="user-info">
        <p>{{ authStore.currentUser?.full_name }}</p>
        <p class="role">{{ authStore.currentUser?.role }}</p>
        <button @click="handleLogout" class="btn-logout">Logout</button>
      </div>
    </aside>
    
    <main class="main-content">
      <header>
        <h1>Dashboard</h1>
      </header>
      
      <div class="stats-grid">
        <div class="stat-card">
          <h3>Total Contracts</h3>
          <p class="stat-value">{{ stats.total_contracts || 0 }}</p>
        </div>
        <div class="stat-card">
          <h3>Contract Value</h3>
          <p class="stat-value">${{ formatNumber(stats.total_contract_value) }}</p>
        </div>
        <div class="stat-card">
          <h3>Active Obligations</h3>
          <p class="stat-value">{{ stats.active_obligations || 0 }}</p>
        </div>
        <div class="stat-card danger">
          <h3>Total Penalties</h3>
          <p class="stat-value">${{ formatNumber(stats.total_penalties) }}</p>
        </div>
      </div>
      
      <div class="recent-activity">
        <h2>Recent Activity</h2>
        <div v-if="notifications.length === 0" class="empty-state">
          No recent activity
        </div>
        <div v-else class="activity-list">
          <div v-for="notif in notifications" :key="notif.id" class="activity-item" :class="notif.type">
            <span class="activity-icon">{{ getIcon(notif.type) }}</span>
            <div class="activity-content">
              <p>{{ notif.message }}</p>
              <span class="activity-time">{{ notif.time }}</span>
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

const stats = ref({})
const notifications = ref([])

const fetchStats = async () => {
  try {
    const response = await api.get('/reports/financial-summary')
    stats.value = response.data
  } catch (err) {
    console.error('Failed to fetch stats:', err)
  }
}

const formatNumber = (num) => {
  if (!num) return '0'
  return new Intl.NumberFormat().format(num)
}

const getIcon = (type) => {
  const icons = {
    obligation_activated: '✓',
    obligation_breached: '!',
    obligation_fulfilled: '✓',
    penalty_applied: '$'
  }
  return icons[type] || '•'
}

const handleLogout = () => {
  authStore.logout()
  router.push('/login')
}

const handleWebSocketMessage = (data) => {
  const messages = {
    obligation_activated: 'New obligation activated',
    obligation_breached: 'Obligation breached',
    obligation_fulfilled: 'Obligation fulfilled',
    penalty_applied: 'Penalty applied'
  }
  
  notifications.value.unshift({
    id: Date.now(),
    type: data.type || 'default',
    message: messages[data.type] || 'Update received',
    time: new Date().toLocaleTimeString()
  })
  
  if (notifications.value.length > 10) {
    notifications.value.pop()
  }
  
  fetchStats()
}

onMounted(() => {
  fetchStats()
  wsService.connect().catch(console.error)
  wsService.on('obligation_activated', handleWebSocketMessage)
  wsService.on('obligation_breached', handleWebSocketMessage)
  wsService.on('obligation_fulfilled', handleWebSocketMessage)
  wsService.on('penalty_applied', handleWebSocketMessage)
})

onUnmounted(() => {
  wsService.disconnect()
})
</script>

<style scoped>
.dashboard {
  display: flex;
  min-height: 100vh;
}

.sidebar {
  width: 250px;
  background: #2c3e50;
  color: white;
  padding: 1rem;
  display: flex;
  flex-direction: column;
}

.logo h2 {
  padding: 1rem 0;
  border-bottom: 1px solid #34495e;
  margin-bottom: 1rem;
}

.sidebar nav {
  flex: 1;
}

.sidebar nav a {
  display: block;
  padding: 0.75rem 1rem;
  color: #bdc3c7;
  text-decoration: none;
  border-radius: 4px;
  margin-bottom: 0.25rem;
}

.sidebar nav a:hover,
.sidebar nav a.router-link-active {
  background: #34495e;
  color: white;
}

.user-info {
  border-top: 1px solid #34495e;
  padding-top: 1rem;
}

.user-info p {
  margin: 0.25rem 0;
}

.role {
  color: #bdc3c7;
  font-size: 0.875rem;
}

.btn-logout {
  margin-top: 1rem;
  width: 100%;
  padding: 0.5rem;
  background: #e74c3c;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

.main-content {
  flex: 1;
  padding: 2rem;
  background: #f5f6fa;
}

.main-content header {
  margin-bottom: 2rem;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 1.5rem;
  margin-bottom: 2rem;
}

.stat-card {
  background: white;
  padding: 1.5rem;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
}

.stat-card h3 {
  margin: 0 0 0.5rem;
  color: #666;
  font-size: 0.875rem;
  font-weight: 500;
}

.stat-value {
  margin: 0;
  font-size: 2rem;
  font-weight: 700;
  color: #2c3e50;
}

.stat-card.danger .stat-value {
  color: #e74c3c;
}

.recent-activity {
  background: white;
  padding: 1.5rem;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
}

.recent-activity h2 {
  margin: 0 0 1rem;
}

.empty-state {
  color: #999;
  text-align: center;
  padding: 2rem;
}

.activity-list {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.activity-item {
  display: flex;
  align-items: flex-start;
  gap: 1rem;
  padding: 0.75rem;
  border-radius: 4px;
  background: #f8f9fa;
}

.activity-icon {
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  font-size: 0.75rem;
}

.activity-item.obligation_activated .activity-icon {
  background: #27ae60;
  color: white;
}

.activity-item.obligation_breached .activity-icon {
  background: #e74c3c;
  color: white;
}

.activity-item.obligation_fulfilled .activity-icon {
  background: #3498db;
  color: white;
}

.activity-content p {
  margin: 0;
}

.activity-time {
  font-size: 0.75rem;
  color: #999;
}
</style>
