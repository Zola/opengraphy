/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./internal/templates/**/*.html",
    "./web/static/js/**/*.js"
  ],
  theme: {
    extend: {
      colors: {
        brand: "#2458ff",
        ink: "#182033",
        muted: "#667085",
        line: "#d8dee8"
      },
      boxShadow: {
        soft: "0 18px 60px rgba(24, 32, 51, 0.10)",
        panel: "0 10px 34px rgba(24, 32, 51, 0.08)"
      }
    }
  },
  plugins: []
};
