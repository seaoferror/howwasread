export interface NotificationMessage {
  tokenMap: Record<string, string>;
  roomName?: string;
  title: string;
  text: string;
  imageURL?: string;
}
