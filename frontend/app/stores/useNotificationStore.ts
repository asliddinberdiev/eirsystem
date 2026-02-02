import { defineStore } from "pinia";
import { sub } from "date-fns";
import type { Notification } from "~/types";

export const useNotificationStore = defineStore("notification", () => {
  // State
  const notifications = ref<Notification[]>([]);
  const loading = ref(false);
  const error = ref<string | null>(null);

  // Mock data (migrated from server/api/notifications.ts)
  const mockNotifications: Notification[] = [
    {
      id: 1,
      unread: true,
      sender: {
        id: 2,
        name: "Jordan Brown",
        email: "jordan.brown@example.com",
        avatar: {
          src: "https://i.pravatar.cc/128?u=2",
        },
        status: "subscribed",
        location: "London, UK",
      },
      body: "sent you a message",
      date: sub(new Date(), { minutes: 7 }).toISOString(),
    },
    {
      id: 2,
      sender: {
        id: 100,
        name: "Lindsay Walton",
        email: "lindsay.walton@example.com",
        status: "subscribed",
        location: "New York, USA",
      },
      body: "subscribed to your email list",
      date: sub(new Date(), { hours: 1 }).toISOString(),
    },
    {
      id: 3,
      unread: true,
      sender: {
        id: 3,
        name: "Taylor Green",
        email: "taylor.green@example.com",
        avatar: {
          src: "https://i.pravatar.cc/128?u=3",
        },
        status: "subscribed",
        location: "Paris, France",
      },
      body: "sent you a message",
      date: sub(new Date(), { hours: 3 }).toISOString(),
    },
    {
      id: 4,
      sender: {
        id: 101,
        name: "Courtney Henry",
        avatar: {
          src: "https://i.pravatar.cc/128?u=4",
        },
        email: "courtney.henry@example.com",
        status: "subscribed",
        location: "Berlin, Germany",
      },
      body: "added you to a project",
      date: sub(new Date(), { hours: 3 }).toISOString(),
    },
    {
      id: 5,
      sender: {
        id: 102,
        name: "Tom Cook",
        avatar: {
          src: "https://i.pravatar.cc/128?u=5",
        },
        email: "tom.cook@example.com",
        status: "subscribed",
        location: "Tokyo, Japan",
      },
      body: "abandonned cart",
      date: sub(new Date(), { hours: 7 }).toISOString(),
    },
    {
      id: 6,
      sender: {
        id: 103,
        name: "Casey Thomas",
        avatar: {
          src: "https://i.pravatar.cc/128?u=6",
        },
        email: "casey.thomas@example.com",
        status: "subscribed",
        location: "Sydney, Australia",
      },
      body: "purchased your product",
      date: sub(new Date(), { days: 1, hours: 3 }).toISOString(),
    },
    {
      id: 7,
      unread: true,
      sender: {
        id: 8,
        name: "Kelly Wilson",
        email: "kelly.wilson@example.com",
        avatar: {
          src: "https://i.pravatar.cc/128?u=8",
        },
        status: "subscribed",
        location: "London, UK",
      },
      body: "sent you a message",
      date: sub(new Date(), { days: 2 }).toISOString(),
    },
    {
      id: 8,
      sender: {
        id: 6,
        name: "Jamie Johnson",
        email: "jamie.johnson@example.com",
        avatar: {
          src: "https://i.pravatar.cc/128?u=9",
        },
        status: "subscribed",
        location: "Sydney, Australia",
      },
      body: "requested a refund",
      date: sub(new Date(), { days: 5, hours: 4 }).toISOString(),
    },
    {
      id: 9,
      unread: true,
      sender: {
        id: 11,
        name: "Morgan Anderson",
        email: "morgan.anderson@example.com",
        status: "subscribed",
        location: "Tokyo, Japan",
      },
      body: "sent you a message",
      date: sub(new Date(), { days: 6 }).toISOString(),
    },
    {
      id: 10,
      sender: {
        id: 104,
        name: "Drew Moore",
        email: "drew.moore@example.com",
        status: "subscribed",
        location: "Paris, France",
      },
      body: "subscribed to your email list",
      date: sub(new Date(), { days: 6 }).toISOString(),
    },
    {
      id: 11,
      sender: {
        id: 7,
        name: "Riley Davis",
        email: "riley.davis@example.com",
        status: "subscribed",
        location: "New York, USA",
      },
      body: "abandonned cart",
      date: sub(new Date(), { days: 7 }).toISOString(),
    },
    {
      id: 12,
      sender: {
        id: 10,
        name: "Jordan Taylor",
        email: "jordan.taylor@example.com",
        status: "subscribed",
        location: "Berlin, Germany",
      },
      body: "subscribed to your email list",
      date: sub(new Date(), { days: 9 }).toISOString(),
    },
    {
      id: 13,
      sender: {
        id: 8,
        name: "Kelly Wilson",
        email: "kelly.wilson@example.com",
        avatar: {
          src: "https://i.pravatar.cc/128?u=8",
        },
        status: "subscribed",
        location: "London, UK",
      },
      body: "subscribed to your email list",
      date: sub(new Date(), { days: 10 }).toISOString(),
    },
    {
      id: 14,
      sender: {
        id: 6,
        name: "Jamie Johnson",
        email: "jamie.johnson@example.com",
        avatar: {
          src: "https://i.pravatar.cc/128?u=9",
        },
        status: "subscribed",
        location: "Sydney, Australia",
      },
      body: "subscribed to your email list",
      date: sub(new Date(), { days: 11 }).toISOString(),
    },
    {
      id: 15,
      sender: {
        id: 11,
        name: "Morgan Anderson",
        email: "morgan.anderson@example.com",
        status: "subscribed",
        location: "Tokyo, Japan",
      },
      body: "purchased your product",
      date: sub(new Date(), { days: 12 }).toISOString(),
    },
    {
      id: 16,
      sender: {
        id: 105,
        name: "Drew Moore",
        avatar: {
          src: "https://i.pravatar.cc/128?u=16",
        },
        email: "drew.moore@example.com",
        status: "subscribed",
        location: "Paris, France",
      },
      body: "subscribed to your email list",
      date: sub(new Date(), { days: 13 }).toISOString(),
    },
    {
      id: 17,
      sender: {
        id: 7,
        name: "Riley Davis",
        email: "riley.davis@example.com",
        status: "subscribed",
        location: "New York, USA",
      },
      body: "subscribed to your email list",
      date: sub(new Date(), { days: 14 }).toISOString(),
    },
    {
      id: 18,
      sender: {
        id: 10,
        name: "Jordan Taylor",
        email: "jordan.taylor@example.com",
        status: "subscribed",
        location: "Berlin, Germany",
      },
      body: "subscribed to your email list",
      date: sub(new Date(), { days: 15 }).toISOString(),
    },
    {
      id: 19,
      sender: {
        id: 8,
        name: "Kelly Wilson",
        email: "kelly.wilson@example.com",
        avatar: {
          src: "https://i.pravatar.cc/128?u=8",
        },
        status: "subscribed",
        location: "London, UK",
      },
      body: "subscribed to your email list",
      date: sub(new Date(), { days: 16 }).toISOString(),
    },
    {
      id: 20,
      sender: {
        id: 6,
        name: "Jamie Johnson",
        email: "jamie.johnson@example.com",
        avatar: {
          src: "https://i.pravatar.cc/128?u=9",
        },
        status: "subscribed",
        location: "Sydney, Australia",
      },
      body: "purchased your product",
      date: sub(new Date(), { days: 17 }).toISOString(),
    },
    {
      id: 21,
      sender: {
        id: 11,
        name: "Morgan Anderson",
        email: "morgan.anderson@example.com",
        status: "subscribed",
        location: "Tokyo, Japan",
      },
      body: "abandonned cart",
      date: sub(new Date(), { days: 17 }).toISOString(),
    },
    {
      id: 22,
      sender: {
        id: 106,
        name: "Drew Moore",
        email: "drew.moore@example.com",
        status: "subscribed",
        location: "Paris, France",
      },
      body: "subscribed to your email list",
      date: sub(new Date(), { days: 18 }).toISOString(),
    },
    {
      id: 23,
      sender: {
        id: 7,
        name: "Riley Davis",
        email: "riley.davis@example.com",
        status: "subscribed",
        location: "New York, USA",
      },
      body: "subscribed to your email list",
      date: sub(new Date(), { days: 19 }).toISOString(),
    },
    {
      id: 24,
      sender: {
        id: 107,
        name: "Jordan Taylor",
        avatar: {
          src: "https://i.pravatar.cc/128?u=24",
        },
        email: "jordan.taylor@example.com",
        status: "subscribed",
        location: "Berlin, Germany",
      },
      body: "subscribed to your email list",
      date: sub(new Date(), { days: 20 }).toISOString(),
    },
    {
      id: 25,
      sender: {
        id: 8,
        name: "Kelly Wilson",
        email: "kelly.wilson@example.com",
        avatar: {
          src: "https://i.pravatar.cc/128?u=8",
        },
        status: "subscribed",
        location: "London, UK",
      },
      body: "subscribed to your email list",
      date: sub(new Date(), { days: 20 }).toISOString(),
    },
    {
      id: 26,
      sender: {
        id: 6,
        name: "Jamie Johnson",
        email: "jamie.johnson@example.com",
        avatar: {
          src: "https://i.pravatar.cc/128?u=9",
        },
        status: "subscribed",
        location: "Sydney, Australia",
      },
      body: "abandonned cart",
      date: sub(new Date(), { days: 21 }).toISOString(),
    },
    {
      id: 27,
      sender: {
        id: 11,
        name: "Morgan Anderson",
        email: "morgan.anderson@example.com",
        status: "subscribed",
        location: "Tokyo, Japan",
      },
      body: "subscribed to your email list",
      date: sub(new Date(), { days: 22 }).toISOString(),
    },
  ];

  // Actions
  async function fetchNotifications() {
    loading.value = true;
    error.value = null;

    try {
      // Simulate API delay
      await new Promise((resolve) => setTimeout(resolve, 100));

      // In production, replace with actual API call:
      // const data = await $fetch('/api/notifications')
      notifications.value = mockNotifications;
    } catch (e) {
      error.value =
        e instanceof Error ? e.message : "Failed to fetch notifications";
      console.error("Error fetching notifications:", e);
    } finally {
      loading.value = false;
    }
  }

  function markAsRead(id: number) {
    const notification = notifications.value.find((n) => n.id === id);
    if (notification) {
      notification.unread = !notification.unread;
    }
  }

  function clearAll() {
    notifications.value = [];
  }

  return {
    // State
    notifications,
    loading,
    error,

    // Actions
    fetchNotifications,
    markAsRead,
    clearAll,
  };
});
