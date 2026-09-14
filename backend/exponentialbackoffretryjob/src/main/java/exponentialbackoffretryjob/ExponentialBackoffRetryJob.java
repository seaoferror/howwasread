package exponentialbackoffretryjob;

import org.apache.flink.api.common.eventtime.WatermarkStrategy;
import org.apache.flink.configuration.Configuration;
import org.apache.flink.connector.kafka.sink.KafkaSink;
import org.apache.flink.connector.kafka.source.KafkaSource;
import org.apache.flink.streaming.api.environment.StreamExecutionEnvironment;
import exponentialbackoffretryjob.config.KafkaConnectorFactory;
import exponentialbackoffretryjob.dto.IncomingEvent;
import exponentialbackoffretryjob.dto.OutgoingEvent;

public class ExponentialBackoffRetryJob {

  public static void main(String[] args) throws Exception {
    Configuration config = new Configuration();
    config.setString("state.backend.type", "rocksdb");
    config.setString("state.backend.incremental", "true");

    StreamExecutionEnvironment env = StreamExecutionEnvironment.getExecutionEnvironment(config);
    env.enableCheckpointing(10000);

    KafkaSource<IncomingEvent> source = KafkaConnectorFactory.createSource(
        "exponential-backoff-retry",
        "exponential-backoff-retry"
    );

    KafkaSink<OutgoingEvent> sink = KafkaConnectorFactory.createSink();

    env.fromSource(source, WatermarkStrategy.noWatermarks(), "Kafka Source")
        .keyBy(IncomingEvent::getPartitionId)
        .process(new ExponentialBackoffRetryFunction())
        .sinkTo(sink);


    env.execute("ExponentialBackoffRetryJob");
  }
}
