import { describe, it, expect, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import FileListView from '@/components/FileListView.vue'

describe('FileListView.vue', () => {
  const mockFiles = [
    {
      id: '1',
      name: 'test-file1.jpg',
      path: '/path/to/test-file1.jpg',
      size: 1024,
      type: 'file',
      mimeType: 'image/jpeg',
      thumbnail: undefined,
      createdAt: new Date('2023-01-01'),
      modifiedAt: new Date('2023-01-02'),
      isSelected: false
    },
    {
      id: '2',
      name: 'test-file2.mp4',
      path: '/path/to/test-file2.mp4',
      size: 2048,
      type: 'file',
      mimeType: 'video/mp4',
      thumbnail: undefined,
      createdAt: new Date('2023-01-03'),
      modifiedAt: new Date('2023-01-04'),
      isSelected: false
    },
    {
      id: '3',
      name: 'test-folder',
      path: '/path/to/test-folder',
      size: 0,
      type: 'directory',
      mimeType: undefined,
      thumbnail: undefined,
      createdAt: new Date('2023-01-05'),
      modifiedAt: new Date('2023-01-06'),
      isSelected: false
    }
  ]

  it('renders properly with files', () => {
    const wrapper = mount(FileListView, {
      props: {
        files: mockFiles
      }
    })

    expect(wrapper.findAll('.list-row')).toHaveLength(3)
    expect(wrapper.text()).toContain('test-file1.jpg')
    expect(wrapper.text()).toContain('test-file2.mp4')
    expect(wrapper.text()).toContain('test-folder')
  })

  it('switches view modes correctly', async () => {
    const wrapper = mount(FileListView, {
      props: {
        files: mockFiles
      }
    })

    // Initially should be in list view
    expect(wrapper.find('.list-view').exists()).toBe(true)
    expect(wrapper.find('.grid-view').exists()).toBe(false)
    expect(wrapper.find('.thumbnail-view').exists()).toBe(false)

    // Switch to grid view
    const gridBtn = wrapper.findAll('.view-mode-btn')[1]
    await gridBtn.trigger('click')
    expect(wrapper.find('.grid-view').exists()).toBe(true)
    expect(wrapper.find('.list-view').exists()).toBe(false)

    // Switch to thumbnail view
    const thumbnailBtn = wrapper.findAll('.view-mode-btn')[2]
    await thumbnailBtn.trigger('click')
    expect(wrapper.find('.thumbnail-view').exists()).toBe(true)
    expect(wrapper.find('.grid-view').exists()).toBe(false)
  })

  it('sorts files by name', async () => {
    const wrapper = mount(FileListView, {
      props: {
        files: mockFiles
      }
    })

    // Initially sorted by modified date (default)
    const initialOrder = wrapper.findAll('.list-row .file-name').map(el => el.text())
    expect(initialOrder).toEqual(['test-folder', 'test-file1.jpg', 'test-file2.mp4'])

    // Change sort field to name
    const sortSelect = wrapper.find('.sort-field')
    await sortSelect.setValue('name')

    // Check that the order has changed
    const sortedOrder = wrapper.findAll('.list-row .file-name').map(el => el.text())
    expect(sortedOrder).toEqual(['test-file1.jpg', 'test-file2.mp4', 'test-folder'])
  })

  it('toggles sort direction', async () => {
    const wrapper = mount(FileListView, {
      props: {
        files: mockFiles
      }
    })

    // Check initial sort direction (should be descending by modified date)
    const initialOrder = wrapper.findAll('.list-row .file-name').map(el => el.text())

    // Toggle sort direction
    const sortDirectionBtn = wrapper.find('.sort-direction-btn')
    await sortDirectionBtn.trigger('click')

    // Check that the order has reversed
    const reversedOrder = wrapper.findAll('.list-row .file-name').map(el => el.text())
    expect(reversedOrder).toEqual(initialOrder.reverse())
  })

  it('formats file sizes correctly', () => {
    const wrapper = mount(FileListView, {
      props: {
        files: mockFiles
      }
    })

    const vm = wrapper.vm as any
    expect(vm.formatFileSize(0)).toBe('0 Bytes')
    expect(vm.formatFileSize(512)).toBe('512.00 Bytes')
    expect(vm.formatFileSize(1024)).toBe('1.00 KB')
    expect(vm.formatFileSize(1024 * 1024)).toBe('1.00 MB')
    expect(vm.formatFileSize(1024 * 1024 * 1024)).toBe('1.00 GB')
  })

  it('displays correct file icons', () => {
    const wrapper = mount(FileListView, {
      props: {
        files: mockFiles
      }
    })

    const vm = wrapper.vm as any
    const imageFile = mockFiles[0]
    const videoFile = mockFiles[1]
    const folder = mockFiles[2]

    expect(vm.getFileIcon(imageFile)).toBe('🖼️')
    expect(vm.getFileIcon(videoFile)).toBe('🎬')
    expect(vm.getFileIcon(folder)).toBe('📁')
  })
})