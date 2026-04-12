export interface MessagingResponse {
  id: string;
  roomId: string;
  fromId: string;
  contentType: string;
  content: string;
}

export interface ChatPreview extends MessagingResponse {
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
  content: string;
  created_at: string;
}
