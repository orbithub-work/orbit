import { createStandaloneToast } from '@chakra-ui/toast'

const { toast } = createStandaloneToast()

export function showNotification(
  title: string,
  status: 'success' | 'error' | 'warning' | 'info' = 'info',
  description?: string
) {
  toast({
    title,
    description,
    status,
    duration: 5000,
    isClosable: true,
    position: 'top-right'
  })
}

export function showSuccessNotification(title: string, description?: string) {
  showNotification(title, 'success', description)
}

export function showErrorNotification(title: string, description?: string) {
  showNotification(title, 'error', description)
}

export function showWarningNotification(title: string, description?: string) {
  showNotification(title, 'warning', description)
}

export function showInfoNotification(title: string, description?: string) {
  showNotification(title, 'info', description)
}
