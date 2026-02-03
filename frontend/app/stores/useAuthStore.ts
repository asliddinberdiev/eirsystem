import { defineStore } from "pinia";
import axios from "axios";
import type {
  IUser,
  ISignInCredentials,
  ISignInResponse,
  IRefreshResponse,
} from "@/types/auth";
import type { Response } from "@/types/api";

export const useAuthStore = defineStore("auth", () => {
  const api = useApi();
  const user = ref<IUser | null>(null);
  const accessToken = ref<string | null>(null);
  const refreshToken = ref<string | null>(null);

  const isLoggedIn = computed(() => !!accessToken.value);

  function setTokens(access: string, refresh: string) {
    accessToken.value = access;
    refreshToken.value = refresh;
  }

  function setUser(userData: IUser) {
    user.value = userData;
  }

  function logout() {
    user.value = null;
    accessToken.value = null;
    refreshToken.value = null;
    navigateTo("/sign-in");
  }

  async function signIn(
    credentials: ISignInCredentials,
  ): Promise<Response<ISignInResponse>> {
    const config = useRuntimeConfig();
    try {
      const response = await api.post<Response<ISignInResponse>>(
        `${config.public.apiBase}/auth/sign-in`,
        credentials,
      );

      const {
        access_token,
        refresh_token,
        user: userData,
      } = response.data.data;

      setTokens(access_token, refresh_token);
      setUser(userData);

      return response.data;
    } catch (error) {
      console.error("[AuthStore] Sign in failed:", error);
      throw error;
    }
  }

  async function refresh(): Promise<void> {
    const config = useRuntimeConfig();
    if (!refreshToken.value) {
      console.warn("[AuthStore] No refresh token available");
      return;
    }

    try {
      const response = await axios.post<IRefreshResponse>(
        `${config.public.apiBase}/auth/refresh`,
        {
          refresh_token: refreshToken.value,
        },
      );
      const { access_token, refresh_token: newRefresh } = response.data;
      setTokens(access_token, newRefresh);
    } catch (error) {
      console.error("[AuthStore] Token refresh failed:", error);
      logout();
      throw error;
    }
  }

  return {
    user,
    accessToken,
    refreshToken,
    isLoggedIn,
    signIn,
    logout,
    refresh,
    setTokens,
    setUser,
  };
});
