<script setup lang="ts">
import type { IBreadcrumbProps } from '@/types/base_component' 

defineProps<IBreadcrumbProps>()
</script>

<template>
  <UBreadcrumb 
    :items="items" 
    :ui="{ list: 'gap-x-2', item: 'text-sm' }" 
    v-bind="$attrs"
  />
</template>