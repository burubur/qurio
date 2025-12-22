<script setup lang="ts">
import { ref } from 'vue'
import { useSourceStore } from './source.store'
import { Plus, Loader2 } from 'lucide-vue-next'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'

const store = useSourceStore()
const url = ref('')
const emit = defineEmits(['submit'])

async function submit() {
  if (!url.value) return
  
  // Basic URL validation
  try {
    new URL(url.value)
  } catch {
    alert('Please enter a valid URL (e.g., https://docs.example.com)')
    return
  }

  await store.addSource({ name: url.value, url: url.value })
  if (!store.error) {
    url.value = ''
    emit('submit')
  }
}
</script>

<template>
  <form @submit.prevent="submit" class="w-full space-y-4">
    <div class="flex w-full items-center space-x-2">
      <Input 
        v-model="url" 
        type="text" 
        placeholder="https://docs.example.com" 
        :disabled="store.isLoading" 
      />
      <Button type="submit" :disabled="store.isLoading">
        <Loader2 v-if="store.isLoading" class="mr-2 h-4 w-4 animate-spin" />
        <Plus v-else class="mr-2 h-4 w-4" />
        Ingest
      </Button>
    </div>
    <div v-if="store.error" class="text-destructive text-sm font-mono flex items-center gap-2">
      <span>!</span>
      {{ store.error }}
    </div>
  </form>
</template>
