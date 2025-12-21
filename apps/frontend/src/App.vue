<script setup lang="ts">
import SourceForm from './features/sources/SourceForm.vue';
import SourceList from './features/sources/SourceList.vue';
import { ref } from 'vue';

const sources = ref<{id: string, url: string}[]>([]);
const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8081';

async function addSource(url: string) {
  try {
    const res = await fetch(`${API_URL}/sources`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ url })
    });
    
    if (!res.ok) throw new Error('Failed to add source');
    
    const data = await res.json();
    sources.value.push(data);
    alert('Source added and queued for ingestion!');
  } catch (e) {
    console.error(e);
    alert('Error adding source');
  }
}
</script>

<template>
  <div class="container">
    <h1>Qurio Admin</h1>
    <div class="section">
      <h2>Add Source</h2>
      <SourceForm @submit="addSource" />
    </div>
    
    <div class="section">
      <h2>Sources</h2>
      <SourceList :sources="sources" />
    </div>
  </div>
</template>

<style scoped>
.container {
  max-width: 800px;
  margin: 0 auto;
  padding: 2rem;
  font-family: Arial, sans-serif;
}
.section {
  margin-bottom: 2rem;
  padding: 1rem;
  border: 1px solid #ddd;
  border-radius: 8px;
}
h1 { color: #2c3e50; }
</style>