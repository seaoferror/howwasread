import { axiosInstance } from "@/api/axios";
import {
  BanParticipantRequest,
  CreateOnlineConversationRequest,
  OnlineConversationFeedResponse,
} from "@/types/conversation";

export async function getOnlineConversations(
  page = 1,
): Promise<OnlineConversationFeedResponse[]> {
  const { data } = await axiosInstance.get(
    `/online/conversation/list?page=${page}`,
  );
  console.log(data);
  return data;
}

export async function createOnlineConversation(
  body: CreateOnlineConversationRequest,
): Promise<{ id: string }> {
  const { data } = await axiosInstance.post(
    "/online/conversation/create",
    body,
  );
  return data;
}

export async function banParticipant(body: BanParticipantRequest) {
  const { data } = await axiosInstance.post("/online/conversation/ban", body);
  return data;
}
