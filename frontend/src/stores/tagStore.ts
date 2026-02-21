// 标签 Store - 标签管理状态
//
// 提供标签的创建、编辑、删除、查询和文件关联功能

import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { apiCall } from '@/services/api'

// 标签项
export interface TagItem {
  id: string
  name: string
  color: string | null
  icon: string | null
  parent_id: string | null
  file_count: number
  created_at: number
}

export interface TagTreeItem extends TagItem {
  children: TagTreeItem[]
}

// 创建标签请求
export interface CreateTagRequest {
  name: string
  color?: string | null
  icon?: string | null
  parent_id?: string | null
}

// 更新标签请求
export interface UpdateTagRequest {
  id: string
  name?: string | null
  color?: string | null
  icon?: string | null
  parent_id?: string | null
}

export const useTagStore = defineStore('tag', () => {
  // 状态
  const tags = ref<TagItem[]>([])
  const tagTree = ref<TagTreeItem[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)

  // 计算属性：获取所有标签（按创建时间降序）
  const sortedTags = computed(() => {
    return [...tags.value].sort((a, b) => b.created_at - a.created_at)
  })

  // 计算属性：获取所有标签（按文件数量降序）
  const tagsByFileCount = computed(() => {
    return [...tags.value].sort((a, b) => b.file_count - a.file_count)
  })

  // 计算属性：获取有文件关联的标签
  const usedTags = computed(() => {
    return tags.value.filter(tag => tag.file_count > 0)
  })

  // 加载所有标签
  async function loadTags() {
    loading.value = true
    error.value = null

    try {
      const result = await apiCall<TagItem[]>('list_tags')
      tags.value = result || []
      try {
        const tree = await apiCall<TagTreeItem[]>('list_tag_tree')
        tagTree.value = tree || []
      } catch {
        tagTree.value = []
      }
    } catch (e) {
      error.value = e instanceof Error ? e.message : '加载标签失败'
      console.error('Failed to load tags:', e)
    } finally {
      loading.value = false
    }
  }

  // 创建标签
  async function createTag(request: CreateTagRequest) {
    loading.value = true
    error.value = null

    try {
      const newTag = await apiCall<TagItem>('create_tag', request)
      tags.value.push(newTag)
      await loadTags()
      return newTag
    } catch (e) {
      error.value = e instanceof Error ? e.message : '创建标签失败'
      console.error('Failed to create tag:', e)
      throw e
    } finally {
      loading.value = false
    }
  }

  // 更新标签
  async function updateTag(request: UpdateTagRequest) {
    loading.value = true
    error.value = null

    try {
      const updatedTag = await apiCall<TagItem>('update_tag', request)

      // 更新本地状态
      const index = tags.value.findIndex(tag => tag.id === request.id)
      if (index !== -1) {
        tags.value[index] = updatedTag
      }
      await loadTags()

      return updatedTag
    } catch (e) {
      error.value = e instanceof Error ? e.message : '更新标签失败'
      console.error('Failed to update tag:', e)
      throw e
    } finally {
      loading.value = false
    }
  }

  // 删除标签
  async function deleteTag(tagId: string) {
    loading.value = true
    error.value = null

    try {
      await apiCall('delete_tag', { tagId })

      // 从本地状态中移除
      tags.value = tags.value.filter(tag => tag.id !== tagId)
    } catch (e) {
      error.value = e instanceof Error ? e.message : '删除标签失败'
      console.error('Failed to delete tag:', e)
      throw e
    } finally {
      loading.value = false
    }
  }

  // 获取文件的标签
  async function getFileTags(fileId: string) {
    loading.value = true
    error.value = null

    try {
      const fileTags = await apiCall<TagItem[]>('get_file_tags', { fileId })
      return fileTags || []
    } catch (e) {
      error.value = e instanceof Error ? e.message : '获取文件标签失败'
      console.error('Failed to get file tags:', e)
      return []
    } finally {
      loading.value = false
    }
  }

  // 批量添加标签到文件
  async function addTagsToFiles(fileIds: string[], tagIds: string[]) {
    loading.value = true
    error.value = null

    try {
      await apiCall('add_tags_to_files', {
        file_ids: fileIds,
        tag_ids: tagIds
      })
    } catch (e) {
      error.value = e instanceof Error ? e.message : '添加标签到文件失败'
      console.error('Failed to add tags to files:', e)
      throw e
    } finally {
      loading.value = false
    }
  }

  // 批量从文件移除标签
  async function removeTagsFromFiles(fileIds: string[], tagIds: string[]) {
    loading.value = true
    error.value = null

    try {
      await apiCall('remove_tags_from_files', {
        file_ids: fileIds,
        tag_ids: tagIds
      })
    } catch (e) {
      error.value = e instanceof Error ? e.message : '从文件移除标签失败'
      console.error('Failed to remove tags from files:', e)
      throw e
    } finally {
      loading.value = false
    }
  }

  // 搜索标签
  async function searchTags(query: string, limit?: number) {
    loading.value = true
    error.value = null

    try {
      const result = await apiCall<TagItem[]>('search_tags', {
        query,
        limit: limit ?? 20
      })
      return result || []
    } catch (e) {
      error.value = e instanceof Error ? e.message : '搜索标签失败'
      console.error('Failed to search tags:', e)
      return []
    } finally {
      loading.value = false
    }
  }

  // 根据ID获取标签
  function getTagById(tagId: string) {
    return tags.value.find(tag => tag.id === tagId)
  }

  // 根据名称获取标签
  function getTagByName(name: string) {
    return tags.value.find(tag => tag.name === name)
  }

  // 获取子标签
  function getChildTags(parentId: string) {
    return tags.value.filter(tag => tag.parent_id === parentId)
  }

  // 清除错误信息
  function clearError() {
    error.value = null
  }

  // 重置状态
  function reset() {
    tags.value = []
    loading.value = false
    error.value = null
  }

  return {
    tags,
    tagTree,
    loading,
    error,
    sortedTags,
    tagsByFileCount,
    usedTags,
    loadTags,
    createTag,
    updateTag,
    deleteTag,
    getFileTags,
    addTagsToFiles,
    removeTagsFromFiles,
    searchTags,
    getTagById,
    getTagByName,
    getChildTags,
    clearError,
    reset
  }
})
