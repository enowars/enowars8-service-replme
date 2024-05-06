const defaultTheme = require('tailwindcss/defaultTheme')


/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./**/*.templ",
  ],
  theme: {
    fontFamily: {
      'mono': [
        '"DejaVuSansM Nerd Font"', ...defaultTheme.fontFamily.mono,
      ],
    },
    extend: {
      gridTemplateColumns: {
        'sessions': 'auto min-content min-content',
      },
    },
  },
  plugins: [],
}

