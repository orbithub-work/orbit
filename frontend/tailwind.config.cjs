/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        primary: {
          DEFAULT: 'var(--color-primary)',
          50: 'var(--color-primary-50)',
          100: 'var(--color-primary-100)',
        },
        background: 'var(--color-background)',
        surface: 'var(--color-surface)',
        text: {
          DEFAULT: 'var(--color-text)',
          secondary: 'var(--color-text-secondary)'
        }
      }
    },
  },
  plugins: [
    require('@tailwindcss/typography'),
    require('daisyui'),
  ],
  daisyui: {
    themes: [
      {
        dark: {
          "primary": "#f97316",
          "primary-content": "#ffffff",
          "secondary": "#4d4d4d",
          "secondary-content": "#cccccc",
          "accent": "#fb923c",
          "accent-content": "#1a1a1a",
          "neutral": "#3d3d3d",
          "neutral-content": "#cccccc",
          "base-100": "#1e1e1e",
          "base-200": "#252526",
          "base-300": "#2d2d2d",
          "base-content": "#cccccc",
          "info": "#06b6d4",
          "success": "#10b981",
          "warning": "#f59e0b",
          "error": "#ef4444",
          "--rounded-box": "0.5rem",
          "--rounded-btn": "0.375rem",
          "--rounded-badge": "1rem",
        },
      },
      {
        light: {
          "primary": "#f97316",
          "primary-content": "#ffffff",
          "secondary": "#e5e7eb",
          "secondary-content": "#374151",
          "accent": "#fb923c",
          "accent-content": "#1a1a1a",
          "neutral": "#374151",
          "neutral-content": "#f9fafb",
          "base-100": "#ffffff",
          "base-200": "#f3f4f6",
          "base-300": "#e5e7eb",
          "base-content": "#1f2937",
          "info": "#06b6d4",
          "success": "#10b981",
          "warning": "#f59e0b",
          "error": "#ef4444",
          "--rounded-box": "0.5rem",
          "--rounded-btn": "0.375rem",
          "--rounded-badge": "1rem",
        },
      },
    ],
    darkTheme: "dark",
  },
}
