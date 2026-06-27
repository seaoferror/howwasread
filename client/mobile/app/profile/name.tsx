import { Pressable, StyleSheet, Text, View } from "react-native";
import { SafeAreaView } from "react-native-safe-area-context";
import { FormProvider, useForm } from "react-hook-form";
import { router, useLocalSearchParams, useNavigation } from "expo-router";
import { useSetName } from "@/hooks/useProfile";
import NameInput from "@/components/profile/NameInput";
import FixedBottomCTA from "@/components/FixedBottomCTA";
import { colors } from "@/constants";
import { useEffect } from "react";
import { Ionicons } from "@expo/vector-icons";
import { setKVStore } from "@/db/storage";

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
            router.push("/conversations");
            return;
          }
          router.push("/account");
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
