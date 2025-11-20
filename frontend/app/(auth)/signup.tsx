import React, { useState } from "react";
import { View, Text, TextInput, TouchableOpacity, Alert } from "react-native";
import { Ionicons } from "@expo/vector-icons";
import { router } from "expo-router";

export default function SignupScreen() {
  const [name, setName] = useState("");
  const [studentId, setStudentId] = useState("");
  const [phone, setPhone] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");

  // 5 biến focus để đổi màu border khi focus
  const [focus, setFocus] = useState("");

  // const handleSignup = () => {
  //   if (!name || !studentId || !phone || !email || !password) {
  //     Alert.alert("Error", "Please fill in all fields.");
  //     return;
  //   }
  //   console.log("New user:", { name, studentId, phone, email });
  //   router.replace("/(tabs)");
  // };

  const handleSignup = async () => {
  try {

    if (password.length < 6) {
      Alert.alert("Error", "Password must be at least 6 characters long.");
      return;
    }
    const response = await fetch(
      "https://querulous-valerie-quanghia-967df8a0.koyeb.app/accounts",
      {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          name,
          email,
          password,
          phone,
          student_id: studentId,  
        }),
      }
    );

    const data = await response.json();
    console.log("Signup response:", data);

    if (!response.ok) {
      Alert.alert("Signup Failed", data.error || "Unknown error");
      return;
    }

    Alert.alert("Success", "Account created successfully!");
    router.replace("/(tabs)");
  } catch (err) {
    console.error(err);
  }
};

  return (
    <View className="flex-1 bg-white px-6 pt-[100px]">
      {/* --- Back button --- */}
      <TouchableOpacity
        onPress={() => router.replace("/")}
        className="absolute top-[50px] left-[20px] flex-row items-center"
      >
        <Ionicons name="arrow-back" size={24} color="#333" />
        <Text className="ml-2 text-[16px] text-[#333] font-medium">Back</Text>
      </TouchableOpacity>

      {/* --- Line separator --- */}
      <View className="absolute top-[85px] left-0 right-0 h-[1px] bg-[#ccc]" />

      {/* --- Title --- */}
      <Text className="text-[28px] font-bold text-center mb-6 text-[#333]">Sign Up</Text>

      {/* --- Name --- */}
      <TextInput
        className="h-[50px] rounded-xl px-4 mb-4 text-[16px] text-[#333] border"
        style={{
          borderColor: focus === "name" ? "#FDBA29" : "#ccc",
          borderWidth: 1,
        }}
        placeholder="Full Name"
        placeholderTextColor="#999"
        value={name}
        onChangeText={setName}
        onFocus={() => setFocus("name")}
        onBlur={() => setFocus("")}
      />

      {/* --- Student ID --- */}
      <TextInput
        className="h-[50px] rounded-xl px-4 mb-4 text-[16px] text-[#333] border"
        style={{
          borderColor: focus === "studentId" ? "#FDBA29" : "#ccc",
          borderWidth: 1,
        }}
        placeholder="Student ID"
        placeholderTextColor="#999"
        value={studentId}
        onChangeText={setStudentId}
        keyboardType="numeric"
        onFocus={() => setFocus("studentId")}
        onBlur={() => setFocus("")}
      />

      {/* --- Phone --- */}
      <TextInput
        className="h-[50px] rounded-xl px-4 mb-4 text-[16px] text-[#333] border"
        style={{
          borderColor: focus === "phone" ? "#FDBA29" : "#ccc",
          borderWidth: 1,
        }}
        placeholder="Phone Number"
        placeholderTextColor="#999"
        value={phone}
        onChangeText={setPhone}
        keyboardType="phone-pad"
        onFocus={() => setFocus("phone")}
        onBlur={() => setFocus("")}
      />

      {/* --- Email --- */}
      <TextInput
        className="h-[50px] rounded-xl px-4 mb-4 text-[16px] text-[#333] border"
        style={{
          borderColor: focus === "email" ? "#FDBA29" : "#ccc",
          borderWidth: 1,
        }}
        placeholder="Email"
        placeholderTextColor="#999"
        value={email}
        onChangeText={setEmail}
        keyboardType="email-address"
        autoCapitalize="none"
        onFocus={() => setFocus("email")}
        onBlur={() => setFocus("")}
      />

      {/* --- Password --- */}
      <TextInput
        className="h-[50px] rounded-xl px-4 mb-6 text-[16px] text-[#333] border"
        style={{
          borderColor: focus === "password" ? "#FDBA29" : "#ccc",
          borderWidth: 1,
        }}
        placeholder="Password"
        placeholderTextColor="#999"
        value={password}
        onChangeText={setPassword}
        secureTextEntry
        onFocus={() => setFocus("password")}
        onBlur={() => setFocus("")}
      />

      {password.length > 0 && password.length < 6 && (
        <Text style={{ color: "red", marginBottom: 10 }}>
          Password must be at least 6 characters.
        </Text>
      )}

      {/* --- Signup button --- */}
      <TouchableOpacity
        onPress={handleSignup}
        className="bg-[#FDBA29] py-4 rounded-xl items-center"
      >
        <Text className="text-white text-[18px] font-semibold">Create Account</Text>
      </TouchableOpacity>

      <View className="flex-row items-center justify-center my-8">
        <View className="flex-1 h-[1px] bg-[#ccc]" />
        <Text className="mx-3 text-[16px] text-[#555] font-medium">OR</Text>
        <View className="flex-1 h-[1px] bg-[#ccc]" />
      </View>

      {/* --- Already have account --- */}
          <View className="mt-3 flex-row justify-center">
      <Text className="text-[#333]">
        Already have an account?
        <Text
          style={{ color: "#FDBA29", fontWeight: "600" }}
          onPress={() => router.push("/login")}
        >
          {" "}Sign in
        </Text>
      </Text>
    </View>
    </View>
  );
}
