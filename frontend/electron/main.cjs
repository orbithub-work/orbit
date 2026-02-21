const { app, BrowserWindow, ipcMain, screen, Tray, Menu, nativeImage, protocol } = require('electron')
const { spawn } = require('child_process')
const path = require('path')
const fs = require('fs')

// Set userData to a local directory for development to avoid permission issues
// This must be done BEFORE app is ready
const isDev = process.env.NODE_ENV === 'development' || !app.isPackaged
if (isDev) {
  // Use a temporary user data directory in the project root to avoid macOS permission issues
  const userDataPath = path.join(__dirname, '..', '..', '.user-data-dev')
  if (!fs.existsSync(userDataPath)) {
    try {
      fs.mkdirSync(userDataPath, { recursive: true })
    } catch (e) {
      console.error('Failed to create user data dir:', e)
    }
  }
  app.setPath('userData', userDataPath)
}
const preloadPath = path.join(__dirname, 'preload.cjs')
const distIndex = path.join(__dirname, '..', 'dist', 'index.html')
const safeModeHtml = path.join(__dirname, 'safe_mode.html')
const coreBaseUrl = process.env.MA_CORE_BASE_URL || 'http://localhost:32000'
const coreHealthUrl = `${coreBaseUrl}/api/v1/health`
const devServerUrl = process.env.MA_DEV_SERVER_URL || 'http://localhost:5176'
const startRoute = normalizeRoute(process.env.MA_START_ROUTE || '/')
const hiddenStart = process.env.MA_ELECTRON_HIDDEN !== '0'
const shouldOpenDevTools = process.env.MA_ELECTRON_DEVTOOLS === '1'

let rendererBase = null
let rendererMode = null
let uiReadyState = null
let devWatcher = null

let mainWindow = null
let onboardingWindow = null
let trayWindow = null
let tray = null
let isQuitting = false
let trayWindowShownAt = 0
let trayWindowLastHiddenAt = 0
let trayHideTimer = null
let coreProcess = null
let coreManagedByElectron = false
let coreShutdownRequested = false

const writeLog = (type, payload = {}) => {
  try {
    const logPath = path.join(app.getPath('userData'), 'electron-debug.log')
    const timestamp = new Date().toISOString()
    const logEntry = `[${timestamp}] ${type} ${JSON.stringify(payload)}\n`
    fs.appendFileSync(logPath, logEntry)
  } catch (e) {
    console.error('Failed to write log:', e)
  }
}

console.log('[electron] boot')
writeLog('boot')

// 规范化路由路径
function normalizeRoute(route) {
  if (!route) return '/'
  return route.startsWith('/') ? route : `/${route}`
}

// 构建带哈希的路由URL
function buildHashUrl(base, route) {
  return `${base}#${normalizeRoute(route)}`
}

// 休眠指定毫秒数
function sleep(ms) {
  return new Promise((resolve) => setTimeout(resolve, ms))
}

// 为Promise添加超时控制
async function withTimeout(promise, ms) {
  let timer = null
  try {
    return await Promise.race([
      promise,
      new Promise((_, reject) => {
        timer = setTimeout(() => reject(new Error('timeout')), ms)
      })
    ])
  } finally {
    if (timer) clearTimeout(timer)
  }
}

// 发起HTTP请求并检查响应是否成功
async function fetchOk(url, options, timeoutMs = 1200) {
  try {
    const res = await withTimeout(fetch(url, options), timeoutMs)
    return !!res && res.ok
  } catch {
    return false
  }
}

async function isDevServerReachable() {
  return fetchOk(`${devServerUrl}/@vite/client`, { method: 'GET' }, 1200)
}

async function isCoreHealthy() {
  return fetchOk(coreHealthUrl, { method: 'GET' }, 1000)
}

async function waitCoreReady(maxWaitMs = 12000) {
  const start = Date.now()
  while (Date.now() - start < maxWaitMs) {
    if (await isCoreHealthy()) return true
    await sleep(300)
  }
  return false
}

