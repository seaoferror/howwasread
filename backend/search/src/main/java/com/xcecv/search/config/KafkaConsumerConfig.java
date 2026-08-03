package com.xcecv.search.config;

import org.apache.kafka.clients.CommonClientConfigs;
import org.apache.kafka.clients.consumer.ConsumerConfig;
import org.apache.kafka.common.config.SaslConfigs;
import org.apache.kafka.common.config.SslConfigs;
import org.apache.kafka.common.serialization.ByteArrayDeserializer;
import org.apache.kafka.common.serialization.StringDeserializer;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.kafka.annotation.EnableKafka;
import org.springframework.kafka.config.ConcurrentKafkaListenerContainerFactory;
import org.springframework.kafka.core.ConsumerFactory;
import org.springframework.kafka.core.DefaultKafkaConsumerFactory;
import org.springframework.kafka.support.serializer.JacksonJsonDeserializer;

import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Paths;
import java.util.HashMap;
import java.util.Map;

@EnableKafka
@Configuration
public class KafkaConsumerConfig {

  @Value("${spring.kafka.bootstrap-servers}")
  private String bootstrapServers;

  @Value("${KAFKA_API_KEY:}")
  private String kafkaApiKey;

  @Value("${KAFKA_API_SECRET:}")
  private String kafkaApiSecret;

  @Value("${KAFKA_USER_CERT_PATH:}")
  private String kafkaUserCertPath;

  @Value("${KAFKA_USER_KEY_PATH:}")
  private String kafkaUserKeyPath;

  @Value("${KAFKA_CA_CERT_PATH:}")
  private String kafkaCACertPath;

  @Bean
  public ConsumerFactory<byte[], Object> consumerFactory() throws IOException {
    Map<String, Object> config = new HashMap<>();

    config.put(ConsumerConfig.BOOTSTRAP_SERVERS_CONFIG, bootstrapServers);
    config.put(ConsumerConfig.GROUP_ID_CONFIG, "search");
    config.put(ConsumerConfig.AUTO_OFFSET_RESET_CONFIG, "earliest");

    // Deserialization Properties
    config.put(ConsumerConfig.KEY_DESERIALIZER_CLASS_CONFIG, ByteArrayDeserializer.class);
    config.put(ConsumerConfig.VALUE_DESERIALIZER_CLASS_CONFIG, StringDeserializer.class);

    config.put(CommonClientConfigs.SECURITY_PROTOCOL_CONFIG, "SSL");

    if(!kafkaApiKey.isEmpty()){
      config.put(CommonClientConfigs.SECURITY_PROTOCOL_CONFIG, "SASL_SSL");
      config.put(SaslConfigs.SASL_MECHANISM, "PLAIN");
      String jaasTemplate = "org.apache.kafka.common.security.plain.PlainLoginModule required username=\"%s\" password=\"%s\";";
      config.put(SaslConfigs.SASL_JAAS_CONFIG, String.format(jaasTemplate, kafkaApiKey, kafkaApiSecret));
    }

    if(!kafkaUserCertPath.isEmpty()) {
      String userCert = new String(Files.readAllBytes(Paths.get(kafkaUserCertPath)));
      String userKey = new String(Files.readAllBytes(Paths.get(kafkaUserKeyPath)));
      config.put(SslConfigs.SSL_KEYSTORE_TYPE_CONFIG, "PEM");
      config.put(SslConfigs.SSL_KEYSTORE_CERTIFICATE_CHAIN_CONFIG, userCert);
      config.put(SslConfigs.SSL_KEYSTORE_KEY_CONFIG, userKey);
    }

    if(!kafkaCACertPath.isEmpty()) {
      String caCert = new String(Files.readAllBytes(Paths.get(kafkaCACertPath)));
      config.put(SslConfigs.SSL_TRUSTSTORE_TYPE_CONFIG, "PEM");
      config.put(SslConfigs.SSL_TRUSTSTORE_CERTIFICATES_CONFIG, caCert);
    }

    return new DefaultKafkaConsumerFactory<>(config);
  }

  @Bean
  public ConcurrentKafkaListenerContainerFactory<byte[], Object> kafkaListenerContainerFactory() throws IOException {
    ConcurrentKafkaListenerContainerFactory<byte[], Object> factory = new ConcurrentKafkaListenerContainerFactory<>();
    factory.setConsumerFactory(consumerFactory());
    return factory;
  }
}