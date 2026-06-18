package com.xcecv.offlineconversation.util;

import java.util.UUID;

public class UUIDUtil {
  public static byte[] uuidToBytes(UUID uuid) {
    long msb = uuid.getMostSignificantBits();
    long lsb = uuid.getLeastSignificantBits();
    byte[] buffer = new byte[16];

    for (int i = 0; i < 8; i++) {
      buffer[i] = (byte) (msb >>> (8 * (7 - i)));
    }
    for (int i = 8; i < 16; i++) {
      buffer[i] = (byte) (lsb >>> (8 * (7 - (i - 8))));
    }
    return buffer;
  }
}
