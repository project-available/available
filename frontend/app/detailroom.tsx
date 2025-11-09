import { useEffect, useState } from "react";
import {
  View,
  Text,
  Image,
  TouchableOpacity,
  ActivityIndicator,
  FlatList,
  Dimensions,
  ScrollView,
} from "react-native";
import { useLocalSearchParams, useRouter } from "expo-router";
import Ionicons from "@expo/vector-icons/Ionicons";
import moment from "moment";

export default function DetailRoom() {
  const router = useRouter();
  const { id } = useLocalSearchParams();
  const [room, setRoom] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  const mockSchedule = {
    days: Array.from({ length: 7 }, (_, i) => {
      const date = moment().add(i, "days");
      return {
        date: date.format("YYYY-MM-DD"),
        weekday: date.format("ddd"),
        day: date.format("D"),
        events:
          i === 1
            ? [
                {
                  id: 1,
                  title: "Meeting",
                  start_time: "09:00 AM",
                  end_time: "10:00 AM",
                },
                {
                  id: 2,
                  title: "Team Meeting",
                  start_time: "02:00 PM",
                  end_time: "03:30 PM",
                },
              ]
            : [],
      };
    }),
  };

  const today = moment().format("YYYY-MM-DD");
  const [selectedDate, setSelectedDate] = useState(
    mockSchedule.days.find((d) => d.date === today)?.date || today
  );
  const selectedDayData = mockSchedule.days.find(
    (d) => d.date === selectedDate
  );

  useEffect(() => {
    if (!id) return;

    const fetchRoom = async () => {
      try {
        const response = await fetch(
          `https://querulous-valerie-quanghia-967df8a0.koyeb.app/rooms/${id}`
        );
        if (!response.ok)
          throw new Error(`HTTP error! status: ${response.status}`);
        const data = await response.json();
        setRoom(data);
      } catch (err) {
        setError(err.message);
      } finally {
        setLoading(false);
      }
    };

    fetchRoom();
  }, [id]);

  // ---------- LOADING ----------
  if (loading) {
    return (
      <View className="flex-1 items-center justify-center bg-[#f8f5f2]">
        <ActivityIndicator size="large" color="#607FBA" />
        <Text className="mt-3 text-gray-500 text-base">
          Đang tải dữ liệu...
        </Text>
      </View>
    );
  }

  // ---------- ERROR ----------
  if (error || !room) {
    return (
      <View className="flex-1 items-center justify-center bg-[#f8f5f2]">
        <Text className="text-gray-500 text-base">
          {error ? `Lỗi: ${error}` : "Room not found"}
        </Text>
      </View>
    );
  }

  const capacity =
    room.customFields?.find((f) => f.key === "Seats")?.value || "Unknown";

  // ---------- HELPER ----------
  const timeToY = (time) => {
    const [hour, minutePart] = time.split(":");
    const minute = parseInt(minutePart);
    const isPM = time.includes("PM") && parseInt(hour) !== 12;
    const hours24 = (parseInt(hour) % 12) + (isPM ? 12 : 0);
    return (hours24 + minute / 60 - 8) * 80 + 9;
  };

  const currentTimeY = () => {
    const now = new Date();
    const hour = now.getHours() + now.getMinutes() / 60;
    return (hour - 8) * 80 + 6;
  };

  return (
    <View className="flex-1 bg-[#f8f5f2]">
      {/* ---------- HEADER ---------- */}
      <View className="flex-row items-center mt-[64px] mx-[24px] gap-3">
        <TouchableOpacity onPress={() => router.back()}>
          <Ionicons name="chevron-back" size={24} color="#000" />
        </TouchableOpacity>
        <Text className="text-[18px] leading-[24px] font-semibold">
          {room.name}
        </Text>
      </View>

      {/* ---------- ROOM INFO ---------- */}
      <View className="mt-[52px] mx-[24px]">
        <Image
          source={{ uri: room.image }}
          className="w-full h-[160px] rounded-3xl"
        />

        <View className="mt-4 flex-row gap-3">
          <View className="flex-1 h-[36px] bg-[#EEEBE5] justify-center items-center rounded-3xl">
            <Text className="text-[14px] font-semibold text-[#2E2F30]">
              {capacity}
            </Text>
          </View>

          <View className="flex-1 h-[36px] bg-[#EEEBE5] justify-center items-center rounded-3xl">
            <Text className="text-[14px] font-semibold text-[#2E2F30]">
              {room.name}
            </Text>
          </View>

          <View className="flex-1 h-[36px] bg-[#EEEBE5] justify-center items-center rounded-3xl">
            <Text className="text-[14px] font-semibold text-[#2E2F30]">
              {room.location}
            </Text>
          </View>
        </View>

        <TouchableOpacity
          onPress={() =>
            router.push({
              pathname: "/booking",
              params: { id: room.roomId },
            })
          }
          className="mt-[24px] h-[36px] border border-[#607FBA]/50 rounded-3xl bg-white flex-row items-center justify-center space-x-2"
        >
          <Ionicons name="calendar-outline" size={16} color="#607FBA" />
          <Text className="ml-1 text-sm leading-5 text-[#607FBA] font-semibold">
            Book Now
          </Text>
        </TouchableOpacity>

        <View className="mt-6 h-[1px] bg-[#DBD8D3]" />
      </View>

      {/* ---------- CALENDAR ---------- */}
      <View className="mt-6 mx-[24px] flex-1">
        {/* Date selection bar */}
        <View className="flex-row items-center mb-3">
          {/* Month-year box, not clickable */}
          <View className="mx-1 w-[64px] h-[64px] rounded-2xl justify-center items-center bg-[#EEEBE5]">
            <Text className="text-[14px] font-semibold text-[#2E2F30]">
              {moment(selectedDate).format("MMM")}
            </Text>
            <Text className="text-[12px] text-[#2E2F30]">
              {moment(selectedDate).format("YYYY")}
            </Text>
          </View>

          {/* List of days */}
          <FlatList
            data={mockSchedule.days}
            horizontal
            showsHorizontalScrollIndicator={false}
            keyExtractor={(item) => item.date}
            renderItem={({ item }) => (
              <TouchableOpacity
                onPress={() => setSelectedDate(item.date)}
                className={`mx-1 w-[48px] h-[64px] rounded-2xl justify-center items-center ${
                  selectedDate === item.date ? "bg-[#FDBA29]" : "bg-[#EEEBE5]"
                }`}
              >
                <Text
                  className={`text-[16px] font-semibold ${
                    selectedDate === item.date ? "text-white" : "text-[#2E2F30]"
                  }`}
                >
                  {item.day}
                </Text>
                <Text
                  className={`text-[12px] ${
                    selectedDate === item.date ? "text-white" : "text-[#2E2F30]"
                  }`}
                >
                  {item.weekday}
                </Text>
              </TouchableOpacity>
            )}
          />
        </View>

        <View className="flex-1 overflow-hidden mt-30">
          {/* Timeline */}
          <ScrollView
            className="relative mt-4"
            contentContainerStyle={{ paddingTop: 8, height: 14 * 80 }}
          >
            {/* Hour rows */}
            {Array.from({ length: 14 }, (_, i) => (
              <View
                key={i}
                className="absolute left-0 right-0 flex-row items-center"
                style={{ top: i * 80 }}
              >
                {/* Hour column */}
                <View className="w-[60px] items-end pr-2">
                  <Text className="text-gray-400 text-sm">{8 + i}AM</Text>
                </View>

                {/* Horizontal line */}
                <View className="flex-1 h-[1px] bg-gray-200" />
              </View>
            ))}

            {/* Current time line */}
            {selectedDate === moment().format("YYYY-MM-DD") && (
              <>
                <View
                  className="absolute left-[65px] right-0 border-t-2 border-[#607FBA]"
                  style={{ top: currentTimeY() + 2 }}
                />

                <View
                  className="absolute bg-[#607FBA] rounded-full"
                  style={{
                    width: 8,
                    height: 8,
                    top: currentTimeY() - 1,
                    left: 60 - 4,
                  }}
                />
              </>
            )}

            {/* Event */}
            {selectedDayData?.events.map((ev) => (
              <View
                key={ev.id}
                className="absolute left-[70px] right-[16px] rounded-2xl justify-center px-4"
                style={{
                  top: timeToY(ev.start_time),
                  height: timeToY(ev.end_time) - timeToY(ev.start_time),
                  backgroundColor: "#607FBA",
                }}
              >
                <Text className="text-white font-semibold">{ev.title}</Text>
                <Text className="text-white text-xs mt-1">
                  {ev.start_time} - {ev.end_time}
                </Text>
              </View>
            ))}
          </ScrollView>
        </View>
      </View>
    </View>
  );
}
