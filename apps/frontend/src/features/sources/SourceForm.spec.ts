import { mount } from '@vue/test-utils'
import { describe, it, expect, vi } from 'vitest'
import { createTestingPinia } from '@pinia/testing'
import SourceForm from './SourceForm.vue'
import { useSourceStore } from './source.store'

describe('SourceForm', () => {
  it('calls addSource on submit', async () => {
    const wrapper = mount(SourceForm, {
      global: {
        plugins: [createTestingPinia({ createSpy: vi.fn })],
      },
    })
    const store = useSourceStore()
    
    const input = wrapper.find('input')
    await input.setValue('https://example.com')
    await wrapper.find('form').trigger('submit')
    
    expect(store.addSource).toHaveBeenCalledWith({ name: 'https://example.com', url: 'https://example.com' })
  })
})
