import { useMutation, useQuery } from "@tanstack/react-query";
import Toast from "react-native-toast-message";
import { getSecureStore, saveSecureStore } from "@/util/secureStore";
import {
  getMyId,
  loginInWithEmail,
  signUpWithEmail,
  requestSMSOTP,
  signInWithApple,
  verifyEmailOTP,
  verifySMSOTP,
} from "@/api/auth";
import { queryKey } from "@/constants";
import { AxiosError } from "axios";

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
      saveSecureStore("verificationId", data.verificationId);
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
        saveSecureStore("verificationId", data.verificationId ?? "");
        const v = await getSecureStore("verificationId");
        console.log(v);
        return;
      }
      if (!data.phoneNumberVerified) {
        saveSecureStore("sessionId", data?.sessionId ?? "");
        return;
      }
      saveSecureStore("accessToken", data?.accessToken ?? "");
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
      saveSecureStore("verificationId", data.verificationId);
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
      saveSecureStore("sessionId", data?.sessionId ?? "");
      console.log(await getSecureStore("sessionId"));
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
        saveSecureStore("accessToken", data?.accessToken ?? "");
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
        saveSecureStore("sessionId", data?.sessionId ?? "");
        return;
      }
      saveSecureStore("accessToken", data?.accessToken ?? "");
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

  return {
    id: data?.id,
    signUpWithEmailMutation,
    loginWithEmailMutation,
    verifyEmailOTPMutation,
    requestSMSOTPMutation,
    verifySMSOTPMutation,
    signInWithAppleMutation,
  };
}
