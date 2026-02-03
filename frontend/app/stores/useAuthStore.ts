import { defineStore } from "pinia";
import axios from "axios";
import type {
  IUser,
  ILoginCredentials,
  IAuthResponse,
  IRefreshResponse,
} from "@/types/auth";

export const useAuthStore = defineStore("auth", () => {
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
    navigateTo("/login");
  }

  async function login(credentials: ILoginCredentials): Promise<IAuthResponse> {
    const config = useRuntimeConfig();
    try {
      // Use raw axios to avoid circular dependency with useApi
      const response = await axios.post<IAuthResponse>(
        `${config.public.apiBase}/auth/login`,
        credentials,
      );

      const { access_token, refresh_token, user: userData } = response.data;

      setTokens(access_token, refresh_token);
      setUser(userData);

      return response.data;
    } catch (error) {
      console.error("[AuthStore] Login failed:", error);
      throw error;
    }
  }

  // Helper action if we need to manually refresh tokens (e.g. on app init if we only have refresh token)
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
    login,
    logout,
    refresh,
    setTokens,
    setUser,
  };
});
