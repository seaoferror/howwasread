import { ReactNode } from "react";
import { router, useFocusEffect } from "expo-router";
import { useAuth } from "@/hooks/useAuth";

interface AuthRouteProps {
  children: ReactNode;
}

export default function AuthRoute({ children }: AuthRouteProps) {
  const { id } = useAuth();

  useFocusEffect(() => {
    // !id && router.replace("/auth");
    !id && router.replace("/conversations");

  });
  return <>{children}</>;
}