function resolveCoreExecutable() {
  const manual = process.env.MA_CORE_EXE
  if (manual && fs.existsSync(manual)) return manual

  const exe = process.platform === 'win32' ? 'core.exe' : 'core'
  const candidates = [
    path.resolve(process.cwd(), exe),
    path.resolve(__dirname, '..', '..', exe),
    path.resolve(__dirname, '..', '..', 'bin', exe),
    path.resolve(process.resourcesPath || '', exe),
    path.resolve(process.resourcesPath || '', 'bin', exe)
  ]

  for (const candidate of candidates) {
    if (candidate && fs.existsSync(candidate)) return candidate
  }
  return null
}

function startCoreByCommand(commandText) {
  coreProcess = spawn(commandText, [], {
    shell: true,
    windowsHide: true,
    stdio: 'ignore',
    env: {
      ...process.env
    }
  })
  coreManagedByElectron = true
  wireCoreProcessEvents()
}

function startCoreByExecutable(executablePath) {
  coreProcess = spawn(executablePath, [], {
    windowsHide: true,
    stdio: 'ignore',
    env: {
      ...process.env
    }
  })
  coreManagedByElectron = true
  wireCoreProcessEvents()
}

function wireCoreProcessEvents() {
  if (!coreProcess) return
  writeLog('core spawned', { pid: coreProcess.pid })
  coreProcess.on('error', (error) => {
    writeLog('core process error', { error: error?.message || String(error) })
  })
  coreProcess.on('exit', (code, signal) => {
    writeLog('core process exit', { code, signal })
    coreProcess = null
  })
}

async function ensureCoreReady() {
  if (await isCoreHealthy()) {
    writeLog('core already running')
    return true
  }

  // Development should always start core externally (script/manual).
  if (!app.isPackaged) {
    writeLog('core managed disabled in development')
    return false
  }

  if (process.env.MA_CORE_MANAGED === '0') {
    writeLog('core managed disabled by env')
    return false
  }

  const coreCmd = process.env.MA_CORE_CMD
  if (coreCmd && coreCmd.trim()) {
    writeLog('starting core by command', { command: coreCmd })
    startCoreByCommand(coreCmd)
  } else {
    const coreExe = resolveCoreExecutable()
    if (!coreExe) {
      writeLog('core executable not found')
      return false
    } else {
      writeLog('starting core by executable', { executable: coreExe })
      startCoreByExecutable(coreExe)
    }
  }

  const ready = await waitCoreReady(12000)
  writeLog('core ready state', { ready })
  return ready
}

async function requestCoreQuit() {
  return fetchOk(`${coreBaseUrl}/api/internal/quit`, { method: 'POST' }, 1200)
}

async function stopManagedCore() {
  if (!coreManagedByElectron || coreShutdownRequested) return
  coreShutdownRequested = true

  await requestCoreQuit()

  const deadline = Date.now() + 3000
  while (coreProcess && coreProcess.exitCode == null && Date.now() < deadline) {
    await sleep(120)
  }

  if (coreProcess && coreProcess.exitCode == null) {
    try {
      coreProcess.kill()
      writeLog('core killed by electron')
    } catch (error) {
      writeLog('core kill failed', { error: error?.message || String(error) })
    }
  }
}

async function postUIReady(ready) {
  if (uiReadyState === ready) return
  uiReadyState = ready
  const url = `${coreBaseUrl}/api/internal/ui/ready`
  try {
    await fetch(url, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ ready })
    })
    writeLog('ui-ready', { ready })
  } catch (error) {
    writeLog('ui-ready failed', { ready, error: error?.message || String(error) })
  }
}

