const { contextBridge, ipcRenderer } = require('electron')

const coreBaseUrl = process.env.MA_CORE_BASE_URL || 'http://localhost:32000'

const sendRendererError = (payload) => {
  ipcRenderer.send('renderer-error', payload)
}

window.addEventListener('error', (event) => {
  sendRendererError({
    type: 'error',
    message: event.message,
    filename: event.filename,
    lineno: event.lineno,
    colno: event.colno,
    stack: event.error?.stack || null
  })
})

window.addEventListener('unhandledrejection', (event) => {
  sendRendererError({
    type: 'unhandledrejection',
    message: event.reason?.message || String(event.reason),
    stack: event.reason?.stack || null
  })
})

contextBridge.exposeInMainWorld('mediaAssistant', {
  coreBaseUrl,
  platform: process.platform,
  sendRendererError,
  window: {
     minimize: () => ipcRenderer.send('window-minimize'),
     maximize: () => ipcRenderer.send('window-maximize'),
     close: () => ipcRenderer.send('window-close')
   },
   onboarding: {
     complete: () => ipcRenderer.send('onboarding-complete')
   },
   tray: {
     sendAction: (action, payload) => ipcRenderer.send('tray-action', { action, payload })
   },
   on: (channel, callback) => {
     const validChannels = ['open-settings']
     if (validChannels.includes(channel)) {
       const subscription = (_event, ...args) => callback(...args)
       ipcRenderer.on(channel, subscription)
       return () => {
         ipcRenderer.removeListener(channel, subscription)
       }
     }
   }
 })
