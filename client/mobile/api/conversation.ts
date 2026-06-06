import { axiosInstance } from "@/api/axios";
import {
  BanParticipantRequest,
  CreateOfflineConversationRequest,
  CreateOnlineConversationRequest,
  OfflineConversationMapResponse,
  OnlineConversationFeedResponse,
} from "@/types/conversation";
import { fetch } from "expo/fetch";

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
  const { data } = await axiosInstance.post("offlineconversation/create", body);
  return data;
}

export async function resolveLocation(
  googleMapsLink: string,
): Promise<{ lat: number; lng: number; placeName: string }> {
  const isLongURL = /https?:\/\/(www\.)?google\.[a-z.]+\/maps/i.test(
    googleMapsLink,
  );
  if (isLongURL) {
    const coords = extractCoords(googleMapsLink);
    const placeName = extractPlaceNameFromPathVariable(googleMapsLink);
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
    redirect: "follow",
  });
  console.log(response.url);
  const placeName = extractPlaceNameFromQueryParam(response.url);
  const htmlText = await response.text();
  const coords = parseCoordsFromHTMLText(htmlText);
  if(!coords) {
    throw new Error("invalid URL")
  }
  return { lat: coords.lat, lng: coords.lng, placeName: placeName };
}

function extractPlaceNameFromQueryParam(longUrl: string): string {
  const match = longUrl.match(/[?&]q=([^&]+)/i);
  if (match && match[1]) {
    const withSpaces = match[1].replace(/\+/g, " ");
    const cleanString = decodeURIComponent(withSpaces);
    const isRawCoordinateString = /^[-]?\d+\.\d+,\s*[-]?\d+\.\d+$/.test(
      cleanString,
    );
    if (isRawCoordinateString) {
      return "Dropped Pin";
    }
    return cleanString.split(",")[0].trim();
  }
  return "";
}

function extractPlaceNameFromPathVariable(longURL: string): string {
  const namePattern = /\/maps\/(?:place|search)\/([^/@]+)/;
  const match = longURL.match(namePattern);
  if (match && match[1]) {
    const spaceCleaned = match[1].replace(/\+/g, " ");
    const decodedName = decodeURIComponent(spaceCleaned);
    const isRawCoordinateString = /^[-]?\d+\.\d+,\s*[-]?\d+\.\d+$/.test(
      decodedName,
    );
    if (isRawCoordinateString) {
      return "Dropped Pin";
    }
    return decodedName;
  }
  return "";
}

function extractCoords(longURL: string): { lat: number; lng: number } | null {
  const atPattern = /@(-?\d+\.\d+),(-?\d+\.\d+)/;
  const atMatch = longURL.match(atPattern);

  if (atMatch && atMatch[1] && atMatch[2]) {
    return {
      lat: parseFloat(atMatch[1]),
      lng: parseFloat(atMatch[2]),
    };
  }
  return null;
}

function parseCoordsFromHTMLText(htmlText: string): {
  lat: number;
  lng: number;
} | null {
  const stateMatch = htmlText.match(
    /\[\[\[[0-9.]+,\s*(-?\d+\.\d+),\s*(-?\d+\.\d+)/,
  );
  if (stateMatch) {
    const lat = parseFloat(stateMatch[2]);
    const lng = parseFloat(stateMatch[1]);
    return { lat, lng };
  }
  return null
}