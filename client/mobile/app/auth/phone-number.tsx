import { StyleSheet, View } from "react-native";
import { SafeAreaView } from "react-native-safe-area-context";
import { colors, time } from "@/constants";
import CountryCodeBox from "@/components/auth/CountryCodeBox";
import PhoneNumberInput from "@/components/auth/PhoneNumberInput";
import { FormProvider, useForm } from "react-hook-form";
import {
  CountryCode,
  isValidPhoneNumber,
  parsePhoneNumberFromString,
} from "libphonenumber-js";
import Toast from "react-native-toast-message";
import { useRequestSMSOTP } from "@/hooks/useAuth";
import { router } from "expo-router";
import { getSecure, setSecure } from "@/db/storage";
import CustomButton from "@/components/CustomButton";

interface FormValue {
  countryCode: CountryCode;
  phoneNumber: string;
}

const canRequestNewSMS = () => {
  const s = getSecure("timeSmsLastSent");
  const t = s ? Number(s) : 0;
  return Date.now() - t > time.TEN_MINUTES;
};

export default function PhoneNumberScreen() {
  const phoneNumberForm = useForm<FormValue>({
    defaultValues: {
      countryCode: "KR",
      phoneNumber: "",
    },
  });

  const requestSMSOTPMutation = useRequestSMSOTP();

  const onSubmit = async (formValues: FormValue) => {
    console.log("start submit");
    if (!canRequestNewSMS()) {
      Toast.show({
        type: "info",
        text1: "Please wait",
        text2: "You can request another code in a few minutes.",
      });
      return;
    }

    const { countryCode, phoneNumber } = formValues;
    const digitsOnly = phoneNumber.replace(/[^\d]/g, "");
    if (!isValidPhoneNumber(digitsOnly, countryCode)) {
      Toast.show({
        type: "error",
        text1: "Invalid phone number",
      });
      return;
    }

    const parsed = parsePhoneNumberFromString(digitsOnly, countryCode);
    const wholeNumber = parsed?.number;
    console.log("wholeNumber: ", wholeNumber);

    if (!wholeNumber) {
      console.log("fail to parse number");
      Toast.show({
        type: "error",
        text1: "Invalid phone number",
      });
      return;
    }

    console.log("execute mutate");
    requestSMSOTPMutation.mutate(
      {
        sessionId: getSecure("sessionId") || null,
        phoneNumber: wholeNumber,
      },
      {
        onSuccess: async () => {
          await setSecure("timeSmsLastSent", String(Date.now()));
          router.push("/auth/otp/sms");
        },
      },
    );
  };

  return (
    <SafeAreaView style={styles.container}>
      <FormProvider {...phoneNumberForm}>
        <View style={styles.content}>
          <View style={styles.phoneRow}>
            <CountryCodeBox />
            <PhoneNumberInput />
          </View>
          <CustomButton
            label={"Send Code"}
            onPress={phoneNumberForm.handleSubmit(onSubmit)}
            disabled={requestSMSOTPMutation.isPending}
          />
        </View>
      </FormProvider>
    </SafeAreaView>
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
    paddingHorizontal: 20,
    paddingTop: 120,
    gap: 50,
  },
  phoneRow: {
    flexDirection: "row",
    width: "100%",
    alignItems: "center",
    gap: 10,
  },
});
