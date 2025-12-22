<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  status: string
}>()

const statusColor = computed(() => {
  switch (props.status.toLowerCase()) {
    case 'indexed':
    case 'completed':
      return 'var(--color-success)'
    case 'processing':
    case 'pending':
      return 'var(--color-primary)'
    case 'failed':
      return 'var(--color-danger)'
    default:
      return 'var(--color-text-muted)'
  }
})
</script>

<template>
  <div class="status-badge" :style="{ color: statusColor, borderColor: statusColor }">
    <div class="dot" :style="{ backgroundColor: statusColor }"></div>
    <span class="label">{{ status }}</span>
  </div>
</template>

<style scoped>
.status-badge {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.25rem 0.75rem;
  border-radius: 12px;
  border: 1px solid currentColor;
  background-color: rgba(255, 255, 255, 0.03);
  font-family: var(--font-mono);
  font-size: 0.75rem;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  font-weight: 500;
}

.dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  box-shadow: 0 0 8px currentColor;
}
</style>
