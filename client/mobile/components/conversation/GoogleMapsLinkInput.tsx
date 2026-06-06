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
          placeholder="https://maps.app.goo.gl/nfMxW4wRAXFd5drq1 or https://www.google.com/maps/place/Dublin,+Ireland/@53.3243941,-6.3282753,12z/data=!3m1!4b1!4m6!3m5!1s0x48670e80ea27ac2f:0xa00c7a9973171a0!8m2!3d53.3498053!4d-6.2603097!16zL20vMDJjZnQ?entry=ttu&g_ep=EgoyMDI2MDYwMS4wIKXMDSoASAFQAw%3D%3D"
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
