<script setup lang="ts">
import type { IInputProps } from "@/types/base_component";

const props = withDefaults(defineProps<IInputProps>(), {
  type: "text",
  placeholder: "",
  disabled: false,
  required: false,
  loading: false,
});

const modelValue = defineModel<string | number>();

const showPassword: Ref<boolean> = ref(false);
</script>

<template>
  <UInput
    v-model="modelValue"
    :type="type === 'password' ? (showPassword ? 'text' : 'password') : type"
    :placeholder="placeholder"
    :disabled="loading || disabled"
    :loading="loading"
    :icon="icon"
    :ui="ui"
    v-bind="$attrs"
    class="w-full"
  >
    <template v-for="(_, name) in $slots" #[name]="slotData">
      <slot :name="name" v-bind="slotData" />
    </template>
    <template v-if="type === 'password'" #trailing>
      <UButton
        color="neutral"
        variant="link"
        :icon="showPassword ? 'i-lucide-eye-off' : 'i-lucide-eye'"
        :aria-label="showPassword ? 'Hide password' : 'Show password'"
        :aria-pressed="showPassword"
        aria-controls="password"
        class="hover:cursor-pointer"
        @click="showPassword = !showPassword"
      />
    </template>
  </UInput>
</template>
