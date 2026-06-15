import { deleteSecure, setSecure } from "@/util/storage";
import { useMutation } from "@tanstack/react-query";
import { AxiosError } from "axios";
import Toast from "react-native-toast-message";
import {
  deleteAccount,
  loginInWithEmail,
  requestSMSOTP,
  signInWithApple,
  signUpWithEmail,
  verifyEmailOTP,
  verifySMSOTP,
} from "@/api/auth";

export function useSignupWithEmail() {
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

export function useLoginWithEmail() {
  return useMutation({
    mutationFn: loginInWithEmail,
    onSuccess: async (data) => {
      if (!data.emailVerified) {
        await setSecure("verificationId", data.verificationId ?? "");
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

export function useRequestSMSOTP() {
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

export function useVerifyEmailOTP() {
  return useMutation({
    mutationFn: verifyEmailOTP,
    onSuccess: async (data) => {
      console.log(data.sessionId);
      await setSecure("sessionId", data?.sessionId ?? "");
    },
    onError: (error: AxiosError) => {
      Toast.show({
        type: "error",
        text1: String(error?.response?.data),
      });
    },
  });
}

export function useVerifySMSOTP() {
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

export function useSignInWithApple() {
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

export function useDeleteAccount() {
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
