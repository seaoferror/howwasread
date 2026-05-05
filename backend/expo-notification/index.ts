import Expo from "expo-server-sdk";
import {
  KafkaJS,
  LibrdKafkaError,
  type TopicPartition,
} from "@confluentinc/kafka-javascript";
import { type Consumer } from "@confluentinc/kafka-javascript/types/kafkajs";
import { NotificationMessage } from "./types/notification";
import { createClient } from "@redis/client";

const { Kafka, ErrorCodes } = KafkaJS;
const expo = new Expo();
const APN_NOTIFICATION = "apn_notification";

function sendPushNotification(
  tokenMap: Record<string, string[]>,
  roomName: string | undefined,
  senderName: string,
  text: string,
  imageURL: string | undefined,
) {}

async function consumerStart() {
  let consumer: Consumer;
  let stopped = false;

  let isPaused = false;
  const togglePauseResume = () => {
    if (stopped) return;

    const assignment = consumer.assignment();
    if (isPaused) {
      console.log(`Resuming partitions ${JSON.stringify(assignment)}`);
      consumer.resume(assignment);
      isPaused = !isPaused;
      return;
    }
    console.log(`Pausing partitions ${JSON.stringify(assignment)}`);
    consumer.pause(assignment);
    isPaused = !isPaused;
  };

  const disconnect = () => {
    process.off("SIGUSR1", togglePauseResume);
    process.off("SIGINT", disconnect);
    process.off("SIGTERM", disconnect);
    stopped = true;
    consumer
      .commitOffsets()
      .finally(() => consumer.disconnect())
      .finally(() => console.log("Disconnected successfully"));
  };

  process.on("SIGUSR1", togglePauseResume);
  process.on("SIGINT", disconnect);
  process.on("SIGTERM", disconnect);

  consumer = new Kafka().consumer({
    "bootstrap.servers": process.env.KAFKA_URL || "127.0.0.1:9092",
    "group.id": APN_NOTIFICATION,
    "auto.offset.reset": "earliest",
    "enable.partition.eof": false,
    "enable.auto.commit": false,
    "group.protocol": "consumer",
    "group.remote.assignor": "uniform",
    rebalance_cb: (err: LibrdKafkaError, assignment: TopicPartition[]) => {
      switch (err.code) {
        case ErrorCodes.ERR__ASSIGN_PARTITIONS:
          console.log(`Assigned partitions ${JSON.stringify(assignment)}`);
          break;
        case ErrorCodes.ERR__REVOKE_PARTITIONS:
          console.log(`Revoked partitions ${JSON.stringify(assignment)}`);
          break;
        default:
          console.error(err);
      }
    },
  });
  await consumer.connect();
  console.log("Connected successfully");
  await consumer.subscribe({ topics: [APN_NOTIFICATION] });

  consumer.run({
    partitionsConsumedConcurrently: 200_000,
    eachMessage: async ({ topic, partition, message }) => {
      const data: NotificationMessage = JSON.parse(
        message?.value?.toString() ?? "",
      );
      sendPushNotification(
        data.tokenMap,
        data.roomName,
        data.senderName,
        data.text,
        data.imageURL,
      );
      await consumer.commitOffsets([
        {
          topic,
          partition,
          offset: (parseInt(message.offset, 10) + 1).toString(),
        },
      ]);
    },
  });
}

consumerStart();

async function ss() {
  const ticket = await expo.sendPushNotificationsAsync([
    {
      to: "token",
      title: "test",
      body: "test",
      badge: 0,
    },
  ]);
  const receiptIds = ticket.filter((t) => t.status === "ok").map((t) => t.id);
  const chunk = expo.chunkPushNotificationReceiptIds(receiptIds);
  expo.getPushNotificationReceiptsAsync(chunk[0]);
}
