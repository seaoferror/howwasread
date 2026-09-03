import type Expo from "expo-server-sdk";
import { makeMessage } from "@/util";

export async function sendPushNotification(
  expo: Expo,
  markNotification: () => Promise<number>,
  tokenMap: Record<string, string>,
  title: string,
  subTitle: string | undefined,
  text: string,
  imageURL: string | undefined,
): Promise<Record<string, string>> {
  const message = makeMessage(title, subTitle, text, imageURL);
  const tickets = await expo.sendPushNotificationsAsync(
    Object.keys(tokenMap).map((token) => {
      message.to = token;
      return message;
    }),
  );
  await markNotification();
  const failedTokenMap: Record<string, string> = {};
  tickets.flatMap((ticket) => {
    if (
      ticket.status != "ok" &&
      ticket.details?.error === "DeviceNotRegistered"
    ) {
      if (ticket.details.expoPushToken) {
        const failedToken = ticket.details.expoPushToken;
        failedTokenMap[failedToken] = tokenMap[failedToken];
      }
    }
  });
  return failedTokenMap;
}
