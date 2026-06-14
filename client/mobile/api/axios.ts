import { create } from "axios";

export const axiosInstance = create({
  adapter: "fetch",
  baseURL: `https://${process.env.EXPO_BASE_URL}`,
  withCredentials: true,
  headers: {
    "Content-Type": "application/json",
  },
});

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
