package com.howwasread.dto;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.io.Serial;
import java.io.Serializable;
import java.util.Map;

@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class NotificationElement implements Serializable {

  @Serial
  private static final long serialVersionUID = 1L;

  private Long scheduledTime;
  private Map<Integer, String> contents;
}
