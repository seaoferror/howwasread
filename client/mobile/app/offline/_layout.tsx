import { Stack } from "expo-router";

export default function OfflineLayout() {
  return (
    <Stack>
      <Stack.Screen
        name="create"
        options={{
          title: "Create your offline conversation",
          headerShown: true,
        }}
      />
    </Stack>
  );
}
