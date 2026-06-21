import { Link, Stack } from "expo-router";
import { MaterialIcons } from "@expo/vector-icons";
import { colors } from "@/constants";

export default function EmailLayout() {
  return (
    <Stack>
      <Stack.Screen
        name="signup"
        options={{
          title: "   Sign up with your email",
          headerShown: true,
          headerLeft: () => (
            <Link href={"/"} replace style={{ paddingRight: 5 }}>
              <MaterialIcons
                name="arrow-back-ios-new"
                size={28}
                color={colors.BLACK}
              />
            </Link>
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
            <Link href={"/"} replace style={{ paddingRight: 5 }}>
              <MaterialIcons
                name="arrow-back-ios-new"
                size={28}
                color={colors.BLACK}
              />
            </Link>
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
    </Stack>
  );
}
