package com.howwasread.dto;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.io.Serial;
import java.io.Serializable;

@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class Content implements Serializable {

  @Serial
  private static final long serialVersionUID = 1L;

  private String title;
  private String body;
//  private String imageURL;
}
