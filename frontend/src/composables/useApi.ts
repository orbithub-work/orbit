/**
 * Vue 3 Composable for API calls
 *
 * 使用方式：
 * ```typescript
 * import { useApi } from '@/composables/useApi';
 *
 * export default {
 *   setup() {
 *     const { loading, error, invoke } = useApi();
 *
 *     const handleIndexFile = async () => {
 *       try {
 *         const result = await invoke('index_file', {
 *           path: '/path/to/file',
 *           project_id: 'proj-1'
 *         });
 *         console.log('File indexed:', result);
 *       } catch (err) {
 *         console.error('Failed to index file:', err);
 *       }
 *     };
 *
 *     return { loading, error, handleIndexFile };
 *   }
 * };
 * ```
 */

import { ref, Ref } from 'vue';
import { apiCall } from '@/services/api';

interface UseApiReturn {
  loading: Ref<boolean>;
  error: Ref<string | null>;
  invoke: <T = any>(command: string, payload?: any) => Promise<T>;
}

export function useApi(): UseApiReturn {
  const loading = ref(false);
  const error = ref<string | null>(null);

  const invoke = async <T = any>(command: string, payload?: any): Promise<T> => {
    loading.value = true;
    error.value = null;

    try {
      const result = await apiCall<T>(command, payload);
      return result;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : String(err);
      error.value = errorMessage;
      throw err;
    } finally {
      loading.value = false;
    }
  };

  return {
    loading,
    error,
    invoke,
  };
}

/**
 * 资产相关的Composable
 */
export function useAssetApi() {
  const { loading, error, invoke } = useApi();

  return {
    loading,
    error,
    async indexFile(path: string, projectId: string) {
      return invoke('index_file', { path, project_id: projectId });
    },
    async getAsset(id: string) {
      return invoke('get_asset', { id });
    },
    async findDuplicates(fingerprint: string) {
      return invoke('find_duplicates', { fingerprint });
    },
    async getChildAssets(parentId: string) {
      return invoke('get_child_assets', { parent_id: parentId });
    },
  };
}

/**
 * 项目相关的Composable
 */
export function useProjectApi() {
  const { loading, error, invoke } = useApi();

  return {
    loading,
    error,
    async createProject(name: string, description?: string) {
      return invoke('create_project', { name, description });
    },
    async getProject(id: string) {
      return invoke('get_project', { id });
    },
    async listProjects() {
      return invoke('list_projects');
    },
    async updateProject(id: string, data: any) {
      return invoke('update_project', { id, ...data });
    },
    async deleteProject(id: string) {
      return invoke('delete_project', { id });
    },
  };
}
