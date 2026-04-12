import { StyleSheet, Text, View } from "react-native";
import { SafeAreaView } from "react-native-safe-area-context";
import { colors } from "@/constants";
import FixedBottomCTA from "@/components/FixedBottomCTA";
import { FormProvider, useForm } from "react-hook-form";
import NameInput from "@/components/profile/NameInput";
import { useMyProfile } from "@/hooks/useMyProfile";
import { router } from "expo-router";

interface FormValue {
  name: string;
}

export default function NameScreen() {
  const nameForm = useForm<FormValue>({
    defaultValues: {
      name: "",
    },
  });
  const { setNameMutation } = useMyProfile();

  const onSubmit = async (formValue: FormValue) => {
    const { name } = formValue;

    setNameMutation.mutate(
      { name: name },
      {
        onSuccess: () => {
          router.push("/conversations");
        },
      },
    );
  };

  return (
    <FormProvider {...nameForm}>
      <SafeAreaView style={styles.container}>
        <Text style={styles.guideLine}>
          We strongly recommend to set a name which is convenient to be called
          by others
        </Text>
        <View style={styles.content}>
          <NameInput />
        </View>
        <FixedBottomCTA
          label="Set this name"
          onPress={nameForm.handleSubmit(onSubmit)}
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
    paddingTop: 100,
    gap: 50,
  },
  guideLine: {
    position: "absolute",
    top: 40,
    paddingHorizontal: 30,
    fontSize: 15,
  },
});
