import { createRouter, createWebHashHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import AppSafe from '@/AppSafe.vue'

// 定义路由组件（使用懒加载）
const DashboardView = () => import('@/views/DashboardView.vue')
const FilesView = () => import('@/views/FilesView.vue')
const CollectionsView = () => import('@/views/CollectionsView.vue')
const ProjectsView = () => import('@/views/ProjectsView.vue')
const SettingsView = () => import('@/views/SettingsView.vue')
const LoginView = () => import('@/views/LoginView.vue')
const NotFoundView = () => import('@/views/NotFoundView.vue')

// 定义路由类型
interface RouteMeta {
  title?: string
  requiresAuth?: boolean
  public?: boolean
}

// 定义路由配置
const routes = [
  {
    path: '/',
    component: AppSafe,
    redirect: '/pool',
    children: [
      {
        path: 'pool',
        name: 'Pool',
        components: {
          sidebar: () => import('@/components/app-shell/PoolSidebar.vue'),
          main: () => import('@/components/app-shell/PoolView.vue')
        },
        meta: { title: '素材库' } satisfies RouteMeta
      },
      {
        path: 'project',
        name: 'Project',
        components: {
          sidebar: () => import('@/components/app-shell/ProjectSidebar.vue'),
          main: () => import('@/components/app-shell/ProjectView.vue')
        },
        meta: { title: '项目库' } satisfies RouteMeta
      },
      {
        path: 'artifact',
        name: 'Artifact',
        components: {
          sidebar: () => import('@/components/app-shell/ArtifactSidebar.vue'),
          main: () => import('@/components/app-shell/ArtifactView.vue')
        },
        meta: { title: '交付库' } satisfies RouteMeta
      },
      {
        path: 'analytics',
        name: 'Analytics',
        components: {
          main: () => import('@/components/analytics/AnalyticsView.vue')
        },
        meta: { title: '数据看板' } satisfies RouteMeta
      }
    ]
  },
  {
    path: '/dashboard',
    name: 'Dashboard',
    component: DashboardView,
    meta: { title: '仪表盘' } satisfies RouteMeta
  },
  {
    path: '/files',
    name: 'Files',
    component: FilesView,
    meta: { title: '文件管理' } satisfies RouteMeta
  },
  {
    path: '/collections/:type?',
    name: 'Collections',
    component: CollectionsView,
    meta: { title: '素材库' } satisfies RouteMeta,
    props: true
  },
  {
    path: '/settings',
    name: 'Settings',
    component: SettingsView,
    meta: { title: '设置' } satisfies RouteMeta
  },
  {
    path: '/projects',
    name: 'Projects',
    component: ProjectsView,
    meta: { title: '项目管理' } satisfies RouteMeta
  },
  {
    path: '/projects/:id',
    name: 'ProjectDetail',
    component: () => import('@/components/ProjectView.vue'),
    meta: { title: '项目详情' } satisfies RouteMeta,
    props: true
  },
  {
    path: '/plugins/:id',
    name: 'PluginDetail',
    component: () => import('@/views/PluginView.vue'),
    meta: { title: '插件' } satisfies RouteMeta
  },
  {
    path: '/login',
    name: 'Login',
    component: LoginView,
    meta: { title: '登录', public: true } satisfies RouteMeta
  },
  {
    path: '/onboarding',
    name: 'Onboarding',
    component: () => import('@/views/OnboardingView.vue'),
    meta: { title: '欢迎使用', public: true } satisfies RouteMeta
  },
  {
    path: '/dock',
    name: 'Dock',
    component: () => import('@/views/DockView.vue'),
    meta: { title: 'Dock', public: true } satisfies RouteMeta
  },
  {
    path: '/tray-menu',
    name: 'TrayMenu',
    component: () => import('@/views/TrayMenuView.vue'),
    meta: { title: 'Tray Menu', public: true } satisfies RouteMeta
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: NotFoundView,
    meta: { title: '页面未找到' } satisfies RouteMeta
  }
]

// 创建路由实例
const router = createRouter({
  history: createWebHashHistory(),
  routes
})

// 导航守卫
router.beforeEach((to, from, next) => {
  // 设置页面标题
  document.title = to.meta?.title ? `${to.meta.title} - 智归档OS` : '智归档OS'

  next()
})

export default router
