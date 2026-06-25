package com.xcecv.offlineconversation.config;

import org.apache.kafka.clients.CommonClientConfigs;
import org.apache.kafka.clients.producer.ProducerConfig;
import org.apache.kafka.common.config.SaslConfigs;
import org.apache.kafka.common.config.SslConfigs;
import org.apache.kafka.common.serialization.ByteArraySerializer;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.kafka.core.DefaultKafkaProducerFactory;
import org.springframework.kafka.core.KafkaTemplate;
import org.springframework.kafka.core.ProducerFactory;
import org.springframework.kafka.support.serializer.JacksonJsonSerializer;

import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Paths;
import java.util.HashMap;
import java.util.Map;
import java.util.UUID;

@Configuration
public class KafkaProducerConfig {

  @Value("${spring.kafka.bootstrap-servers}")
  private String bootstrapServers;

  @Value("${KAFKA_API_KEY}")
  private String kafkaApiKey;

  @Value("${KAFKA_API_SECRET}")
  private String kafkaApiSecret;

  @Bean
  public ProducerFactory<byte[], Object> producerFactory() throws IOException {
    Map<String, Object> config = new HashMap<>();
    config.put(ProducerConfig.BOOTSTRAP_SERVERS_CONFIG, bootstrapServers);
    config.put(ProducerConfig.CLIENT_ID_CONFIG, "producer_offlineconversation" + UUID.randomUUID());

    config.put(ProducerConfig.KEY_SERIALIZER_CLASS_CONFIG, ByteArraySerializer.class);
    config.put(ProducerConfig.VALUE_SERIALIZER_CLASS_CONFIG, JacksonJsonSerializer.class);

    config.put(ProducerConfig.COMPRESSION_TYPE_CONFIG, "snappy");
    config.put(ProducerConfig.ACKS_CONFIG, "1");
    config.put(ProducerConfig.ENABLE_IDEMPOTENCE_CONFIG, false);
    config.put(ProducerConfig.RETRIES_CONFIG, 5);
    config.put(ProducerConfig.RETRY_BACKOFF_MS_CONFIG, 500);
    config.put(ProducerConfig.MAX_IN_FLIGHT_REQUESTS_PER_CONNECTION, 5);

    config.put(CommonClientConfigs.SECURITY_PROTOCOL_CONFIG, "SASL_SSL");
    config.put(SaslConfigs.SASL_MECHANISM, "PLAIN");

    String jaasTemplate = "org.apache.kafka.common.security.plain.PlainLoginModule required username=\"%s\" password=\"%s\";";
    config.put(SaslConfigs.SASL_JAAS_CONFIG, String.format(jaasTemplate, kafkaApiKey, kafkaApiSecret));

//    config.put(ProducerConfig.DELIVERY_TIMEOUT_MS_CONFIG, 2500);
//    String caCert = new String(Files.readAllBytes(Paths.get("/cert/kafka/cluster/ca.crt")));
//    String userCert = new String(Files.readAllBytes(Paths.get("/cert/kafka/user/user.crt")));
//    String userKey = new String(Files.readAllBytes(Paths.get("/cert/kafka/user/user.key")));
//    config.put(SslConfigs.SSL_TRUSTSTORE_TYPE_CONFIG, "PEM");
//    config.put(SslConfigs.SSL_TRUSTSTORE_CERTIFICATES_CONFIG, caCert);
//    config.put(SslConfigs.SSL_KEYSTORE_TYPE_CONFIG, "PEM");
//    config.put(SslConfigs.SSL_KEYSTORE_CERTIFICATE_CHAIN_CONFIG, userCert);
//    config.put(SslConfigs.SSL_KEYSTORE_KEY_CONFIG, userKey);

    return new DefaultKafkaProducerFactory<>(config);
  }

  @Bean
  public KafkaTemplate<byte[], Object> kafkaTemplate() throws IOException {
    return new KafkaTemplate<>(producerFactory());
  }
}