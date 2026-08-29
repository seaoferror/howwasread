package com.howwasread;

import org.apache.flink.connector.kafka.sink.KafkaRecordSerializationSchema;
import org.apache.flink.connector.kafka.sink.KafkaSink;
import org.apache.flink.connector.kafka.source.KafkaSource;
import org.apache.flink.connector.kafka.source.enumerator.initializer.OffsetsInitializer;
import org.apache.flink.api.common.serialization.SimpleStringSchema;

/**
 * Factory for creating Kafka source and sink connectors.
 * TODO: Replace SimpleStringSchema with proper serialization/deserialization
 *       for NotificationEvent (e.g., JSON using Jackson or Avro).
 */
public class KafkaConnectorFactory {

    private static final String BOOTSTRAP_SERVERS =
            System.getenv().getOrDefault("KAFKA_BOOTSTRAP_SERVERS", "localhost:9092");

    @SuppressWarnings("unchecked")
    public static KafkaSource<NotificationEvent> createSource(String topic, String groupId) {
        // TODO: Replace SimpleStringSchema with a proper DeserializationSchema<NotificationEvent>
        return (KafkaSource<NotificationEvent>) (KafkaSource<?>) KafkaSource.<String>builder()
                .setBootstrapServers(BOOTSTRAP_SERVERS)
                .setTopics(topic)
                .setGroupId(groupId)
                .setStartingOffsets(OffsetsInitializer.earliest())
                .setValueOnlyDeserializer(new SimpleStringSchema())
                .build();
    }

    @SuppressWarnings("unchecked")
    public static KafkaSink<NotificationEvent> createSink(String topic) {
        // TODO: Replace SimpleStringSchema with a proper SerializationSchema<NotificationEvent>
        return (KafkaSink<NotificationEvent>) (KafkaSink<?>) KafkaSink.<String>builder()
                .setBootstrapServers(BOOTSTRAP_SERVERS)
                .setRecordSerializer(
                        KafkaRecordSerializationSchema.builder()
                                .setTopic(topic)
                                .setValueSerializationSchema(new SimpleStringSchema())
                                .build()
                )
                .build();
    }
}
