package com.xcecv.offlineconversation.config;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.json.JsonMapper;
import glide.api.GlideClusterClient;
import glide.api.models.configuration.*;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Paths;

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
  public GlideClusterClient glideClusterClient() throws IOException {
    byte[] caCertBytes = Files.readAllBytes(Paths.get("/cert/valkey/ca.crt"));

    TlsAdvancedConfiguration tlsConfig = TlsAdvancedConfiguration.builder()
        .rootCertificates(caCertBytes)
        .build();

    AdvancedGlideClusterClientConfiguration advancedConfig = AdvancedGlideClusterClientConfiguration.builder()
        .tlsAdvancedConfiguration(tlsConfig)
        .build();

    GlideClusterClientConfiguration config = GlideClusterClientConfiguration.builder()
        .address(NodeAddress.builder().host(host).port(port).build())
        .useTLS(true)
        .advancedConfiguration(advancedConfig)
        .credentials(ServerCredentials.builder()
            .username(username)
            .password(password)
            .build()
        )
        .requestTimeout(2000)
        .build();
    return GlideClusterClient.createClient(config).join();
  }

  @Bean
  public ObjectMapper objectMapper() {
    return JsonMapper.builder()
        .findAndAddModules()
        .build();
  }
}