function startDevWatcher() {
  if (devWatcher) return
  devWatcher = setInterval(async () => {
    if (rendererMode !== 'safe') {
      clearInterval(devWatcher)
      devWatcher = null
      return
    }

    const devOk = await isDevServerReachable()
    if (!devOk) return

    rendererBase = devServerUrl
    rendererMode = 'dev'
    writeLog('renderer-mode switched', { mode: 'dev' })
    await postUIReady(true)

    if (mainWindow && !mainWindow.isDestroyed()) {
      mainWindow.loadURL(buildHashUrl(rendererBase, startRoute))
    }
    if (trayWindow && !trayWindow.isDestroyed()) {
      trayWindow.loadURL(buildHashUrl(rendererBase, '/tray-menu'))
    }

    clearInterval(devWatcher)
    devWatcher = null
  }, 3000)
}

async function resolveRendererBase() {
  if (rendererBase && rendererMode) {
    return { base: rendererBase, mode: rendererMode }
  }

  // In development, prefer dev server over dist
  if (await isDevServerReachable()) {
    rendererBase = devServerUrl
    rendererMode = 'dev'
    return { base: rendererBase, mode: rendererMode }
  }

  if (fs.existsSync(distIndex)) {
    rendererBase = `file://${distIndex.replace(/\\/g, '/')}`
    rendererMode = 'dist'
    return { base: rendererBase, mode: rendererMode }
  }

  rendererBase = `file://${safeModeHtml.replace(/\\/g, '/')}`
  rendererMode = 'safe'
  startDevWatcher()
  return { base: rendererBase, mode: rendererMode }
}

function forceShowMainWindow() {
  if (!mainWindow || mainWindow.isDestroyed()) return
  mainWindow.show()
  mainWindow.focus()
  if (!mainWindow.isMaximized()) {
    mainWindow.center()
  }
  mainWindow.setAlwaysOnTop(true)
  setTimeout(() => {
    if (mainWindow && !mainWindow.isDestroyed()) {
      mainWindow.setAlwaysOnTop(false)
    }
  }, 150)
}

async function createWindowAsync() {
  const resolved = await resolveRendererBase()
  await postUIReady(resolved.mode !== 'safe')

  // 检查是否首次启动
  const isFirstLaunch = await checkFirstLaunch()
  
  if (isFirstLaunch) {
    // 创建引导窗口
    await createOnboardingWindow(resolved)
  } else {
    // 创建主窗口
    await createMainWindow(resolved)
  }
}

async function checkFirstLaunch() {
  try {
    const res = await fetch(`${coreBaseUrl}/api/v1/system/first-launch`)
    const data = await res.json()
    return data?.data?.is_first_launch || false
  } catch (e) {
    return false
  }
}

async function createOnboardingWindow(resolved) {
  onboardingWindow = new BrowserWindow({
    width: 600,
    height: 550,
    center: true,
    resizable: false,
    frame: false,
    backgroundColor: '#1a1a1a',
    webPreferences: {
      preload: preloadPath,
      contextIsolation: true
    }
  })

  onboardingWindow.on('closed', () => {
    onboardingWindow = null
  })

  if (shouldOpenDevTools) {
    onboardingWindow.webContents.openDevTools({ mode: 'detach' })
  }

  if (resolved.mode === 'safe') {
    await onboardingWindow.loadURL(`file://${safeModeHtml.replace(/\\/g, '/')}`)
  } else {
    await onboardingWindow.loadURL(buildHashUrl(resolved.base, '/onboarding'))
  }
}

async function createMainWindow(resolved) {
  mainWindow = new BrowserWindow({
    width: 1200,
    height: 800,
    show: !hiddenStart,
    backgroundColor: '#1e1e1e',
    frame: false,
    autoHideMenuBar: true,
    webPreferences: {
      preload: preloadPath,
      contextIsolation: true
    }
  })

  mainWindow.on('close', (event) => {
    if (isQuitting) return
    event.preventDefault()
    mainWindow.hide()
    writeLog('mainWindow hidden')
  })

  mainWindow.webContents.on('did-fail-load', (_event, errorCode, errorDescription, validatedURL) => {
    writeLog('main did-fail-load', { errorCode, errorDescription, validatedURL })
  })

  if (shouldOpenDevTools) {
    mainWindow.webContents.openDevTools({ mode: 'detach' })
  }

  mainWindow.webContents.on('did-finish-load', () => {
    writeLog('main did-finish-load')
    if (!hiddenStart) {
      forceShowMainWindow()
    }
  })

  if (resolved.mode === 'safe') {
    await mainWindow.loadURL(`file://${safeModeHtml.replace(/\\/g, '/')}`)
  } else {
    await mainWindow.loadURL(buildHashUrl(resolved.base, startRoute))
  }
}

