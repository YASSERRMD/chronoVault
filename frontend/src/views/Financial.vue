<template>
  <div class="financial-page">
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
        <h1>Financial Impact</h1>
      </header>
      
      <div class="summary-cards">
        <div class="summary-card">
          <h3>Total Contract Value</h3>
          <p class="value">${{ formatNumber(summary.total_contract_value) }}</p>
        </div>
        <div class="summary-card">
          <h3>Total Penalties</h3>
          <p class="value danger">${{ formatNumber(summary.total_penalties) }}</p>
        </div>
        <div class="summary-card">
          <h3>Active Obligations</h3>
          <p class="value">{{ summary.active_obligations || 0 }}</p>
        </div>
      </div>
      
      <div class="reports-grid">
        <div class="report-card">
          <h2>Risk Exposure</h2>
          <div v-if="riskExposure.length === 0" class="empty-state">No data</div>
          <table v-else>
            <thead>
              <tr>
                <th>Contract</th>
                <th>Value</th>
                <th>Potential Penalty</th>
                <th>Risk Level</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="r in riskExposure" :key="r.contract_id">
                <td>{{ r.contract_title }}</td>
                <td>${{ formatNumber(r.contract_value) }}</td>
                <td>${{ formatNumber(r.potential_penalty) }}</td>
                <td><span class="risk-badge" :class="r.risk_level">{{ r.risk_level }}</span></td>
              </tr>
            </tbody>
          </table>
        </div>
        
        <div class="report-card">
          <h2>Penalty Tracking</h2>
          <div v-if="penalties.length === 0" class="empty-state">No data</div>
          <table v-else>
            <thead>
              <tr>
                <th>Contract</th>
                <th>Total Penalty</th>
                <th>Breached</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="p in penalties" :key="p.contract_id">
                <td>{{ p.contract_title }}</td>
                <td>${{ formatNumber(p.total_penalty) }}</td>
                <td>{{ p.breached_count }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
      
      <div class="report-card yearly">
        <h2>Yearly Impact</h2>
        <div class="year-selector">
          <button v-for="year in ['2024', '2025', '2026']" :key="year" @click="selectedYear = year; fetchYearlyImpact()" :class="{ active: selectedYear === year }">{{ year }}</button>
        </div>
        <div v-if="yearlyImpact.length === 0" class="empty-state">No data for {{ selectedYear }}</div>
        <div v-else class="chart-bars">
          <div v-for="item in yearlyImpact" :key="item.month" class="bar-container">
            <div class="bar" :style="{ height: getBarHeight(item.penalties) + '%' }"></div>
            <span class="bar-label">{{ item.month }}</span>
          </div>
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

const summary = ref({})
const riskExposure = ref([])
const penalties = ref([])
const yearlyImpact = ref([])
const selectedYear = ref('2024')

const fetchSummary = async () => {
  try {
    const response = await api.get('/reports/financial-summary')
    summary.value = response.data
  } catch (err) { console.error(err) }
}

const fetchRiskExposure = async () => {
  try {
    const response = await api.get('/reports/risk-exposure')
    riskExposure.value = response.data
  } catch (err) { console.error(err) }
}

const fetchPenaltyTracking = async () => {
  try {
    const response = await api.get('/reports/penalty-tracking')
    penalties.value = response.data
  } catch (err) { console.error(err) }
}

const fetchYearlyImpact = async () => {
  try {
    const response = await api.get(`/reports/yearly-impact?year=${selectedYear.value}`)
    yearlyImpact.value = response.data
  } catch (err) { console.error(err) }
}

const formatNumber = (num) => num ? new Intl.NumberFormat().format(num) : '0'

const getBarHeight = (penalties) => {
  const max = Math.max(...yearlyImpact.value.map(y => parseFloat(y.penalties) || 0), 1)
  return ((parseFloat(penalties) || 0) / max) * 100
}

const handleLogout = () => {
  authStore.logout()
  router.push('/login')
}

onMounted(() => {
  fetchSummary()
  fetchRiskExposure()
  fetchPenaltyTracking()
  fetchYearlyImpact()
})
</script>

<style scoped>
.financial-page { display: flex; min-height: 100vh; }
.sidebar { width: 250px; background: #2c3e50; color: white; padding: 1rem; display: flex; flex-direction: column; }
.logo h2 { padding: 1rem 0; border-bottom: 1px solid #34495e; margin-bottom: 1rem; }
.sidebar nav { flex: 1; }
.sidebar nav a { display: block; padding: 0.75rem 1rem; color: #bdc3c7; text-decoration: none; border-radius: 4px; margin-bottom: 0.25rem; }
.sidebar nav a:hover, .sidebar nav a.router-link-active { background: #34495e; color: white; }
.btn-logout { width: 100%; padding: 0.5rem; background: #e74c3c; color: white; border: none; border-radius: 4px; cursor: pointer; }
.main-content { flex: 1; padding: 2rem; background: #f5f6fa; }
.summary-cards { display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: 1.5rem; margin-bottom: 2rem; }
.summary-card { background: white; padding: 1.5rem; border-radius: 8px; box-shadow: 0 2px 8px rgba(0,0,0,0.1); }
.summary-card h3 { margin: 0 0 0.5rem; color: #666; font-size: 0.875rem; }
.summary-card .value { margin: 0; font-size: 2rem; font-weight: 700; color: #2c3e50; }
.summary-card .value.danger { color: #e74c3c; }
.reports-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(400px, 1fr)); gap: 1.5rem; margin-bottom: 2rem; }
.report-card { background: white; padding: 1.5rem; border-radius: 8px; box-shadow: 0 2px 8px rgba(0,0,0,0.1); }
.report-card h2 { margin: 0 0 1rem; font-size: 1.25rem; }
table { width: 100%; border-collapse: collapse; }
th, td { padding: 0.75rem; text-align: left; border-bottom: 1px solid #eee; }
th { font-weight: 600; color: #666; font-size: 0.875rem; }
.risk-badge { padding: 0.25rem 0.5rem; border-radius: 4px; font-size: 0.75rem; font-weight: 600; }
.risk-badge.high { background: #fadbd8; color: #e74c3c; }
.risk-badge.medium { background: #fdebd0; color: #d68910; }
.risk-badge.low { background: #d5f5e3; color: #27ae60; }
.yearly { margin-top: 0; }
.year-selector { display: flex; gap: 0.5rem; margin-bottom: 1.5rem; }
.year-selector button { padding: 0.5rem 1rem; border: 1px solid #ddd; background: white; border-radius: 4px; cursor: pointer; }
.year-selector button.active { background: #667eea; color: white; border-color: #667eea; }
.chart-bars { display: flex; align-items: flex-end; gap: 0.5rem; height: 200px; padding-top: 1rem; }
.bar-container { flex: 1; display: flex; flex-direction: column; align-items: center; height: 100%; }
.bar { width: 100%; background: #667eea; border-radius: 4px 4px 0 0; min-height: 4px; }
.bar-label { font-size: 0.75rem; color: #666; margin-top: 0.5rem; }
.empty-state { padding: 2rem; text-align: center; color: #999; }
</style>
