package com.howwasread;

import org.apache.flink.streaming.api.functions.KeyedProcessFunction;
import org.apache.flink.util.Collector;

/**
 * Delays notification events using Flink's timer mechanism.
 * TODO: Implement your actual delay logic using event-time or processing-time timers.
 */
public class DelayNotificationFunction
        extends KeyedProcessFunction<String, NotificationEvent, NotificationEvent> {

    private static final long serialVersionUID = 1L;

    // Default delay: 5 minutes (in milliseconds)
    private static final long DELAY_MS = 5 * 60 * 1000L;

    @Override
    public void processElement(
            NotificationEvent event,
            KeyedProcessFunction<String, NotificationEvent, NotificationEvent>.Context ctx,
            Collector<NotificationEvent> out) throws Exception {

        // TODO: Register a processing-time timer instead of passing through immediately
        // For now, pass events through as a placeholder
        long fireTimestamp = ctx.timerService().currentProcessingTime() + DELAY_MS;
        ctx.timerService().registerProcessingTimeTimer(fireTimestamp);

        // TODO: Store the event in Flink state so it can be emitted in onTimer()
        out.collect(event);
    }

    @Override
    public void onTimer(
            long timestamp,
            KeyedProcessFunction<String, NotificationEvent, NotificationEvent>.OnTimerContext ctx,
            Collector<NotificationEvent> out) throws Exception {

        // TODO: Retrieve the stored event from state and emit it here
        // This is where the delayed notification should be emitted
    }
}
