import {
  useInfiniteQuery,
  useMutation,
  useQueries,
  useQuery,
} from "@tanstack/react-query";
import {
  banParticipant,
  createOfflineConversation,
  createOnlineConversation,
  getOfflineConversationDetail,
  getOnlineConversations,
  joinOfflineConversation,
  mapOfflineConversation,
  quitOfflineConversation,
} from "@/api/conversation";
import { queryKey } from "@/constants";
import { AxiosError } from "axios";
import Toast from "react-native-toast-message";
import queryClient from "@/api/queryClient";
import { OfflineConversationMapResponse } from "@/types/conversation";

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

export function useBanParticipant() {
  return useMutation({
    mutationFn: banParticipant,
    onError: (error: AxiosError) => {
      console.log(error?.response?.data);
      Toast.show({
        type: "error",
        text1: String(error?.response?.data),
      });
    },
  });
}

export function useMapOfflineConversations({
  resolution,
  h3Indexes,
}: {
  resolution: number;
  h3Indexes: string[];
}): OfflineConversationMapResponse[] {
  const queries = useQueries({
    queries: h3Indexes.map((h3Index) => {
      return {
        queryKey: [
          queryKey.CONVERSATION,
          queryKey.MAP_OFFLINE_CONVERSATION,
          h3Index,
        ],
        queryFn: () => mapOfflineConversation({ resolution, h3Index }),
      };
    }),
  });

  return queries.flatMap((query) => {
    if (query.data) {
      return query.data;
    }
    return [];
  });
}

export function useCreateOfflineConversation() {
  return useMutation({
    mutationFn: createOfflineConversation,
    onSuccess: async (data, variables) => {
      await queryClient.invalidateQueries({
        queryKey: [
          queryKey.CONVERSATION,
          queryKey.MAP_OFFLINE_CONVERSATION,
          variables?.h3Res5,
        ],
      });
      await queryClient.invalidateQueries({
        queryKey: [
          queryKey.CONVERSATION,
          queryKey.MAP_OFFLINE_CONVERSATION,
          variables?.h3Res7,
        ],
      });
    },
  });
}

export function useGetOfflineConversationDetail(id: string) {
  const { data } = useQuery({
    queryFn: () => getOfflineConversationDetail(id),
    queryKey: [
      queryKey.CONVERSATION,
      queryKey.GET_OFFLINE_CONVERSATION_DETAIL,
      id,
    ],
    enabled: !!id,
  });
  return { data };
}

export function useJoinOfflineConversation() {
  return useMutation({
    mutationFn: joinOfflineConversation,
    onSuccess: async (data, variables) => {
      await queryClient.invalidateQueries({
        queryKey: [
          queryKey.CONVERSATION,
          queryKey.GET_OFFLINE_CONVERSATION_DETAIL,
          variables.conversationId,
        ],
      });
    },
  });
}

export function useQuitOfflineConversation() {
  return useMutation({
    mutationFn: quitOfflineConversation,
    onSuccess: async (data, variables) => {
      await queryClient.invalidateQueries({
        queryKey: [
          queryKey.CONVERSATION,
          queryKey.GET_OFFLINE_CONVERSATION_DETAIL,
          variables.conversationId,
        ],
      });
    },
  });
}
