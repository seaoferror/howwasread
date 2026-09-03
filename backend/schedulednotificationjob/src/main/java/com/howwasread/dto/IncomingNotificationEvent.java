package com.howwasread.dto;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.io.Serial;
import java.io.Serializable;
import java.util.Map;
import java.util.UUID;

@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class IncomingNotificationEvent implements Serializable {

  @Serial
  private static final long serialVersionUID = 1L;

  private UUID partitionId;
  private String type;
  private Long scheduledTime;
  private UUID keyId;
  private Map<Integer, String> contents;
}
