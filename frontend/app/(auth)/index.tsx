import { View, Text, TouchableOpacity, StyleSheet } from "react-native";
import { router } from "expo-router";

export default function WelcomeScreen() {
  return (
    <View style={styles.container}>

      <TouchableOpacity
        style={[styles.button, { backgroundColor: "#FDBA29" }]}
        onPress={() => router.push("/(auth)/signup")}
      >
        <Text style={styles.buttonText}>Create Account</Text>
      </TouchableOpacity>

      <TouchableOpacity
        style={[styles.button, { backgroundColor: "#fff" }]}
        onPress={() => router.push("/(auth)/login")}
      >
        <Text style={styles.buttonText}>Sign Ip</Text>
      </TouchableOpacity>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: "#fff",
    justifyContent: "center",
    alignItems: "center",
    padding: 24,
  },
  title: { fontSize: 32, fontWeight: "700", marginBottom: 8 },
  subtitle: { fontSize: 16, color: "#666", marginBottom: 40 },
  button: {
    width: "80%",
    paddingVertical: 14,
    borderRadius: 10,
    marginBottom: 16,
    alignItems: "center",
  },
  buttonText: { color: "#000", fontSize: 18, fontWeight: "600" },
});