async function createTrayWindow() {
  const { base, mode } = await resolveRendererBase()
  if (mode === 'safe') {
    await postUIReady(false)
    return
  }
  await postUIReady(true)

  trayWindow = new BrowserWindow({
    width: 248,
    height: 236,
    show: false,
    frame: false,
    resizable: false,
    alwaysOnTop: true,
    skipTaskbar: true,
    transparent: true,
    backgroundColor: '#00000000',
    webPreferences: {
      preload: preloadPath,
      contextIsolation: true
    }
  })

  trayWindow.loadURL(buildHashUrl(base, '/tray-menu'))

  trayWindow.on('close', (event) => {
    if (isQuitting) return
    event.preventDefault()
    trayWindow.hide()
    trayWindowLastHiddenAt = Date.now()
  })

  trayWindow.on('blur', () => {
    if (!trayWindow || trayWindow.isDestroyed() || !trayWindow.isVisible()) return
    
    // Slight delay to prevent immediate closing on interaction
    if (Date.now() - trayWindowShownAt < 100) return

    trayWindow.hide()
    trayWindowLastHiddenAt = Date.now()
  })

  trayWindow.on('closed', () => {
    trayWindow = null
    if (trayHideTimer) {
      clearTimeout(trayHideTimer)
      trayHideTimer = null
    }
  })
}

function resolveTrayIconPath() {
  if (process.platform === 'darwin') {
    // macOS specific icons
    const macCandidates = [
      path.resolve(process.cwd(), 'build', 'appicon.png'),
      path.resolve(__dirname, '..', 'build', 'appicon.png'),
      path.resolve(__dirname, '..', '..', 'build', 'appicon.png'),
      path.resolve(process.resourcesPath || '', 'build', 'appicon.png'),
      path.resolve(process.resourcesPath || '', 'appicon.png')
    ]
    for (const candidate of macCandidates) {
      if (candidate && fs.existsSync(candidate)) return candidate
    }
    // Fallback to searching in parent directories if running from dev
    const devFallback = path.resolve(__dirname, '..', '..', 'build', 'appicon.png')
    if (fs.existsSync(devFallback)) return devFallback
  }

  const candidates = [
    path.resolve(process.cwd(), 'build', 'windows', 'icon.ico'),
    path.resolve(__dirname, '..', 'build', 'icon.ico'),
    path.resolve(__dirname, '..', '..', 'build', 'windows', 'icon.ico'),
    path.resolve(process.resourcesPath || '', 'build', 'windows', 'icon.ico'),
    path.resolve(process.resourcesPath || '', 'icon.ico')
  ]

  for (const candidate of candidates) {
    if (candidate && fs.existsSync(candidate)) return candidate
  }
  return null
}

