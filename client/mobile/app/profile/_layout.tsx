import { Stack } from "expo-router";

export default function ProfileLayout() {
  return (
    <Stack>
      <Stack.Screen
        name="name"
        options={{
          title: "Set your name",
          headerShown: true,
        }}
      />
      <Stack.Screen
        name="term-of-use"
        options={{
          headerShown: false,
        }}
      />
    </Stack>
  );
}
