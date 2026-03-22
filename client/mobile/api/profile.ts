import { axiosInstance } from "@/api/axios";
import { GetMyProfileResponse } from "@/types/profile";

export async function getMyProfile(): Promise<GetMyProfileResponse> {
  const { data } = await axiosInstance.get("/profile");
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
