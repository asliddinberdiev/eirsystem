<script setup lang="ts">
definePageMeta({ layout: "auth" });

const { schema, state, onSubmit } = useAuthPage();
</script>

<template>
  <UCard>
    <template #header>
      <h1 class="text-2xl font-bold text-center">EIR SYSTEM</h1>
    </template>

    <UForm :schema="schema" :state="state" class="space-y-4" @submit="onSubmit">
      <UFormField name="username">
        <BaseInput
          v-model="state.username"
          size="xl"
          placeholder="Username kiriting"
          icon="i-lucide-user"
        />
      </UFormField>

      <UFormField name="password">
        <BaseInput
          v-model="state.password"
          size="xl"
          placeholder="Parolni kiriting"
          icon="i-lucide-lock"
          type="password"
        />
      </UFormField>
      <BaseBtn type="submit" color="primary" size="xl" block>Kirish</BaseBtn>
    </UForm>
  </UCard>
</template>
