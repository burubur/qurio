import { mount } from '@vue/test-utils'
import { describe, it, expect, vi } from 'vitest'
import { createTestingPinia } from '@pinia/testing'
import SourceForm from './SourceForm.vue'
import { useSourceStore } from './source.store'

describe('SourceForm', () => {
  it('calls addSource on submit with advanced config', async () => {
    const wrapper = mount(SourceForm, {
      global: {
        plugins: [createTestingPinia({ createSpy: vi.fn })],
      },
    })
    const store = useSourceStore()
    
    // Set URL
    const input = wrapper.find('input[type="text"]')
    await input.setValue('https://example.com')

    // Open Advanced
    await wrapper.find('button[type="button"]').trigger('click')
    
    // Set Depth
    const depthInput = wrapper.find('input[type="number"]')
    await depthInput.setValue(2)

    // Set Exclusions
    const textarea = wrapper.find('textarea')
    await textarea.setValue(`/login
/admin`)

    // Submit
    await wrapper.find('form').trigger('submit')
    
    expect(store.addSource).toHaveBeenCalledWith({ 
      name: 'https://example.com', 
      url: 'https://example.com',
      max_depth: 2,
      exclusions: ['/login', '/admin']
    })
  })
})
