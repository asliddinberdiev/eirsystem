<script setup lang="ts">
import type { DropdownMenuItem } from '@nuxt/ui'

defineProps<{
  collapsed?: boolean
}>()

const roles = ref([
  { label: 'Super Admin' },
  { label: 'Tenant Owner' },
  { label: 'Admin' },
  { label: 'Doctor' },
  { label: 'Nurse' },
  { label: 'Technician' },
  { label: 'Reception' }
])
const selectedRole = ref(roles.value[0])

const items = computed<DropdownMenuItem[][]>(() => {
  return [roles.value.map(role => ({
    ...role,
    onSelect() {
      selectedRole.value = role
    }
  }))]
})
</script>

<template>
  <UDropdownMenu :items="items" :content="{ align: 'center', collisionPadding: 12 }"
    :ui="{ content: collapsed ? 'w-40' : 'w-(--reka-dropdown-menu-trigger-width)' }">
    <UButton v-bind="{
      ...selectedRole,
      label: collapsed ? undefined : selectedRole?.label,
      trailingIcon: collapsed ? undefined : 'i-lucide-chevrons-up-down'
    }" color="neutral" variant="ghost" block :square="collapsed" class="data-[state=open]:bg-elevated"
      :class="[!collapsed && 'py-2']" :ui="{
        trailingIcon: 'text-dimmed'
      }" />
  </UDropdownMenu>
</template>
