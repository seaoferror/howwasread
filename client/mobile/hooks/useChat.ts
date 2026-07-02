import { useInfiniteQuery, useMutation, useQuery } from "@tanstack/react-query";
import {
  checkBlock,
  generatePresignedURL,
  getChatParticipants,
  getChatRoomInfo,
  sendMessage,
} from "@/api/chat";
import { AxiosError } from "axios";
import Toast from "react-native-toast-message";
import { queryKey } from "@/constants";
import { findMessagesByRoomId } from "@/db/message";
import { SQLiteDatabase } from "expo-sqlite";
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

export function useCheckBlock(id: string) {
  return useQuery({
    queryFn: () => checkBlock(id),
    queryKey: [queryKey.CHAT, queryKey.CHECK_BLOCK],
    enabled: getKVStore("type" + id) === "personal",
  });
}

export function useGetChatParticipants(roomId: string) {
  return useQuery({
    queryFn: () => getChatParticipants(roomId),
    queryKey: [queryKey.CHAT, queryKey.GET_CHAT_PARTICIPANT_IDS],
  });
}
