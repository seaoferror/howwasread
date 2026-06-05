import { Controller, useFormContext } from "react-hook-form";
import InputField from "@/components/InputField";

export default function GoogleMapsLinkInput() {
  const { control } = useFormContext();
  return (
    <Controller
      name="googleMapsLink"
      control={control}
      render={({ field: { onChange, value } }) => (
        <InputField
          variant="standard"
          label="Google Maps link(required)"
          placeholder="https://maps.app.goo.gl/nfMxW4wRAXFd5drq1"
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
