import { KafkaMessage } from "@confluentinc/kafka-javascript/types/kafkajs";
import { NotificationMessage } from "@/types";
import type { ExpoPushMessage } from "expo-server-sdk";

export function extractData(message: KafkaMessage) {
  const value: NotificationMessage = JSON.parse(
    message?.value?.toString() ?? "",
  );
  const key = Buffer.concat([
    Buffer.from("apn"),
    message?.key?.subarray(0, 16) ?? Buffer.from(""),
  ]);
  const member = message?.key?.subarray(16, 17) ?? Buffer.from("");
  return { value, key, member};
}

export function makeMessage(
  title: string,
  subTitle: string | undefined,
  text: string,
  imageURL: string | undefined,
) {
  const message: ExpoPushMessage = {
    to: "",
    title: title,
    body: text,
    richContent: { image: imageURL },
  };
  if (subTitle) {
    message.subtitle = subTitle;
  }
  return message;
}