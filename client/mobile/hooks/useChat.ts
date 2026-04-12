import { useMutation, useQuery } from "@tanstack/react-query";
import { getChatRoomInfo, sendLike } from "@/api/chat";
import { AxiosError } from "axios";
import Toast from "react-native-toast-message";
import { queryKey } from "@/constants";

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

export function useGetChatRoomInfo(roomId: string) {
  const { data } = useQuery({
    queryFn: () => getChatRoomInfo(roomId),
    queryKey: [queryKey.CHAT, queryKey.GET_CHAT_ROOM_INFO, roomId],
  })
  return { data }
}
