import { ReactNode, useEffect } from "react";
import { router } from "expo-router";
import { useGetMyProfile } from "@/hooks/useProfile";
import { getSecure, setSecure } from "@/util/storage";
import { useSQLiteContext } from "expo-sqlite";
import { deleteAllMessages } from "@/db/message";

interface AuthRouteProps {
  children: ReactNode;
}


export default function AuthRoute({ children }: AuthRouteProps) {
  const { data: profile, isError, isLoading } = useGetMyProfile();
  const db = useSQLiteContext();
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
    if (getSecure("id") !== profile?.id) {
        deleteAllMessages(db);
      setSecure("id", profile?.id ?? "");
    }
    if (!profile?.name) {
      router.replace("/profile/name?newcomer=true");
      return;
    }
    router.replace("/conversations");
  }, [profile, isError, isLoading]);
  return <>{children}</>;
}
