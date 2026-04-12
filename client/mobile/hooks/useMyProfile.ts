import { useMutation, useQuery } from "@tanstack/react-query";
import { getMyProfile, getProfile, setBirthYear, setName } from "@/api/profile";
import Toast from "react-native-toast-message";
import { queryKey } from "@/constants";

function useGetMyProfile() {
  const { data } = useQuery({
    queryFn: getMyProfile,
    queryKey: [queryKey.PROFILE, queryKey.GET_MY_PROFILE],
  });

  return { data };
}

export function useGetProfile(id: string) {
  const { data } = useQuery({
    queryFn: () => getProfile(id),
    queryKey: [queryKey.PROFILE, queryKey.GET_PROFILE, id]
  });

  return { data }
}

function useSetName() {
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

function useSetBirthYear() {
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

export function useProfile() {
  const { data } = useGetMyProfile();

  const setNameMutation = useSetName();
  const setBirthYearMutation = useSetBirthYear();

  return {
    profile: {
      id: data?.id ?? "",
      name: data?.name ?? "",
    },
    setNameMutation: setNameMutation,
    setBirthYearMutation,
  };
}
