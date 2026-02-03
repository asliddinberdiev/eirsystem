export const useAppToast = () => {
  const toast = useToast();

  function successToast(title: string, description: string) {
    toast.add({
      title,
      description,
      color: "success",
    });
  }

  function errorToast(title: string, description: string) {
    toast.add({
      title,
      description,
      color: "error",
    });
  }

  function warningToast(title: string, description: string) {
    toast.add({
      title,
      description,
      color: "warning",
    });
  }

  function infoToast(title: string, description: string) {
    toast.add({
      title,
      description,
      color: "info",
    });
  }

  return {
    toast,
    successToast,
    errorToast,
    warningToast,
    infoToast,
  };
};
