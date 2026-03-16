import { deleteItemAsync, getItemAsync, setItemAsync } from "expo-secure-store";

export async function saveSecureStore(
  key: string,
  value: string,
): Promise<void> {
  await setItemAsync(key, value);
}

export async function getSecureStore(key: string): Promise<string> {
  const value = (await getItemAsync(key)) ?? "";
  return value;
}

export async function deleteSecureStore(key: string): Promise<void> {
  await deleteItemAsync(key);
}
