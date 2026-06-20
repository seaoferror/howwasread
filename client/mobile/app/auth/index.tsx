import { SafeAreaView } from "react-native-safe-area-context";
import { Platform, StyleSheet, View } from "react-native";
import CustomButton from "@/components/CustomButton";
import { Link, router, useFocusEffect } from "expo-router";
import { colors } from "@/constants";
import AppleSignInButton from "@/components/auth/AppleSignInButton";
import { deleteSecure } from "@/util/storage";
import {
  GoogleSignin,
  GoogleSigninButton,
} from "@react-native-google-signin/google-signin";
import { useSignInWithGoogle } from "@/hooks/useAuth";
import { Image } from "expo-image";

const AppIcon = require("@/assets/images/icon_transparent_background.png");

export default function AuthScreen() {
  const googleSignInMutation = useSignInWithGoogle();
  useFocusEffect(() => {
    deleteSecure("verificationId");
    deleteSecure("sessionId");
  });

  const handleGoogleSignInButtonPress = async () => {
    const response = await GoogleSignin.signIn();
    console.log(response);
    if (response.idToken) {
      googleSignInMutation.mutate(
        {
          idToken: response.idToken,
        },
        {
          onSuccess: (data) => {
            if (data.sessionId) {
              router.push("/auth/phone-number");
            }
            if (data.accessToken) {
              router.push("/");
            }
          },
        },
      );
    }
  };

  return (
    <SafeAreaView style={styles.container}>
      <View style={styles.buttonContainer}>
        <Image style={styles.icon} source={AppIcon} />
        <CustomButton
          label={"Start with your Phone number"}
          onPress={() => router.push("/auth/phone-number")}
        />
        <GoogleSigninButton
          size={GoogleSigninButton.Size.Wide}
          color={GoogleSigninButton.Color.Dark}
          onPress={handleGoogleSignInButtonPress}
        />
        {Platform.OS === "ios" && <AppleSignInButton />}
        <View style={styles.emailContainer}>
          <CustomButton
            label={"Start with your Email"}
            onPress={() => router.push("/auth/email/signup")}
          />
          <Link href={"/auth/email/login"} style={styles.signupText}>
            Do you have account? Login with your email
          </Link>
        </View>
      </View>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: colors.SAND_110,
    borderTopWidth: StyleSheet.hairlineWidth,
    borderColor: colors.GRAY_700,
  },
  buttonContainer: {
    flex: 1,
    width: "100%",
    alignItems: "center",
    paddingHorizontal: 30,
    marginTop: 90,
    gap: 40,
  },
  icon: {
    width: 160,
    height: 120,
  },
  signupText: {
    textAlign: "center",
    textDecorationLine: "underline",
    fontSize: 15,
    marginTop: 20,
  },
  emailContainer: {
    flex: 1,
    width: "100%",
  },
});
