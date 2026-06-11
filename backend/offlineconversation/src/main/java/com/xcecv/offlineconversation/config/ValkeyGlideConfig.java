package com.xcecv.offlineconversation.config;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.json.JsonMapper;
import glide.api.GlideClient;
import glide.api.models.configuration.GlideClientConfiguration;
import glide.api.models.configuration.NodeAddress;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

@Configuration
public class ValkeyGlideConfig {

  @Value("${valkey.host:localhost}")
  private String host;

  @Value("${valkey.port:6379}")
  private int port;

  @Bean(destroyMethod = "close")
  public GlideClient glideClient() {
    GlideClientConfiguration config = GlideClientConfiguration.builder()
        .address(NodeAddress.builder().host(host).port(port).build())
        // .useTLS(true) // Enable if using TLS
        // .credentials(Credentials.builder().username("user").password("pass").build())
        .requestTimeout(2000)
        .build();
    return GlideClient.createClient(config).join();
  }

  @Bean
  public ObjectMapper objectMapper() {
    return JsonMapper.builder()
        .findAndAddModules()
        .build();
  }
}