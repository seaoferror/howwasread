export interface NotificationMessage {
  tokenMap: Record<string, string[]>;
  roomName?: string;
  senderName: string;
  text: string;
  imageURL?: string;
}
