import { z } from "zod";
import type { FormSubmitEvent } from "@nuxt/ui";

export const useAuthPage = () => {
  const authStore = useAuthStore();
  const { successToast, errorToast } = useAppToast();

  const schema = z.object({
    username: z
      .string("Username kiritilmadi!")
      .min(3, "Eng kamida 3 ta belgi bo'lishi kerak")
      .max(50, "Eng ko'pi 50 ta belgi bo'lishi kerak"),
    password: z
      .string("Password kiritilmadi!")
      .min(6, "Eng kamida 6 ta belgi bo'lishi kerak"),
  });

  type Schema = z.output<typeof schema>;

  const state = reactive<Partial<Schema>>({
    username: undefined,
    password: undefined,
  });

  async function onSubmit(event: FormSubmitEvent<Schema>) {
    try {
      await authStore.login(event.data);
      successToast("Kirish muvaffaqiyatli amalga oshirildi", "");
    } catch (error) {
      errorToast("Kirishda xatolik", "Username yoki password noto'g'ri");
    }
  }

  return {
    schema,
    state,
    onSubmit,
  };
};
