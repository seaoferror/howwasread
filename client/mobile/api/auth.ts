import { axiosInstance } from "@/api/axios";
import {
  LoginWithEmailResponse,
  SendSMSOTPRequest,
  SignInWithAppleRequest,
  SignInWithEmailRequest,
  SignInWithGoogleRequest,
  SignInWithThirdPartyResponse,
  VerifyEmailOTPRequest,
  VerifyEmailOTPResponse,
  VerifySMSOTPRequest,
  VerifySMSOTPResponse,
} from "@/types/auth";

export async function signUpWithEmail(
  body: SignInWithEmailRequest,
): Promise<{ verificationId?: string }> {
  console.log("post email sign up");
  const { data } = await axiosInstance.post("/auth/email/create", body);
  return data;
}

export async function loginInWithEmail(
  body: SignInWithEmailRequest,
): Promise<LoginWithEmailResponse> {
  const { data } = await axiosInstance.post("/auth/email/login", body);
  console.log(data);
  return data;
}

export async function requestEmailOTP(body: {
  id: string;
}): Promise<{ verificationId: string }> {
  const { data } = await axiosInstance.post("/auth/email/otp/send", body);
  return data;
}

export async function requestSMSOTP(
  body: SendSMSOTPRequest,
): Promise<{ verificationId: string }> {
  const { data } = await axiosInstance.post("/auth/sms/otp/send", body);
  return data;
}

export async function verifyEmailOTP(
  body: VerifyEmailOTPRequest,
): Promise<VerifyEmailOTPResponse> {
  const { data } = await axiosInstance.post("/auth/email/otp/verify", body);
  return data;
}

export async function verifySMSOTP(
  body: VerifySMSOTPRequest,
): Promise<VerifySMSOTPResponse> {
  const { data } = await axiosInstance.post("/auth/sms/otp/verify", body);
  return data;
}

export async function signInWithApple(
  body: SignInWithAppleRequest,
): Promise<SignInWithThirdPartyResponse> {
  const { data } = await axiosInstance.post("/auth/email/apple", body);
  return data;
}

export async function signInWithGoogle(
  body: SignInWithGoogleRequest,
): Promise<SignInWithThirdPartyResponse> {
  const { data } = await axiosInstance.post("/auth/email/google", body);
  return data;
}

export async function logout() {
  const { data } = await axiosInstance.post("/auth/account/logout");
  return data;
}

export async function deleteAccount() {
  const { data } = await axiosInstance.delete("/auth/account/delete");
  return data;
}

export async function forgetPassword(body: {
  email: string;
}): Promise<{ verificationId: string }> {
  const { data } = await axiosInstance.post("/auth/email/password/forget", body);
  return data;
}

export async function setNewPassword(body: {
  password: string;
  sessionId: string;
}): Promise<void> {
  const { data } = await axiosInstance.put(
    "/auth/email/password/set-new",
    body,
  );
  return data
}