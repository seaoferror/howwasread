import { StyleSheet, View } from "react-native";
import { SafeAreaView } from "react-native-safe-area-context";
import { colors } from "@/constants";
import FixedBottomCTA from "@/components/FixedBottomCTA";
import { FormProvider, useForm } from "react-hook-form";
import { useAuth } from "@/hooks/useAuth";
import { router } from "expo-router";
import { getSecureAsync } from "@/util/storage";
import Toast from "react-native-toast-message";
import OTPInput from "@/components/auth/OTPInput";

interface FormValue {
  otp: string;
}

export default function SMSScreen() {
  const { verifySMSOTPMutation } = useAuth();
  const smsOTPForm = useForm<FormValue>({
    defaultValues: {
      otp: "",
    },
  });

  const onSubmit = async (formValue: FormValue) => {
    const { otp } = formValue;
    const verificationId = await getSecureAsync("verificationId");
    if (!verificationId) {
      console.error("fail to get verification Id");
      Toast.show({
        type: "error",
        text1: "you can't send code",
      });
      return;
    }
    console.log("execute post sms otp mutate");
    verifySMSOTPMutation.mutate(
      { otp, verificationId, sessionId: await getSecureAsync("sessionId") },
      {
        onSuccess: () => router.replace("/"),
      },
    );
  };

  return (
    <FormProvider {...smsOTPForm}>
      <SafeAreaView style={styles.container}>
        <View style={styles.content}>
          <OTPInput />
        </View>
        <FixedBottomCTA
          label="Confirm"
          onPress={smsOTPForm.handleSubmit(onSubmit)}
        />
      </SafeAreaView>
    </FormProvider>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: colors.SAND_110,
    borderTopWidth: StyleSheet.hairlineWidth,
    borderColor: colors.GRAY_700,
  },
  content: {
    flex: 1,
    width: "100%",
    paddingHorizontal: 100,
    paddingTop: 120,
  },
});
