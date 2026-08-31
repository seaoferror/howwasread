package com.howwasread;

import com.howwasread.dto.IncomingNotificationEvent;
import com.howwasread.dto.OutgoingNotificationEvent;
import com.howwasread.config.KafkaConnectorFactory;
import org.apache.flink.api.common.eventtime.WatermarkStrategy;
import org.apache.flink.configuration.Configuration;
import org.apache.flink.connector.kafka.sink.KafkaSink;
import org.apache.flink.connector.kafka.source.KafkaSource;
import org.apache.flink.streaming.api.environment.StreamExecutionEnvironment;

public class ScheduledNotificationJob {

  public static void main(String[] args) throws Exception {
    Configuration config = new Configuration();
    config.setString("state.backend.type", "rocksdb");
    config.setString("state.backend.incremental", "true");

    StreamExecutionEnvironment env = StreamExecutionEnvironment.getExecutionEnvironment(config);
    env.enableCheckpointing(10000);

    KafkaSource<IncomingNotificationEvent> source = KafkaConnectorFactory.createSource(
        "incoming-notifications",
        "flink-notification-group"
    );

    KafkaSink<OutgoingNotificationEvent> sink = KafkaConnectorFactory.createSink(
        "delayed-notifications-out"
    );

    // Build the pipeline
    env.fromSource(source, WatermarkStrategy.noWatermarks(), "Kafka Source")
        .keyBy(event -> event.getConversationId().toString())
        .process(new ScheduledNotificationFunction())
        .sinkTo(sink);

    // Execute the job
    env.execute("Notification Delay Pipeline");
  }
}
