import { axiosInstance } from "@/api/axios";
import {
  BanParticipantRequest,
  CreateOfflineConversationRequest,
  CreateOnlineConversationRequest,
  OfflineConversationDetailResponse,
  OfflineConversationMapResponse,
  OfflineConversationSearchResponse,
  OnlineConversationDetailResponse,
  OnlineConversationFeedResponse,
} from "@/types/conversation";

export async function getOnlineConversations(
  page = 1,
): Promise<OnlineConversationFeedResponse[]> {
  const { data } = await axiosInstance.get(`/onlineconversation/list`, {
    params: {
      page,
      time: new Date().toISOString(),
    },
  });
  return data;
}

export async function getOnlineConversationDetail(
  id: string,
): Promise<OnlineConversationDetailResponse> {
  const { data } = await axiosInstance.get(`/onlineconversation/detail`, {
    params: {
      id: id,
    },
  });
  return data;
}

export async function registerOnlineConversation(body: { id: string }) {
  const { data } = await axiosInstance.post(
    `/onlineconversation/register`,
    body,
  );
  return data;
}

export async function deregisterOnlineConversation(body: { id: string }) {
  const { data } = await axiosInstance.post(
    `/onlineconversation/deregister`,
    body,
  );
  return data;
}

export async function createOnlineConversation(
  body: CreateOnlineConversationRequest,
): Promise<{ id: string }> {
  const { data } = await axiosInstance.post("/onlineconversation/create", body);
  return data;
}

export async function banParticipant(body: BanParticipantRequest) {
  const { data } = await axiosInstance.post("/onlineconversation/ban", body);
  return data;
}

export async function mapOfflineConversations({
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
      time: new Date().toISOString(),
    },
  });
  return data;
}

export async function searchOfflineConversations({
  input,
  resolution,
  h3Indexes,
  page = 1,
}: {
  input: string;
  resolution: number;
  h3Indexes: string[];
  page: number;
}): Promise<OfflineConversationSearchResponse[]> {
  const { data } = await axiosInstance.get(`/search/conversation/offline`, {
    params: {
      input,
      resolution,
      h3Indexes: h3Indexes.join(","),
      page,
      time: new Date().toISOString(),
    },
  });
  return data;
}

export async function searchOnlineConversations({
  input,
  page = 1,
}: {
  input: string;
  page: number;
}): Promise<OnlineConversationFeedResponse[]> {
  const { data } = await axiosInstance.get(`/search/conversation/online`, {
    params: {
      input,
      page,
      time: new Date().toISOString(),
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
  console.log(data);
  return data;
}

export async function joinOfflineConversation(body: {
  conversationId: string;
}) {
  const { data } = await axiosInstance.patch("/offlineconversation/join", body);
  return data;
}

export async function quitOfflineConversation(body: {
  conversationId: string;
}) {
  const { data } = await axiosInstance.patch("/offlineconversation/quit", body);
  return data;
}

export async function blockConversation(body: { id: string }) {
  const { data } = await axiosInstance.post("/chat/block/conversation", body);
  return data;
}

export async function getBlockedConversations(): Promise<{ id: string }[]> {
  const { data } = await axiosInstance.get("/chat/block/conversations");
  return data;
}

export async function reportOnlineConversation(body: { id: string }) {
  const { data } = await axiosInstance.post("/onlineconversation/report", body);
  return data;
}

export async function reportOfflineConversation(body: { id: string }) {
  const { data } = await axiosInstance.post(
    "/offlineconversation/report",
    body,
  );
  return data;
}
