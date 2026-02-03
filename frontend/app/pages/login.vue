<script setup lang="ts">
import { z } from 'zod'
import { isAxiosError } from 'axios'
import type { FormSubmitEvent } from '#ui/types'

definePageMeta({
    layout: 'auth'
})

const authStore = useAuthStore()
const toast = useToast()

const schema = z.object({
    email: z.string().email('Invalid email address'),
    password: z.string().min(6, 'Password must be at least 6 characters')
})

type Schema = z.output<typeof schema>

const state = reactive({
    email: '',
    password: ''
})

const loading = ref(false)

async function onSubmit(event: FormSubmitEvent<Schema>) {
    loading.value = true
    try {
        await authStore.login(event.data)
        toast.add({
            title: 'Welcome back!',
            description: 'You have successfully signed in.',
            icon: 'i-heroicons-check-circle'
        })
        await navigateTo('/')
    } catch (error: unknown) {
        console.error('[Login] Authentication failed:', error)

        let errorMessage = 'Invalid email or password'
        if (isAxiosError(error) && error.response?.data?.message) {
            errorMessage = error.response.data.message
        }

        toast.add({
            title: 'Authentication Failed',
            description: errorMessage,
            color: 'error',
            icon: 'i-heroicons-exclamation-circle'
        })
    } finally {
        loading.value = false
    }
}
</script>

<template>
    <UCard>
        <template #header>
            <h1 class="text-xl font-bold text-gray-900 dark:text-white">
                Sign In
            </h1>
            <p class="text-sm text-gray-500 dark:text-gray-400 mt-1">
                Access your eirsystem account
            </p>
        </template>

        <UForm :schema="schema" :state="state" class="grid grid-rows-1 gap-4 items-center" @submit="onSubmit">
            <UFormGroup label="Email" name="email">
                <UInput v-model="state.email" type="email" placeholder="you@example.com" autofocus class="w-full" />
            </UFormGroup>

            <UFormGroup label="Password" name="password">
                <UInput v-model="state.password" type="password" placeholder="••••••••" class="w-full" />
            </UFormGroup>

            <UButton type="submit" block :loading="loading">
                Sign In
            </UButton>
        </UForm>
    </UCard>
</template>
