import { queryKey } from "@/constants";
import { deleteSecure, getSecureAsync, setSecure } from "@/util/storage";
import { useMutation, useQuery } from "@tanstack/react-query";
import { AxiosError } from "axios";
import Toast from "react-native-toast-message";
import {
  deleteAccount,
  getMyId,
  loginInWithEmail,
  requestSMSOTP,
  signInWithApple,
  signUpWithEmail,
  verifyEmailOTP,
  verifySMSOTP,
} from "@/api/auth";

function useGetMyId() {
  const { data } = useQuery({
    queryFn: getMyId,
    queryKey: [queryKey.AUTH, queryKey.GET_MY_ID],
  });

  return { data };
}

function useSignupWithEmail() {
  return useMutation({
    mutationFn: signUpWithEmail,
    onSuccess: async (data) => {
      setSecure("verificationId", data.verificationId);
      console.log("success to save verification Id");
    },
    onError: (error: AxiosError) => {
      Toast.show({
        type: "error",
        text1: String(error?.response?.data),
      });
    },
  });
}

function useLoginWithEmail() {
  return useMutation({
    mutationFn: loginInWithEmail,
    onSuccess: async (data) => {
      if (!data.emailVerified) {
        setSecure("verificationId", data.verificationId ?? "");
        const v = await getSecureAsync("verificationId");
        console.log(v);
        return;
      }
      if (!data.phoneNumberVerified) {
        setSecure("sessionId", data?.sessionId ?? "");
        return;
      }
      setSecure("accessToken", data?.accessToken ?? "");
    },
    onError: (error: AxiosError) => {
      Toast.show({
        type: "error",
        text1: String(error?.response?.data),
      });
    },
  });
}

function useRequestSMSOTP() {
  return useMutation({
    mutationFn: requestSMSOTP,
    onSuccess: (data) => {
      setSecure("verificationId", data.verificationId);
      console.log("success to save verificationId");
    },
    onError: (error: AxiosError) => {
      Toast.show({
        type: "error",
        text1: String(error?.response?.data),
      });
    },
  });
}

function useVerifyEmailOTP() {
  return useMutation({
    mutationFn: verifyEmailOTP,
    onSuccess: async (data) => {
      console.log(data.sessionId);
      setSecure("sessionId", data?.sessionId ?? "");
      console.log(await getSecureAsync("sessionId"));
    },
    onError: (error: AxiosError) => {
      Toast.show({
        type: "error",
        text1: String(error?.response?.data),
      });
    },
  });
}

function useVerifySMSOTP() {
  return useMutation({
    mutationFn: verifySMSOTP,
    onSuccess: (data) => {
      if (data.phoneNumberVerified)
        setSecure("accessToken", data?.accessToken ?? "");
    },
    onError: (error: AxiosError) => {
      Toast.show({
        type: "error",
        text1: String(error?.response?.data),
      });
    },
  });
}

function useSignInWithApple() {
  return useMutation({
    mutationFn: signInWithApple,
    onSuccess: async (data) => {
      if (!data.phoneNumberVerified) {
        await setSecure("sessionId", data?.sessionId ?? "");
        return;
      }
      await setSecure("accessToken", data?.accessToken ?? "");
    },
    onError: (error: AxiosError) => {
      Toast.show({
        type: "error",
        text1: String(error?.response?.data),
      });
    },
  });
}

function useDeleteAccount() {
  return useMutation({
    mutationFn: deleteAccount,
    onSuccess: async () => {
      Toast.show({
        type: "success",
        text1: `Bye, see you.`,
      });
      await deleteSecure("accessToken");
    },
    onError: (error: AxiosError) => {
      Toast.show({
        type: "error",
        text1: String(error?.response?.data),
      });
    },
  });
}

export function useAuth() {
  const { data } = useGetMyId();
  const signUpWithEmailMutation = useSignupWithEmail();
  const loginWithEmailMutation = useLoginWithEmail();
  const verifyEmailOTPMutation = useVerifyEmailOTP();
  const requestSMSOTPMutation = useRequestSMSOTP();
  const verifySMSOTPMutation = useVerifySMSOTP();
  const signInWithAppleMutation = useSignInWithApple();
  const deleteAccountMutation = useDeleteAccount();

  return {
    id: data?.id,
    signUpWithEmailMutation,
    loginWithEmailMutation,
    verifyEmailOTPMutation,
    requestSMSOTPMutation,
    verifySMSOTPMutation,
    signInWithAppleMutation,
    deleteAccountMutation,
  };
}
