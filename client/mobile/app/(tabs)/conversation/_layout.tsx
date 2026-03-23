import { Stack } from "expo-router";

export default function ConversationLayout() {
  return (
    <Stack>
      <Stack.Screen
        name="index"
        options={{
          title: "Conversation",
          headerShown: true,
          animation: "slide_from_left",
        }}
      />
    </Stack>
  );
}
