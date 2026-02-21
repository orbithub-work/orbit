import { defineStore } from 'pinia';
import { ref } from 'vue';
import { themes, type ThemeDefinition } from '../styles/themes';

export const useThemeStore = defineStore('theme', () => {
  const currentThemeId = ref(localStorage.getItem('theme-id') || 'eagle-dark');
  const currentTheme = ref<ThemeDefinition>(
    themes.find(t => t.id === currentThemeId.value) || themes[0]
  );

  const setTheme = (themeId: string) => {
    const theme = themes.find(t => t.id === themeId);
    if (theme) {
      currentThemeId.value = themeId;
      currentTheme.value = theme;
      localStorage.setItem('theme-id', themeId);
      applyTheme(theme);
    }
  };

  const applyTheme = (theme: ThemeDefinition) => {
    const root = document.documentElement;
    Object.entries(theme.variables).forEach(([key, value]) => {
      root.style.setProperty(key, value);
    });
    
    // 特殊处理流光球颜色
    theme.glowColors.forEach((color, index) => {
      root.style.setProperty(`--glow-${index + 1}`, color);
    });
  };

  // 初始化应用主题
  applyTheme(currentTheme.value);

  return {
    currentThemeId,
    currentTheme,
    themes,
    setTheme
  };
});
