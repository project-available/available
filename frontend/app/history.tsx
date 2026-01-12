import {
  View,
  Text,
  TouchableOpacity,
  FlatList,
  ActivityIndicator,
} from "react-native";
import Ionicons from "@expo/vector-icons/Ionicons";
import { useRouter } from "expo-router";
import { useEffect, useState } from "react";
import AsyncStorage from "@react-native-async-storage/async-storage";

interface Booking {
  id: number;
  account_id: number;
  room_id: number;
  start: string;
  end: string;
  status: string;
  phone_booking: string;
}

export default function History() {
  const router = useRouter();
  const [bookings, setBookings] = useState<Booking[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [cancellingId, setCancellingId] = useState<number | null>(null);

  const handleCancelBooking = async (bookingId: number) => {
    try {
      setCancellingId(bookingId);
      const token = await AsyncStorage.getItem("access_token");

      if (!token) {
        throw new Error("Not logged in. Please login first.");
      }

      const url = `https://querulous-valerie-quanghia-967df8a0.koyeb.app/bookings/update/${bookingId}`;

      const response = await fetch(url, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({
          status: "cancelled",
        }),
      });

      if (!response.ok) {
        throw new Error(`Failed to cancel booking: ${response.status}`);
      }

      // Update local state
      setBookings((prevBookings) =>
        prevBookings.map((booking) =>
          booking.id === bookingId
            ? { ...booking, status: "cancelled" }
            : booking
        )
      );

      console.log("✅ Booking cancelled successfully");
    } catch (err) {
      console.error("❌ Error cancelling booking:", err);
      alert(err instanceof Error ? err.message : "Failed to cancel booking");
    } finally {
      setCancellingId(null);
    }
  };

  useEffect(() => {
    const fetchBookings = async () => {
      try {
        const accountData = await AsyncStorage.getItem("account");
        const token = await AsyncStorage.getItem("access_token");

        console.log("📦 Account:", accountData);
        console.log("🔑 Token:", token ? "exists" : "missing");
        console.log("🔑 Full Token:", token);
        if (!accountData || !token) {
          throw new Error("Not logged in. Please login first.");
        }

        const account = JSON.parse(accountData);
        console.log("👤 Account ID:", account.id);

        const url = `https://querulous-valerie-quanghia-967df8a0.koyeb.app/bookings/${account.id}`;
        console.log("📡 Fetching:", url);

        const response = await fetch(url, {
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
          },
        });

        console.log("📥 Response status:", response.status);

        if (response.status === 401) {
          throw new Error("Session expired. Please login again.");
        }

        if (!response.ok) {
          const errorText = await response.text();
          console.error("❌ Error response:", errorText);
          throw new Error(`HTTP error! status: ${response.status}`);
        }

        const data = await response.json();
        console.log("✅ Bookings loaded:", data.length);

        // Sort by ID descending (newest booking first)
        const sortedData = data.sort((a: Booking, b: Booking) => b.id - a.id);

        setBookings(sortedData);
      } catch (err: any) {
        console.error("❌ Fetch bookings error:", err);
        setError(err.message);
      } finally {
        setLoading(false);
      }
    };

    fetchBookings();
  }, []);

  const getStatusColor = (status: string) => {
    switch (status.toLowerCase()) {
      case "confirmed":
      case "successful":
        return "#69B76B";
      case "cancelled":
        return "#EC1861";
      default:
        return "#000";
    }
  };

  const getStatusText = (status: string) => {
    switch (status.toLowerCase()) {
      case "confirmed":
        return "Successful";
      case "cancelled":
        return "Cancelled";
      case "pending":
        return "Booked";
      default:
        return status;
    }
  };

  const formatDateTime = (start: string, end: string) => {
    const startDate = new Date(start);
    const endDate = new Date(end);

    const dateStr = startDate.toLocaleDateString("en-GB", {
      day: "2-digit",
      month: "short",
      year: "numeric",
    });

    const startTime = startDate.toLocaleTimeString("en-US", {
      hour: "2-digit",
      minute: "2-digit",
      hour12: true,
    });

    const endTime = endDate.toLocaleTimeString("en-US", {
      hour: "2-digit",
      minute: "2-digit",
      hour12: true,
    });

    return `${dateStr} | ${startTime} - ${endTime}`;
  };

  const getTimeAgo = (start: string) => {
    const now = new Date();
    const bookingDate = new Date(start);
    const diffMs = now.getTime() - bookingDate.getTime();
    const diffHours = Math.floor(diffMs / (1000 * 60 * 60));

    if (diffHours < 1) return "Just now";
    if (diffHours < 24) return `${diffHours}hrs ago`;
    const diffDays = Math.floor(diffHours / 24);
    return `${diffDays}d ago`;
  };

  if (loading) {
    return (
      <View className="flex-1 bg-[#f8f5f2] items-center justify-center">
        <ActivityIndicator size="large" color="#607FBA" />
      </View>
    );
  }

  return (
    <View className="flex-1 bg-[#f8f5f2]">
      {/* Custom Header */}
      <View className="flex-row items-center mt-[64px] ml-[24px] gap-3">
        <TouchableOpacity onPress={() => router.back()}>
          <Ionicons name="chevron-back" size={24} color="#000" />
        </TouchableOpacity>

        <Text className="text-[18px] leading-[24px] font-semibold">
          History
        </Text>
      </View>

      {/* Horizontal Line */}
      <View className="mt-[20px] h-[1px] bg-gray-300 mx-[24px]" />

      {/* Booking List */}
      {error ? (
        <View className="flex-1 items-center justify-center px-6">
          <Ionicons name="alert-circle-outline" size={64} color="#EC1861" />
          <Text className="text-red-500 text-center mt-4 text-base">
            {error}
          </Text>
          {error.includes("login") && (
            <TouchableOpacity
              onPress={() => router.push("/(auth)/login")}
              className="mt-6 bg-[#607FBA] px-6 py-3 rounded-[16px]"
            >
              <Text className="text-white font-semibold">Go to Login</Text>
            </TouchableOpacity>
          )}
        </View>
      ) : (
        <FlatList
          data={bookings}
          keyExtractor={(item) => item.id.toString()}
          contentContainerStyle={{ paddingTop: 16, paddingBottom: 20 }}
          renderItem={({ item }) => {
            const now = new Date();
            const startTime = new Date(item.start);
            const canCancel = 
              item.status.toLowerCase() !== "cancelled" && now < startTime;

            return (
              <View className="bg-[#EEEBE5] border border-[#404040] mx-[8px] mb-[20px] rounded-[16px] px-[20px] py-[18px]">
                <View className="flex-row justify-between">
                  {/* Left Side */}
                  <View className="flex-1">
                    {/* Room Name */}
                    <Text className="text-[14px] font-bold text-[#5E5E5E]">
                      Meeting Room - {item.room_id}
                    </Text>

                    {/* Status */}
                    <Text
                      className="text-[12px] mt-[4px]"
                      style={{ color: getStatusColor(item.status) }}
                    >
                      {getStatusText(item.status)}
                    </Text>

                    {/* Date & Time */}
                    <View className="flex-row items-center mt-[8px]">
                      <Ionicons
                        name="calendar-outline"
                        size={20}
                        color="#5E5E5E"
                      />
                      <Text className="text-[12px] font-bold text-[#5E5E5E] ml-[8px]">
                        {formatDateTime(item.start, item.end)}
                      </Text>
                    </View>
                  </View>

                  {/* Right Side */}
                  <View className="items-end justify-between ml-[24px]">
                    {/* Time Ago */}
                    <Text className="text-[10px] text-[#5E5E5E]">
                      {getTimeAgo(item.start)}
                    </Text>

                    {/* Cancel Button (only if booked and before start time) */}
                    {canCancel && (
                      <TouchableOpacity
                        className="w-[76px] h-[20px] bg-[#EC1861] border border-white rounded-[8px] items-center justify-center"
                        onPress={() => handleCancelBooking(item.id)}
                        disabled={cancellingId === item.id}
                      >
                        <Text className="text-white text-[10px] font-semibold">
                          {cancellingId === item.id ? "..." : "Cancel"}
                        </Text>
                      </TouchableOpacity>
                    )}
                  </View>
                </View>
              </View>
            );
          }}
          ListEmptyComponent={
            <View className="flex-1 items-center justify-center mt-[100px]">
              <Text className="text-gray-500">No booking history</Text>
            </View>
          }
        />
      )}
    </View>
  );
}
