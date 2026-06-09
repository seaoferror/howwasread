import { Controller, useFormContext } from "react-hook-form";
import InputField from "@/components/InputField";

export default function LocationInput() {
  const { control } = useFormContext();
  return (
    <Controller
      name="location"
      control={control}
      render={({ field: { onChange, value } }) => (
        <InputField
          variant="standard"
          label="location(optional)"
          placeholder="Starbucks texas houston"
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
