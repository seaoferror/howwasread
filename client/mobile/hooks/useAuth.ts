import { deleteSecure, setSecure } from "@/db/storage";
import { useMutation } from "@tanstack/react-query";
import { AxiosError } from "axios";
import Toast from "react-native-toast-message";
import {
  deleteAccount,
  loginInWithEmail,
  logout,
  requestSMSOTP,
  signInWithApple,
  signInWithGoogle,
  signUpWithEmail,
  verifyEmailOTP,
  verifySMSOTP,
} from "@/api/auth";

export function useSignupWithEmail() {
  return useMutation({
    mutationFn: signUpWithEmail,
    onSuccess: async (data) => {
      if (data.verificationId) {
        await setSecure("verificationId", data.verificationId);
      }
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
      if (data.verificationId) {
        await setSecure("verificationId", data.verificationId);
        return;
      }
      if (data.sessionId) {
        await setSecure("sessionId", data.sessionId);
        return;
      }
      if (data.accessToken) {
        await setSecure("accessToken", data.accessToken);
      }
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
    onSuccess: async (data) => {
      await setSecure("verificationId", data.verificationId);
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
      if (data.sessionId) {
        await setSecure("sessionId", data?.sessionId);
      }
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
    onSuccess: async (data) => {
      if (data.accessToken) {
        await setSecure("accessToken", data.accessToken);
      }
    },
    onError: (error: AxiosError) => {
      Toast.show({
        type: "error",
        text1: String(error?.response?.data),
      });
    },
  });
}

export function useSignInWithGoogle() {
  return useMutation({
    mutationFn: signInWithGoogle,
    onSuccess: async (data) => {
      if (data.sessionId) {
        await setSecure("sessionId", data.sessionId);
        return;
      }
      if (data.accessToken) {
        await setSecure("accessToken", data.accessToken);
      }
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
      if (data.sessionId) {
        await setSecure("sessionId", data.sessionId);
        return;
      }
      if (data.accessToken) {
        await setSecure("accessToken", data.accessToken);
      }
    },
    onError: (error: AxiosError) => {
      Toast.show({
        type: "error",
        text1: String(error?.response?.data),
      });
    },
  });
}

export function useLogout() {
  return useMutation({
    mutationFn: logout,
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
