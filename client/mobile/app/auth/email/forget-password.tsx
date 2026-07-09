import { StyleSheet, View } from "react-native";
import { FormProvider, useForm } from "react-hook-form";
import { router } from "expo-router";
import EmailInput from "@/components/auth/EmailInput";
import { colors } from "@/constants";
import CustomButton from "@/components/CustomButton";
import { useForgetPassword } from "@/hooks/useAuth";

interface FormValue {
  email: string;
}

export default function ForgetPasswordScreen() {
  const forgetPasswordMutation = useForgetPassword();

  const forgetPasswordForm = useForm<FormValue>({
    defaultValues: {
      email: "",
    },
  });
  const onSubmit = (formValues: FormValue) => {
    const { email } = formValues;

    forgetPasswordMutation.mutate(
      {
        email,
      },
      {
        onSuccess: async (data) => {
          if (data.verificationId) {
            router.replace("/auth/otp/password");
          }
        },
      },
    );
  };
  return (
    <FormProvider {...forgetPasswordForm}>
      <View style={styles.container}>
        <View style={styles.content}>
          <EmailInput />
        </View>
        <CustomButton
          label="login"
          onPress={forgetPasswordForm.handleSubmit(onSubmit)}
          disabled={forgetPasswordMutation.isPending}
        />
      </View>
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
    margin: 16,
    gap: 16,
    paddingHorizontal: 20,
    backgroundColor: colors.SAND_110,
  },
});
