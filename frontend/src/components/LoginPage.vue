<script setup lang="ts">
import { ref } from "vue";
import { router } from "../router";

let username = ref("");
let password = ref("");

const login = async () => {
  const response = await fetch(`${import.meta.env.VITE_API_URL}api/login`, {
    method: "POST",
    body: JSON.stringify({
      username: username.value,
      password: password.value,
    }),
    credentials: "include",
  });

  if (!response.ok) {
    return;
  }
  router.push("/")
};
</script>
<template>
  <div class="flex items-center justify-center p-4">
    <div class="w-full max-w-sm rounded-xl bg-white p-8 flex-col flex gap-3">
      <input
        type="text"
        placeholder="username"
        class="rounded-lg border border-slate-300"
        v-model="username"
      />
      <input
        type="text"
        placeholder="password"
        class="rounded-lg border border-slate-300"
        v-model="password"
      />
      <button
        class="rounded-lg bg-blue-600 px-4 py-3 font-medium text-white hover:cursor-pointer"
        @click="login"
      >
        Login
      </button>
    </div>
  </div>
</template>
<!-- <template>
  <div class="flex min-h-screen items-center justify-center bg-slate-100 p-4">
    <div class="w-full max-w-sm rounded-xl bg-white p-8 shadow-lg">
      <h1 class="mb-6 text-center text-2xl font-semibold text-slate-800">
        Login
      </h1>

      <form class="flex flex-col gap-4">
        <input
          type="text"
          placeholder="Username"
          class="rounded-lg border border-slate-300 px-4 py-3 outline-none transition focus:border-blue-500 focus:ring-2 focus:ring-blue-200"
        />

        <input
          type="password"
          placeholder="Passwort"
          class="rounded-lg border border-slate-300 px-4 py-3 outline-none transition focus:border-blue-500 focus:ring-2 focus:ring-blue-200"
        />

        <button
          type="submit"
          class="rounded-lg bg-blue-600 px-4 py-3 font-medium text-white transition hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-400 focus:ring-offset-2"
        >
          Anmelden
        </button>
      </form>
    </div>
  </div>
</template> -->
