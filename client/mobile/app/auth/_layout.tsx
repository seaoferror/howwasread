import { router, Stack } from "expo-router";
import { Pressable } from "react-native";
import { Ionicons } from "@expo/vector-icons";

export default function AuthLayout() {
  return (
    <Stack>
      <Stack.Screen
        name="index"
        options={{
          title: "",
          headerShown: false,
          headerLeft: () => (
            <Pressable onPress={() => router.back()}>
              <Ionicons name="chevron-back" size={28} color="black" />
            </Pressable>
          ),
        }}
      />
      <Stack.Screen
        name="phone-number"
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
          headerShown: false,
        }}
      />
      <Stack.Screen
        name="otp"
        options={{
          headerShown: false,
        }}
      />
    </Stack>
  );
}
