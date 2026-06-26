import { ReactNode, useEffect } from "react";
import { router } from "expo-router";
import { getKVStore, getSecure } from "@/db/storage";
import { useGetMyProfile } from "@/hooks/useProfile";

interface AuthRouteProps {
  children: ReactNode;
}

export default function AuthRoute({ children }: AuthRouteProps) {
  const { data: profile, isLoading } = useGetMyProfile();

  useEffect(() => {
    if(isLoading) {
      return
    }
    if(!!(profile?.name ?? getKVStore("myName"))){
      router.replace("/conversations");
      return
    }
    if(!!getSecure("accessToken")) {
      router.replace("/profile/name")
      return
    }
    router.replace("/auth");
  }, [isLoading, profile?.name]);

  return <>{children}</>;
}
