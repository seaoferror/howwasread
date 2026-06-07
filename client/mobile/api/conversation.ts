import { axiosInstance } from "@/api/axios";
import {
  BanParticipantRequest,
  CreateOfflineConversationRequest,
  CreateOnlineConversationRequest,
  OfflineConversationMapResponse,
  OnlineConversationFeedResponse,
} from "@/types/conversation";
import {
  extractCoordsFromURL,
  extractPlaceNameFromPathVariable,
  extractPlaceNameFromQueryParam,
} from "@/util/url";

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

export async function resolveLocation(
  googleMapsLink: string,
): Promise<{ lat: number; lng: number; placeName: string }> {
  const isLongURL = /https?:\/\/(www\.)?google\.[a-z.]+\/maps/i.test(
    googleMapsLink,
  );
  if (isLongURL) {
    const placeName = extractPlaceNameFromPathVariable(googleMapsLink);
    const coords = extractCoordsFromURL(googleMapsLink);
    if (!coords) {
      throw new Error("invalid URL");
    }
    return {
      lat: coords.lat,
      lng: coords.lng,
      placeName: placeName,
    };
  }
  const response = await fetch(googleMapsLink, {
    method: "GET",
    redirect: "manual",
  });
  console.log(response.url);
  const placeName = extractPlaceNameFromQueryParam(response.url);
  const gu = new URL(googleMapsLink);
  if (gu.searchParams.get("g_st") === "ic") {
    const coords = await captureCoordinatesFromNetwork(googleMapsLink);
    if (!coords) {
      throw new Error("invalid URL");
    }
    return { lat: coords.lat, lng: coords.lng, placeName: placeName };
  }
  const coords = extractCoordsFromURL(response.url);
  if (!coords) {
    throw new Error("invalid URL");
  }
  return { lat: coords.lat, lng: coords.lng, placeName: placeName };
}