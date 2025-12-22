import { defineStore } from 'pinia'
import { ref } from 'vue'

export interface Source {
  id: string
  name: string
  url?: string
  status?: string
  lastSyncedAt?: string
}

export const useSourceStore = defineStore('sources', () => {
  const sources = ref<Source[]>([])
  const isLoading = ref(false)
  const error = ref<string | null>(null)

  async function fetchSources() {
    isLoading.value = true
    error.value = null
    try {
      const res = await fetch('/api/sources')
      if (!res.ok) {
        throw new Error(`Failed to fetch sources: ${res.statusText}`)
      }
      sources.value = await res.json()
    } catch (e: any) {
      error.value = e.message || 'Unknown error'
      console.error('Failed to fetch sources', e)
    } finally {
      isLoading.value = false
    }
  }

  async function addSource(source: Omit<Source, 'id'>) {
    isLoading.value = true
    error.value = null
    try {
      const res = await fetch('/api/sources', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(source),
      })
      if (!res.ok) {
        throw new Error(`Failed to add source: ${res.statusText}`)
      }
      const newSource = await res.json()
      sources.value.push(newSource)
    } catch (e: any) {
      error.value = e.message || 'Unknown error'
      console.error('Failed to add source', e)
    } finally {
      isLoading.value = false
    }
  }

  return { 
    sources, 
    isLoading, 
    error, 
    fetchSources, 
    addSource 
  }
})