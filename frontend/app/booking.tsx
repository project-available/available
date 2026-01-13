import {
  View,
  Text,
  TouchableOpacity,
  TextInput,
  Modal,
  Alert,
} from "react-native";
import Ionicons from "@expo/vector-icons/Ionicons";
import { useLocalSearchParams, useRouter } from "expo-router";
import AsyncStorage from "@react-native-async-storage/async-storage";
import { useEffect, useState } from "react";
import DateTimePicker from "@react-native-community/datetimepicker";
import { API_BASE_URL } from "../utils/apiConfig";

export default function Booking() {
  const router = useRouter();
  const [student, setStudent] = useState<any>(null);
  const [phone, setPhone] = useState("");
  const now = new Date();
  const [selectedDate, setSelectedDate] = useState(now);
  const [startTime, setStartTime] = useState(now);
  const [endTime, setEndTime] = useState(
    new Date(now.getTime() + 60 * 60 * 1000)
  );
  const [showDatePicker, setShowDatePicker] = useState(false);
  const [showStartPicker, setShowStartPicker] = useState(false);
  const [showEndPicker, setShowEndPicker] = useState(false);
  const { id } = useLocalSearchParams();

  // Popup state
  const [showPopup, setShowPopup] = useState(false);
  const [popupType, setPopupType] = useState<
    "success" | "fail" | "invalid" | null
  >(null);

  // Load account
  useEffect(() => {
    if (!id) return;
    const loadAccount = async () => {
      const token = await AsyncStorage.getItem("access_token");
      const data = await AsyncStorage.getItem("account");

      if (!token || !data) {
        // Not logged in, redirect to login
        router.replace("/(auth)/login");
        return;
      }

      console.log("📦 Account stored:", data);
      const acc = JSON.parse(data);
      setStudent(acc);
      setPhone(acc.phone);
    };
    loadAccount();
  }, []);

  const combineDateTime = (date: Date, time: Date) => {
    const combined = new Date(date);
    combined.setHours(time.getHours());
    combined.setMinutes(time.getMinutes());
    combined.setSeconds(0);
    combined.setMilliseconds(0);
    return combined.toISOString();
  };

  const handleBooking = async (bookingData: any) => {
    try {
      const token = await AsyncStorage.getItem("access_token");
      const response = await fetch(`${API_BASE_URL}/bookings`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify(bookingData),
      });

      if (!response.ok)
        throw new Error(`HTTP error! status: ${response.status}`);

      const data = await response.json();
      console.log("Booking success:", data);
      return data;
    } catch (err) {
      console.error("Booking failed:", err);
      return null;
    }
  };

  return (
    <View className="flex-1 bg-[#f8f5f2]">
      {/* Header */}
      <View className="flex-row items-center mt-[64px] ml-[24px] gap-3">
        <TouchableOpacity onPress={() => router.back()}>
          <Ionicons name="chevron-back" size={24} color="#000" />
        </TouchableOpacity>
        <Text className="text-[18px] font-semibold">New Booking</Text>
      </View>
      <View className="w-full mt-[20px] h-[1px] bg-[#DBD8D3]" />

      {/* Form */}
      <View className="flex-col justify-center mt-[22px] mx-6 gap-4">
        {/* Full name */}
        <View className="h-[56px] bg-white rounded-2xl border border-[#404040] justify-center px-4">
          <View className="flex-row justify-between items-center">
            <Text className="text-base text-[#8C8C8C]">Full Name:</Text>
            <Text className="text-base font-semibold">{student?.name}</Text>
          </View>
        </View>

        {/* MSSV */}
        <View className="h-[56px] bg-white rounded-2xl border border-[#404040] justify-center px-4">
          <View className="flex-row justify-between items-center">
            <Text className="text-base text-[#8C8C8C]">MSSV:</Text>
            <Text className="text-base font-semibold">
              {student?.student_id}
            </Text>
          </View>
        </View>

        {/* Phone (editable) */}
        <View className="h-[56px] bg-white rounded-2xl border border-[#404040] justify-center px-4">
          <View className="flex-row justify-between items-center">
            <Text className="text-base text-[#8C8C8C]">Phone:</Text>
            <TextInput
              value={phone}
              onChangeText={setPhone}
              keyboardType="phone-pad"
              className="text-base font-semibold text-right w-[140px]"
            />
          </View>
        </View>

        {/* Date & Time */}
        <View className="bg-white rounded-2xl border border-[#404040] p-4">
          <View className="flex-row items-center mb-4">
            <Ionicons name="time-outline" size={16} color="#000" />
            <Text className="ml-3 text-base font-semibold">Date & Time</Text>
          </View>
          <View className="h-[1px] bg-[#DBD8D3]" />

          {/* Date */}
          <TouchableOpacity
            onPress={() => setShowDatePicker(true)}
            className="mt-4 flex-row justify-between items-center"
          >
            <Text className="text-base text-[#8C8C8C]">Date</Text>
            <Text className="text-base font-semibold">
              {selectedDate.toLocaleDateString("en-GB", {
                day: "2-digit",
                month: "short",
                year: "numeric",
              })}
            </Text>
          </TouchableOpacity>
          <View className="mt-4 h-[1px] bg-[#DBD8D3]" />

          {/* Start Time */}
          <TouchableOpacity
            onPress={() => setShowStartPicker(true)}
            className="mt-4 flex-row justify-between items-center"
          >
            <Text className="text-base text-[#8C8C8C]">Start Time</Text>
            <Text className="text-base font-semibold">
              {startTime.toLocaleTimeString([], {
                hour: "2-digit",
                minute: "2-digit",
              })}
            </Text>
          </TouchableOpacity>

          <View className="mt-4 h-[1px] bg-[#DBD8D3]" />

          {/* End Time */}
          <TouchableOpacity
            onPress={() => setShowEndPicker(true)}
            className="mt-4 flex-row justify-between items-center"
          >
            <Text className="text-base text-[#8C8C8C]">End Time</Text>
            <Text className="text-base font-semibold">
              {endTime.toLocaleTimeString([], {
                hour: "2-digit",
                minute: "2-digit",
              })}
            </Text>
          </TouchableOpacity>

          {/* Pickers */}
          {showDatePicker && (
            <DateTimePicker
              value={selectedDate}
              mode="date"
              display="default"
              onChange={(e, date) => {
                setShowDatePicker(false);
                if (date) setSelectedDate(date);
              }}
            />
          )}
          {showStartPicker && (
            <DateTimePicker
              value={startTime}
              mode="time"
              display="default"
              onChange={(e, date) => {
                setShowStartPicker(false);
                if (date) {
                  setStartTime(date);
                  setEndTime(new Date(date.getTime() + 60 * 60 * 1000));
                }
              }}
            />
          )}
          {showEndPicker && (
            <DateTimePicker
              value={endTime}
              mode="time"
              display="default"
              onChange={(e, date) => {
                setShowEndPicker(false);
                if (date) setEndTime(date);
              }}
            />
          )}
        </View>

        {/* Book Room */}
        <TouchableOpacity
          onPress={async () => {
            console.log("Preparing booking data...");
            const startDateTime = combineDateTime(selectedDate, startTime);
            const endDateTime = combineDateTime(selectedDate, endTime);
            console.log("Start:", startDateTime);
            console.log("End:", endDateTime);
            console.log("-----------------------------------------");
            // Validate: start time must be less than end time
            if (new Date(startDateTime) >= new Date(endDateTime)) {
              setPopupType("invalid");
              setShowPopup(true);
              return;
            }

            const bookingData = {
              account_id: student.id,
              room_id: parseInt(id as string, 10),
              start: startDateTime,
              end: endDateTime,
              phone_booking: phone,
            };

            const res = await handleBooking(bookingData);

            if (res) {
              setPopupType("success");
              setShowPopup(true);
            } else {
              setPopupType("fail");
              setShowPopup(true);
            }
          }}
          className="mt-4 h-[56px] border border-[#607FBA]/50 rounded-[40px] bg-[#607FBA] items-center justify-center"
        >
          <Text className="text-sm text-white font-semibold">Book Room</Text>
        </TouchableOpacity>
      </View>

      {/* Custom Popup */}
      <Modal visible={showPopup} transparent animationType="fade">
        <View className="flex-1 items-center justify-center bg-black/50">
          <View className="w-[327px] h-[264px] bg-white rounded-[24px] items-center p-6">
            {/* Icon */}
            <View
              className="w-[64px] h-[64px] rounded-full items-center justify-center mt-2 mb-4"
              style={{
                backgroundColor:
                  popupType === "success"
                    ? "#FFF8E1"
                    : popupType === "invalid"
                    ? "#FFF3E0"
                    : "#FFECEC",
              }}
            >
              <Ionicons
                name={
                  popupType === "success"
                    ? "checkmark"
                    : popupType === "invalid"
                    ? "alert-circle"
                    : "close"
                }
                size={40}
                color={
                  popupType === "success"
                    ? "#FFD54F"
                    : popupType === "invalid"
                    ? "#FF9800"
                    : "#E57373"
                }
              />
            </View>

            {/* Title */}
            <Text className="text-[18px] font-semibold mb-2">
              {popupType === "success"
                ? "Booking Sent"
                : popupType === "invalid"
                ? "Invalid Time"
                : "Booking Failed"}
            </Text>

            {/* Message */}
            <Text className="text-base text-[#666] text-center mb-6">
              {popupType === "success"
                ? "The booking was sent successfully"
                : popupType === "invalid"
                ? "Start time must be earlier than end time"
                : "There was an error sending your booking"}
            </Text>

            {/* OK Button */}
            <TouchableOpacity
              onPress={() => {
                setShowPopup(false);
                if (popupType === "success") router.push("/(tabs)");
              }}
              className="w-full h-[56px] bg-[#607FBA] rounded-[16px] items-center justify-center"
              style={{ marginHorizontal: 16 }}
            >
              <Text className="text-white text-[16px] font-semibold">Ok</Text>
            </TouchableOpacity>
          </View>
        </View>
      </Modal>
    </View>
  );
}
