import { ReactNode, useEffect } from "react";
import { router } from "expo-router";
import { useGetMyProfile } from "@/hooks/useProfile";

interface AuthRouteProps {
  children: ReactNode;
}

export default function AuthRoute({ children }: AuthRouteProps) {
  const { data: profile, isError, isLoading } = useGetMyProfile();
  // const networkState = useNetworkState();

  useEffect(() => {
    console.log(profile, isError, isLoading);
    if (isLoading) {
      return;
    }
    if (isError) {
      router.replace("/auth");
      return;
    }
    if (!profile?.name) {
      router.replace("/profile/name?newcomer=true");
      return;
    }
    router.replace("/conversations");
  }, [profile, isError, isLoading]);
  return <>{children}</>;
}
