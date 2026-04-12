import { Stack } from "expo-router";

export default function ChatLayout() {
  return (
    <Stack>
      <Stack.Screen
        name="[id]"
        options={{
          title: "",
          headerShown: false,
          animation: "none",
        }}
      />
    </Stack>
  );
}
