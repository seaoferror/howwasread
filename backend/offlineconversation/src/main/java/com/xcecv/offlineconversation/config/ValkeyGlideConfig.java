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

  @Value("${VALKEY_CA_CERT_PATH:}")
  private String caCertPath;

  @Bean(destroyMethod = "close")
  public GlideClusterClient glideClient() throws IOException {
    AdvancedGlideClusterClientConfiguration advancedConfig = null;
    if (!caCertPath.isEmpty()) {
      byte[] caCertBytes = Files.readAllBytes(Paths.get(caCertPath));
      TlsAdvancedConfiguration tlsConfig = TlsAdvancedConfiguration.builder()
          .rootCertificates(caCertBytes)
          .build();
      advancedConfig = AdvancedGlideClusterClientConfiguration.builder()
          .tlsAdvancedConfiguration(tlsConfig)
          .build();
    }

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