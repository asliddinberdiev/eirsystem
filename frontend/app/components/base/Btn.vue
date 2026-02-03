<script setup lang="ts">
import type { IBtnProps } from "@/types/base_component";

withDefaults(defineProps<IBtnProps>(), {
  type: "button",
  block: false,
  loading: false,
  disabled: false,
  variant: "solid",
  color: "primary",
});
</script>

<template>
  <UButton
    :type="type"
    :to="to"
    :loading="loading"
    :disabled="loading || disabled"
    :variant="variant"
    :color="color"
    :icon="icon"
    :leading-icon="leadingIcon"
    :trailing-icon="trailingIcon"
    :label="label"
    :block="block"
    class="transition-transform active:scale-[0.98]"
    v-bind="$attrs"
  >
    <slot />
  </UButton>
</template>
