import { RegisterNotificationRequest } from "@/types/token";
import { axiosInstance } from "@/api/axios";

export async function registerNotification(body: RegisterNotificationRequest) {
  console.log(body);
  const { data } = await axiosInstance.post("/notification/register", body);
  return data;
}
