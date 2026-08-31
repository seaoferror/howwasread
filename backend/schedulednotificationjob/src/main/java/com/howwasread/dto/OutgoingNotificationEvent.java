package com.howwasread.dto;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.io.Serial;
import java.io.Serializable;
import java.util.List;
import java.util.UUID;


@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class OutgoingNotificationEvent implements Serializable {

    @Serial
    private static final long serialVersionUID = 1L;

    private UUID conversationId;
    private List<UUID> memberIds;
    private String writtenBy;
    private long scheduledTime;
    private String notificationType;
}
