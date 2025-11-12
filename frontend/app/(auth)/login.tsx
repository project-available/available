import React, { useState } from "react";
import { View, Text, TextInput, TouchableOpacity, StyleSheet } from "react-native";
import { Ionicons } from "@expo/vector-icons";
import { router } from "expo-router";

export default function LoginScreen({ navigation }: any) {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");

  const handleLogin = () => {
    console.log("Email:", email);
    console.log("Password:", password);
  };

  return (
    <View style={styles.container}>
      <TouchableOpacity onPress={() => router.back()} style={styles.backButton}>
        <Ionicons name="arrow-back" size={24} color="#333" />
        <Text style={styles.backText}>Back</Text>
      </TouchableOpacity>

      <View style={styles.seperator} />

      <Text style={styles.title}>Sign In</Text>

      <TextInput
        style={styles.input}
        placeholder="Email"
        placeholderTextColor="#999"
        value={email}
        onChangeText={setEmail}
        keyboardType="email-address"
        autoCapitalize="none"
      />

      <TextInput
        style={styles.input}
        placeholder="Password"
        placeholderTextColor="#999"
        value={password}
        onChangeText={setPassword}
        secureTextEntry
      />

      <TouchableOpacity style={styles.signInButton} onPress={handleLogin}>
        <Text style={styles.signInText}>Sign In</Text>
      </TouchableOpacity>
      <View style={styles.orContainer}>
        <View style={styles.line} />
        <Text style={styles.orText}>OR</Text>
        <View style={styles.line} />
      </View>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: "#fff",
    paddingHorizontal: 24,
    paddingTop: 100, 
  },
  backButton: {
    position: "absolute",
    top: 50,
    left: 20,
    flexDirection: "row",    
  },
  backText: {
    marginLeft: 6,  
    fontSize: 16,
    color: "#333",
    fontWeight: "500",
  },
  seperator: {
    position: "absolute",
    top: 85,            
    left: 0,
    right: 0,
    height: 1,
    backgroundColor: "#ccc", 
  },
  title: {
    fontSize: 28,
    fontWeight: "700",
    marginBottom: 32,
    textAlign: "left",
    color: "#333",
  },
  input: {
    height: 50,
    borderColor: "#ccc",
    borderWidth: 1,
    borderRadius: 12,
    paddingHorizontal: 16,
    marginBottom: 20,
    fontSize: 16,
    color: "#333",
  },
  signInButton: {
    backgroundColor: "#FDBA29",
    paddingVertical: 14,
    borderRadius: 12,
    alignItems: "center",
  },
  signInText: {
    color: "#fff",
    fontSize: 18,
    fontWeight: "600",
  },
  orContainer: {
    flexDirection: "row",       
    alignItems: "center",       
    justifyContent: "center",
    marginVertical: 30,        
  },
  line: {
    flex: 1,                    
    height: 1,
    backgroundColor: "#ccc",    
  },
  orText: {
    marginHorizontal: 10,       
    fontSize: 16,
    color: "#555",
    fontWeight: "500",
  },
});
