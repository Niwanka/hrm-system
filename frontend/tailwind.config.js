/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  darkMode: 'class',
  theme: {
    extend: {
      colors: {
        dark: {
          bg: '#09090b',
          card: '#121215',
          border: '#27272a',
          hover: '#1e1e24',
        },
        emerald: {
          400: '#34d399',
          500: '#10b981',
          600: '#059669',
          700: '#047857',
          glow: 'rgba(16, 185, 129, 0.35)',
        }
      },
      fontFamily: {
        sans: ['Inter', 'sans-serif'],
      },
    },
  },
  plugins: [],
}
