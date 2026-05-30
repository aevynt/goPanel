/** @type {import('tailwindcss').Config} */
export default {
  content: [
    './index.html',
    './src/**/*.{vue,ts}',
  ],
  theme: {
    extend: {
      colors: {
        surface: '#1f1f1d',
        'surface-raised': '#2a2a28',
        'text-primary': '#e8e6dc',
        'text-secondary': '#9c9a92',
        'text-tertiary': '#73726c',
        accent: '#c96442',
        'accent-hover': '#d4733f',
        'accent-muted': 'rgba(201, 100, 66, 0.15)',
        border: '#3d3d3a',
        'border-light': '#2a2a28',
        danger: '#ef4444',
        success: '#22c55e',
        warning: '#f59e0b',
      },
      fontFamily: {
        heading: ['Playfair Display', 'Georgia', 'serif'],
        body: ['DM Sans', 'system-ui', '-apple-system', 'sans-serif'],
        mono: ['JetBrains Mono', 'Fira Code', 'monospace'],
      },
      borderRadius: {
        DEFAULT: '8px',
        sm: '4px',
        lg: '12px',
      },
    },
  },
  plugins: [],
}
