export interface ThemeDefinition {
  id: string;
  name: string;
  description: string;
  variables: Record<string, string>;
  glowColors: string[];
}

export const themes: ThemeDefinition[] = [
  {
    id: 'eagle-dark',
    name: 'Eagle 经典',
    description: '朴素、高效、专注内容的深色风格',
    glowColors: ['rgba(0, 122, 204, 0.1)', 'rgba(0, 0, 0, 0)', 'rgba(0, 0, 0, 0)'],
    variables: {
      '--color-primary': '#007acc', // Eagle Blue
      '--color-secondary': '#4d4d4d', // Dark Gray
      '--color-accent': '#007acc',
      '--color-purple': '#007acc',
      '--grad-primary': 'linear-gradient(180deg, #007acc 0%, #0063a5 100%)', // Subtle gradient
      '--color-bg-base': '#1e1e1e', // Main Content
      '--color-bg-sidebar': '#252526', // Sidebar
      '--color-bg-header': '#2d2d2d', // Header / Toolbar
      '--color-bg-surface': '#2d2d2d', // Cards
      '--glass-border-bright': 'rgba(255, 255, 255, 0.05)',
      '--color-text-primary': '#cccccc', // Soft White
      '--color-text-secondary': '#858585', // Muted Text
      '--color-border': '#3e3e42', // VS Code Border
      '--glass-blur': '0px', // No blur
      '--shadow-sm': 'none',
      '--shadow-lg': '0 2px 8px rgba(0, 0, 0, 0.3)',
      '--color-active-bg': 'rgba(0, 122, 204, 0.15)',
      '--color-selected': 'rgba(0, 122, 204, 0.25)',
      '--color-hover': 'rgba(255, 255, 255, 0.05)',
      '--color-surface': '#2d2d2d',
    }
  },
  {
    id: 'pro-dark',
    name: '专业深色',
    description: '沉浸式暗色调，专注内容创作',
    glowColors: ['rgba(59, 130, 246, 0.15)', 'rgba(99, 102, 241, 0.15)', 'rgba(255, 255, 255, 0.05)'],
    variables: {
      '--color-primary': '#3b82f6', // Blue 500
      '--color-secondary': '#64748b', // Slate 500
      '--color-accent': '#f59e0b', // Amber 500
      '--color-purple': '#8b5cf6', // Violet 500
      '--grad-primary': 'linear-gradient(135deg, #3b82f6 0%, #2563eb 100%)',
      '--color-bg-base': '#18181b', // Zinc 950
      '--color-bg-sidebar': '#27272a', // Zinc 900
      '--color-bg-header': '#27272a', // Zinc 900
      '--color-bg-surface': '#27272a', // Zinc 900
      '--glass-border-bright': 'rgba(255, 255, 255, 0.1)',
      '--color-text-primary': '#f4f4f5', // Zinc 100
      '--color-text-secondary': '#a1a1aa', // Zinc 400
      '--color-border': 'rgba(255, 255, 255, 0.08)',
    }
  },
  {
    id: 'pro-light',
    name: '专业浅色',
    description: '明亮清爽，适合日间工作',
    glowColors: ['rgba(59, 130, 246, 0.1)', 'rgba(99, 102, 241, 0.1)', 'rgba(0, 0, 0, 0.05)'],
    variables: {
      '--color-primary': '#3b82f6', // Blue 500
      '--color-secondary': '#64748b', // Slate 500
      '--color-accent': '#f59e0b', // Amber 500
      '--color-purple': '#8b5cf6', // Violet 500
      '--grad-primary': 'linear-gradient(135deg, #3b82f6 0%, #2563eb 100%)',
      '--color-bg-base': '#ffffff', // White
      '--color-bg-sidebar': '#f8fafc', // Slate 50
      '--color-bg-header': '#ffffff', // White
      '--color-bg-surface': '#ffffff', // White
      '--glass-border-bright': 'rgba(0, 0, 0, 0.06)',
      '--color-text-primary': '#18181b', // Zinc 950
      '--color-text-secondary': '#71717a', // Zinc 500
      '--color-border': 'rgba(0, 0, 0, 0.08)',
    }
  },
  {
    id: 'cyberpunk',
    name: '极致赛博',
    description: '黑透材质 + 霓虹彩虹流光',
    glowColors: ['rgba(255, 0, 85, 0.25)', 'rgba(0, 242, 255, 0.2)', 'rgba(188, 0, 255, 0.2)', 'rgba(240, 255, 0, 0.15)'],
    variables: {
      '--color-primary': '#ff0055',
      '--color-secondary': '#00f2ff',
      '--color-accent': '#f0ff00',
      '--color-purple': '#bc00ff',
      '--grad-primary': 'linear-gradient(135deg, #ff0055 0%, #bc00ff 50%, #00f2ff 100%)',
      '--color-bg-base': '#050506',
      '--color-bg-sidebar': 'rgba(10, 10, 12, 0.85)',
      '--color-bg-header': 'rgba(10, 10, 12, 0.7)',
      '--color-bg-surface': 'rgba(20, 20, 25, 0.6)',
      '--glass-border-bright': 'rgba(255, 0, 85, 0.4)',
      '--color-text-primary': '#ffffff',
      '--color-text-secondary': '#a0a0b0',
      '--color-border': 'rgba(255, 255, 255, 0.08)',
    }
  },
  {
    id: 'tech-cyan',
    name: '高科技冷流',
    description: '黑透 + 荧光蓝/绿',
    glowColors: ['rgba(0, 242, 255, 0.3)', 'rgba(0, 255, 157, 0.2)', 'rgba(0, 100, 255, 0.2)'],
    variables: {
      '--color-primary': '#00f2ff',
      '--color-secondary': '#00ff9d',
      '--color-accent': '#00ff9d',
      '--color-purple': '#0064ff',
      '--grad-primary': 'linear-gradient(135deg, #00f2ff 0%, #00ff9d 100%)',
      '--color-bg-base': '#080a0f',
      '--color-bg-sidebar': 'rgba(8, 10, 15, 0.9)',
      '--color-bg-header': 'rgba(8, 10, 15, 0.75)',
      '--color-bg-surface': 'rgba(15, 20, 30, 0.6)',
      '--glass-border-bright': 'rgba(0, 242, 255, 0.4)',
    }
  },
  {
    id: 'mystic-purple',
    name: '神秘幽魅',
    description: '黑透 + 紫色/紫红流光',
    glowColors: ['rgba(188, 0, 255, 0.3)', 'rgba(255, 0, 153, 0.2)', 'rgba(100, 0, 255, 0.2)'],
    variables: {
      '--color-primary': '#bc00ff',
      '--color-secondary': '#ff0099',
      '--color-accent': '#ff0099',
      '--color-purple': '#bc00ff',
      '--grad-primary': 'linear-gradient(135deg, #bc00ff 0%, #ff0099 100%)',
      '--color-bg-base': '#0a0510',
      '--color-bg-sidebar:': 'rgba(10, 5, 16, 0.9)',
      '--color-bg-header': 'rgba(10, 5, 16, 0.75)',
      '--color-bg-surface': 'rgba(20, 10, 30, 0.6)',
      '--glass-border-bright': 'rgba(188, 0, 255, 0.4)',
    }
  },
  {
    id: 'speed-red',
    name: '速度动感',
    description: '黑透 + 熔岩红流光',
    glowColors: ['rgba(255, 50, 0, 0.3)', 'rgba(255, 100, 0, 0.2)', 'rgba(150, 0, 0, 0.2)'],
    variables: {
      '--color-primary': '#ff3200',
      '--color-secondary': '#ff6400',
      '--color-accent': '#ffc800',
      '--color-purple': '#960000',
      '--grad-primary': 'linear-gradient(135deg, #ff3200 0%, #ff6400 100%)',
      '--color-bg-base': '#0f0505',
      '--color-bg-sidebar': 'rgba(15, 5, 5, 0.9)',
      '--color-bg-header': 'rgba(15, 5, 5, 0.75)',
      '--color-bg-surface': 'rgba(25, 10, 10, 0.6)',
      '--glass-border-bright': 'rgba(255, 50, 0, 0.4)',
    }
  },
  {
    id: 'phantom-white',
    name: '高级渐变',
    description: '黑透 + 幻彩白流光',
    glowColors: ['rgba(255, 255, 255, 0.25)', 'rgba(200, 220, 255, 0.2)', 'rgba(255, 200, 255, 0.15)'],
    variables: {
      '--color-primary': '#ffffff',
      '--color-secondary': '#e2e8f0',
      '--color-accent': '#ffffff',
      '--color-purple': '#94a3b8',
      '--grad-primary': 'linear-gradient(135deg, #ffffff 0%, #cbd5e1 100%)',
      '--color-bg-base': '#020617',
      '--color-bg-sidebar': 'rgba(15, 23, 42, 0.9)',
      '--color-bg-header': 'rgba(15, 23, 42, 0.75)',
      '--color-bg-surface': 'rgba(30, 41, 59, 0.6)',
      '--glass-border-bright': 'rgba(255, 255, 255, 0.4)',
    }
  }
];
