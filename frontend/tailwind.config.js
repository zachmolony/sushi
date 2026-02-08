/** @type {import('tailwindcss').Config} */
export default {
  content: ["./index.html", "./src/**/*.{svelte,js,ts,jsx,tsx}"],
  theme: {
    extend: {
      colors: {
        base: {
          900: "#111827",
          800: "#141c28",
          700: "#1b2636",
          600: "#1e293b",
        },
        accent: {
          DEFAULT: "#50a0ff",
          dim: "rgba(80,160,255,0.25)",
          hover: "rgba(80,160,255,0.35)",
          border: "rgba(80,160,255,0.4)",
          glow: "rgba(80,160,255,0.15)",
        },
        surface: {
          DEFAULT: "rgba(255,255,255,0.08)",
          hover: "rgba(255,255,255,0.14)",
          dim: "rgba(255,255,255,0.03)",
          border: "rgba(255,255,255,0.1)",
        },
        fav: "rgba(255,200,50,0.9)",
        success: "rgba(80,200,120,0.6)",
      },
      fontFamily: {
        sans: [
          "Nunito",
          "-apple-system",
          "BlinkMacSystemFont",
          "Segoe UI",
          "Roboto",
          "sans-serif",
        ],
      },
    },
  },
  plugins: [],
};
