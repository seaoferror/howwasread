import { Stack } from "expo-router";

export default function ChatsLayout() {
  return (
    <Stack>
      <Stack.Screen
        name="index"
        options={{
          title: "Chat",
          headerShown: true,
        }}
      />
    </Stack>
  );
}
