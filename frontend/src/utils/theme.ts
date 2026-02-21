// Theme configuration interface
export interface ThemeConfig {
  name: 'light' | 'dark' | 'auto';
  colors: {
    primary: string;
    secondary: string;
    background: string;
    surface: string;
    text: string;
    textSecondary: string;
    border: string;
    error: string;
    warning: string;
    success: string;
    info: string;
  };
}

// Default light theme
export const lightTheme: ThemeConfig = {
  name: 'light',
  colors: {
    primary: '#2563eb',
    secondary: '#64748b',
    background: '#f8fafc',
    surface: '#ffffff',
    text: '#1e293b',
    textSecondary: '#64748b',
    border: '#e2e8f0',
    error: '#ef4444',
    warning: '#f59e0b',
    success: '#22c55e',
    info: '#3b82f6'
  }
};

// Default dark theme
export const darkTheme: ThemeConfig = {
  name: 'dark',
  colors: {
    primary: '#3b82f6',
    secondary: '#94a3b8',
    background: '#0f172a',
    surface: '#1e293b',
    text: '#f1f5f9',
    textSecondary: '#94a3b8',
    border: '#334155',
    error: '#f87171',
    warning: '#fbbf24',
    success: '#4ade80',
    info: '#60a5fa'
  }
};

// Theme manager class
class ThemeManager {
  private currentTheme: ThemeConfig = lightTheme;
  private observer: MutationObserver | null = null;

  constructor() {
    this.init();
  }

  private init(): void {
    // Initialize theme based on user preference or system setting
    this.applyTheme(this.getStoredTheme() || this.detectSystemTheme());
    
    // Watch for system theme changes
    this.setupSystemThemeListener();
  }

  private detectSystemTheme(): ThemeConfig {
    const isDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
    return isDark ? darkTheme : lightTheme;
  }

  private setupSystemThemeListener(): void {
    const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)');
    mediaQuery.addEventListener('change', (e) => {
      if (this.currentTheme.name === 'auto') {
        const newTheme = e.matches ? darkTheme : lightTheme;
        this.applyTheme(newTheme);
      }
    });
  }

  public setTheme(themeName: 'light' | 'dark' | 'auto'): void {
    let theme: ThemeConfig;
    
    if (themeName === 'auto') {
      theme = this.detectSystemTheme();
      theme.name = 'auto'; // Ensure the name is correctly set
    } else {
      theme = themeName === 'dark' ? darkTheme : lightTheme;
    }
    
    this.currentTheme = theme;
    this.applyTheme(theme);
    this.storeTheme(themeName);
  }

  private applyTheme(theme: ThemeConfig): void {
    const root = document.documentElement;
    
    // Apply CSS variables
    Object.entries(theme.colors).forEach(([key, value]) => {
      root.style.setProperty(`--color-${key}`, value);
    });
    
    // Apply theme class
    root.className = theme.name === 'dark' ? 'dark' : '';
  }

  private storeTheme(themeName: 'light' | 'dark' | 'auto'): void {
    localStorage.setItem('theme', themeName);
  }

  private getStoredTheme(): 'light' | 'dark' | 'auto' | null {
    const stored = localStorage.getItem('theme');
    return stored as 'light' | 'dark' | 'auto' || null;
  }

  public getCurrentTheme(): ThemeConfig {
    return this.currentTheme;
  }

  public toggleTheme(): void {
    const nextTheme = this.currentTheme.name === 'light' 
      ? 'dark' 
      : this.currentTheme.name === 'dark' 
        ? 'auto' 
        : 'light';
        
    this.setTheme(nextTheme);
  }
}

// Create a singleton instance
export const themeManager = new ThemeManager();

// Expose theme manager globally for use in components
declare global {
  interface Window {
    themeManager: ThemeManager;
  }
}

if (typeof window !== 'undefined') {
  window.themeManager = themeManager;
}