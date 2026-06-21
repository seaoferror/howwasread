import { Stack } from "expo-router";
import { QueryClientProvider } from "@tanstack/react-query";
import queryClient from "@/api/queryClient";
import Toast from "react-native-toast-message";
import { ActionSheetProvider } from "@expo/react-native-action-sheet";
import { SQLiteProvider, useSQLiteContext } from "expo-sqlite";
import { stringify as uuidStringify } from "uuid";
import { MessagingResponse } from "@/types/chat";
import { useEffect, useRef } from "react";
import { getSecure, getSecureAsync, setSecure } from "@/util/storage";
import { Platform } from "react-native";
import { getRecentMessages } from "@/api/chat";
import {
  deleteMessagesBeforeQuit,
  initDB,
  saveRecentMessage,
} from "@/db/message";
import { randomUUID } from "expo-crypto";
import {
  getDevicePushTokenAsync,
  getExpoPushTokenAsync,
  requestPermissionsAsync,
} from "expo-notifications";
import { registerNotification } from "@/api/notification";
import { setAudioModeAsync } from "expo-audio";
import { useGetMyProfile } from "@/hooks/useProfile";
import { GoogleSignin } from "@react-native-google-signin/google-signin";

GoogleSignin.configure({
  webClientId: process.env.EXPO_PUBLIC_GOOGLE_SIGN_IN_WEB_CLIENT_ID,
});

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
  const { data: profile } = useGetMyProfile();
  const ws = useRef<WebSocket>(null);

  useEffect(() => {
    console.log(getSecure("accessToken"));
    const connectMessaging = async () => {
      await setAudioModeAsync({
        allowsRecording: true,
        playsInSilentMode: true,
      });
      let deviceId = await getSecureAsync("deviceId");
      if (!deviceId) {
        deviceId = randomUUID();
        await setSecure("deviceId", deviceId);
      }
      console.log(deviceId);
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
        `SELECT id FROM message ORDER BY id DESC LIMIT 1`,
      );
      console.log(row);
      let cursor = "00000000-0000-7000-8000-000000000000";
      if (row) {
        cursor = uuidStringify(row.id);
        console.log("last inserted ID:", cursor);
      }
      const messages = await getRecentMessages(cursor);
      if (messages && messages.length > 0) {
        const deleteQueue: MessagingResponse[] = [];
        await Promise.all(
          messages.map((m) => {
            if (m.contentType === "quit" && m.fromId === profile?.id) {
              deleteQueue.push(m);
              return
            }
            return saveRecentMessage(db, m);
          }),
        );
        await Promise.all(
          deleteQueue.map((m) => {
            return deleteMessagesBeforeQuit(db, m);
          })
        )
      }

      ws.current = new WebSocket(
        `wss://${process.env.EXPO_PUBLIC_API_URL}/chat/messaging/connect`,
        undefined,
        {
          headers: {
            Authorization: `Bearer ${await getSecureAsync("accessToken")}`,
            // "X-User-Id": `${Platform.OS === "ios" ? localDevId.ios : localDevId.android}`,
            "Device-Id": await getSecureAsync("deviceId"),
          },
        },
      );
      ws.current.onmessage = async (event) => {
        const m: MessagingResponse = JSON.parse(event.data);
        if (m.contentType === "quit" && m.fromId === profile?.id) {
          await deleteMessagesBeforeQuit(db, m);
          return
        }
        await saveRecentMessage(db, m);
      };
    };
    if (profile?.name) connectMessaging();
    return () => {
      ws.current?.close();
    };
  }, [profile]);
  return (
    <Stack>
      <Stack.Screen name="(init)" options={{ headerShown: false }} />
      <Stack.Screen name="(tabs)" options={{ headerShown: false }} />
      <Stack.Screen name="auth" options={{ headerShown: false }} />
      <Stack.Screen name="profile" options={{ headerShown: false }} />
      <Stack.Screen name="chat" options={{ headerShown: false }} />
      <Stack.Screen name="offline" options={{ headerShown: false }} />
      <Stack.Screen name="online" options={{ headerShown: false }} />
    </Stack>
  );
}
