import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import FileIcon from '@/components/file-manager/FileIcon.vue'

describe('FileIcon', () => {
  describe('renders correctly', () => {
    it('should render image icon', () => {
      const wrapper = mount(FileIcon, {
        props: { fileType: 'image' }
      })
      
      expect(wrapper.find('.file-icon--image').exists()).toBe(true)
      expect(wrapper.find('.icon-svg').exists()).toBe(true)
    })

    it('should render video icon', () => {
      const wrapper = mount(FileIcon, {
        props: { fileType: 'video' }
      })
      
      expect(wrapper.find('.file-icon--video').exists()).toBe(true)
      expect(wrapper.find('.icon-svg').exists()).toBe(true)
    })

    it('should render audio icon', () => {
      const wrapper = mount(FileIcon, {
        props: { fileType: 'audio' }
      })
      
      expect(wrapper.find('.file-icon--audio').exists()).toBe(true)
      expect(wrapper.find('.icon-svg').exists()).toBe(true)
    })

    it('should render document icon', () => {
      const wrapper = mount(FileIcon, {
        props: { fileType: 'document' }
      })
      
      expect(wrapper.find('.file-icon--document').exists()).toBe(true)
      expect(wrapper.find('.icon-svg').exists()).toBe(true)
    })

    it('should render folder emoji for folder type', () => {
      const wrapper = mount(FileIcon, {
        props: { fileType: 'folder' }
      })
      
      expect(wrapper.find('.file-icon--folder').exists()).toBe(true)
      expect(wrapper.find('.icon-svg').exists()).toBe(true)
    })

    it('should render default icon for unknown type', () => {
      const wrapper = mount(FileIcon, {
        props: { fileType: 'unknown' }
      })
      
      expect(wrapper.find('.icon-svg').exists()).toBe(true)
    })
  })

  describe('props', () => {
    it('should accept custom size', () => {
      const wrapper = mount(FileIcon, {
        props: { fileType: 'image', size: 32 }
      })
      
      expect(wrapper.html()).toContain('32px')
    })

    it('should use default size when not provided', () => {
      const wrapper = mount(FileIcon, {
        props: { fileType: 'image' }
      })
      
      expect(wrapper.html()).toContain('24px')
    })
  })
})
