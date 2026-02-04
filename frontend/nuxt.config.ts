export default defineNuxtConfig({
  compatibilityDate: "2026-02-03",
  ssr: true,
  devtools: {
    enabled: true,
  },
  css: ["~/assets/css/main.css"],

  modules: ["@nuxt/eslint", "@nuxt/ui", "@vueuse/nuxt", "@pinia/nuxt"],

  runtimeConfig: {
    public: {
      apiBase: process.env.NUXT_PUBLIC_API_BASE,
      accessExp: parseInt(process.env.NUXT_PUBLIC_ACCESS_EXP || "15"),
      refreshExp: parseInt(process.env.NUXT_PUBLIC_REFRESH_EXP || "10080"),
    },
  },

  eslint: {
    config: {
      stylistic: {
        commaDangle: "never",
        braceStyle: "1tbs",
      },
    },
  },

  vite: {
    server: {
      allowedHosts: true,
    },
  },
});
