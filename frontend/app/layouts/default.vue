<template>
  <div class="min-h-screen flex flex-col bg-white dark:bg-gray-900 text-gray-900 dark:text-gray-100 transition-colors">
    <header class="border-b border-gray-200 dark:border-gray-800 p-4 sticky top-0 bg-white/75 dark:bg-gray-900/75 backdrop-blur z-50">
      <div class="container mx-auto flex items-center justify-between">
        <h1 class="text-xl font-bold text-primary-500">eirsystem</h1>
        <nav class="flex items-center gap-4">
          <NuxtLink to="/" class="hover:text-primary-500">Home</NuxtLink>
           <!-- Add more links here -->
          <ColorModeToggle />
        </nav>
      </div>
    </header>

    <main class="flex-grow container mx-auto p-4">
      <slot />
    </main>

    <footer class="border-t border-gray-200 dark:border-gray-800 p-4 text-center text-sm text-gray-500">
      &copy; {{ new Date().getFullYear() }} eirsystem. All rights reserved.
    </footer>
  </div>
</template>
