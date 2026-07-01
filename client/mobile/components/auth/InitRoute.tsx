import { ReactNode, useEffect } from "react";
import { router } from "expo-router";
import { getKVStore } from "@/db/storage";

interface AuthRouteProps {
  children: ReactNode;
}

export default function InitRoute({ children }: AuthRouteProps) {

  useEffect(() => {
    if(!!getKVStore("myName")){
      router.replace("/conversations");
      return
    }
    if(!!getKVStore("myId")) {
      router.replace("/profile/name")
      return
    }
    router.replace("/auth");
  }, []);

  return <>{children}</>;
}
