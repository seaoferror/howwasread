import { APN_NOTIFICATION } from "@/constants";
import { createClient } from "@redis/client";
import {
  KafkaJS,
  LibrdKafkaError,
  type TopicPartition,
} from "@confluentinc/kafka-javascript";
import cassandra from "cassandra-driver";

const { Kafka, ErrorCodes } = KafkaJS;

export async function createRedisClient() {
  const redis = createClient({
    url: process.env.REDIS_ADDRESS
      ? "redis://" + process.env.REDIS_ADDRESS
      : "redis://127.0.0.1:6379",
  });

  await redis.connect();
  console.log("connect to redis");

  return redis;
}

export async function createKafkaConsumer() {
  const consumer = new Kafka().consumer({
    "bootstrap.servers":
      process.env.KAFKA_ADDRESS ||
      "kafka-cluster-kafka-bootstrap.kafka-system.svc.cluster.local:9093",
    "security.protocol": "ssl",
    "ssl.ca.location": "kafka-ca.crt",
    "ssl.certificate.location": "kafka-user.crt",
    "ssl.key.location": "kafka-user.key",

    "group.id": APN_NOTIFICATION,
    "auto.offset.reset": "earliest",
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
  console.log("connect to Kafka");
  await consumer.subscribe({ topics: [APN_NOTIFICATION] });

  let isPaused = false;
  const togglePauseResume = () => {
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

  process.on("SIGUSR1", togglePauseResume);

  return consumer;
}

export async function createCassandraClient() {
  const cloud = { secureConnectBundle: "/astradb/astradb-secure-connect.zip" };
  const authProvider = new cassandra.auth.PlainTextAuthProvider(
    "token",
    process.env["ASTRA_DB_TOKEN"] || "",
  );
  const client = new cassandra.Client({
    cloud: cloud,
    authProvider: authProvider,
    localDataCenter: process.env.CASSANDRA_DATACENTER || "datacenter1",
    keyspace: "default",
    queryOptions: {
      consistency: cassandra.types.consistencies.quorum,
    },
    socketOptions: {
      readTimeout: 60_000,
    },
  });

  await client.connect();
  console.log("connect to cassandra");

  return client;
}
