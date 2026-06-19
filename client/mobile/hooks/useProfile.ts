import { useMutation, useQuery } from "@tanstack/react-query";
import { getMyProfile, getProfile, setBirthYear, setName } from "@/api/profile";
import Toast from "react-native-toast-message";
import { queryKey } from "@/constants";

export function useGetMyProfile() {
  return useQuery({
    queryFn: getMyProfile,
    queryKey: [queryKey.PROFILE, queryKey.GET_MY_PROFILE],
  },);

}

export function useGetProfile(id: string) {
  return useQuery({
    queryFn: () => getProfile(id),
    queryKey: [queryKey.PROFILE, queryKey.GET_PROFILE, id]
  });
}

export function useSetName() {
  return useMutation({
    mutationFn: setName,
    onSuccess: async () => {},
    onError: (error) => {
      Toast.show({
        type: "error",
        text1: error.message,
      });
    },
  });
}

export function useSetBirthYear() {
  return useMutation({
    mutationFn: setBirthYear,
    onSuccess: async () => {},
    onError: (error) => {
      Toast.show({
        type: "error",
        text1: error.message,
      });
    },
  });
}