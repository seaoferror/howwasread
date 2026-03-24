import { axiosInstance } from "@/api/axios";

export async function sendLike(body: { toId: string }) {
  const { data } = await axiosInstance.post("/chat/like", body);
  return data;
}
