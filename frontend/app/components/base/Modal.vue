<script setup lang="ts">
import type { IModalProps } from '@/types/base_component'

defineProps<IModalProps>()
const isOpen = defineModel<boolean>()
</script>

<template>
  <UModal v-model="isOpen" :prevent-close="preventClose">
    <UCard v-if="isOpen" :ui="{
      root: 'ring-0 divide-y divide-gray-100 dark:divide-gray-800',
      header: 'padding-x-4 padding-y-4 sm:px-6',
      body: 'p-0 sm:p-0',
      footer: 'padding-x-4 padding-y-4 sm:px-6'
    }">

      <template #header>
        <div class="flex items-center justify-between">
          <div>
            <h3 v-if="title" class="text-base font-semibold leading-6 text-gray-900 dark:text-white">
              {{ title }}
            </h3>
            <p v-if="description" class="mt-1 text-sm text-gray-500">
              {{ description }}
            </p>
          </div>
          <UButton color="primary" variant="ghost" icon="i-heroicons-x-mark-20-solid" class="-my-1"
            @click="isOpen = false" />
        </div>
      </template>

      <div class="p-4">
        <slot />
      </div>

      <template v-if="$slots.footer" #footer>
        <div class="flex justify-end gap-2">
          <slot name="footer" />
        </div>
      </template>

    </UCard>
  </UModal>
</template>