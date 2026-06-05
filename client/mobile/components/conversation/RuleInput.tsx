import { Controller, useFormContext } from "react-hook-form";
import InputField from "@/components/InputField";

export default function RuleInput() {
  const { control } = useFormContext();
  return (
    <Controller
      name="rule"
      control={control}
      render={({ field: { onChange, value } }) => (
        <InputField
          variant="standard"
          label="Rule(optional)"
          placeholder={
            "1. Respect each other\n2. Try to speak more\n3. Try to yield and listen"
          }
          inputMode="text"
          returnKeyType="default"
          submitBehavior="newline"
          value={value}
          onChangeText={onChange}
          multiline
        />
      )}
    />
  );
}
