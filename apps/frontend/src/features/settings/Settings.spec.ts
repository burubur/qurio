import { mount } from '@vue/test-utils'
import { createTestingPinia } from '@pinia/testing'
import { describe, it, expect, vi } from 'vitest'
import Settings from './Settings.vue'
import { useSettingsStore } from './settings.store'

describe('Settings.vue', () => {
  it('fetches settings on mount', () => {
    const wrapper = mount(Settings, {
      global: {
        plugins: [createTestingPinia({ createSpy: vi.fn })],
      },
    })
    const store = useSettingsStore()
    expect(store.fetchSettings).toHaveBeenCalled()
  })

  it('calls updateSettings when save button is clicked', async () => {
    const wrapper = mount(Settings, {
      global: {
        plugins: [createTestingPinia({ createSpy: vi.fn })],
      },
    })
    const store = useSettingsStore()
    
    await wrapper.find('button').trigger('click')
    expect(store.updateSettings).toHaveBeenCalled()
  })
})
