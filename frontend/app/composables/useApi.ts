import axios from "axios";
import type { InternalAxiosRequestConfig, AxiosError } from "axios";

// Mutex for token refresh
let isRefreshing = false;
let failedQueue: Array<{
  resolve: (token: string) => void;
  reject: (error: any) => void;
}> = [];

const processQueue = (error: any, token: string | null = null) => {
  failedQueue.forEach((prom) => {
    if (error) {
      prom.reject(error);
    } else {
      prom.resolve(token!);
    }
  });
  failedQueue = [];
};

export const useApi = () => {
  const config = useRuntimeConfig();
  const authStore = useAuthStore();

  const api = axios.create({
    baseURL: config.public.apiBase,
    headers: {
      "Content-Type": "application/json",
    },
  });

  // Request Interceptor
  api.interceptors.request.use(
    (config: InternalAxiosRequestConfig) => {
      const token = authStore.accessToken;
      if (token) {
        config.headers.Authorization = `Bearer ${token}`;
      }
      return config;
    },
    (error) => {
      return Promise.reject(error);
    },
  );

  // Response Interceptor
  api.interceptors.response.use(
    (response) => response,
    async (error: AxiosError) => {
      const originalRequest = error.config as InternalAxiosRequestConfig & {
        _retry?: boolean;
      };

      if (error.response?.status === 401 && !originalRequest._retry) {
        if (isRefreshing) {
          return new Promise<string>((resolve, reject) => {
            failedQueue.push({ resolve, reject });
          })
            .then((token) => {
              originalRequest.headers.Authorization = `Bearer ${token}`;
              return api(originalRequest);
            })
            .catch((err) => {
              return Promise.reject(err);
            });
        }

        originalRequest._retry = true;
        isRefreshing = true;

        try {
          // Attempt to refresh token
          // Note: We avoid using 'api' instance here to prevent infinite loops if the refresh fails with 401 again
          // unless the endpoint is public? Usually refresh endpoint is accessible with refresh token.
          // We assume authStore.refreshToken is available.

          if (!authStore.refreshToken) {
            throw new Error("No refresh token available");
          }

          /* 
             Using a fresh axios call for refresh to avoid interceptor issues, 
             or we could use the same instance if we ensure the refresh URL is excluded from the 401 check,
             but safely we use a raw call or helper.
          */
          const response = await axios.post(
            `${config.public.apiBase}/auth/refresh`,
            {
              refresh_token: authStore.refreshToken,
            },
          );

          const { access_token, refresh_token } = response.data;

          // Update store
          authStore.setTokens(access_token, refresh_token);

          // Process queued requests
          processQueue(null, access_token);

          // Retry original request
          originalRequest.headers.Authorization = `Bearer ${access_token}`;
          return api(originalRequest);
        } catch (err) {
          console.error("[useApi] Token refresh failed:", err);
          processQueue(err, null);
          authStore.logout();
          return Promise.reject(err);
        } finally {
          isRefreshing = false;
        }
      }

      return Promise.reject(error);
    },
  );

  return api;
};
