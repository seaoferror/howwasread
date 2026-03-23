import { Stack } from "expo-router";
import "react-native-reanimated";
import { QueryClientProvider } from "@tanstack/react-query";
import queryClient from "@/api/queryClient";
import Toast from "react-native-toast-message";

export default function RootLayout() {
  return (
    <QueryClientProvider client={queryClient}>
      <RootNavigator />
      <Toast />
    </QueryClientProvider>
  );
}

function RootNavigator() {
  return (
    <Stack>
      <Stack.Screen name="(init)" options={{ headerShown: false }} />
      <Stack.Screen name="(tabs)" options={{ headerShown: false }} />
      <Stack.Screen name="auth" options={{ headerShown: false }} />
      <Stack.Screen name="newcomer" options={{ headerShown: false }} />
      <Stack.Screen name="online" options={{ headerShown: false }} />
    </Stack>
  );
}
