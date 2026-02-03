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
      apiBase:
        process.env.NUXT_PUBLIC_API_BASE || "http://localhost:8080/api/v1",
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
