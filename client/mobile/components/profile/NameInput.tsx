import { Controller, useFormContext } from "react-hook-form";
import InputField from "@/components/InputField";
import { StyleSheet } from "react-native";

export default function NameInput() {
  const { control } = useFormContext();

  return (
    <Controller
      name="name"
      control={control}
      rules={{
        validate: (data: string) => {
          if (!data || data.trim().length <= 1) {
            return "name should be more than 1 characters";
          }
          if (!/^[A-Za-z ]+$/.test(data)) {
            return "use latin alphabets only";
          }
        },
      }}
      render={({ field: { onChange, value }, fieldState: { error } }) => (
        <InputField
          style={styles.text}
          customHeight={55}
          variant="outlined"
          autoFocus
          label="Your name"
          placeholder=""
          returnKeyType="done"
          submitBehavior="blurAndSubmit"
          value={value}
          onChangeText={onChange}
          error={error?.message}
        />
      )}
    />
  );
}

const styles = StyleSheet.create({
  text: {
    fontSize: 24,
  },
});
