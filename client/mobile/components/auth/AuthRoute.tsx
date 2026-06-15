import { ReactNode } from "react";
import { router, useFocusEffect } from "expo-router";
import { useAuth } from "@/hooks/useAuth";
import { useMyProfile } from "@/hooks/useMyProfile";

interface AuthRouteProps {
  children: ReactNode;
}

export default function AuthRoute({ children }: AuthRouteProps) {
  const { profile } = useMyProfile();

  useFocusEffect(() => {
    console.log(profile.id)
    if (!profile.id) {
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
