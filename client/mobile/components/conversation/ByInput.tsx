import { Controller, useFormContext } from "react-hook-form";
import InputField from "@/components/InputField";

export default function ByInput() {
  const { control } = useFormContext();
  return (
    <Controller
      name="by"
      control={control}
      render={({ field: { onChange, value } }) => (
        <InputField
          variant="standard"
          label="written by(required)"
          placeholder="Shakespeare Soseki Pound"
          inputMode="text"
          returnKeyType="done"
          submitBehavior="blurAndSubmit"
          value={value}
          onChangeText={onChange}
        />
      )}
    />
  );
}
