<script setup lang="ts" generic="T extends Record<string, any> | string | number">
import type { ISelectProps } from '@/types/base_component'

const props = withDefaults(defineProps<ISelectProps>(), {
    searchable: true,
    placeholder: 'Tanlang...',
    searchPlaceholder: 'Qidirish...',
    disabled: false,
    error: undefined,
})

const selected = defineModel<T | string | number>()

const isError = computed(() => !!props.error)
</script>

<template>
    <div class="w-full">
        <label v-if="label" class="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-1">
            {{ label }}
        </label>

        <USelectMenu v-model="selected" :items="options" :searchable="searchable"
            :searchable-placeholder="searchPlaceholder" :placeholder="placeholder" :disabled="disabled"
            :loading="loading" :color="isError ? 'error' : 'primary'" class="w-full" v-bind="$attrs">
            <template v-for="(_, name) in $slots" #[name]="slotData">
                <slot :name="name" v-bind="slotData" />
            </template>
        </USelectMenu>

        <p v-if="typeof error === 'string' && error" class="mt-1 text-xs text-red-500 font-medium">
            {{ error }}
        </p>
    </div>
</template>