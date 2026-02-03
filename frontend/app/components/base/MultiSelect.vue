<script setup lang="ts" generic="T extends Record<string, any> | string | number">
import type { IMultiSelectProps } from '@/types/base_component'

const props = withDefaults(defineProps<IMultiSelectProps>(), {
    searchable: true,
    placeholder: 'Tanlang...',
    optionAttribute: 'label',
    disabled: false
})

const selected = defineModel<T[] | string[] | number[]>({ default: [] })
</script>

<template>
    <div class="w-full">
        <label v-if="label" class="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-1">
            {{ label }}
        </label>

        <USelectMenu multiple v-model="selected" :items="options" :placeholder="placeholder" :searchable="searchable"
            :disabled="disabled" class="w-full">
        </USelectMenu>

        <p v-if="typeof error === 'string' && error" class="mt-1 text-xs text-red-500 font-medium">
            {{ error }}
        </p>
    </div>
</template>