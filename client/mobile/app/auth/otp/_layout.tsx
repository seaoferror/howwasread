import { router, Stack } from "expo-router";
import { Pressable } from "react-native";
import { Ionicons } from "@expo/vector-icons";

export default function OTPLayout() {
  return (
    <Stack>
      <Stack.Screen
        name="sms"
        options={{
          title: "Verify your phone number",
          headerShown: true,
          headerLeft: () => (
            <Pressable onPress={() => router.back()}>
              <Ionicons name="chevron-back" size={28} color="black" />
            </Pressable>
          ),
        }}
      />
      <Stack.Screen
        name="email"
        options={{
          title: "Verify your email",
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
