import { View, Text, Image, TouchableOpacity } from "react-native";
import { Feather } from "@expo/vector-icons";

export default function Profile() {
  const mockData = {
    name: "Nhi",
    mssv: "2212416",
    avatar: "https://i.pravatar.cc/100?img=9",
    email: "email@hcmut.edu.vn",
    phone: "+84 82 987 2221",
  };

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
          source={{ uri: mockData.avatar }}
          className="w-[100px] h-[100px] rounded-[24px] border-[4px] border-white"
        />
        <Text className="text-white font-semibold text-[16px] leading-[24px] mt-[8px]">
          {mockData.name}
        </Text>
        <Text className="text-white text-[12px] leading-[16px] mt-[2px]">
          MSSV: {mockData.mssv}
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
            {mockData.email}
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
            {mockData.phone}
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
        <TouchableOpacity className="flex-1 bg-[#EEEBe5] rounded-2xl h-[56px] flex-row items-center ml-[6px]">
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
