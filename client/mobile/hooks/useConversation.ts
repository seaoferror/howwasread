import { useInfiniteQuery, useMutation, useQuery } from "@tanstack/react-query";
import {
  getOnlineConversations,
  createOnlineConversation,
} from "@/api/conversation";
import { queryKey } from "@/constants";
import { AxiosError } from "axios";
import Toast from "react-native-toast-message";
import queryClient from "@/api/queryClient";

export function useGetInfiniteOnlineConversations() {
  return useInfiniteQuery({
    queryFn: ({ pageParam }) => getOnlineConversations(pageParam),
    queryKey: [queryKey.CONVERSATION, queryKey.GET_ONLINE_CONVERSATIONS],
    initialPageParam: 1,
    getNextPageParam: (lastPage, allPages) => {
      const lastPost = lastPage[lastPage.length - 1];
      return lastPost ? allPages.length + 1 : undefined;
    },
  });
}

export function useCreateOnlineConversation() {
  return useMutation({
    mutationFn: createOnlineConversation,
    onSuccess: async () => {
      await queryClient.invalidateQueries({
        queryKey: [queryKey.CONVERSATION, queryKey.GET_ONLINE_CONVERSATIONS],
      });
    },
    onError: (error: AxiosError) => {
      console.log(error?.response?.data);
      Toast.show({
        type: "error",
        text1: String(error?.response?.data),
      });
    },
  });
}
