import { axiosInstance } from "@/api/axios";
import { fetch } from "expo/fetch";
import {
  GeneratePresignedURLResponse,
  GetChatRoomInfoResponse,
  MessagingResponse,
  SendMessagingRequest,
  UploadToS3Request,
} from "@/types/chat";
import { File } from "expo-file-system";

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

export async function sendMessage(
  body: SendMessagingRequest,
): Promise<{ id: string }> {
  const { data } = await axiosInstance.post("/chat/messaging/send", body);
  return data;
}

export async function generatePresignedURL(body: {
  contentType: string;
  n: number;
}): Promise<GeneratePresignedURLResponse[]> {
  const { data } = await axiosInstance.post(`/chat/messaging/presigned`, body);
  return data;
}

export async function uploadToS3({
  awsPresignedURL,
  awsFields,
  localFileURI,
  mimeType,
  filename,
}: UploadToS3Request) {
  const formData = new FormData();
  Object.keys(awsFields).forEach((key) => {
    formData.append(key, awsFields[key]);
  });
  formData.append("Content-type", mimeType);
  formData.append("file", new File(localFileURI), filename);
  const data = await fetch(awsPresignedURL, {
    method: "POST",
    body: formData,
    headers: {
      "Content-Type": "multipart/form-data",
    },
  });
  console.log(data.status);
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
      contentType,
      filename,
    },
  });
  return data;
}

export async function checkBlock(id: string): Promise<{ didBlock: boolean }> {
  const { data } = await axiosInstance.get(`/chat/block/check?id=${id}`);
  return data;
}

export async function getChatParticipants(
  roomId: string,
): Promise<{ id: string }[]> {
  const { data } = await axiosInstance.get(
    `/chat/participants?roomId=${roomId}`,
  );
  console.log(data)
  return data;
}
