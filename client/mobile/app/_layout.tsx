import { Stack, useFocusEffect } from "expo-router";
import { QueryClientProvider } from "@tanstack/react-query";
import queryClient from "@/api/queryClient";
import Toast from "react-native-toast-message";
import { ActionSheetProvider } from "@expo/react-native-action-sheet";
import { SQLiteProvider, useSQLiteContext } from "expo-sqlite";
import { parse as uuidParse, stringify as uuidStringify } from "uuid";
import { baseUrl, localDevId } from "@/api/axios";
import { MessagingResponse } from "@/types/chat";
import { useCallback, useRef } from "react";
import { useAuth } from "@/hooks/useAuth";
import { getSecureAsync, setSecure } from "@/util/storage";
import { Platform } from "react-native";
import { getTimestamp } from "@/util/time";
import { getRecentMessages } from "@/api/chat";
import { initDB, saveRecentMessage } from "@/db/message";
import { randomUUID } from "expo-crypto";
import {
  getDevicePushTokenAsync,
  getExpoPushTokenAsync,
  requestPermissionsAsync,
} from "expo-notifications";
import { registerNotification } from "@/api/notification";
import { setAudioModeAsync } from "expo-audio";

declare const WebSocket: {
  prototype: WebSocket;
  new (
    url: string,
    protocols?: string | string[] | null,
    options?: {
      headers?: { [header: string]: string };
      [key: string]: any;
    } | null,
  ): WebSocket;
};

export default function RootLayout() {
  return (
    <SQLiteProvider
      databaseName="db"
      onInit={initDB}
      options={{ enableChangeListener: true }}
    >
      <ActionSheetProvider>
        <QueryClientProvider client={queryClient}>
          <RootNavigator />
          <Toast />
        </QueryClientProvider>
      </ActionSheetProvider>
    </SQLiteProvider>
  );
}

function RootNavigator() {
  const db = useSQLiteContext();
  const { id } = useAuth();
  const ws = useRef<WebSocket>(null);


  useFocusEffect(
    useCallback(() => {
      const connectMessaging = async () => {
        await setAudioModeAsync({
          allowsRecording: true,
          playsInSilentMode: true,
        });
        let deviceId = await getSecureAsync("deviceId");
        console.log(deviceId);
        if (!deviceId) {
          deviceId = randomUUID();
          console.log(deviceId);
          await setSecure("deviceId", deviceId);
        }
        try {
          await requestPermissionsAsync({
            ios: {
              allowAlert: true,
              allowBadge: true,
              allowSound: true,
            },
          });
          await registerNotification({
            deviceId: deviceId,
            os: Platform.OS,
            devicePushToken:
              Platform.OS === "android"
                ? (await getDevicePushTokenAsync()).data
                : (await getExpoPushTokenAsync()).data,
          });
        } catch (error) {
          console.log(error);
        }

        const row = await db.getFirstAsync<{ id: Uint8Array }>(
          `SELECT id FROM message ORDER BY rowid DESC LIMIT 1`,
        );
        console.log(row);
        let cursor = "00000000-0000-7000-8000-000000000000";
        if (row) {
          cursor = uuidStringify(row.id);
          console.log("last inserted ID:", cursor);
        }
        const messages = await getRecentMessages(cursor);
        if (messages && messages.length > 0) {
          await Promise.all(
            messages.map((m) => {
              const timestamp = getTimestamp(m.id);
              const roomId = uuidParse(m.roomId);
              return saveRecentMessage(db, m, roomId, timestamp);
            }),
          );
        }

        ws.current = new WebSocket(
          `ws://${baseUrl.ios}:8080/chat/messaging/connect`,
          undefined,
          {
            headers: {
              Authorization: `Bearer ${await getSecureAsync("accessToken")}`,
              "X-User-Id": `${Platform.OS === "ios" ? localDevId.ios : localDevId.android}`,
              "Device-Id": await getSecureAsync("deviceId"),
            },
          },
        );
        ws.current.onmessage = async (event) => {
          const data: MessagingResponse = JSON.parse(event.data);
          const timestamp = getTimestamp(data.id);
          const roomId = uuidParse(data.roomId);
          await saveRecentMessage(db, data, roomId, timestamp);
        };
      };
      // if (id)
      connectMessaging();
      return () => {
        ws.current?.close();
      };
    }, [db]),
  );
  return (
    <Stack>
      <Stack.Screen name="(init)" options={{ headerShown: false }} />
      <Stack.Screen name="(tabs)" options={{ headerShown: false }} />
      <Stack.Screen name="auth" options={{ headerShown: false }} />
      <Stack.Screen name="profile" options={{ headerShown: false }} />
      <Stack.Screen name="online" options={{ headerShown: false }} />
    </Stack>
  );
}
