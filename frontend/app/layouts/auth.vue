<script setup lang="ts">
/*
  Auth Layout
  - Centered minimalist design
  - Gray background
  - Optimized for login/signup pages
*/
</script>

<template>
    <div class="min-h-screen flex items-center justify-center bg-gray-50 dark:bg-gray-950 p-4">
        <div class="w-full max-w-md">
            <slot />
        </div>
    </div>
</template>
