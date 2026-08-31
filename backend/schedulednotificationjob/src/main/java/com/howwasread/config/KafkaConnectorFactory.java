package com.howwasread.config;

import com.howwasread.dto.IncomingNotificationEvent;
import com.howwasread.dto.OutgoingNotificationEvent;
import com.howwasread.mapper.NotificationEventDeserializationSchema;
import com.howwasread.mapper.NotificationEventSerializationSchema;
import org.apache.flink.connector.kafka.sink.KafkaRecordSerializationSchema;
import org.apache.flink.connector.kafka.sink.KafkaSink;
import org.apache.flink.connector.kafka.source.KafkaSource;
import org.apache.flink.connector.kafka.source.enumerator.initializer.OffsetsInitializer;

public class KafkaConnectorFactory {

    private static final String BOOTSTRAP_SERVERS =
        System.getenv().getOrDefault("KAFKA_BOOTSTRAP_SERVERS", "localhost:9092");

    public static KafkaSource<IncomingNotificationEvent> createSource(String topic, String groupId) {
        // Builder is now explicitly typed to <IncomingNotificationEvent>
        return KafkaSource.<IncomingNotificationEvent>builder()
            .setBootstrapServers(BOOTSTRAP_SERVERS)
            .setTopics(topic)
            .setGroupId(groupId)
            .setStartingOffsets(OffsetsInitializer.earliest())
            // Use a custom schema instead of SimpleStringSchema
            .setValueOnlyDeserializer(new NotificationEventDeserializationSchema())
            .build();
    }

    public static KafkaSink<OutgoingNotificationEvent> createSink(String topic) {
        // Builder is now explicitly typed to <OutgoingNotificationEvent>
        return KafkaSink.<OutgoingNotificationEvent>builder()
            .setBootstrapServers(BOOTSTRAP_SERVERS)
            .setRecordSerializer(
                KafkaRecordSerializationSchema.builder()
                    .setTopic(topic)
                    // Use a custom schema instead of SimpleStringSchema
                    .setValueSerializationSchema(new NotificationEventSerializationSchema())
                    .build()
            )
            .build();
    }
}