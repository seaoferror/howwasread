import { axiosInstance } from "@/api/axios";
import { GetChatRoomInfoResponse, MessagingResponse } from "@/types/chat";

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

export async function getChatRoomInfo(roomId: string) : Promise<GetChatRoomInfoResponse> {
  const { data } = await axiosInstance.get(`/chat/room/info?id=${roomId}`)
  return data;
}