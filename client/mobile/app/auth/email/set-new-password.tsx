import { StyleSheet, View } from "react-native";
import { FormProvider, useForm } from "react-hook-form";
import PasswordInput from "@/components/auth/PasswordInput";
import PasswordConfirmInput from "@/components/auth/PasswordConfirmInput";
import { colors } from "@/constants";
import { router } from "expo-router";
import CustomButton from "@/components/CustomButton";
import { getSecureAsync } from "@/db/storage";
import { useSetNewPassword } from "@/hooks/useAuth";

interface FormValue {
  password: string;
  passwordConfirm: string;
}

export default function SetNewPasswordScreen() {
  const setNewPasswordMutation = useSetNewPassword();

  const passwordForm = useForm<FormValue>({
    defaultValues: {
      password: "",
      passwordConfirm: "",
    },
  });

  const onSubmit = async (formValues: FormValue) => {
    const { password } = formValues;

    setNewPasswordMutation.mutate(
      {
        sessionId: await getSecureAsync("sessionId"),
        password,
      },
      {
        onSuccess: () => {
          router.replace("/auth/email/login");
        },
      },
    );
  };
  return (
    <FormProvider {...passwordForm}>
      <View style={styles.container}>
        <View style={styles.content}>
          <PasswordInput submitBehavior="submit" />
          <PasswordConfirmInput />
        </View>
        <CustomButton
          label="Sign up"
          onPress={passwordForm.handleSubmit(onSubmit)}
          disabled={setNewPasswordMutation.isPending}
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
    backgroundColor: colors.SAND_110,
  },
});
