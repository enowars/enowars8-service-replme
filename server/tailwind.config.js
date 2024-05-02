/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [ "./**/*.templ" ],
  theme: {
    extend: {
      gridTemplateColumns: {
        'sessions': 'auto min-content min-content',
      }
    },
  },
  plugins: [],
}

