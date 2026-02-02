<script setup lang="ts">
const { locale, setLocale, locales } = useI18n()

// Computed property to format locales for Nuxt UI Select or Dropdown
// Assuming we want a simple select for now, or we can build a custom one.
// We'll use a simple select for simplicity and robustness.
const availableLocales = computed(() => {
  return (locales.value as any[]).map(l => ({
    label: l.name || l.code,
    value: l.code
  }))
})

function onLocaleChange(newLocale: string) {
  setLocale(newLocale as any)
}
</script>

<template>
  <div class="flex items-center gap-2">
    <!-- Using standard select for maximum compatibility if USelect is tricky, 
         but trying USelect first if available in the project -->
    <select
      :value="locale"
      @change="onLocaleChange(($event.target as HTMLSelectElement).value)"
      class="bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded px-3 py-1 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500"
    >
      <option v-for="l in availableLocales" :key="l.value" :value="l.value">
        {{ l.label }}
      </option>
    </select>
  </div>
</template>
