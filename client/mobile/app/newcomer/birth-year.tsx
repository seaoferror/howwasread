import { StyleSheet, Text, View } from "react-native";
import { SafeAreaView } from "react-native-safe-area-context";
import { colors } from "@/constants";
import FixedBottomCTA from "@/components/FixedBottomCTA";
import { FormProvider, useForm } from "react-hook-form";
import BirthYearBox from "@/components/profile/BirthYearBox";
import { useMyProfile } from "@/hooks/useMyProfile";
import { router } from "expo-router";

interface FormValue {
  birthYear: number;
}

export default function BirthYearScreen() {
  const birthYearForm = useForm<FormValue>({
    defaultValues: { birthYear: 2000 },
  });
  const { setBirthYearMutation } = useMyProfile();

  const onSubmit = async (formValue: FormValue) => {
    const { birthYear } = formValue;

    setBirthYearMutation.mutate(
      { birthYear },
      {
        onSuccess: () => {
          // router.push("/newcomer/ethnicity");
        },
      },
    );
  };

  return (
    <FormProvider {...birthYearForm}>
      <SafeAreaView style={styles.container}>
        <Text style={styles.guideLine}>
          When were you born? This is non-modifiable
        </Text>
        <View style={styles.content}>
          <BirthYearBox />
        </View>
        <FixedBottomCTA
          label="Set birth year"
          onPress={birthYearForm.handleSubmit(onSubmit)}
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
