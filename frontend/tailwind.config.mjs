/** @type {import('tailwindcss').Config} */
export default {
  content: ['./src/**/*.{js,jsx,ts,tsx}'],
  theme: {
    colors: {
      background: 'var(--background)',
      foreground: 'var(--foreground)',
      'muted-background': 'var(--muted-background)',
      'muted-foreground': 'var(--muted-foreground)',
      border: 'var(--border)',
      white: 'var(--white)',
    },
  },
}
