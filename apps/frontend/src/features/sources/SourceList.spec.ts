import { mount } from '@vue/test-utils'
import { describe, it, expect } from 'vitest'
import SourceList from './SourceList.vue'

describe('SourceList', () => {
  it('displays sources', () => {
    const sources = [{ id: '1', url: 'https://example.com' }]
    const wrapper = mount(SourceList, { props: { sources } })
    expect(wrapper.text()).toContain('https://example.com')
  })
})
