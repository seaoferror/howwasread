import { Controller, useFormContext } from "react-hook-form";

import InputField from "@/components/InputField";

export default function LengthInput() {
  const { control } = useFormContext();

  return (
    <Controller
      name="length"
      control={control}
      rules={{
        validate: (data: string) => {
          if (data.trim().length === 0) {
            return "length is required";
          }

          const length = Number(data);
          if (Number.isNaN(length) || length <= 19 || length >= 181) {
            return "length should be a number between 20 and 180";
          }
        },
      }}
      render={({ field: { onChange, value }, fieldState: { error } }) => (
        <InputField
          variant="standard"
          label="Length(minute)"
          placeholder="100"
          inputMode="numeric"
          maxLength={3}
          returnKeyType="done"
          submitBehavior="blurAndSubmit"
          value={String(value ?? "")}
          onChangeText={(text) => onChange(text.replace(/[^\d]/g, ""))}
          error={error?.message}
        />
      )}
    />
  );
}
