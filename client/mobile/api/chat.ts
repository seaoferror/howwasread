import { axiosInstance } from "@/api/axios";
import {
  GeneratePresignedURLResponse,
  GetChatRoomInfoResponse,
  MessagingResponse,
  SendMessagingRequest,
  UploadToS3Request,
} from "@/types/chat";
import axios from "axios";

export async function sendLike(body: { toId: string }) {
  const { data } = await axiosInstance.post("/chat/like", body);
  return data;
}

export async function getRecentMessages(
  cursor: string,
): Promise<MessagingResponse[]> {
  const { data } = await axiosInstance.get(
    `/chat/messaging/recent?cursor=${cursor}`,
  );
  console.log(data);
  return data;
}

export async function getChatRoomInfo(
  roomId: string,
): Promise<GetChatRoomInfoResponse> {
  const { data } = await axiosInstance.get(`/chat/room/info?id=${roomId}`);
  return data;
}

export async function sendMessaging(
  body: SendMessagingRequest,
): Promise<{ id: string }> {
  const { data } = await axiosInstance.post("/chat/messaging/send", body);
  return data;
}

export async function generatePresignedURL(body: {
  contentType: string;
}): Promise<GeneratePresignedURLResponse> {
  const { data } = await axiosInstance.post(`/chat/messaging/presigned`, body);
  return data;
}

export async function uploadToS3({
  awsPresignedURL,
  awsFields,
  localFileURI,
  mimeType,
}: UploadToS3Request) {
  const formData = new FormData();
  Object.keys(awsFields).forEach((key) => {
    formData.append(key, awsFields[key]);
  });
  formData.append("file", {
    uri: localFileURI,
    name: "lol",
    type: mimeType,
  } as any);
  const { data } = await axios.post(awsPresignedURL, formData, {
    headers: {
      "Content-Type": "multipart/form-data",
    },
  });
  return data;
}

export async function getSignedURL({
  contentType,
  filename,
}: {
  contentType: string;
  filename: string;
}): Promise<{ url: string }> {
  const { data } = await axiosInstance.get(`/chat/messaging/signed`, {
    params: {
      contentType, // equivalent to contentType: contentType
      filename, // equivalent to filename: filename
    },
  });
  return data;
}
