import {
  deleteItemAsync,
  getItem,
  getItemAsync,
  setItemAsync,
} from "expo-secure-store";
import Storage from "expo-sqlite/kv-store";

export async function saveSecureStore(
  key: string,
  value: string,
): Promise<void> {
  await setItemAsync(key, value);
}

export async function getSecureAsync(key: string): Promise<string> {
  const value = await getItemAsync(key);
  return value ?? "";
}

export function getSecure(key: string): string {
  const value = getItem(key)
  return value ?? ""
}

export async function deleteSecure(key: string): Promise<void> {
  await deleteItemAsync(key);
}

export function getKVStore(key: string): string {
  const value = Storage.getItemSync(key);
  return value ?? "";
}

export function setKVStore(key: string, value: string) {
  Storage.setItemSync(key, value);
}
