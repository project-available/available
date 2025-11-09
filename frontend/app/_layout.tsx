import { Stack } from "expo-router";
import "./globals.css";
import AsyncStorage from "@react-native-async-storage/async-storage";
import { useEffect } from "react";

export default function RootLayout() {
  useEffect(() => {
    const initLogin = async () => {
      try {
        console.log("🔐 Login...");

        const res = await fetch(
          "https://querulous-valerie-quanghia-967df8a0.koyeb.app/accounts/login",
          {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({
              email: "nhibui@hcmut.edu.vn",
              password: "123456",
            }),
          }
        );

        if (!res.ok) {
          throw new Error(`Login failed: ${res.status}`);
        }

        const data = await res.json();

        // Lưu account và token
        await AsyncStorage.setItem("account", JSON.stringify(data.account));
        await AsyncStorage.setItem("access_token", data.access_token);

        console.log("✅ Login successful, token has been saved!");
      } catch (err) {
        console.error("❌ Auto-login error:", err);
      }
    };

    initLogin();
  }, []);

  return (
    <Stack screenOptions={{ headerShown: false }}>
      <Stack.Screen name="(tabs)" />
      <Stack.Screen name="history" />
      <Stack.Screen name="detailroom" />
      <Stack.Screen name="booking" />
    </Stack>
  );
}
