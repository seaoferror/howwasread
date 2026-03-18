import { Stack } from "expo-router";

export default function OnlineLayout() {
  return (
    <Stack>
      <Stack.Screen
        name="index"
        options={{
          title: "Online conversation",
          headerShown: true,
          animation: "slide_from_left",
        }}
      />
      <Stack.Screen
        name="create"
        options={{
          title: "Create your online conversation",
          headerShown: true,
        }}
      />
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
