import { useEffect, useState, useRef } from "react";
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
  const [bookings, setBookings] = useState([]);
  const [loadingBookings, setLoadingBookings] = useState(false);
  const [bookingError, setBookingError] = useState(null);
  const flatListRef = useRef(null);


  const days = Array.from({ length: 7 }, (_, i) => {
    const date = moment().add(i, "days");
    return {
      date: date.format("YYYY-MM-DD"),
      weekday: date.format("ddd"),
      day: date.format("D"),
    };
  });


  const today = moment().format("YYYY-MM-DD");
  const [selectedDate, setSelectedDate] = useState(today);


  // Auto scroll to today when component mounts
  useEffect(() => {
    if (flatListRef.current && days.length > 0) {
      const todayIndex = days.findIndex(d => d.date === today);
      if (todayIndex >= 0) {
        setTimeout(() => {
          flatListRef.current?.scrollToIndex({
            index: todayIndex,
            animated: true,
            viewPosition: 0.5
          });
        }, 100);
      }
    }
  }, [loading]); // Run after loading is complete


  // Fetch bookings for selected date with retry logic
  useEffect(() => {
    if (!id) return;


    const fetchBookingsWithRetry = async (retryCount = 0) => {
      try {
        setLoadingBookings(true);
        setBookingError(null);
       
        const url = `https://querulous-valerie-quanghia-967df8a0.koyeb.app/rooms/${id}/bookings?date=${selectedDate}`;
        console.log(`🔍 Fetching bookings (attempt ${retryCount + 1}):`, url);
       
        const response = await fetch(url);
       
        // Handle 500 error with retry
        if (response.status === 500) {
          // Get error details from backend
          let errorDetail = '';
          try {
            const errorData = await response.json();
            errorDetail = errorData.error || errorData.message || JSON.stringify(errorData);
          } catch {
            const errorText = await response.text();
            errorDetail = errorText || 'Unknown error';
          }
         
          // Prepared statement/connection errors need retry
          const isConnectionError = errorDetail.includes('prepared statement') ||
                                   errorDetail.includes('bind message') ||
                                   errorDetail.includes('connection') ||
                                   errorDetail.includes('pq:');
         
          if (retryCount < 3 && isConnectionError) {
            console.warn(`⚠️ API connection error (${errorDetail.substring(0, 50)}), retrying attempt ${retryCount + 2}...`);
            await new Promise(resolve => setTimeout(resolve, 500 * (retryCount + 1))); // Delay 500ms, 1s, 1.5s
            return fetchBookingsWithRetry(retryCount + 1);
          } else {
            console.error('❌ API trả về 500:', errorDetail);
            if (isConnectionError) {
              setBookingError('Lỗi kết nối database. Vui lòng thử lại.');
            } else {
              setBookingError(`Lỗi server: ${errorDetail.substring(0, 100)}`);
            }
            setBookings([]);
            return;
          }
        }
       
        if (response.status === 404) {
          console.log('ℹ️ No bookings for this date');
          setBookings([]);
          return;
        }
       
        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }
       
        const data = await response.json();
       
        // Check if data is an array and not empty
        if (!Array.isArray(data)) {
          console.warn('⚠️ API không trả về array:', data);
          setBookings([]);
          return;
        }
       
        // Convert bookings to display format
        const events = data.map((booking) => ({
          id: booking.id,
          title: `Booking #${booking.id}`,
          start_time: moment(booking.start).format("hh:mm A"),
          end_time: moment(booking.end).format("hh:mm A"),
          status: booking.status,
        }));
       
        console.log('✅ Loaded', events.length, 'bookings');
        setBookings(events);
      } catch (err) {
        console.error('❌ Error fetching bookings:', err);
        setBookingError('Không thể tải lịch. Vui lòng thử lại sau.');
        setBookings([]);
      } finally {
        setLoadingBookings(false);
      }
    };


    fetchBookingsWithRetry();
  }, [id, selectedDate]);


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
            ref={flatListRef}
            data={days}
            horizontal
            showsHorizontalScrollIndicator={false}
            keyExtractor={(item) => item.date}
            onScrollToIndexFailed={(info) => {
              const wait = new Promise(resolve => setTimeout(resolve, 100));
              wait.then(() => {
                flatListRef.current?.scrollToIndex({ index: info.index, animated: true });
              });
            }}
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


            {/* Loading indicator for bookings */}
            {loadingBookings && (
              <View className="absolute left-[70px] right-[16px] top-[100px] items-center">
                <ActivityIndicator size="small" color="#607FBA" />
              </View>
            )}


            {/* Error message */}
            {!loadingBookings && bookingError && (
              <View className="absolute left-[70px] right-[16px] top-[100px] items-center bg-red-50 p-3 rounded-xl">
                <Ionicons name="alert-circle-outline" size={24} color="#E74C3C" />
                <Text className="text-red-600 text-xs text-center mt-2">{bookingError}</Text>
              </View>
            )}


            {/* Events from API */}
            {!loadingBookings && bookings.map((ev) => (
              <View
                key={ev.id}
                className="absolute left-[70px] right-[16px] rounded-2xl justify-center px-4"
                style={{
                  top: timeToY(ev.start_time),
                  height: timeToY(ev.end_time) - timeToY(ev.start_time),
                  backgroundColor: ev.status === 'confirmed' ? "#607FBA" : "#95A5C6",
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
