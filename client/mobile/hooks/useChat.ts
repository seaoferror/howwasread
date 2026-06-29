import {
  useInfiniteQuery,
  useMutation,
  useQuery,
} from "@tanstack/react-query";
import {
  checkBlock,
  generatePresignedURL,
  getChatRoomInfo,
  sendMessage,
} from "@/api/chat";
import { AxiosError } from "axios";
import Toast from "react-native-toast-message";
import { queryKey } from "@/constants";
import {
  deleteMessagesBeforeQuit,
  findMessagesByRoomId,
  findNewMessage,
  findPreview,
} from "@/db/message";
import * as SQLite from "expo-sqlite";
import { SQLiteDatabase, useSQLiteContext } from "expo-sqlite";
import { stringify as uuidStringify } from "uuid";
import { useEffect, useState } from "react";
import { Message } from "@/types/chat";
import { getKVStore } from "@/db/storage";

export function useGetInfiniteMessages(db: SQLiteDatabase, roomId: string) {
  return useInfiniteQuery({
    queryFn: ({ pageParam }) => findMessagesByRoomId(db, roomId, pageParam),
    queryKey: [queryKey.CHAT, queryKey.FIND_MESSAGES_BY_ROOM_ID, roomId],
    initialPageParam: 1,
    getNextPageParam: (lastPage, allPages) => {
      const lastPost = lastPage[lastPage.length - 1];
      return lastPost ? allPages.length + 1 : undefined;
    },
    staleTime: 0,
    gcTime: 0,
  });
}

export function useGetChatRoomInfo(roomId: string) {
  return useQuery({
    queryFn: () => getChatRoomInfo(roomId),
    queryKey: [queryKey.CHAT, queryKey.GET_CHAT_ROOM_INFO, roomId],
  });
}

export function useSendMessage() {
  return useMutation({
    mutationFn: sendMessage,
    onError: (error: AxiosError) => {
      console.log(error?.response?.data);
      Toast.show({
        type: "error",
        text1: String(error?.response?.data),
      });
    },
  });
}

export function useGeneratePresignedURL() {
  return useMutation({
    mutationFn: generatePresignedURL,
    onError: (error: AxiosError) => {
      console.log(error?.response?.data);
      Toast.show({
        type: "error",
        text1: String(error?.response?.data),
      });
    },
  });
}

export function usePreview() {
  const db = useSQLiteContext();
  const [preview, setPreview] = useState<Message[]>([]);

  useEffect(() => {
    const wrapper = async () => {
      const previewRaw = await findPreview(db);

      const p = previewRaw.map((row) => ({
        id: uuidStringify(row.id),
        roomId: uuidStringify(row.room_id),
        fromId: uuidStringify(row.from_id),
        contentType: row.content_type,
        contents: JSON.parse(row.contents),
        createdAt: row.created_at,
      }));
      setPreview(p);

      SQLite.addDatabaseChangeListener(async (event) => {
        const newMessageRaw = await findNewMessage(db, event.rowId);

        if (!newMessageRaw) return;

        const newPreviewItem: Message = {
          id: uuidStringify(newMessageRaw.id),
          roomId: uuidStringify(newMessageRaw.room_id),
          fromId: uuidStringify(newMessageRaw.from_id),
          contentType: newMessageRaw.content_type,
          contents: JSON.parse(newMessageRaw.contents),
          createdAt: newMessageRaw.created_at,
        };
        if (newPreviewItem.contentType === "quit") {
          await deleteMessagesBeforeQuit(db, newPreviewItem);
          setPreview((prev) =>
            prev.filter((item) => item.roomId !== newPreviewItem.roomId),
          );
          return;
        }

        setPreview((prev) => {
          const filtered = prev.filter(
            (item) => item.roomId !== newPreviewItem.roomId,
          );
          return [newPreviewItem, ...filtered];
        });
      });
    };
    wrapper();

    return () => {};
  }, []);
  return preview;
}

export function useCheckBlock(id: string) {
  return useQuery({
    queryFn: () => checkBlock(id),
    queryKey: [queryKey.CHAT, queryKey.CHECK_BLOCK],
    enabled: getKVStore("type" + id) === "personal",
  });
}
