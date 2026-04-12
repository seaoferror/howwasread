import { axiosInstance } from "@/api/axios";
import { GetProfileResponse } from "@/types/profile";

export async function getMyProfile(): Promise<GetProfileResponse> {
  const { data } = await axiosInstance.get("/profile/my");
  return data;
}

export async function getProfile(id: string): Promise<GetProfileResponse> {
  const { data } = await axiosInstance.get(`/profile?id=${id}`)
  return data;
}

export async function setName(body: { name: string }) {
  const { data } = await axiosInstance.put("/profile/name", body);
  return data;
}

export async function setBirthYear(body: { birthYear: number }) {
  const { data } = await axiosInstance.post("/profile/birth-year", body);
  return data;
}