function buildApplicationMenu() {
  const template = [
    {
      label: app.getName(),
      submenu: [
        {
          label: '关于智归档OS',
          role: 'about'
        },
        { type: 'separator' },
        {
          label: '设置...',
          accelerator: 'Cmd+,',
          click: () => {
            if (mainWindow && !mainWindow.isDestroyed()) {
              forceShowMainWindow()
              mainWindow.webContents.send('open-settings')
            }
          }
        },
        { type: 'separator' },
        {
          label: '隐藏智归档OS',
          accelerator: 'Cmd+H',
          role: 'hide'
        },
        {
          label: '隐藏其他',
          accelerator: 'Cmd+Option+H',
          role: 'hideothers'
        },
        {
          label: '显示全部',
          role: 'unhide'
        },
        { type: 'separator' },
        {
          label: '退出智归档OS',
          accelerator: 'Cmd+Q',
          click: () => {
            void requestQuit()
          }
        }
      ]
    },
    {
      label: '编辑',
      submenu: [
        {
          label: '撤销',
          accelerator: 'Cmd+Z',
          role: 'undo'
        },
        {
          label: '重做',
          accelerator: 'Shift+Cmd+Z',
          role: 'redo'
        },
        { type: 'separator' },
        {
          label: '剪切',
          accelerator: 'Cmd+X',
          role: 'cut'
        },
        {
          label: '复制',
          accelerator: 'Cmd+C',
          role: 'copy'
        },
        {
          label: '粘贴',
          accelerator: 'Cmd+V',
          role: 'paste'
        },
        {
          label: '全选',
          accelerator: 'Cmd+A',
          role: 'selectall'
        }
      ]
    },
    {
      label: '视图',
      submenu: [
        {
          label: '重新加载',
          accelerator: 'Cmd+R',
          role: 'reload'
        },
        {
          label: '强制重新加载',
          accelerator: 'Shift+Cmd+R',
          role: 'forceReload'
        },
        {
          label: '开发者工具',
          accelerator: 'Alt+Cmd+I',
          role: 'toggleDevTools'
        },
        { type: 'separator' },
        {
          label: '实际大小',
          accelerator: 'Cmd+0',
          role: 'resetZoom'
        },
        {
          label: '放大',
          accelerator: 'Cmd+Plus',
          role: 'zoomIn'
        },
        {
          label: '缩小',
          accelerator: 'Cmd+-',
          role: 'zoomOut'
        },
        { type: 'separator' },
        {
          label: '全屏',
          accelerator: 'Ctrl+Cmd+F',
          role: 'togglefullscreen'
        }
      ]
    },
    {
      label: '窗口',
      submenu: [
        {
          label: '最小化',
          accelerator: 'Cmd+M',
          role: 'minimize'
        },
        {
          label: '关闭',
          accelerator: 'Cmd+W',
          role: 'close'
        },
        { type: 'separator' },
        {
          label: '全部置于顶层',
          role: 'front'
        }
      ]
    },
    {
      label: '帮助',
      submenu: [
        {
          label: '学习更多',
          click: () => {
            require('electron').shell.openExternal('https://example.com')
          }
        }
      ]
    }
  ]

  return Menu.buildFromTemplate(template)
}

function buildNativeTrayMenu() {
  return Menu.buildFromTemplate([
    {
      label: '显示主界面',
      click: () => forceShowMainWindow()
    },
    { type: 'separator' },
    {
      label: '退出',
      click: () => {
        void requestQuit()
      }
    }
  ])
}

function detectTaskbarSide(display) {
  const bounds = display.bounds
  const area = display.workArea

  if (area.y > bounds.y) return 'top'
  if (area.x > bounds.x) return 'left'
  if (area.x === bounds.x && area.width < bounds.width) return 'right'
  return 'bottom'
}

function getTrayAnchor(anchor) {
  if (anchor?.trayBounds) {
    return {
      x: anchor.trayBounds.x + Math.floor(anchor.trayBounds.width / 2),
      y: anchor.trayBounds.y + Math.floor(anchor.trayBounds.height / 2)
    }
  }
  if (anchor?.cursor && typeof anchor.cursor.x === 'number' && typeof anchor.cursor.y === 'number') {
    return anchor.cursor
  }
  return screen.getCursorScreenPoint()
}

