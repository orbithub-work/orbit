import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import Breadcrumb from '@/components/file-manager/Breadcrumb.vue'

describe('Breadcrumb', () => {
  describe('renders correctly', () => {
    it('should render home button', () => {
      const wrapper = mount(Breadcrumb, {
        props: { currentPath: '/' }
      })
      
      expect(wrapper.find('.breadcrumb-item--home').exists()).toBe(true)
    })

    it('should render path segments', () => {
      const wrapper = mount(Breadcrumb, {
        props: { currentPath: '/folder1/folder2/folder3' }
      })
      
      const items = wrapper.findAll('.breadcrumb-item:not(.breadcrumb-item--home)')
      expect(items).toHaveLength(3)
      expect(items[0].text()).toBe('folder1')
      expect(items[1].text()).toBe('folder2')
      expect(items[2].text()).toBe('folder3')
    })

    it('should render separators between items', () => {
      const wrapper = mount(Breadcrumb, {
        props: { currentPath: '/folder1/folder2' }
      })
      
      const separators = wrapper.findAll('.breadcrumb-separator')
      expect(separators).toHaveLength(3)
    })

    it('should handle empty path', () => {
      const wrapper = mount(Breadcrumb, {
        props: { currentPath: '/' }
      })
      
      const items = wrapper.findAll('.breadcrumb-item:not(.breadcrumb-item--home)')
      expect(items).toHaveLength(0)
    })
  })

  describe('events', () => {
    it('should emit navigate with root path when home button clicked', async () => {
      const wrapper = mount(Breadcrumb, {
        props: { currentPath: '/folder1/folder2' }
      })
      
      await wrapper.find('.breadcrumb-item--home').trigger('click')
      
      expect(wrapper.emitted('navigate')).toBeTruthy()
      expect(wrapper.emitted('navigate')![0]).toEqual(['/'])
    })

    it('should emit navigate with correct path when segment clicked', async () => {
      const wrapper = mount(Breadcrumb, {
        props: { currentPath: '/folder1/folder2/folder3' }
      })
      
      const items = wrapper.findAll('.breadcrumb-item:not(.breadcrumb-item--home)')
      await items[1].trigger('click')
      
      expect(wrapper.emitted('navigate')).toBeTruthy()
      expect(wrapper.emitted('navigate')![0]).toEqual(['/folder1/folder2'])
    })

    it('should emit navigate with full path when last segment clicked', async () => {
      const wrapper = mount(Breadcrumb, {
        props: { currentPath: '/folder1/folder2' }
      })
      
      const items = wrapper.findAll('.breadcrumb-item:not(.breadcrumb-item--home)')
      await items[1].trigger('click')
      
      expect(wrapper.emitted('navigate')![0]).toEqual(['/folder1/folder2'])
    })
  })

  describe('path handling', () => {
    it('should handle backslash separators', () => {
      const wrapper = mount(Breadcrumb, {
        props: { currentPath: 'C:\\folder1\\folder2' }
      })
      
      const items = wrapper.findAll('.breadcrumb-item:not(.breadcrumb-item--home)')
      expect(items).toHaveLength(2)
      expect(items[0].text()).toBe('folder1')
      expect(items[1].text()).toBe('folder2')
    })

    it('should handle forward slash separators', () => {
      const wrapper = mount(Breadcrumb, {
        props: { currentPath: '/folder1/folder2' }
      })
      
      const items = wrapper.findAll('.breadcrumb-item:not(.breadcrumb-item--home)')
      expect(items).toHaveLength(2)
      expect(items[0].text()).toBe('folder1')
      expect(items[1].text()).toBe('folder2')
    })

    it('should trim leading and trailing separators', () => {
      const wrapper = mount(Breadcrumb, {
        props: { currentPath: '///folder1/folder2///' }
      })
      
      const items = wrapper.findAll('.breadcrumb-item:not(.breadcrumb-item--home)')
      expect(items).toHaveLength(2)
    })
  })
})