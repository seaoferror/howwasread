import { useMutation } from "@tanstack/react-query";
import { AxiosError } from "axios";
import Toast from "react-native-toast-message";
import { registerNotification } from "@/api/notification";

export function useRegisterNotification() {
  return useMutation({
    mutationFn: registerNotification,
    onError: (error: AxiosError) => {
      console.log(error?.response?.data);
      Toast.show({
        type: "error",
        text1: String(error?.response?.data),
      });
    },
  })
}
