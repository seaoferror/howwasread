import { axiosInstance } from "@/api/axios";
import {
  ConversationFeedResponse,
  CreateOnlineConversationRequest,
  GetConversationResponse,
} from "@/types/conversation";

export async function getOnlineConversations(
  page = 1,
): Promise<ConversationFeedResponse[]> {
  const { data } = await axiosInstance.get(
    `/online/conversation/list?page=${page}`,
  );
  console.log(data);
  return data;
}

export async function getOnlineConversation(
  id: string,
): Promise<GetConversationResponse> {
  const { data } = await axiosInstance.get(`/online/conversation?id=${id}`);
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
