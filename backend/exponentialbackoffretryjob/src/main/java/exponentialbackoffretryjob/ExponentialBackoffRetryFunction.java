package exponentialbackoffretryjob;

import com.fasterxml.jackson.databind.ObjectMapper;
import org.apache.flink.api.common.functions.OpenContext;
import org.apache.flink.api.common.state.ListState;
import org.apache.flink.api.common.state.ListStateDescriptor;
import org.apache.flink.api.common.state.ValueState;
import org.apache.flink.api.common.state.ValueStateDescriptor;
import org.apache.flink.streaming.api.functions.KeyedProcessFunction;
import org.apache.flink.util.Collector;
import exponentialbackoffretryjob.dto.IncomingEvent;
import exponentialbackoffretryjob.dto.OutgoingEvent;

import java.io.Serial;
import java.time.Instant;

public class ExponentialBackoffRetryFunction extends KeyedProcessFunction<String, IncomingEvent, OutgoingEvent> {

  @Serial
  private static final long serialVersionUID = 1L;

  private transient ListState<String> reasons;
  private transient ValueState<Long> backoff;
  private transient ValueState<Long> multiplier;
  private transient ValueState<Long> cap;
  private transient ValueState<Integer> currentFailure;
  private transient ValueState<Integer> maxFailure;
  private transient ValueState<String> topic;
  private transient ValueState<String> originalTopic;
  private transient ValueState<String> type;
  private transient ValueState<byte[]> value;
  private transient ObjectMapper objectMapper;

  @Override
  public void open(OpenContext openContext) throws Exception {
    ListStateDescriptor<String> reasonsDescriptor =
        new ListStateDescriptor<>("reasons", String.class);
    reasons = getRuntimeContext().getListState(reasonsDescriptor);

    ValueStateDescriptor<Long> backoffDescriptor =
        new ValueStateDescriptor<>("backoff", Long.class);
    backoff = getRuntimeContext().getState(backoffDescriptor);

    ValueStateDescriptor<Long> multiplierDescriptor =
        new ValueStateDescriptor<>("multiplier", Long.class);
    multiplier = getRuntimeContext().getState(multiplierDescriptor);

    ValueStateDescriptor<Long> capDescriptor =
        new ValueStateDescriptor<>("cap", Long.class);
    cap = getRuntimeContext().getState(capDescriptor);

    ValueStateDescriptor<Integer> currentFailureDescriptor =
        new ValueStateDescriptor<>("current-failure", Integer.class);
    currentFailure = getRuntimeContext().getState(currentFailureDescriptor);

    ValueStateDescriptor<Integer> maxFailureDescriptor =
        new ValueStateDescriptor<>("max-failure", Integer.class);
    maxFailure = getRuntimeContext().getState(maxFailureDescriptor);

    ValueStateDescriptor<String> topicDescriptor =
        new ValueStateDescriptor<>("topic", String.class);
    topic = getRuntimeContext().getState(topicDescriptor);

    ValueStateDescriptor<String> originalTopicDescriptor =
        new ValueStateDescriptor<>("original-topic", String.class);
    originalTopic = getRuntimeContext().getState(originalTopicDescriptor);

    ValueStateDescriptor<String> typeDescriptor =
        new ValueStateDescriptor<>("type", String.class);
    type = getRuntimeContext().getState(typeDescriptor);

    ValueStateDescriptor<byte[]> valueDescriptor =
        new ValueStateDescriptor<>("value", byte[].class);
    value = getRuntimeContext().getState(valueDescriptor);

    objectMapper = new ObjectMapper();
  }

  @Override
  public void processElement(IncomingEvent event, Context ctx, Collector<OutgoingEvent> out) throws Exception {
    reasons.add(event.getReason());
    int c = currentFailure.value() + 1;
    if (c >= maxFailure.value()) {
      originalTopic.update(topic.value());
      topic.update("dlq");
      ctx.timerService().registerProcessingTimeTimer(0);
      return;
    }
    currentFailure.update(c);
    if (event.getTopic() != null) {
      topic.update(event.getTopic());
    }
    if (event.getType() != null) {
      type.update(event.getType());
    }
    if (event.getCap() != null) {
      cap.update(event.getCap());
    }
    if (event.getValue() != null) {
      value.update(event.getValue());
    }
    if (event.getMultiplier() != null) {
      multiplier.update(event.getMultiplier());
    }
    if (event.getMaxFailure() != null) {
      maxFailure.update(event.getMaxFailure());
    }
    if (event.getBackoff() != null) {
      backoff.update(event.getBackoff());
    }
    var b = Math.min(backoff.value() * multiplier.value(), cap.value());
    backoff.update(b);
    ctx.timerService().registerProcessingTimeTimer(Instant.now().toEpochMilli() + b);
  }

  @Override
  public void onTimer(long timestamp, KeyedProcessFunction<String, IncomingEvent, OutgoingEvent>.OnTimerContext ctx, Collector<OutgoingEvent> out) throws Exception {
    var t = topic.value();
    if (t.equals("dlq")) {
      out.collect(OutgoingEvent.builder()
          .topic("dlq")
          .type(type.value())
          .originalTopic(originalTopic.value())
          .rawReasons(objectMapper.writeValueAsBytes(reasons.get()))
          .value(value.value())
          .build());
      reasons.clear();
      backoff.clear();
      multiplier.clear();
      cap.clear();
      currentFailure.clear();
      maxFailure.clear();
      topic.clear();
      originalTopic.clear();
      type.clear();
      value.clear();
      return;
    }
    out.collect(OutgoingEvent.builder()
        .topic(t)
        .type(type.value())
        .value(value.value())
        .build());
  }
}
