import { Stack } from "expo-router";
import "./globals.css";
import * as Sentry from "@sentry/react-native";
import { API_BASE_URL } from "../utils/apiConfig";

Sentry.init({
  dsn: "https://bddadb2f6afa0ebb0d96d8f1c8e93d42@o4510268002664448.ingest.de.sentry.io/4510697430646864",

  // Adds more context data to events (IP address, cookies, user, etc.)
  // For more information, visit: https://docs.sentry.io/platforms/react-native/data-management/data-collected/
  sendDefaultPii: true,

  // Enable Logs
  enableLogs: true,

  // Set tracesSampleRate to 1.0 to capture 100% of transactions for performance monitoring.
  // We recommend adjusting this value in production.
  tracesSampleRate: 1.0,

  // Capture distributed traces
  tracePropagationTargets: [
    "localhost",
    API_BASE_URL,
    /^https:\/\/querulous-valerie-quanghia-967df8a0\.koyeb\.app/,
  ],

  // Configure Session Replay
  replaysSessionSampleRate: 0.1,
  replaysOnErrorSampleRate: 1,
  integrations: [
    Sentry.mobileReplayIntegration(),
    Sentry.feedbackIntegration(),
  ],

  // uncomment the line below to enable Spotlight (https://spotlightjs.com)
  // spotlight: __DEV__,
});

export default Sentry.wrap(function RootLayout() {
  return (
    <Stack screenOptions={{ headerShown: false }}>
      <Stack.Screen name="(auth)" />
      <Stack.Screen name="(tabs)" />
      <Stack.Screen name="history" />
      <Stack.Screen name="detailroom" />
      <Stack.Screen name="booking" />
    </Stack>
  );
});
