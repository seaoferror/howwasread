import { RegisterNotificationRequest } from "@/types/token";
import { localDevInstance } from "@/api/axios";

export async function registerNotification(body: RegisterNotificationRequest) {
  console.log(body);
  const { data } = await localDevInstance.post("/notification/register", body);
  return data;
}
