import { Stack } from "expo-router";

export default function ConversationsLayout() {
  return (
    <Stack>
      <Stack.Screen
        name="index"
        options={{
          title: "Conversations",
          headerShown: true,
          animation: "slide_from_left",
        }}
      />
    </Stack>
  );
}
