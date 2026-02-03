export const useIndexStore = defineStore("index", () => {
  const loading: Ref<boolean> = ref(false);

  return {
    loading,
  };
});
