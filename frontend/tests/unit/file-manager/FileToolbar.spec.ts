import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import FileToolbar from '@/components/file-manager/FileToolbar.vue'

describe('FileToolbar', () => {
  describe('renders correctly', () => {
    it('should render all toolbar buttons', () => {
      const wrapper = mount(FileToolbar, {
        props: {
          canGoBack: false,
          canGoForward: false,
          currentViewMode: 'list'
        }
      })
      
      expect(wrapper.find('.toolbar-btn').exists()).toBe(true)
      expect(wrapper.find('.search-box').exists()).toBe(true)
      expect(wrapper.find('.view-toggle').exists()).toBe(true)
    })

    it('should render search input', () => {
      const wrapper = mount(FileToolbar)
      
      expect(wrapper.find('.search-input').exists()).toBe(true)
    })

    it('should render filter dropdown', () => {
      const wrapper = mount(FileToolbar)
      
      expect(wrapper.find('.filter-select').exists()).toBe(true)
    })

    it('should render view toggle buttons', () => {
      const wrapper = mount(FileToolbar)
      
      const viewBtns = wrapper.findAll('.view-btn')
      expect(viewBtns).toHaveLength(2)
    })
  })

  describe('navigation buttons', () => {
    it('should disable back button when cannot go back', () => {
      const wrapper = mount(FileToolbar, {
        props: { canGoBack: false }
      })
      
      const backBtn = wrapper.findAll('.toolbar-btn')[0]
      expect(backBtn.attributes('disabled')).toBeDefined()
    })

    it('should enable back button when can go back', () => {
      const wrapper = mount(FileToolbar, {
        props: { canGoBack: true }
      })
      
      const backBtn = wrapper.findAll('.toolbar-btn')[0]
      expect(backBtn.attributes('disabled')).toBeUndefined()
    })

    it('should disable forward button when cannot go forward', () => {
      const wrapper = mount(FileToolbar, {
        props: { canGoForward: false }
      })
      
      const forwardBtn = wrapper.findAll('.toolbar-btn')[1]
      expect(forwardBtn.attributes('disabled')).toBeDefined()
    })

    it('should emit go-back event', async () => {
      const wrapper = mount(FileToolbar, {
        props: { canGoBack: true }
      })
      
      await wrapper.findAll('.toolbar-btn')[0].trigger('click')
      
      expect(wrapper.emitted('go-back')).toBeTruthy()
    })

    it('should emit go-forward event', async () => {
      const wrapper = mount(FileToolbar, {
        props: { canGoForward: true }
      })
      
      await wrapper.findAll('.toolbar-btn')[1].trigger('click')
      
      expect(wrapper.emitted('go-forward')).toBeTruthy()
    })

    it('should emit refresh event', async () => {
      const wrapper = mount(FileToolbar)
      
      await wrapper.findAll('.toolbar-btn')[2].trigger('click')
      
      expect(wrapper.emitted('refresh')).toBeTruthy()
    })
  })

  describe('search', () => {
    it('should update search query', async () => {
      const wrapper = mount(FileToolbar)
      
      const input = wrapper.find('.search-input')
      await input.setValue('test query')
      
      expect(input.element.value).toBe('test query')
    })

    it('should show clear button when search has value', async () => {
      const wrapper = mount(FileToolbar)
      
      expect(wrapper.find('.search-clear').exists()).toBe(false)
      
      const input = wrapper.find('.search-input')
      await input.setValue('test')
      await input.trigger('input')
      
      // Wait for debounce
      await new Promise(resolve => setTimeout(resolve, 350))
      
      expect(wrapper.find('.search-clear').exists()).toBe(true)
    })

    it('should clear search when clear button clicked', async () => {
      const wrapper = mount(FileToolbar)
      
      const input = wrapper.find('.search-input')
      await input.setValue('test')
      await input.trigger('input')
      
      await new Promise(resolve => setTimeout(resolve, 350))
      
      await wrapper.find('.search-clear').trigger('click')
      
      expect(wrapper.emitted('search')).toBeTruthy()
      expect(wrapper.emitted('search')![0]).toEqual([''])
    })
  })

  describe('view toggle', () => {
    it('should highlight active view mode', () => {
      const wrapper = mount(FileToolbar, {
        props: { currentViewMode: 'list' }
      })
      
      const viewBtns = wrapper.findAll('.view-btn')
      expect(viewBtns[0].classes()).toContain('view-btn--active')
      expect(viewBtns[1].classes()).not.toContain('view-btn--active')
    })

    it('should emit view-change event when clicking view button', async () => {
      const wrapper = mount(FileToolbar)
      
      const viewBtns = wrapper.findAll('.view-btn')
      await viewBtns[1].trigger('click')
      
      expect(wrapper.emitted('view-change')).toBeTruthy()
      expect(wrapper.emitted('view-change')![0]).toEqual(['grid'])
    })
  })

  describe('filter', () => {
    it('should have default filter option', () => {
      const wrapper = mount(FileToolbar)
      
      const select = wrapper.find('.filter-select')
      expect(select.element.value).toBe('')
    })

    it('should emit filter-change event when filter changes', async () => {
      const wrapper = mount(FileToolbar)
      
      const select = wrapper.find('.filter-select')
      await select.setValue('image')
      
      expect(wrapper.emitted('filter-change')).toBeTruthy()
      expect(wrapper.emitted('filter-change')![0]).toEqual(['image'])
    })
  })
})