import { Link, router, Stack } from "expo-router";
import { Ionicons, MaterialIcons } from "@expo/vector-icons";
import { colors } from "@/constants";
import { Pressable } from "react-native";

export default function EmailLayout() {
  return (
    <Stack>
      <Stack.Screen
        name="signup"
        options={{
          title: "   Sign up with your email",
          headerShown: true,
          headerLeft: () => (
            <Pressable onPress={() => router.replace("/auth")}>
              <Ionicons name="chevron-back" size={28} color="black" />
            </Pressable>
          ),
          headerRight: () => (
            <Link
              href={"/auth/email/login"}
              style={{ fontSize: 20, color: colors.BLUE_500 }}
            >
              Login
            </Link>
          ),
        }}
      />
      <Stack.Screen
        name="login"
        options={{
          title: "   Login",
          headerShown: true,
          headerLeft: () => (
            <Pressable onPress={() => router.replace("/auth")}>
              <Ionicons name="chevron-back" size={28} color="black" />
            </Pressable>
          ),
          headerRight: () => (
            <Link
              href={"/auth/email/signup"}
              style={{ fontSize: 20, color: colors.BLUE_500 }}
            >
              Sign up
            </Link>
          ),
        }}
      />
      <Stack.Screen
        name="forget-password"
        options={{
          title: "   Verify your email",
          headerShown: true,
          headerLeft: () => (
            <Pressable onPress={() => router.back()}>
              <Ionicons name="chevron-back" size={28} color="black" />
            </Pressable>
          ),
        }}
      />
      <Stack.Screen
        name="set-new-password"
        options={{
          title: "   Set your new password",
          headerShown: true,
          headerLeft: () => (
            <Pressable onPress={() => router.replace("/auth")}>
              <Ionicons name="chevron-back" size={28} color="black" />
            </Pressable>
          ),
        }}
      />
    </Stack>
  );
}
