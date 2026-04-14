import { Stack } from "expo-router";
import "react-native-reanimated";
import { QueryClientProvider } from "@tanstack/react-query";
import queryClient from "@/api/queryClient";
import Toast from "react-native-toast-message";
import { ActionSheetProvider } from "@expo/react-native-action-sheet";
import { SQLiteProvider, useSQLiteContext } from "expo-sqlite";
import { parse as uuidParse, stringify as uuidStringify } from "uuid";
import { baseUrl, localDevId } from "@/api/axios";
import { MessagingResponse } from "@/types/chat";
import { useEffect, useRef } from "react";
import { useAuth } from "@/hooks/useAuth";
import { getSecureAsync } from "@/util/storage";
import { Platform } from "react-native";
import { getTimestamp } from "@/util/time";
import { getRecentMessages } from "@/api/chat";
import { checkIfFirstOfDay, initDB } from "@/db/message";

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

  useEffect(() => {
    const connectMessaging = async () => {
      const row = await db.getFirstAsync<{ id: Uint8Array }>(
        `SELECT id FROM message ORDER BY rowid DESC LIMIT 1`,
      );
      console.log(row);
      let cursor = "00000000-0000-7000-8000-000000000000";
      if (row) {
        cursor = uuidStringify(row.id);
        console.log("The very last inserted ID is:", cursor);
      }
      const messages = await getRecentMessages(cursor);
      console.log(messages);
      if (messages && messages.length > 0) {
        for (const m of messages) {
          const timestamp = getTimestamp(m.id);
          console.log(m.roomId);
          const roomId = uuidParse(m.roomId);
          await db.runAsync(
            `INSERT OR IGNORE INTO message (id, room_id, from_id, content_type, content, created_at, is_day_first)
               VALUES (?, ?, ?, ?, ?, ?, ?);`,
            uuidParse(m.id),
            roomId,
            uuidParse(m.fromId),
            m.contentType,
            m.content,
            timestamp,
            await checkIfFirstOfDay(db, roomId, timestamp),
          );
        }
      }

      ws.current = new WebSocket(
        `ws://${baseUrl.ios}:8080/chat/messaging/connect`,
        undefined,
        {
          headers: {
            Authorization: `Bearer ${await getSecureAsync("accessToken")}`,
            "X-User-Id": `${Platform.OS === "ios" ? localDevId.ios : localDevId.android}`,
          },
        },
      );
      ws.current.onmessage = async (event) => {
        const data: MessagingResponse = JSON.parse(event.data);
        const timestamp = getTimestamp(data.id);
        const roomId = uuidParse(data.roomId);
        await db.runAsync(
          `INSERT OR IGNORE INTO message (id, room_id, from_id, content_type, content, created_at)
           VALUES (?, ?, ?, ?, ?, ?);`,
          uuidParse(data.id),
          roomId,
          uuidParse(data.fromId),
          data.contentType,
          data.content,
          timestamp,
          await checkIfFirstOfDay(db, roomId, timestamp),
        );
      };
    };
    // if (id)//just for development, this should be non commented when production
      connectMessaging();
    return () => {
      ws.current?.close();
    };
  }, [db, id]);
  return (
    <Stack>
      <Stack.Screen name="(init)" options={{ headerShown: false }} />
      <Stack.Screen name="(tabs)" options={{ headerShown: false }} />
      <Stack.Screen name="auth" options={{ headerShown: false }} />
      <Stack.Screen name="newcomer" options={{ headerShown: false }} />
      <Stack.Screen name="online" options={{ headerShown: false }} />
    </Stack>
  );
}
