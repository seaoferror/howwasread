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

    KafkaSource <IncomingNotificationEvent> source = KafkaConnectorFactory.createSource(
        "scheduled-notification",
        "scheduled-notification"
    );

    KafkaSink<OutgoingNotificationEvent> sink = KafkaConnectorFactory.createSink(
        "notification"
    );

    env.fromSource(source, WatermarkStrategy.noWatermarks(), "Kafka Source")
        .keyBy(event -> event.getPartitionId().toString())
        .process(new ScheduledNotificationFunction())
        .sinkTo(sink);

    env.execute("Notification Delay Pipeline");
  }
}
