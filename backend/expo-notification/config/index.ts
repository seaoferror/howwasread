import { GlideClusterClient } from "@valkey/valkey-glide";
import {
  KafkaJS,
  LibrdKafkaError,
  type TopicPartition,
} from "@confluentinc/kafka-javascript";
import cassandra from "cassandra-driver";
import { readFileSync } from "node:fs";

const { Kafka, ErrorCodes } = KafkaJS;

export async function createValkeyClient() {
  const caCertBuffer = readFileSync("/cert/valkey/ca.crt");
  return GlideClusterClient.createClient({
    addresses: [
      {
        host: process.env.VALKEY_HOST || "localhost",
        port: 6379,
      },
    ],
    credentials:
      process.env.PROFILE === "production"
        ? {
            username: process.env.VALKEY_USERNAME,
            password: process.env.VALKEY_PASSWORD || "",
          }
        : undefined,
    useTLS: process.env.PROFILE === "production" ? true : undefined,
    advancedConfiguration: {
      tlsAdvancedConfiguration: {
        rootCertificates: caCertBuffer
      }
    },
    compression: {
      enabled: true,
    },
  });
}

export async function createKafkaConsumer() {
  const consumer = new Kafka().consumer({
    "bootstrap.servers": process.env.KAFKA_ADDRESS || "localhost:9092",
    "security.protocol":
      process.env.PROFILE === "production" ? "ssl" : undefined,
    "ssl.ca.location":
      process.env.PROFILE === "production" ? "/cert/kafka/cluster/ca.crt" : undefined,
    "ssl.certificate.location":
      process.env.PROFILE === "production" ? "/cert/kafka/user/user.crt" : undefined,
    "ssl.key.location":
      process.env.PROFILE === "production" ? "/cert/kafka/user/user.key" : undefined,

    "group.id": "apn_notification",
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
  await consumer.subscribe({ topics: ["apn-notification"] });

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
  const authProvider = new cassandra.auth.PlainTextAuthProvider(
    "token",
    process.env.ASTRADB_TOKEN || "",
  );
  const client = new cassandra.Client({
    cloud: {
      secureConnectBundle: "/cert/astradb/astradb-secure-connect.zip",
    },
    authProvider: authProvider,
    localDataCenter: process.env.PROFILE ? undefined : "datacenter1",
    keyspace: process.env.PROFILE,
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
