import { defineStore } from "pinia";
import type { Member } from "~/types";

export const useMemberStore = defineStore(
  "member",
  () => {
    // State
    const members: Ref<Member[]> = ref([]);
    const loading = ref(false);
    const error = ref<string | null>(null);

    // Mock data (migrated from server/api/members.ts)
    const mockMembers: Member[] = [
      {
        name: "Anthony Fu",
        username: "antfu",
        role: "member",
        avatar: {
          src: "https://ipx.nuxt.com/f_auto,s_192x192/gh_avatar/antfu",
        },
      },
      {
        name: "Baptiste Leproux",
        username: "larbish",
        role: "member",
        avatar: {
          src: "https://ipx.nuxt.com/f_auto,s_192x192/gh_avatar/larbish",
        },
      },
      {
        name: "Benjamin Canac",
        username: "benjamincanac",
        role: "owner",
        avatar: {
          src: "https://ipx.nuxt.com/f_auto,s_192x192/gh_avatar/benjamincanac",
        },
      },
      {
        name: "Céline Dumerc",
        username: "celinedumerc",
        role: "member",
        avatar: {
          src: "https://ipx.nuxt.com/f_auto,s_192x192/gh_avatar/celinedumerc",
        },
      },
      {
        name: "Daniel Roe",
        username: "danielroe",
        role: "member",
        avatar: {
          src: "https://ipx.nuxt.com/f_auto,s_192x192/gh_avatar/danielroe",
        },
      },
      {
        name: "Farnabaz",
        username: "farnabaz",
        role: "member",
        avatar: {
          src: "https://ipx.nuxt.com/f_auto,s_192x192/gh_avatar/farnabaz",
        },
      },
      {
        name: "Ferdinand Coumau",
        username: "FerdinandCoumau",
        role: "member",
        avatar: {
          src: "https://ipx.nuxt.com/f_auto,s_192x192/gh_avatar/FerdinandCoumau",
        },
      },
      {
        name: "Hugo Richard",
        username: "hugorcd",
        role: "owner",
        avatar: {
          src: "https://ipx.nuxt.com/f_auto,s_192x192/gh_avatar/hugorcd",
        },
      },
      {
        name: "Pooya Parsa",
        username: "pi0",
        role: "member",
        avatar: { src: "https://ipx.nuxt.com/f_auto,s_192x192/gh_avatar/pi0" },
      },
      {
        name: "Sarah Moriceau",
        username: "SarahM19",
        role: "member",
        avatar: {
          src: "https://ipx.nuxt.com/f_auto,s_192x192/gh_avatar/SarahM19",
        },
      },
      {
        name: "Sébastien Chopin",
        username: "Atinux",
        role: "owner",
        avatar: {
          src: "https://ipx.nuxt.com/f_auto,s_192x192/gh_avatar/atinux",
        },
      },
    ];

    // Actions
    async function fetchMembers() {
      loading.value = true;
      error.value = null;

      try {
        // Simulate API delay
        await new Promise((resolve) => setTimeout(resolve, 100));

        // In production, replace with actual API call:
        // const data = await $fetch('/api/members')
        members.value = mockMembers;
      } catch (e) {
        error.value =
          e instanceof Error ? e.message : "Failed to fetch members";
        console.error("Error fetching members:", e);
      } finally {
        loading.value = false;
      }
    }

    function addMember(member: Member) {
      members.value.push(member);
    }

    function removeMember(username: string) {
      const index = members.value.findIndex((m) => m.username === username);
      if (index !== -1) {
        members.value.splice(index, 1);
      }
    }

    return {
      // State
      members,
      loading,
      error,

      // Actions
      fetchMembers,
      addMember,
      removeMember,
    };
  },
);
