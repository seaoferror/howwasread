import { createContext, ReactNode, useContext, useEffect, useRef } from "react";
import { baseUrl } from "@/api/axios";
import { MessagingResponse } from "@/types/chat";
import { useSQLiteContext } from "expo-sqlite";
import { parse as uuidParse } from "uuid";

const ChatContext = createContext<WebSocket | null>(null);

export function ChatProvider({ children }: { children: ReactNode }) {
  const ws = useRef<WebSocket>(null);
  const db = useSQLiteContext();

  useEffect(() => {
    ws.current = new WebSocket(`ws://${baseUrl.ios}:8080/chat/connect`);
    ws.current.onmessage = async (event) => {
      const data: MessagingResponse = JSON.parse(event.data);
      await db.runAsync(
        `INSERT INTO message (id, room_id, from_id, content_type, content)
                         VALUES (?, ?, ?, ?, ?);`,
        uuidParse(data.id),
        data.roomId ? uuidParse(data.roomId) : null,
        data.fromId,
        data.contentType,
        data.content,
      );
    };

    return () => {
      ws.current?.close();
    };
  }, []);
  return (
    <ChatContext.Provider value={ws.current}>{children}</ChatContext.Provider>
  );
}

export function useChat() {
  const context = useContext(ChatContext);
  if (!context) throw new Error("useChat must be used within a provider");
  return context;
}
