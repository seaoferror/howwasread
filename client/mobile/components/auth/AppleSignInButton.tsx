import { StyleSheet, View } from "react-native";
import {
  AppleAuthenticationButton,
  AppleAuthenticationButtonStyle,
  AppleAuthenticationButtonType,
  AppleAuthenticationScope,
  signInAsync,
} from "expo-apple-authentication";
import { useSignInWithApple } from "@/hooks/useAuth";
import { randomUUID } from "expo-crypto";
import { router } from "expo-router";

export default function AppleSignInButton() {
  const signInWithAppleMutation = useSignInWithApple();

  const onSignIn = async () => {
    const rawNonce = randomUUID();
    const credential = await signInAsync({
      requestedScopes: [AppleAuthenticationScope.EMAIL],
      nonce: rawNonce,
    });
    const idt = credential.identityToken;
    if (!idt) {
      return;
    }
    signInWithAppleMutation.mutate(
      {
        isFirstSignIn: !!credential.email,
        identityToken: idt,
      },
      {
        onSuccess: async (data) => {
          if (data.sessionId) {
            router.push("/auth/phone-number");
            return;
          }
          if (data.accessToken){
            router.replace("/");
          }
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
    width: 312,
    height: 44,
  },
});
