import plugin from "tailwindcss/plugin";

const skeletonPlugin = plugin(({ addComponents, theme }) => {
  const borderRadius = theme("borderRadius.sm");

  addComponents({
    ".skeleton-text-2xl": {
      height: theme("spacing.6"),
      marginTop: theme("spacing.1"),
      marginBottom: theme("spacing.1"),
      borderRadius,
    },
    ".skeleton-text-xl": {
      height: theme("spacing.5"),
      marginTop: theme("spacing.1"),
      marginBottom: theme("spacing.1"),
      borderRadius,
    },
    ".skeleton-text-lg": {
      height: "1.125rem",
      marginTop: 5,
      marginBottom: 5,
      borderRadius,
    },
    ".skeleton-text-base": {
      height: theme("spacing.4"),
      marginTop: theme("spacing.1"),
      marginBottom: theme("spacing.1"),
      borderRadius,
    },
    ".skeleton-text-sm": {
      height: "0.875rem",
      marginTop: 3,
      marginBottom: 3,
      borderRadius,
    },
    ".skeleton-text-xs": {
      height: theme("spacing.3"),
      marginTop: 2,
      marginBottom: 2,
      borderRadius,
    },
  });
});

/** @type {import('tailwindcss').Config} */
export default {
  content: ["./index.html", "./src/**/*.{js,ts,jsx,tsx}"],
  theme: {
    fontFamily: {
      sans: "var(--font-family)",
    },
    extend: {
      colors: {
        brand: "rgba(var(--brand-color) / <alpha-value>)",
        "dc-purple": {
          50: `rgba(var(--dc-purple-50) / <alpha-value>)`,
          100: `rgba(var(--dc-purple-100) / <alpha-value>)`,
          200: `rgba(var(--dc-purple-200) / <alpha-value>)`,
          300: `rgba(var(--dc-purple-300) / <alpha-value>)`,
          400: `rgba(var(--dc-purple-400) / <alpha-value>)`,
          500: `rgba(var(--dc-purple-500) / <alpha-value>)`,
          600: `rgba(var(--dc-purple-600) / <alpha-value>)`,
          700: `rgba(var(--dc-purple-700) / <alpha-value>)`,
          800: `rgba(var(--dc-purple-800) / <alpha-value>)`,
          900: `rgba(var(--dc-purple-900) / <alpha-value>)`,
          950: `rgba(var(--dc-purple-950) / <alpha-value>)`,
        },
        "dc-green": {
          50: `rgba(var(--dc-green-50) / <alpha-value>)`,
          100: `rgba(var(--dc-green-100) / <alpha-value>)`,
          200: `rgba(var(--dc-green-200) / <alpha-value>)`,
          300: `rgba(var(--dc-green-300) / <alpha-value>)`,
          400: `rgba(var(--dc-green-400) / <alpha-value>)`,
          500: `rgba(var(--dc-green-500) / <alpha-value>)`,
          600: `rgba(var(--dc-green-600) / <alpha-value>)`,
          700: `rgba(var(--dc-green-700) / <alpha-value>)`,
          800: `rgba(var(--dc-green-800) / <alpha-value>)`,
          900: `rgba(var(--dc-green-900) / <alpha-value>)`,
          950: `rgba(var(--dc-green-950) / <alpha-value>)`,
        },
        "dc-yellow": {
          50: `rgba(var(--dc-yellow-50) / <alpha-value>)`,
          100: `rgba(var(--dc-yellow-100) / <alpha-value>)`,
          200: `rgba(var(--dc-yellow-200) / <alpha-value>)`,
          300: `rgba(var(--dc-yellow-300) / <alpha-value>)`,
          400: `rgba(var(--dc-yellow-400) / <alpha-value>)`,
          500: `rgba(var(--dc-yellow-500) / <alpha-value>)`,
          600: `rgba(var(--dc-yellow-600) / <alpha-value>)`,
          700: `rgba(var(--dc-yellow-700) / <alpha-value>)`,
          800: `rgba(var(--dc-yellow-800) / <alpha-value>)`,
          900: `rgba(var(--dc-yellow-900) / <alpha-value>)`,
          950: `rgba(var(--dc-yellow-950) / <alpha-value>)`,
        },
        "dc-pink": {
          50: `rgba(var(--dc-pink-50) / <alpha-value>)`,
          100: `rgba(var(--dc-pink-100) / <alpha-value>)`,
          200: `rgba(var(--dc-pink-200) / <alpha-value>)`,
          300: `rgba(var(--dc-pink-300) / <alpha-value>)`,
          400: `rgba(var(--dc-pink-400) / <alpha-value>)`,
          500: `rgba(var(--dc-pink-500) / <alpha-value>)`,
          600: `rgba(var(--dc-pink-600) / <alpha-value>)`,
          700: `rgba(var(--dc-pink-700) / <alpha-value>)`,
          800: `rgba(var(--dc-pink-800) / <alpha-value>)`,
          900: `rgba(var(--dc-pink-900) / <alpha-value>)`,
          950: `rgba(var(--dc-pink-950) / <alpha-value>)`,
        },
        "dc-red": {
          50: `rgba(var(--dc-red-50) / <alpha-value>)`,
          100: `rgba(var(--dc-red-100) / <alpha-value>)`,
          200: `rgba(var(--dc-red-200) / <alpha-value>)`,
          300: `rgba(var(--dc-red-300) / <alpha-value>)`,
          400: `rgba(var(--dc-red-400) / <alpha-value>)`,
          500: `rgba(var(--dc-red-500) / <alpha-value>)`,
          600: `rgba(var(--dc-red-600) / <alpha-value>)`,
          700: `rgba(var(--dc-red-700) / <alpha-value>)`,
          800: `rgba(var(--dc-red-800) / <alpha-value>)`,
          900: `rgba(var(--dc-red-900) / <alpha-value>)`,
          950: `rgba(var(--dc-red-950) / <alpha-value>)`,
        },
        "dc-neutral": {
          50: `rgba(var(--dc-neutral-50) / <alpha-value>)`,
          100: `rgba(var(--dc-neutral-100) / <alpha-value>)`,
          200: `rgba(var(--dc-neutral-200) / <alpha-value>)`,
          300: `rgba(var(--dc-neutral-300) / <alpha-value>)`,
          400: `rgba(var(--dc-neutral-400) / <alpha-value>)`,
          500: `rgba(var(--dc-neutral-500) / <alpha-value>)`,
          600: `rgba(var(--dc-neutral-600) / <alpha-value>)`,
          700: `rgba(var(--dc-neutral-700) / <alpha-value>)`,
          800: `rgba(var(--dc-neutral-800) / <alpha-value>)`,
          850: `rgba(var(--dc-neutral-850) / <alpha-value>)`,
          900: `rgba(var(--dc-neutral-900) / <alpha-value>)`,
          950: `rgba(var(--dc-neutral-950) / <alpha-value>)`,
          1000: `rgba(var(--dc-neutral-1000) / <alpha-value>)`,
        },
      },
      borderRadius: {
        sm: "3px",
        md: "5px",
      },
    },
  },
  plugins: [skeletonPlugin],
};
