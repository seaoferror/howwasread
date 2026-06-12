package com.xcecv.offlineconversation.config;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.json.JsonMapper;
import com.zaxxer.hikari.util.Credentials;
import glide.api.GlideClient;
import glide.api.models.configuration.GlideClientConfiguration;
import glide.api.models.configuration.NodeAddress;
import glide.api.models.configuration.ServerCredentials;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

@Configuration
public class ValkeyGlideConfig {

  @Value("${valkey.host}")
  private String host;

  @Value("${valkey.port}")
  private int port;

  @Value("${valkey.username}")
  private String username;

  @Value("${valkey.password}")
  private String password;


  @Bean(destroyMethod = "close")
  public GlideClient glideClient() {
    GlideClientConfiguration config = GlideClientConfiguration.builder()
        .address(NodeAddress.builder().host(host).port(port).build())
        .useTLS(true)
        .credentials(ServerCredentials.builder()
            .username(username)
            .password(password)
            .build()
        )
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