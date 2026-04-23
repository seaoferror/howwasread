import {
  useInfiniteQuery,
  useMutation,
  useQueries,
  useQuery,
} from "@tanstack/react-query";
import {
  generatePresignedURL,
  getChatRoomInfo,
  getSignedURL,
  sendLike,
  sendMessage,
} from "@/api/chat";
import { AxiosError } from "axios";
import Toast from "react-native-toast-message";
import { queryKey } from "@/constants";
import { findMessagesByRoomId } from "@/db/message";
import { SQLiteDatabase } from "expo-sqlite";
import { data } from "browserslist";

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
    queryKey: [queryKey.CHAT, queryKey.FIND_MESSAGES_BY_ROOM_ID, roomId],
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

export function useGetSignedURLs({
  contentType,
  contents,
}: {
  contentType: string;
  contents: string[];
}) {
  const queries = useQueries({
    queries: contents.map((filename) => {
      return {
        queryKey: [
          queryKey.CHAT,
          queryKey.GET_SIGNED_URL,
          contentType,
          filename,
        ],
        queryFn: () =>
          getSignedURL({
            contentType,
            filename: filename,
          }),
        enabled: contentType !== "text" && !!filename,
        staleTime: 1000 * 60 * 60,
      };
    }),
  });

  return queries.map((result) => result.data)
}
