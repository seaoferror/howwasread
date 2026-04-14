import { useInfiniteQuery, useMutation, useQuery } from "@tanstack/react-query";
import { getChatRoomInfo, sendLike, sendMessaging } from "@/api/chat";
import { AxiosError } from "axios";
import Toast from "react-native-toast-message";
import { queryKey } from "@/constants";
import { findMessagesByRoomId } from "@/db/message";
import { SQLiteDatabase } from "expo-sqlite";

export function useSendLike() {
  return useMutation({
    mutationFn: sendLike,
    onError: (error: AxiosError) => {
      console.log(error?.response?.data);
      Toast.show({
        type: "error",
        text1: String(error?.response?.data),
      });
    },
  });
}

export function useGetInfiniteMessages(db: SQLiteDatabase, roomId: string) {
  return useInfiniteQuery({
    queryFn: ({ pageParam }) => findMessagesByRoomId(db, roomId, pageParam),
    queryKey: [
      queryKey.CONVERSATION,
      queryKey.FIND_MESSAGES_BY_ROOM_ID,
      roomId,
    ],
    initialPageParam: 1,
    getNextPageParam: (lastPage, allPages) => {
      const lastPost = lastPage[lastPage.length - 1];
      return lastPost ? allPages.length + 1 : undefined;
    },
  });
}

export function useGetChatRoomInfo(roomId: string) {
  const { data } = useQuery({
    queryFn: () => getChatRoomInfo(roomId),
    queryKey: [queryKey.CHAT, queryKey.GET_CHAT_ROOM_INFO, roomId],
  });
  return { data };
}

export function useSendMessaging() {
  return useMutation({
    mutationFn: sendMessaging,
    onError: (error: AxiosError) => {
      console.log(error?.response?.data);
      Toast.show({
        type: "error",
        text1: String(error?.response?.data),
      });
    },
  })
}
