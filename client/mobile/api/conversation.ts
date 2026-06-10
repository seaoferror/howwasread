import { axiosInstance } from "@/api/axios";
import {
  BanParticipantRequest,
  CreateOfflineConversationRequest,
  CreateOnlineConversationRequest,
  OfflineConversationDetailResponse,
  OfflineConversationMapResponse,
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

export async function mapOfflineConversation({
  resolution,
  h3Index,
}: {
  resolution: number;
  h3Index: string;
}): Promise<OfflineConversationMapResponse[]> {
  const { data } = await axiosInstance.get(`/offlineconversation/map`, {
    params: {
      resolution,
      h3Index,
    },
  });
  return data;
}

export async function createOfflineConversation(
  body: CreateOfflineConversationRequest,
) {
  const { data } = await axiosInstance.post(
    "/offlineconversation/create",
    body,
  );
  return data;
}

export async function getOfflineConversationDetail(
  id: string,
): Promise<OfflineConversationDetailResponse> {
  const { data } = await axiosInstance.get(
    `/offlineconversation/detail?conversationId=${id}`,
  );
  return data;
}

export async function joinOfflineConversation(body: { conversationId: string }) {
  const { data } = await axiosInstance.patch("/offlineconversation/join", body);
  return data;
}

export async function quitOfflineConversation(body: { conversationId: string }) {
  const { data } = await axiosInstance.patch("/offlineconversation/quit", body);
  return data;
}
