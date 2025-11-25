import { useEffect, useState } from "react";
import { Redirect } from "expo-router";
import AsyncStorage from "@react-native-async-storage/async-storage";

export default function RootRedirect() {
  const [hasSeenOnboarding, setHasSeenOnboarding] = useState<boolean | null>(null);

  useEffect(() => {
    const checkOnboarding = async () => {
      const value = await AsyncStorage.getItem("hasSeenOnboarding");
      setHasSeenOnboarding(value === "true");
    };

    checkOnboarding();
  }, []);

  if (hasSeenOnboarding === null) {
    return null;
  }

  return hasSeenOnboarding
    ? <Redirect href="/(auth)" />
    : <Redirect href="/onboarding/screen1" />;
}
