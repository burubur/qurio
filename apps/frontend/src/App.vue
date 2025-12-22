<script setup lang="ts">
import SourceForm from './features/sources/SourceForm.vue';
import SourceList from './features/sources/SourceList.vue';
</script>

<template>
  <div class="app-container">
    <header class="app-header">
      <div class="logo">Qurio</div>
      <nav>
        <!-- Future: Add nav links -->
      </nav>
    </header>

    <main class="main-content">
      <div class="card">
        <h2>Add New Source</h2>
        <p class="description">Enter a URL to ingest content into your knowledge base.</p>
        <SourceForm />
      </div>

      <div class="card">
        <h2>Your Sources</h2>
        <SourceList />
      </div>
    </main>
  </div>
</template>

<style>
/* Global Resets & Basics */
body {
  margin: 0;
  font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, 'Open Sans', 'Helvetica Neue', sans-serif;
  background-color: #f5f7fa;
  color: #2c3e50;
}

* {
  box-sizing: border-box;
}
</style>

<style scoped>
.app-container {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.app-header {
  background-color: #fff;
  padding: 1rem 2rem;
  box-shadow: 0 2px 4px rgba(0,0,0,0.05);
  display: flex;
  align-items: center;
}

.logo {
  font-size: 1.5rem;
  font-weight: 800;
  color: #3498db;
  letter-spacing: -0.5px;
}

.main-content {
  flex: 1;
  max-width: 900px;
  width: 100%;
  margin: 0 auto;
  padding: 2rem;
  display: flex;
  flex-direction: column;
  gap: 2rem;
}

.card {
  background: white;
  border-radius: 8px;
  padding: 2rem;
  box-shadow: 0 4px 6px rgba(0,0,0,0.02);
  border: 1px solid #eee;
}

h2 {
  margin-top: 0;
  margin-bottom: 0.5rem;
  font-size: 1.25rem;
  font-weight: 600;
}

.description {
  color: #7f8c8d;
  margin-bottom: 1.5rem;
  font-size: 0.95rem;
}
</style>
