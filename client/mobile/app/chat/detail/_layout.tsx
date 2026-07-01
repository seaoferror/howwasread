import { router, Stack } from "expo-router";
import { Pressable } from "react-native";
import { Ionicons } from "@expo/vector-icons";

export default function ChatDetailLayout() {
  return (
    <Stack>
      <Stack.Screen
        name="[id]"
        options={{
          title: "",
          headerShown: true,
          headerLeft: () => (
            <Pressable onPress={() => router.back()}>
              <Ionicons name="chevron-back" size={28} color="black" />
            </Pressable>
          ),
          headerTransparent: true
        }}
      />
    </Stack>
  );
}
