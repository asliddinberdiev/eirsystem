import type { UseFetchOptions } from '#app'

export function useApi<T>(url: string, options: UseFetchOptions<T> = {}) {
  const config = useRuntimeConfig()
  const authStore = useAuthStore()

  const defaults: UseFetchOptions<T> = {
    baseURL: config.public.apiBase as string,
    key: url,

    onRequest({ options }) {
      if (authStore.token) {
        options.headers = options.headers || {}
        // @ts-expect-error - headers type mismatch in some generic cases but valid here
        options.headers.Authorization = `Bearer ${authStore.token}`
      }
    },

    onResponseError({ response }) {
      if (response.status === 401) {
        authStore.logout()
        // Optional: redirect to login
        navigateTo('/login')
      }
    }
  }

  // Merge defaults with user options.
  // Note: deep merge might be better for complex headers, but shallow merge of top-level props + manual header handling is often sufficient.
  // We use defu or simplified spread here.
  const params = {
    ...defaults,
    ...options,
    headers: {
      ...defaults.headers,
      ...options.headers,
    }
  }

  return useFetch(url, params)
}
