import {
  View,
  Text,
  Image,
  TouchableOpacity,
  ActivityIndicator,
} from "react-native";
import { Feather } from "@expo/vector-icons";
import { useEffect, useState } from "react";
import AsyncStorage from "@react-native-async-storage/async-storage";
import { router } from "expo-router";

interface AccountData {
  id: number;
  name: string;
  role: string;
  email: string;
  phone: string;
  student_id: string;
}

export default function Profile() {
  const [accountData, setAccountData] = useState<AccountData | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    const fetchProfile = async () => {
      try {
        const token = await AsyncStorage.getItem("access_token");
        const accountString = await AsyncStorage.getItem("account");

        if (!token || !accountString) {
          throw new Error("Not logged in. Please login first!");
        }

        const account = JSON.parse(accountString);
        const studentId = account.student_id;

        const url =
          `https://querulous-valerie-quanghia-967df8a0.koyeb.app/accounts/${studentId}`;

        const response = await fetch(url, {
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
          },
        });

        if (response.status === 401) {
          await AsyncStorage.removeItem("access_token");
          router.replace("/(auth)/login");
          return;
        }

        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }

        const data: AccountData = await response.json();
        setAccountData(data);
      } catch (err) {
        console.error("❌ Error fetching profile:", err);
        setError(err instanceof Error ? err.message : "Failed to load profile");
      } finally {
        setLoading(false);
      }
    };

    fetchProfile();
  }, []);

  const handleLogout = async () => {
    try {
      await AsyncStorage.removeItem("access_token");
      await AsyncStorage.removeItem("account");
      router.replace("/");
    } catch (err) {
      console.error("❌ Error logging out:", err);
    }
  };

  if (loading) {
    return (
      <View className="flex-1 bg-[#f8f5f2] items-center justify-center">
        <ActivityIndicator size="large" color="#2A429A" />
      </View>
    );
  }

  if (error || !accountData) {
    return (
      <View className="flex-1 bg-[#f8f5f2] items-center justify-center px-6">
        <Feather name="alert-circle" size={48} color="#EC1861" />
        <Text className="text-center text-red-500 mt-4">
          {error || "Failed to load profile"}
        </Text>
      </View>
    );
  }

  return (
    <View className="flex-1 bg-[#f8f5f2]">
      {/* --- Title --- */}
      <Text className="mt-[66px] ml-[24px] text-[18px] leading-[24px] font-semibold">
        Profile Details
      </Text>

      {/* --- Horizontal Divider --- */}
      <View className="mt-[40px] h-[1px] bg-gray-300 mx-[24px]" />

      {/* --- Profile Info Card --- */}
      <View className="mx-[24px] mt-[24px] bg-[#2A429A] rounded-2xl items-center py-[24px]">
        <Image
          source={{ uri: "https://i.pravatar.cc/100?img=9" }}
          className="w-[100px] h-[100px] rounded-[24px] border-[4px] border-white"
        />
        <Text className="text-white font-semibold text-[16px] leading-[24px] mt-[8px]">
          {accountData.name}
        </Text>
        <Text className="text-white text-[12px] leading-[16px] mt-[2px]">
          MSSV: {accountData.student_id}
        </Text>
      </View>

      {/* --- Email Section --- */}
      <View className="h-[56px] bg-[#EEEBe5] mx-[24px] mt-[24px] rounded-2xl flex-row items-center">
        <View className="ml-[8px] my-[8px] w-[40px] h-[40px] bg-white rounded-[12px] items-center justify-center">
          <Feather name="mail" size={24} color="black" />
        </View>
        <View className="ml-[12px]">
          <Text className="text-[12px] leading-[16px] text-[#333]">Email</Text>
          <Text className="text-[14px] leading-[20px] font-semibold text-[#000]">
            {accountData.email}
          </Text>
        </View>
      </View>

      {/* --- Personal Phone Section --- */}
      <View className="h-[56px] bg-[#EEEBe5] mx-[24px] mt-[12px] rounded-2xl flex-row items-center">
        <View className="ml-[8px] my-[8px] w-[40px] h-[40px] bg-white rounded-[12px] items-center justify-center">
          <Feather name="phone" size={24} color="black" />
        </View>
        <View className="ml-[12px]">
          <Text className="text-[12px] leading-[16px] text-[#333]">
            Personal Phone
          </Text>
          <Text className="text-[14px] leading-[20px] font-semibold text-[#000]">
            {accountData.phone}
          </Text>
        </View>
      </View>

      {/* --- Edit & Logout Buttons (2 columns) --- */}
      <View className="flex-row mx-[24px] mt-[12px] justify-between">
        {/* Edit Button */}
        <TouchableOpacity className="flex-1 bg-[#EEEBe5] rounded-2xl h-[56px] flex-row items-center mr-[6px]">
          {/* Icon Container */}
          <View className="ml-[8px] my-[8px] w-[40px] h-[40px] bg-white rounded-[12px] items-center justify-center">
            <Feather name="edit-3" size={24} color="#EC1861" />
          </View>
          {/* Label */}
          <Text className="text-[14px] leading-[20px] font-semibold text-[#000] ml-[12px]">
            Edit
          </Text>
        </TouchableOpacity>

        {/* Logout Button */}
        <TouchableOpacity 
          className="flex-1 bg-[#EEEBe5] rounded-2xl h-[56px] flex-row items-center ml-[6px]"
          onPress={handleLogout}
        >
          <View className="ml-[8px] my-[8px] w-[40px] h-[40px] bg-white rounded-[12px] items-center justify-center">
            <Feather name="log-out" size={24} color="#EC1861" />
          </View>
          <Text className="text-[14px] leading-[20px] font-semibold text-[#000] ml-[12px]">
            Logout
          </Text>
        </TouchableOpacity>
      </View>
    </View>
  );
}
