import {
  createCassandraClient,
  createKafkaConsumer,
  createRedisClient,
} from "./config";
import { extractData } from "./util";
import { sendPushNotification } from "./service/notification";
import Expo from "expo-server-sdk";
import { removePushTokenByIdAndDeviceId } from "./db";

async function main() {
  const expo = new Expo();
  const consumer = await createKafkaConsumer();
  const redis = await createRedisClient();
  const cassandra = await createCassandraClient();

  const disconnect = () => {
    process.off("SIGINT", disconnect);
    process.off("SIGTERM", disconnect);
    consumer
      .commitOffsets()
      .finally(() => consumer.disconnect())
      .finally(() => redis.destroy())
      .finally(() => cassandra.shutdown())
      .finally(() => console.log("Disconnected successfully"));
  };
  process.on("SIGINT", disconnect);
  process.on("SIGTERM", disconnect);

  consumer.run({
    partitionsConsumedConcurrently: 6,
    eachMessage: async ({ message }) => {
      const { value, redisKey, redisMember } = extractData(message);
      const did = await redis.sIsMember(redisKey, redisMember);
      if (did) {
        return;
      }
      const failedIdPairs = await sendPushNotification(
        expo,
        () => redis.sAdd(redisKey, redisMember),
        value.tokenMap,
        value.roomName,
        value.senderName,
        value.text,
        value.imageURL,
      );
      if (failedIdPairs.length > 0) {
        await Promise.all(failedIdPairs.map((failedIdPair) => {
          return removePushTokenByIdAndDeviceId(
            cassandra,
            failedIdPair[0],
            failedIdPair[1],
          );
        }))
      }
    },
  });
}

main();

