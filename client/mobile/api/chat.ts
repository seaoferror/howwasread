import { axiosInstance } from "@/api/axios";
import { MessagingResponse } from "@/types/chat";
import { data } from "browserslist";

export async function sendLike(body: { toId: string }) {
  const { data } = await axiosInstance.post("/chat/like", body);
  return data;
}

export async function getRecentMessages(cursor: string) : Promise<MessagingResponse[]> {
  const { data } = await axiosInstance.get(
    `/chat/messaging/recent?cursor=${cursor}`,
  );
  return data;
}
