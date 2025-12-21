import { mount } from '@vue/test-utils'
import { describe, it, expect } from 'vitest'
import SourceForm from './SourceForm.vue'

describe('SourceForm', () => {
  it('emits submit event with url', async () => {
    const wrapper = mount(SourceForm)
    const input = wrapper.find('input')
    await input.setValue('https://example.com')
    await wrapper.find('form').trigger('submit')
    
    expect(wrapper.emitted('submit')?.[0]).toEqual(['https://example.com'])
  })
})
