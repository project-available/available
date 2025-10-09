import { View, Text, TouchableOpacity } from "react-native";
import { Ionicons } from "@expo/vector-icons";
import { useRouter } from "expo-router";

export default function HomeTab() {
  const router = useRouter();

  return (
    <View className="flex-1 bg-[#f8f5f2] px-[24px] pt-[66px]">
      {/* ---------- HEADER ---------- */}
      <View className="flex-row items-center justify-between bg-white rounded-full px-4 py-2 h-[64px]">
        {/* Avatar */}
        <View className="w-[52px] h-[52px] rounded-full bg-[#f8f5f2] items-center justify-center">
          <Ionicons name="person-circle-outline" size={52} color="#000" />
        </View>

        {/* Hello text */}
        <View className="flex-1 ml-3">
          <Text className="text-[12px] leading-4">Hello!!!</Text>
          <Text className="text-[16px] font-semibold leading-6">Nhi</Text>
        </View>

        {/* Notification button */}
        <TouchableOpacity
          className="w-10 h-10 rounded-full bg-[#f8f5f2] items-center justify-center"
          onPress={() => router.push("/notification")}
        >
          <Ionicons name="notifications-outline" size={24} color="black" />
        </TouchableOpacity>
      </View>

      {/* ---------- TITLE ---------- */}
      <Text className="text-xl font-bold mt-8 mb-4">Services</Text>

      {/* ---------- SERVICE OPTIONS ---------- */}
      <View className="flex-col">
        {/* --- Available --- */}
        <TouchableOpacity
          className="bg-white rounded-2xl h-[146px] justify-start px-[16px] py-[24px] mb-3"
          onPress={() => router.push("/room")}
        >
          <Ionicons
            name="calendar-outline"
            size={24}
            color="black"
            className="mb-[24px]"
          />
          <Text className="text-base font-bold leading-6 mb-[2px]">
            Available
          </Text>
          <Text className="text-xs text-gray-500 leading-4">
            Checking Availability
          </Text>
        </TouchableOpacity>

        {/* --- Your bookings --- */}
        <TouchableOpacity
          className="bg-white rounded-2xl h-[146px] justify-start px-[16px] py-[24px]"
          onPress={() => router.push("../history")}
        >
          <Ionicons
            name="list-outline"
            size={24}
            color="black"
            className="mb-[24px]"
          />
          <Text className="text-base font-bold leading-6 mb-[2px]">
            Your bookings
          </Text>
          <Text className="text-xs text-gray-500 leading-4">
            View booking history
          </Text>
        </TouchableOpacity>
      </View>
    </View>
  );
}
