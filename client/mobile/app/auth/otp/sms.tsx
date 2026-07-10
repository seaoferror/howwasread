import { StyleSheet, View } from "react-native";
import { SafeAreaView } from "react-native-safe-area-context";
import { colors } from "@/constants";
import { FormProvider, useForm } from "react-hook-form";
import { useVerifySMSOTP } from "@/hooks/useAuth";
import { router } from "expo-router";
import { getSecureAsync, setKVStore } from "@/db/storage";
import Toast from "react-native-toast-message";
import OTPInput from "@/components/auth/OTPInput";
import { getMyProfile } from "@/api/profile";
import CustomButton from "@/components/CustomButton";

interface FormValue {
  otp: string;
}

export default function SMSOTPScreen() {
  const verifySMSOTPMutation = useVerifySMSOTP();
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
      {
        otp,
        verificationId,
        sessionId: (await getSecureAsync("sessionId")) || null,
      },
      {
        onSuccess: async (data) => {
          if (data.accessToken) {
            const my = await getMyProfile();
            setKVStore("myId", my.id);
            if (!my.name) {
              router.replace("/profile/name");
              return;
            }
            setKVStore("myName", my.name);
            router.replace("/conversations");
          }
        },
      },
    );
  };

  return (
    <FormProvider {...smsOTPForm}>
      <SafeAreaView style={styles.container}>
        <View style={styles.content}>
          <OTPInput />
          <CustomButton
            label="Confirm"
            onPress={smsOTPForm.handleSubmit(onSubmit)}
            disabled={verifySMSOTPMutation.isPending}
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
