import { Controller, useFormContext } from "react-hook-form";

import InputField from "@/components/InputField";

export default function CapacityInput() {
  const { control } = useFormContext();

  return (
    <Controller
      name="capacity"
      control={control}
      rules={{
        validate: (data: number | undefined) => {
          if (data === undefined) {
            return "capacity is required";
          }
          if (Number.isNaN(data) || data >= 7 || data <= 1) {
            return "capacity should be a number between 2 and 6";
          }
        },
      }}
      render={({ field: { onChange, value }, fieldState: { error } }) => (
        <InputField
          variant="standard"
          label="Capacity"
          placeholder="5"
          inputMode="numeric"
          maxLength={3}
          returnKeyType="done"
          submitBehavior="blurAndSubmit"
          value={String(value ?? "")}
          onChangeText={(text) => {
            const digits = text.replace(/[^\d]/g, "");
            onChange(digits.length === 0 ? undefined : Number(digits));
          }}
          error={error?.message}
        />
      )}
    />
  );
}
