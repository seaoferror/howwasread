package com.howwasread;

import org.apache.flink.api.common.eventtime.WatermarkStrategy;
import org.apache.flink.connector.kafka.sink.KafkaSink;
import org.apache.flink.connector.kafka.source.KafkaSource;
import org.apache.flink.streaming.api.environment.StreamExecutionEnvironment;

public class NotificationDelayJob {

  public static void main(String[] args) throws Exception {
    StreamExecutionEnvironment env = StreamExecutionEnvironment.getExecutionEnvironment();
    env.enableCheckpointing(10000);

    // Initialize connectors using your clean factory
    KafkaSource<NotificationEvent> source = KafkaConnectorFactory.createSource(
        "incoming-notifications",
        "flink-notification-group"
    );

    KafkaSink<NotificationEvent> sink = KafkaConnectorFactory.createSink(
        "delayed-notifications-out"
    );

    // Build the pipeline
    env.fromSource(source, WatermarkStrategy.noWatermarks(), "Kafka Source")
        .keyBy(NotificationEvent::getUserId)
        .process(new DelayNotificationFunction())
        .sinkTo(sink);

    // Execute the job
    env.execute("Notification Delay Pipeline");
  }
}
