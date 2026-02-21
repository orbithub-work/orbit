import { defineStore } from 'pinia'
import { ref } from 'vue'
import { apiCall } from '@/services/api'

export interface Collection {
  id: string
  name: string
  description?: string
  itemCount: number
  createdAt: string
  updatedAt: string
}

export interface CollectionItem {
  collectionId: string
  fileId: string
  addedAt: string
}

export const useCollectionStore = defineStore('collection', () => {
  const collections = ref<Collection[]>([])
  const currentCollection = ref<Collection | null>(null)
  const collectionFiles = ref<any[]>([])
  const loading = ref(false)

  async function loadCollections() {
    loading.value = true
    try {
      const result = await apiCall('list_collections') as Collection[]
      collections.value = result
    } catch (error) {
      console.error('Failed to load collections:', error)
      collections.value = []
    } finally {
      loading.value = false
    }
  }

  async function createCollection(name: string, description?: string) {
    loading.value = true
    try {
      const result = await apiCall('create_collection', {
        request: { name, description }
      }) as Collection
      collections.value.push(result)
      return result
    } catch (error) {
      console.error('Failed to create collection:', error)
      throw error
    } finally {
      loading.value = false
    }
  }

  async function deleteCollection(collectionId: string) {
    loading.value = true
    try {
      await apiCall('delete_collection', { collectionId })
      collections.value = collections.value.filter(c => c.id !== collectionId)
      if (currentCollection.value?.id === collectionId) {
        currentCollection.value = null
        collectionFiles.value = []
      }
    } catch (error) {
      console.error('Failed to delete collection:', error)
      throw error
    } finally {
      loading.value = false
    }
  }

  async function renameCollection(collectionId: string, newName: string) {
    loading.value = true
    try {
      await apiCall('rename_collection', { collectionId, newName })
      const collection = collections.value.find(c => c.id === collectionId)
      if (collection) {
        collection.name = newName
        collection.updatedAt = new Date().toISOString()
      }
    } catch (error) {
      console.error('Failed to rename collection:', error)
      throw error
    } finally {
      loading.value = false
    }
  }

  async function addFileToCollection(collectionId: string, fileId: string) {
    loading.value = true
    try {
      await apiCall('add_file_to_collection', { collectionId, fileId })
      const collection = collections.value.find(c => c.id === collectionId)
      if (collection) {
        collection.itemCount++
        collection.updatedAt = new Date().toISOString()
      }
      if (currentCollection.value?.id === collectionId) {
        await loadCollectionFiles(collectionId)
      }
    } catch (error) {
      console.error('Failed to add file to collection:', error)
      throw error
    } finally {
      loading.value = false
    }
  }

  async function removeFileFromCollection(collectionId: string, fileId: string) {
    loading.value = true
    try {
      await apiCall('remove_file_from_collection', { collectionId, fileId })
      const collection = collections.value.find(c => c.id === collectionId)
      if (collection) {
        collection.itemCount = Math.max(0, collection.itemCount - 1)
        collection.updatedAt = new Date().toISOString()
      }
      if (currentCollection.value?.id === collectionId) {
        await loadCollectionFiles(collectionId)
      }
    } catch (error) {
      console.error('Failed to remove file from collection:', error)
      throw error
    } finally {
      loading.value = false
    }
  }

  async function loadCollectionFiles(collectionId: string) {
    loading.value = true
    try {
      const result = await apiCall('get_collection_files', { collectionId }) as any[]
      collectionFiles.value = result
    } catch (error) {
      console.error('Failed to load collection files:', error)
      collectionFiles.value = []
    } finally {
      loading.value = false
    }
  }

  function setCurrentCollection(collection: Collection | null) {
    currentCollection.value = collection
    if (collection) {
      loadCollectionFiles(collection.id)
    } else {
      collectionFiles.value = []
    }
  }

  return {
    collections,
    currentCollection,
    collectionFiles,
    loading,
    loadCollections,
    createCollection,
    deleteCollection,
    renameCollection,
    addFileToCollection,
    removeFileFromCollection,
    loadCollectionFiles,
    setCurrentCollection
  }
})
