import { Stack } from "expo-router";
import "react-native-reanimated";
import { QueryClientProvider } from "@tanstack/react-query";
import queryClient from "@/api/queryClient";
import Toast from "react-native-toast-message";
import { ActionSheetProvider } from "@expo/react-native-action-sheet";
import {
  type SQLiteDatabase,
  SQLiteProvider,
  useSQLiteContext,
} from "expo-sqlite";
import { parse as uuidParse } from "uuid";
import { baseUrl, localDevId } from "@/api/axios";
import { MessagingResponse } from "@/types/chat";
import { useEffect } from "react";
import { useAuth } from "@/hooks/useAuth";
import { getSecure } from "@/util/storage";
import { Platform } from "react-native";
import { getTimestamp } from "@/util/time";

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

async function initDB(db: SQLiteDatabase) {
  await db.execAsync(`
    PRAGMA journal_mode = 'wal';
    CREATE TABLE IF NOT EXISTS message (
        id BLOB PRIMARY KEY NOT NULL,
        roomId BLOB,
        fromId BLOB NOT NULL,
        content_type TEXT NOT NULL,
        content TEXT NOT NULL,
        created_at TEXT NOT NULL
    );
    `);
}

function RootNavigator() {
  const db = useSQLiteContext();
  const { id } = useAuth();

  useEffect(() => {
    const ws = new WebSocket(
      `ws://${baseUrl.ios}:8080/chat/connect`,
      undefined,
      {
        headers: {
          Authorization: `Bearer ${getSecure("accessToken")}`,
          "X-User-Id": `${Platform.OS === "ios" ? localDevId.ios : localDevId.android}`,
        },
      },
    );
    ws.onmessage = async (event) => {
      const data: MessagingResponse = JSON.parse(event.data);
      await db.runAsync(
        `INSERT INTO message (id, roomId, fromId, content_type, content, created_at)
                         VALUES (?, ?, ?, ?, ?, ?);`,
        uuidParse(data.id),
        data.roomId ? uuidParse(data.roomId) : null,
        uuidParse(data.fromId),
        data.contentType,
        data.content,
        getTimestamp(data.id),
      );
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
