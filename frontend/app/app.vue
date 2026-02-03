<template>
  <UApp :toaster="{ position: 'top-right', expand: false }">
    <NuxtLoadingIndicator />

    <NuxtLayout>
      <NuxtPage />
    </NuxtLayout>
  </UApp>
</template>
