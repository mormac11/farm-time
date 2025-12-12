/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        // Rustic brown - deer hide, leather
        farm: {
          50: '#faf6f1',
          100: '#f0e6d8',
          200: '#e0ccb0',
          300: '#c9a87a',
          400: '#b8915a',
          500: '#a67c4e',
          600: '#8b6342',
          700: '#6f4e36',
          800: '#5c4030',
          900: '#4a352a',
        },
        // Forest green - hunting, woods
        forest: {
          50: '#f3f6f4',
          100: '#e0e9e2',
          200: '#c2d4c7',
          300: '#9ab8a2',
          400: '#6f9679',
          500: '#4e7a5a',
          600: '#3d6247',
          700: '#334f3b',
          800: '#2b4032',
          900: '#24352a',
        },
        // Rust/orange - tractors, autumn
        rust: {
          50: '#fdf6f3',
          100: '#fae8e0',
          200: '#f5cfc0',
          300: '#edae96',
          400: '#e28460',
          500: '#d4683f',
          600: '#c25434',
          700: '#a1432c',
          800: '#843929',
          900: '#6c3226',
        },
        // Cream background
        cream: {
          50: '#fefdfb',
          100: '#fdf9f3',
          200: '#f9f1e4',
          300: '#f3e4ce',
          400: '#e9d0ac',
        },
      },
    },
  },
  plugins: [],
}
