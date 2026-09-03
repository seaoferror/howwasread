export interface NotificationMessage {
  tokenMap: Record<string, string>;
  title: string;
  subTitle?: string;
  text: string;
  imageURL?: string;
}
