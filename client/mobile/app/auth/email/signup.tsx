import { StyleSheet, View } from "react-native";
import { FormProvider, useForm } from "react-hook-form";
import EmailInput from "@/components/auth/EmailInput";
import PasswordInput from "@/components/auth/PasswordInput";
import PasswordConfirmInput from "@/components/auth/PasswordConfirmInput";
import { useSignupWithEmail } from "@/hooks/useAuth";
import { colors } from "@/constants";
import { router } from "expo-router";
import CustomButton from "@/components/CustomButton";

interface FormValue {
  email: string;
  password: string;
  passwordConfirm: string;
}

export default function SignupScreen() {
  const signUpWithEmailMutation = useSignupWithEmail();

  const emailSignupForm = useForm<FormValue>({
    defaultValues: {
      email: "",
      password: "",
      passwordConfirm: "",
    },
  });

  const onSubmit = (formValues: FormValue) => {
    const { email, password } = formValues;

    signUpWithEmailMutation.mutate(
      {
        email,
        password,
      },
      {
        onSuccess: (data) => {
          if (data.verificationId) {
            router.replace("/auth/otp/email");
          }
        },
      },
    );
  };
  return (
    <FormProvider {...emailSignupForm}>
      <View style={styles.container}>
        <View style={styles.content}>
          <EmailInput />
          <PasswordInput submitBehavior="submit" />
          <PasswordConfirmInput />
          <View style={{marginTop: 30}}>
            <CustomButton
              label="Sign up"
              onPress={emailSignupForm.handleSubmit(onSubmit)}
              disabled={signUpWithEmailMutation.isPending}
            />
          </View>
        </View>
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
    backgroundColor: colors.SAND_110,
  },
});
