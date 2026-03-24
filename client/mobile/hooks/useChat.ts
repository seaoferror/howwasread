import { useMutation } from "@tanstack/react-query";
import { sendLike } from "@/api/chat";
import { AxiosError } from "axios";
import Toast from "react-native-toast-message";

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
