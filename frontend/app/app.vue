<script setup lang="ts">
import { ref, onMounted } from 'vue'

const data = ref<any>(null)
const error = ref<string | null>(null)

onMounted(async () => {
  try {
    const response = await fetch('http://localhost:8080/health')
    if (!response.ok) throw new Error('Network response was not ok')
    
    data.value = await response.json()
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'An error occurred'
    console.error('Fetch error:', err)
  }
})
</script>

<template>
  <div>
    <h1>Backend Response:</h1>
    <pre v-if="data">{{ data }}</pre>
    <p v-else-if="error" style="color: red;">Error: {{ error }}</p>
    <p v-else>Loading...</p>
  </div>
</template>