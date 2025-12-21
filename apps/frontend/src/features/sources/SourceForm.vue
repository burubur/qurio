<template>
  <form @submit.prevent="submit">
    <input v-model="url" type="text" placeholder="Enter URL" />
    <button type="submit">Add</button>
  </form>
</template>
<script setup lang="ts">
import { ref } from 'vue'
const url = ref('')
const emit = defineEmits(['submit'])
function submit() {
  emit('submit', url.value)
}
</script>
