package com.xcecv.offlineconversation.util;

import java.util.UUID;

public class UUIDUtil {

  private static final char[] HEX_ARRAY = "0123456789abcdef".toCharArray();

  public static String bytesToHex(byte[] bytes) {
    if (bytes == null) {
      return null;
    }
    char[] hexChars = new char[bytes.length * 2];
    for (int j = 0; j < bytes.length; j++) {
      int v = bytes[j] & 0xFF;
      hexChars[j * 2] = HEX_ARRAY[v >>> 4];
      hexChars[j * 2 + 1] = HEX_ARRAY[v & 0x0F];
    }
    return new String(hexChars);
  }
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
