import { axiosInstance } from "@/api/axios";
import { GetProfileResponse } from "@/types/profile";

export async function getMyProfile(): Promise<GetProfileResponse> {
  const { data } = await axiosInstance.get("/chat/profile/my");
  return data;
}

export async function getProfile(id: string): Promise<GetProfileResponse> {
  const { data } = await axiosInstance.get(`/chat/profile?id=${id}`)
  return data;
}

export async function setName(body: { name: string }) {
  const { data } = await axiosInstance.put("/chat/profile/name", body);
  return data;
}

export async function setBirthYear(body: { birthYear: number }) {
  const { data } = await axiosInstance.post("/chat/profile/birth-year", body);
  return data;
}
