import { StyleSheet, View } from "react-native";
import { SafeAreaView } from "react-native-safe-area-context";
import { colors } from "@/constants";
import { FormProvider, useForm } from "react-hook-form";
import { useVerifyEmailOTP } from "@/hooks/useAuth";
import { router } from "expo-router";
import { getSecureAsync } from "@/db/storage";
import Toast from "react-native-toast-message";
import OTPInput from "@/components/auth/OTPInput";
import CustomButton from "@/components/CustomButton";

interface FormValue {
  otp: string;
}

export default function PasswordOTPScreen() {
  const verifyEmailOTPMutation = useVerifyEmailOTP();

  const emailOTPForm = useForm<FormValue>({
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
    console.log("execute verify email otp mutate");
    verifyEmailOTPMutation.mutate(
      {
        verificationId,
        otp,
      },
      {
        onSuccess: () => router.replace("/auth/email/set-new-password"),
      },
    );
  };

  return (
    <FormProvider {...emailOTPForm}>
      <SafeAreaView style={styles.container}>
        <View style={styles.content}>
          <OTPInput />
          <CustomButton
            label="Confirm"
            onPress={emailOTPForm.handleSubmit(onSubmit)}
            disabled={verifyEmailOTPMutation.isPending}
          />
        </View>
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
    gap: 50,
  },
});
