import axios from "axios";
import { Platform } from "react-native";

export const baseUrl = {
  android: process.env.EXPO_PUBLIC_API_URL,
  ios: process.env.EXPO_PUBLIC_API_URL,
};

export const localDevId = {
  android: "b7e2c1af-4f3d-4e2a-1c85-2f6b7a1e5d3c",
  ios: "b7e2c1af-4f3d-4e2a-1c85-2f6b7a1e5d4c",
};

export const axiosInstance = axios.create({
  adapter: "fetch",
  baseURL: `http://${Platform.OS === "ios" ? baseUrl.ios : baseUrl.android}:8080`,
  withCredentials: true,
  headers: {
    "Content-Type": "application/json",
    "X-User-Id": `${Platform.OS === "ios" ? localDevId.ios : localDevId.android}`,
  },
});

export const localDevInstance = axios.create({
  adapter: "fetch",
  baseURL: `http://${Platform.OS === "ios" ? baseUrl.ios : baseUrl.android}:8081`,
  withCredentials: true,
  headers: {
    "Content-Type": "application/json",
    "X-User-Id": `${Platform.OS === "ios" ? localDevId.ios : localDevId.android}`,
  },
});
