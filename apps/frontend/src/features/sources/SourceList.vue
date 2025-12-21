<template>
  <div v-for="source in sources" :key="source.id">
    {{ source.url }}
  </div>
</template>
<script setup lang="ts">
defineProps<{ sources: { id: string, url: string }[] }>()
</script>