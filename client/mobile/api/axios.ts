import { create } from "axios";
import { getSecureAsync, setSecure } from "@/util/storage";

export const axiosInstance = create({
  adapter: "fetch",
  baseURL: `https://${process.env.EXPO_PUBLIC_API_URL}`,
  withCredentials: true,
  headers: {
    "Content-Type": "application/json",
  },
});

async function refreshAccessToken(): Promise<{ accessToken: string }> {
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
          const { accessToken } = await refreshAccessToken();
          await setSecure("accessToken", accessToken);
          originalRequest.headers.Authorization = `Bearer ${accessToken}`;
          return axiosInstance(originalRequest);
        } catch (err) {
          console.error("Refresh token failed", err);
        }
      }
      return Promise.reject(error);
    }
  },
);

// export const localDevInstance = create({
//   adapter: "fetch",
//   baseURL: `http://${Platform.OS === "ios" ? baseUrl.ios : baseUrl.android}:8078`,
//   withCredentials: true,
//   headers: {
//     "Content-Type": "application/json",
//     "X-User-Id": `${Platform.OS === "ios" ? localDevId.ios : localDevId.android}`,
//   },
// });

// export const baseUrl = {
//   android: process.env.EXPO_PUBLIC_API_URL,
//   ios: process.env.EXPO_PUBLIC_API_URL,
// };
//
// export const localDevId = {
//   android: "019e0e84-f358-71a2-8b3c-d4e5f6012345",
//   ios: "019e0e84-f358-7d8e-9fa0-b1c2d3e4f506",
// };
