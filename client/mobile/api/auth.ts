import { localDevInstance } from "@/api/axios";
import {
  LoginWithEmailResponse,
  SendSMSOTPRequest,
  SignInWithAppleRequest,
  SignInWithAppleResponse,
  SignInWithEmailRequest,
  VerifyEmailOTPRequest,
  VerifyEmailOTPResponse,
  VerifySMSOTPRequest,
  VerifySMSOTPResponse,
} from "@/types/auth";

export async function signUpWithEmail(
  body: SignInWithEmailRequest,
): Promise<{ verificationId: string }> {
  console.log("post email sign up");
  const { data } = await localDevInstance.post("/auth/email/create", body);
  return data;
}

export async function loginInWithEmail(
  body: SignInWithEmailRequest,
): Promise<LoginWithEmailResponse> {
  const { data } = await localDevInstance.post("/auth/email/login", body);
  console.log(data);
  return data;
}

export async function requestEmailOTP(body: {
  id: string;
}): Promise<{ verificationId: string }> {
  const { data } = await localDevInstance.post("/auth/email/otp/send", body);
  return data;
}

export async function requestSMSOTP(
  body: SendSMSOTPRequest,
): Promise<{ verificationId: string }> {
  const { data } = await localDevInstance.post("/auth/sms/otp/send", body);
  return data;
}

export async function verifyEmailOTP(
  body: VerifyEmailOTPRequest,
): Promise<VerifyEmailOTPResponse> {
  const { data } = await localDevInstance.post("/auth/email/otp/verify", body);
  return data;
}

export async function verifySMSOTP(
  body: VerifySMSOTPRequest,
): Promise<VerifySMSOTPResponse> {
  const { data } = await localDevInstance.post("/auth/sms/otp/verify", body);
  return data;
}

export async function signInWithApple(
  body: SignInWithAppleRequest,
): Promise<SignInWithAppleResponse> {
  const { data } = await localDevInstance.post("/auth/email/apple", body);
  return data;
}

export async function getMyId(): Promise<{ id: string }> {
  const { data } = await localDevInstance.get("/auth/my-id");
  return data;
}

export async function deleteAccount() {
  const { data } = await localDevInstance.post("/auth/delete");
  return data;
}