function computeTrayWindowPosition(winBounds, anchor) {
  const trayBounds = anchor?.trayBounds
  const cursor = getTrayAnchor(anchor)
  const display = screen.getDisplayNearestPoint(cursor)
  const area = display.workArea
  const side = detectTaskbarSide(display)
  const gap = -11

  let x, y

  if (trayBounds) {
    // Horizontal alignment: Center on icon
    x = trayBounds.x + Math.round((trayBounds.width - winBounds.width) / 2)
    
    // Vertical alignment based on taskbar side
    if (side === 'bottom') {
      // Anchor to bottom of work area (top of taskbar)
      y = area.y + area.height - winBounds.height - gap
    } else if (side === 'top') {
      // Anchor to top of work area (bottom of taskbar)
      y = area.y + gap
    } else if (side === 'left') {
      x = area.x + gap
      y = trayBounds.y + Math.round((trayBounds.height - winBounds.height) / 2)
    } else { // right
      x = area.x + area.width - winBounds.width - gap
      y = trayBounds.y + Math.round((trayBounds.height - winBounds.height) / 2)
    }
  } else {
    // Fallback if no tray bounds
    x = cursor.x - winBounds.width / 2
    if (side === 'bottom') {
      y = area.y + area.height - winBounds.height - gap
    } else {
      y = cursor.y - winBounds.height - gap
    }
  }

  // Ensure fully within work area (clamping) - ONLY for X axis to prevent horizontal overflow
  // We trust the Y calculation above (especially for negative gaps/overlaps)
  if (x < area.x) x = area.x
  if (x + winBounds.width > area.x + area.width) x = area.x + area.width - winBounds.width
  
  // Removed strict Y clamping to allow negative gaps (overlapping taskbar)
  // if (y < area.y) y = area.y
  // if (y + winBounds.height > area.y + area.height) y = area.y + area.height - winBounds.height

  return { x: Math.round(x), y: Math.round(y) }
}

async function showTrayWindow(anchor) {
  const resolved = await resolveRendererBase()
  if (resolved.mode === 'safe') {
    await postUIReady(false)
    if (tray) {
      tray.popUpContextMenu(buildNativeTrayMenu())
    }
    return
  }

  await postUIReady(true)

  if (!trayWindow || trayWindow.isDestroyed()) {
    await createTrayWindow()
    if (!trayWindow || trayWindow.isDestroyed()) return
  }

  const now = Date.now()
  if (now - trayWindowLastHiddenAt < 260) {
    return
  }

  if (trayWindow.isVisible()) {
    trayWindow.hide()
    trayWindowLastHiddenAt = Date.now()
    return
  }

  if (trayHideTimer) {
    clearTimeout(trayHideTimer)
    trayHideTimer = null
  }

  const winBounds = trayWindow.getBounds()
  const pos = computeTrayWindowPosition(winBounds, anchor)
  trayWindow.setPosition(pos.x, pos.y, false)
  trayWindow.setAlwaysOnTop(true, 'status')
  trayWindow.setVisibleOnAllWorkspaces(true, { visibleOnFullScreen: true })
  trayWindow.show()
  trayWindow.focus()
  trayWindowShownAt = Date.now()
}

function hideTrayWindow() {
  if (!trayWindow || trayWindow.isDestroyed()) return
  trayWindow.hide()
  trayWindowLastHiddenAt = Date.now()
}

function createTray() {
  const iconPath = resolveTrayIconPath()
  writeLog('create-tray', { iconPath })
  
  // Create native image and resize for macOS tray (usually 16x16 or 22x22 points)
  let icon = nativeImage.createEmpty()
  if (iconPath) {
    icon = nativeImage.createFromPath(iconPath)
    if (process.platform === 'darwin') {
      icon = icon.resize({ width: 22, height: 22 })
      icon.setTemplateImage(true) // Adapts to light/dark mode
    }
  }

  tray = new Tray(icon)
  tray.setToolTip('Media Assistant - 智归档')

  tray.setIgnoreDoubleClickEvents(true) 

  tray.on('click', (event, bounds, position) => {
    writeLog('tray-click', { bounds, position })
    const cursor = screen.getCursorScreenPoint()
    void showTrayWindow({
        trayBounds: bounds,
        cursor: cursor
    })
  })

  tray.on('double-click', () => {
    writeLog('tray-double-click')
    hideTrayWindow()
    forceShowMainWindow()
  })

  tray.on('right-click', (event, bounds) => {
    writeLog('tray-right-click', { bounds })
    const cursor = screen.getCursorScreenPoint()
    void showTrayWindow({
        trayBounds: bounds,
        cursor: cursor
    })
  })
}

