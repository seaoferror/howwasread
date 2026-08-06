import { Stack } from "expo-router";
import { QueryClientProvider } from "@tanstack/react-query";
import queryClient from "@/api/queryClient";
import Toast from "react-native-toast-message";
import { ActionSheetProvider } from "@expo/react-native-action-sheet";
import { SQLiteProvider, useSQLiteContext } from "expo-sqlite";
import { MessagingResponse } from "@/types/chat";
import { useEffect, useRef } from "react";
import {
  getKVStore,
  getSecureAsync,
  setKVStore,
  setSecure,
} from "@/db/storage";
import { Platform } from "react-native";
import { getRecentMessages } from "@/api/chat";
import {
  deleteMessage,
  deleteMessagesBeforeQuit,
  downloadFiles,
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
import { useGetMyProfile } from "@/hooks/useProfile";
import { GoogleSignin } from "@react-native-google-signin/google-signin";
import * as Sentry from "@sentry/react-native";
import { setAudioModeAsync } from "expo-audio";

Sentry.init({
  dsn: "https://4d3c893e931b2ee0c936c646d22033bc@o4511734928572416.ingest.de.sentry.io/4511734931783760",

  // Adds more context data to events (IP address, cookies, user, etc.)
  // For more information, visit: https://docs.sentry.io/platforms/react-native/data-management/data-collected/
  sendDefaultPii: true,

  // Enable Logs
  enableLogs: true,

  // Configure Session Replay
  replaysSessionSampleRate: 0.1,
  replaysOnErrorSampleRate: 1,
  integrations: [
    Sentry.mobileReplayIntegration(),
    Sentry.feedbackIntegration(),
  ],

  // uncomment the line below to enable Spotlight (https://spotlightjs.com)
  // spotlight: __DEV__,
});

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

export default Sentry.wrap(function RootLayout() {
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
});

function RootNavigator() {
  const db = useSQLiteContext();
  const { data: profile } = useGetMyProfile();
  const ws = useRef<WebSocket>(null);

  useEffect(() => {
    if (!profile?.name) {
      return;
    }
    const connectMessaging = async () => {
      console.log("myId", profile.id);
      setKVStore("myId", profile.id);
      setKVStore("myName", profile.name);
      let deviceId = await getSecureAsync("deviceId");
      if (!deviceId) {
        deviceId = randomUUID();
        await setSecure("deviceId", deviceId);
      }
      console.log(deviceId);
      try {
        await setAudioModeAsync({
          allowsRecording: true,
          allowsBackgroundRecording: true,
          shouldPlayInBackground: true,
          shouldRouteThroughEarpiece: true,
          playsInSilentMode: true,
        });
        await requestPermissionsAsync({
          ios: {
            allowAlert: true,
            allowBadge: true,
            allowSound: true,
          },
        });
        await registerNotification({
          os: Platform.OS,
          devicePushToken:
            Platform.OS === "android"
              ? (await getDevicePushTokenAsync()).data
              : (await getExpoPushTokenAsync()).data,
        });
      } catch (error) {
        console.log(error);
      }

      let cursor = getKVStore("recentMessageId");
      if (!cursor) {
        cursor = "00000000-0000-7000-8000-000000000000";
      }
      const messages = await getRecentMessages(cursor);
      if (messages && messages.length > 0) {
        setKVStore("recentMessageId", messages[0].id);
        const myQuitMessages = messages.filter((m) => m.contentType === "quit");
        const myDeleteMessages = messages.filter(
          (m) => m.contentType === "delete",
        );
        const validMessagesToSave = messages.filter((m) => {
          if (m.contentType === "quit" || m.contentType === "delete") {
            return false;
          }
          const wipedByQuit = myQuitMessages.some(
            (q) => q.roomId === m.roomId && q.id > m.id,
          );
          const wipedByDelete = myDeleteMessages.some(
            (d) => d.contents[0] === m.id,
          );
          return !wipedByQuit && !wipedByDelete;
        });
        for (const q of myQuitMessages) {
          await deleteMessagesBeforeQuit(db, q);
        }
        for (const d of myDeleteMessages) {
          await deleteMessage(db, d.contents[0]);
        }
        for (const m of validMessagesToSave) {
          if (
            m.contentType === "audio" ||
            m.contentType === "image" ||
            m.contentType === "video"
          ) {
            await downloadFiles(m);
          }
          await saveRecentMessage(db, m);
        }
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
        console.log(m);
        setKVStore("recentMessageId", m.id);
        if (
          m.contentType === "audio" ||
          m.contentType === "image" ||
          m.contentType === "video"
        ) {
          await downloadFiles(m);
        }
        await saveRecentMessage(db, m);
      };
    };
    connectMessaging();
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
