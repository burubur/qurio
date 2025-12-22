<script setup lang="ts">
import Sidebar from './Sidebar.vue'
</script>

<template>
  <div class="app-layout">
    <Sidebar />
    <main class="main-content">
      <slot />
    </main>
  </div>
</template>

<style scoped>
.app-layout {
  display: flex;
  min-height: 100vh;
}

.main-content {
  margin-left: 260px; /* Width of sidebar */
  flex: 1;
  padding: 3rem;
  background-color: #0b1120; /* Slightly darker than void for content contrast, or keep void */
  background: radial-gradient(circle at top right, #1e293b 0%, #0f172a 40%);
}
</style>