async function requestQuit() {
  if (isQuitting) return
  isQuitting = true
  hideTrayWindow()
  await stopManagedCore()
  app.quit()
}

const singleInstance = app.requestSingleInstanceLock()
if (!singleInstance) {
  console.log('[electron] another instance running, quitting')
  app.quit()
} else {
  console.log('[electron] single instance lock acquired')
  app.on('second-instance', () => {
    forceShowMainWindow()
  })

  app.whenReady().then(async () => {
    console.log('[electron] app ready')
    writeLog('ready')
    process.env.MA_CORE_BASE_URL = coreBaseUrl

    // Register custom protocol for assets
    protocol.registerFileProtocol('assets', (request, callback) => {
      const url = request.url.replace('assets://', '')
      const parts = url.split('/')
      
      if (parts[0] === 'thumbnails' && parts[1]) {
        const assetId = parts[1]
        // Thumbnails are stored in {projectRoot}/data/cache/thumbnails/{assetId}.jpg
        const projectRoot = path.join(__dirname, '..', '..')
        const thumbnailPath = path.join(projectRoot, 'data', 'cache', 'thumbnails', `${assetId}.jpg`)
        callback({ path: thumbnailPath })
      } else {
        callback({ error: -6 }) // FILE_NOT_FOUND
      }
    })

    const applicationMenu = buildApplicationMenu()
    Menu.setApplicationMenu(applicationMenu)

    const coreReady = await ensureCoreReady()
    writeLog('core-ready', { coreReady })

    ipcMain.on('renderer-error', (_event, payload) => {
      writeLog('renderer-error', payload)
    })

    ipcMain.on('window-minimize', () => {
      if (mainWindow && !mainWindow.isDestroyed()) mainWindow.minimize()
    })

    ipcMain.on('window-maximize', () => {
      if (!mainWindow || mainWindow.isDestroyed()) return
      if (mainWindow.isMaximized()) mainWindow.unmaximize()
      else mainWindow.maximize()
    })

    ipcMain.on('window-close', () => {
      if (mainWindow && !mainWindow.isDestroyed()) {
        mainWindow.close()
      }
    })

    ipcMain.on('onboarding-complete', async () => {
      writeLog('onboarding-complete')
      
      // 关闭引导窗口
      if (onboardingWindow && !onboardingWindow.isDestroyed()) {
        onboardingWindow.close()
      }
      
      // 创建主窗口
      const resolved = await resolveRendererBase()
      await createMainWindow(resolved)
      
      // 显示主窗口
      if (mainWindow && !mainWindow.isDestroyed()) {
        mainWindow.show()
        mainWindow.focus()
      }
    })

    ipcMain.on('tray-action', (_event, data) => {
      const action = typeof data === 'string' ? data : data?.action
      const payload = typeof data === 'object' ? data?.payload : undefined
      writeLog('tray-action', { action })

      if (action === 'show-tray-menu') {
        void showTrayWindow(payload)
        return
      }
      if (action === 'hide-tray-menu') {
        hideTrayWindow()
        return
      }
      if (action === 'show-main') {
        hideTrayWindow()
        forceShowMainWindow()
        return
      }
      if (action === 'quit') {
        void requestQuit()
        return
      }
    })

    await createWindowAsync()
    await createTrayWindow()
    createTray()
  })

  app.on('activate', () => {
    if (!mainWindow || mainWindow.isDestroyed()) {
      void createWindowAsync()
    } else {
      forceShowMainWindow()
    }
  })

  app.on('before-quit', () => {
    isQuitting = true
  })

  app.on('window-all-closed', () => {
    // Keep app running in tray.
  })

  app.on('will-quit', () => {
    if (devWatcher) {
      clearInterval(devWatcher)
      devWatcher = null
    }
    void stopManagedCore()
  })
}
