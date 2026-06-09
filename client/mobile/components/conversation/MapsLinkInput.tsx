import { Controller, useFormContext } from "react-hook-form";
import InputField from "@/components/InputField";

export default function MapsLinkInput() {
  const { control } = useFormContext();
  return (
    <Controller
      name="mapsLink"
      control={control}
      render={({ field: { onChange, value } }) => (
        <InputField
          variant="standard"
          label="Maps share link(required)"
          placeholder="https://maps.app.goo.gl/G3hMpRVRoh3ZXue39 or https://www.google.com/maps/place/Seorae+Island/@37.5079857,126.9899704,17z/data=!3m1!4b1!4m6!3m5!1s0x357ca184dcefac4f:0x4efa5c5519a2579b!8m2!3d37.5079857!4d126.9899704!16s%2Fg%2F1vn16gj7!18m1!1e1?entry=ttu&g_ep=EgoyMDI2MDYwMy4xIKXMDSoASAFQAw%3D%3D"
          inputMode="text"
          multiline={true}
          returnKeyType="done"
          submitBehavior="blurAndSubmit"
          value={value}
          onChangeText={onChange}
        />
      )}
    />
  );
}
