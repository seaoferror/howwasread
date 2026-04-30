export interface MessagingResponse {
  id: string;
  roomId: string;
  fromId: string;
  contentType: string;
  contents: string[];
}

export interface Message extends MessagingResponse {
  createdAt: string;
}

export interface GetChatRoomInfoResponse {
  name: string;
  type: string;
}

export interface MessageEntity {
  id: Uint8Array;
  room_id: Uint8Array;
  from_id: Uint8Array;
  content_type: string;
  contents: string; //TEXT type column is used with JSON.stringify
  created_at: string;
}

export interface SendMessagingRequest {
  toIdType: string;
  toId: string;
  contentType: string;
  contents: string[];
}

export interface GeneratePresignedURLResponse {
  filename: string;
  url: string;
  fields: Record<string, string>;
}

export interface UploadToS3Request {
  awsPresignedURL: string;
  awsFields: Record<string, string>;
  mimeType: string;
  localFileURI: string;
  filename: string;
}
