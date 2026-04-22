/** @type {import('tailwindcss').Config} */
export default {
  content: ["./index.html", "./src/**/*.{ts,tsx}"],
  theme: {
    extend: {
      fontFamily: {
        sans: ["Avenir Next", "Segoe UI", "Helvetica Neue", "sans-serif"],
      },
      colors: {
        ink: "#101828",
        steel: "#475467",
        sand: "#f8f5ef",
        flare: "#d97706",
        pine: "#0f766e"
      },
      boxShadow: {
        panel: "0 18px 40px rgba(16, 24, 40, 0.08)"
      }
    },
  },
  plugins: [],
};
