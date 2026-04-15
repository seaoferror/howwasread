import { getSecureAsync, setSecure } from "@/util/storage";
import { axiosInstance } from "@/api/axios";

async function refreshAccessToken() {
  try {
    const { data } = await axiosInstance.post("/auth/refresh-token");
    return data;
  } catch (err: any) {
    const message =
      err.response?.data?.message || "Failed to refresh access token";
    throw new Error(message);
  }
}


axiosInstance.interceptors.request.use(async (config) => {
  const token = await getSecureAsync("accessToken");
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }

  return config;
});

axiosInstance.interceptors.response.use(
  (res) => res,
  async (error) => {
    if (error.response?.status) {
      const originalRequest = error.config;

      if (
        error.response?.status === 401 &&
        !originalRequest._retry &&
        !originalRequest.url.includes("/refresh-token")
      ) {
        originalRequest._retry = true;
        try {
          const { newToken } = await refreshAccessToken();
          await setSecure("accessToken", newToken);
          originalRequest.headers.Authorization = `Bearer ${newToken}`;
          return axiosInstance(originalRequest);
        } catch (err) {
          console.error("Refresh token failed", err);
        }
      }
      return Promise.reject(error);
    }
  },
);
