import React, { useState } from "react";
import { View, Text, TextInput, TouchableOpacity, Alert } from "react-native";
import { Ionicons } from "@expo/vector-icons";
import { router } from "expo-router";
import AsyncStorage from "@react-native-async-storage/async-storage";


export default function LoginScreen() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");

  const handleLogin = async () => {
    if (!email || !password) {
      Alert.alert("Error", "Please enter both email and password.");
      return;
    }

    try {
      const response = await fetch(
        "https://querulous-valerie-quanghia-967df8a0.koyeb.app/accounts/login",
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({
            email,
            password,
          }),
        }
      );

      const data = await response.json();
      console.log("Login response:", data);

      if (!response.ok) {
        Alert.alert("Login Failed", data.error || "Unknown error");
        return;
      }

    await AsyncStorage.setItem("accessToken", data.access_token);

    await AsyncStorage.setItem("user", JSON.stringify(data.account));

      Alert.alert("Success", "Welcome!");
      router.replace("/(tabs)");
    } catch (error) {
      console.error(error);
      Alert.alert("Error", "Something went wrong. Try again.");
    }
  };

  return (
    <View className="flex-1 bg-white px-6 pt-[100px]">
      <TouchableOpacity
        onPress={() => router.back()}
        className="absolute top-[50px] left-[20px] flex-row items-center"
      >
        <Ionicons name="arrow-back" size={24} color="#333" />
        <Text className="ml-2 text-[16px] text-[#333] font-medium">Back</Text>
      </TouchableOpacity>

      <View className="absolute top-[85px] left-0 right-0 h-[1px] bg-[#ccc]" />

      <Text className="text-[28px] font-bold mb-8 text-[#333]">Sign In</Text>

      <TextInput
        className="h-[50px] border border-[#ccc] rounded-xl px-4 mb-5 text-[16px] text-[#333]"
        placeholder="Email"
        placeholderTextColor="#999"
        value={email}
        onChangeText={setEmail}
        keyboardType="email-address"
        autoCapitalize="none"
      />

      <TextInput
        className="h-[50px] border border-[#ccc] rounded-xl px-4 mb-5 text-[16px] text-[#333]"
        placeholder="Password"
        placeholderTextColor="#999"
        value={password}
        onChangeText={setPassword}
        secureTextEntry
      />

      <TouchableOpacity
        onPress={handleLogin}
        className="bg-[#FDBA29] py-3 rounded-xl items-center"
      >
        <Text className="text-white text-[18px] font-semibold">Sign In</Text>
      </TouchableOpacity>

      <View className="flex-row items-center justify-center my-8">
        <View className="flex-1 h-[1px] bg-[#ccc]" />
        <Text className="mx-3 text-[16px] text-[#ccc] font-medium">or</Text>
        <View className="flex-1 h-[1px] bg-[#ccc]" />
      </View>
    </View>
  );
}
