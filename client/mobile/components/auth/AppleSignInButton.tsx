import { StyleSheet, View } from "react-native";
import {
  AppleAuthenticationButton,
  AppleAuthenticationButtonStyle,
  AppleAuthenticationButtonType,
  AppleAuthenticationScope,
  signInAsync,
} from "expo-apple-authentication";
import { useAuth } from "@/hooks/useAuth";
import { randomUUID } from "expo-crypto";
import { router } from "expo-router";

export default function AppleSignInButton() {
  const { signInWithAppleMutation } = useAuth();

  const onSignIn = async () => {
    const rawNonce = randomUUID();
    const credential = await signInAsync({
      requestedScopes: [AppleAuthenticationScope.EMAIL],
      nonce: rawNonce,
    });
    const idt = credential.identityToken ?? "";
    if (!idt) {
      return;
    }
    signInWithAppleMutation.mutate(
      {
        user: credential.user,
        email: credential.email,
        identityToken: idt,
        nonce: rawNonce,
      },
      {
        onSuccess: async (data) => {
          if (!data.phoneNumberVerified) {
            router.push("/auth/phone-number");
            return;
          }
          router.replace("/");
        },
      },
    );
  };

  return (
    <View>
      <AppleAuthenticationButton
        buttonType={AppleAuthenticationButtonType.SIGN_IN}
        buttonStyle={AppleAuthenticationButtonStyle.BLACK}
        cornerRadius={5}
        style={styles.appleButton}
        onPress={onSignIn}
      />
    </View>
  );
}

const styles = StyleSheet.create({
  appleButton: {
    width: 300,
    height: 44,
  },
});
