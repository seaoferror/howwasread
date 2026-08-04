import { Pressable, StyleSheet, Text, View } from "react-native";
import { SafeAreaView } from "react-native-safe-area-context";
import { FormProvider, useForm } from "react-hook-form";
import { router, useLocalSearchParams, useNavigation } from "expo-router";
import { useSetName } from "@/hooks/useProfile";
import NameInput from "@/components/profile/NameInput";
import { colors } from "@/constants";
import { useEffect } from "react";
import { Ionicons } from "@expo/vector-icons";
import { setKVStore } from "@/db/storage";
import CustomButton from "@/components/CustomButton";

interface FormValue {
  name: string;
}

export default function NameScreen() {
  const { newcomer } = useLocalSearchParams();
  const nameForm = useForm<FormValue>({
    defaultValues: {
      name: "",
    },
  });
  const setNameMutation = useSetName();
  const navigation = useNavigation();

  useEffect(() => {
    if (!newcomer) {
      navigation.setOptions({
        headerLeft: () => (
          <Pressable onPress={() => router.back()}>
            <Ionicons name="chevron-back" size={28} color="black" />
          </Pressable>
        ),
      });
    }
  }, [newcomer]);

  const onSubmit = async (formValue: FormValue) => {
    const { name } = formValue;

    setNameMutation.mutate(
      { name: name },
      {
        onSuccess: () => {
          setKVStore("myName", name);
          if (newcomer) {
            router.replace("/conversations");
            return;
          }
          router.replace("/account");
        },
      },
    );
  };

  return (
    <FormProvider {...nameForm}>
      <SafeAreaView style={styles.container}>
        <Text style={styles.guideLine}>
          {
            "We strongly recommend to set a name \n which is convenient to be called by others"
          }
        </Text>
        <View style={styles.content}>
          <NameInput />
        </View>
        <CustomButton
          label="Set this name"
          onPress={nameForm.handleSubmit(onSubmit)}
          disabled={setNameMutation.isPending}
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
    paddingHorizontal: 40,
  },
  content: {
    width: "100%",
    paddingHorizontal: 57,
    marginTop: 58,
    marginBottom: 97,
  },
  guideLine: {
    fontSize: 14,
    textAlign: "center",
  },
});
