package com.xcecv.offlineconversation.config;

import org.opensearch.data.client.osc.OpenSearchConfiguration;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.data.elasticsearch.client.ClientConfiguration;


@Configuration
public class OpenSearchConfig extends OpenSearchConfiguration {

  @Value("${spring.elasticsearch.uris}")
  private String uri;

  @Value("${spring.elasticsearch.username}")
  private String username;

  @Value("${spring.elasticsearch.password}")
  private String password;

  @Override
  @Bean
  public ClientConfiguration clientConfiguration() {
    return ClientConfiguration.builder()
        .connectedTo(uri)
        .usingSsl()
        .withBasicAuth(username, password)
        .build();
  }
}