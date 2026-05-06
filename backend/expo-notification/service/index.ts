import type Expo from "expo-server-sdk";
import { makeMessage } from "../util";

export async function sendPushNotification(
  expo: Expo,
  markNotification: () => Promise<number>,
  tokenMap: Record<string, string[]>,
  roomName: string | undefined,
  senderName: string,
  text: string,
  imageURL: string | undefined,
) {
  const message = makeMessage(senderName, text, imageURL, roomName);
  const tickets = await expo.sendPushNotificationsAsync(
    Object.keys(tokenMap).map((token) => {
      message.to = token;
      return message;
    }),
  );
  await markNotification();
  return tickets.flatMap((ticket) => {
    if (
      ticket.status !== "ok" &&
      ticket.details &&
      ticket.details.error === "DeviceNotRegistered"
    ) {
      const failedToken = ticket.details.expoPushToken;
      return failedToken && tokenMap[failedToken]
        ? [tokenMap[failedToken]]
        : [];
    }
    return [];
  });
}
