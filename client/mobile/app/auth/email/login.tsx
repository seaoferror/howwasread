import { StyleSheet, View } from "react-native";
import { FormProvider, useForm } from "react-hook-form";
import FixedBottomCTA from "@/components/FixedBottomCTA";
import EmailInput from "@/components/auth/EmailInput";
import PasswordInput from "@/components/auth/PasswordInput";
import { useAuth } from "@/hooks/useAuth";
import { colors } from "@/constants";
import { router } from "expo-router";

interface FormValue {
  email: string;
  password: string;
  passwordConfirm: string;
}

export default function LoginScreen() {
  const { loginWithEmailMutation } = useAuth();

  const emailLoginForm = useForm<FormValue>({
    defaultValues: {
      email: "",
      password: "",
    },
  });

  const onSubmit = (formValues: FormValue) => {
    const { email, password } = formValues;

    loginWithEmailMutation.mutate(
      {
        email,
        password,
      },
      {
        onSuccess: (data) => {
          console.log("data.emailVerified:", data?.emailVerified);
          console.log("data.phoneNumberVerified:", data?.phoneNumberVerified);
          if (!data.emailVerified) {
            router.replace("/auth/otp/email");
            return;
          }
          if (!data.phoneNumberVerified) {
            router.replace("/auth/phone-number");
            return;
          }
          if (data.accessToken) {
            router.replace("/");
          }
        },
      },
    );
  };
  return (
    <FormProvider {...emailLoginForm}>
      <View style={styles.container}>
        <View style={styles.content}>
          <EmailInput />
          <PasswordInput />
        </View>
        <FixedBottomCTA
          label="login"
          onPress={emailLoginForm.handleSubmit(onSubmit)}
          disabled={loginWithEmailMutation.isPending}
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
