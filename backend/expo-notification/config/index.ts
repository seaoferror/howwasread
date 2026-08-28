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
  const ca = process.env.VALKEY_CA_CERT_PATH;
  let caCertBuffer;
  if(ca) {  
    caCertBuffer = readFileSync(ca);
  }
  return GlideClusterClient.createClient({
    addresses: [
      {
        host: process.env.VALKEY_HOST || "localhost",
        port: Number(process.env.VALKEY_PORT) || 6379,
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
    "bootstrap.servers": process.env.KAFKA_ADDRESS,
    "security.protocol": process.env.KAFKA_API_KEY ? "sasl_ssl" : "ssl",
    "sasl.mechanism": process.env.KAFKA_API_KEY ? "PLAIN" : undefined,
    "sasl.username": process.env.KAFKA_API_KEY || undefined,
    "sasl.password": process.env.KAFKA_API_SECRET || undefined,
    "ssl.ca.location": process.env.KAFKA_CA_CERT_PATH || undefined,
    "ssl.certificate.location": process.env.KAFKA_USER_CERT_PATH || undefined,
    "ssl.key.location": process.env.KAFKA_USER_KEY_PATH || undefined,

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
    process.env.K8SSANDRA_USERNAME ?? "",
    process.env.K8SSANDRA_PASSWORD ?? "",
  );

  const client = new cassandra.Client({
    contactPoints: [
      process.env.K8SSANDRA_HOST ??
        "k8ssandra-cluster-dc1-service.k8ssandra.svc.cluster.local",
    ],
    localDataCenter: "dc1",
    authProvider: authProvider,
    keyspace: process.env.PROFILE,
    sslOptions: {
      rejectUnauthorized: false,
    },
    queryOptions: {
      consistency: cassandra.types.consistencies.localOne,
    },
    socketOptions: {
      readTimeout: 60_000,
    },
  });
  await client.connect();
  console.log("connect to cassandra");

  return client;
}
