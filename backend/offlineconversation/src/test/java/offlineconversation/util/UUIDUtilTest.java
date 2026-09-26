package offlineconversation.util;

import org.junit.jupiter.api.Test;

import java.nio.ByteBuffer;
import java.util.UUID;

import static org.assertj.core.api.Assertions.assertThat;

class UUIDUtilTest {

  @Test
  void uuidToBytes_isSixteenBigEndianBytes() {
    UUID uuid = UUID.fromString("01234567-89ab-cdef-0123-456789abcdef");

    byte[] bytes = UUIDUtil.uuidToBytes(uuid);

    assertThat(bytes).hasSize(16);
    assertThat(bytes[0]).isEqualTo((byte) 0x01);
    assertThat(bytes[15]).isEqualTo((byte) 0xef);
    ByteBuffer buffer = ByteBuffer.wrap(bytes);
    assertThat(new UUID(buffer.getLong(), buffer.getLong())).isEqualTo(uuid);
  }
}
