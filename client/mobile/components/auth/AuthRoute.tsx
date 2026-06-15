import { ReactNode } from "react";
import { router, useFocusEffect } from "expo-router";
import { useAuth } from "@/hooks/useAuth";
import { useMyProfile } from "@/hooks/useMyProfile";

interface AuthRouteProps {
  children: ReactNode;
}

export default function AuthRoute({ children }: AuthRouteProps) {
  const { id } = useAuth();
  const { profile } = useMyProfile();

  useFocusEffect(() => {
    console.log(id)
    if (!id) {
      router.replace("/auth");
      return
    }
    if (!profile.name) {
      router.replace("/profile/name?newcomer=true");
      return;
    }
    router.replace("/conversations")
  });
  return <>{children}</>;
}
