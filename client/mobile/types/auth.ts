export interface SignInWithEmailRequest {
  email: string;
  password: string;
}

export interface VerifyEmailOTPRequest {
  otp: string;
  verificationId: string;
}

export interface VerifyEmailOTPResponse {
  sessionId?: string;
}

export interface SendSMSOTPRequest {
  phoneNumber: string;
  sessionId: string | null;
}

export interface VerifySMSOTPRequest {
  otp: string;
  verificationId: string;
  sessionId: string | null;
}

export interface VerifySMSOTPResponse {
  accessToken?: string;
}

export interface LoginWithEmailResponse {
  verificationId?: string;
  sessionId?: string;
  accessToken?: string;
}

export interface SignInWithAppleRequest {
  identityToken: string;
}

export interface SignInWithGoogleRequest {
  idToken: string;
}

export interface SignInWithThirdPartyResponse {
  sessionId?: string;
  accessToken?: string;
}