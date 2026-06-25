import { router, Stack } from "expo-router";
import { Pressable } from "react-native";
import { Ionicons } from "@expo/vector-icons";

export default function ProfileLayout() {
  return (
    <Stack>
      <Stack.Screen
        name="name"
        options={{
          title: "Set your name",
          headerShown: true,
          headerLeft: () => (
            <Pressable onPress={() => router.back()}>
              <Ionicons name="chevron-back" size={28} color="black" />
            </Pressable>
          ),
        }}
      />
    </Stack>
  );
}
