<script setup lang="ts">
import { ref } from "vue";

let selectedFiles = ref<File[]>([]);
const emit = defineEmits(["close"]);

const handleFileSelect = (e: Event) => {
  const input = e.target as HTMLInputElement;
  const filesAsArray = Array.from(input?.files || []);
  selectedFiles.value = selectedFiles.value.concat(filesAsArray);
};

const uploadFiles = async () => {
  const many = 4;
  const queue = [...selectedFiles.value];

  const worker = async () => {
    while (queue.length > 0) {
      let file = queue.shift();

      if (!file) return;

      await uploadFile(file);
    }
  };

  const workers = Array.from({ length: Math.min(many, queue.length) }, () =>
    worker(),
  );
  await Promise.all(workers);
  close();
};

const uploadFile = async (file: File) => {
  const url = `${import.meta.env.VITE_API_URL}api/upload`;

  const formData = new FormData();
  formData.append("files[]", file);

  const response = await fetch(url, {
    method: "POST",
    body: formData,
  });

  if (!response.ok) {
    throw new Error(`Could not upload ${file.name}`);
  }

  console.log(await response.text());
};
const close = () => {
  selectedFiles.value = [];
  emit("close");
};
</script>
<template>
  <div
    class="fixed inset-0 z-50 flex items-center justify-center bg-black/70 backdrop-blur-sm"
  >
    <button
      class="fixed top-3 right-4 text-white text-xl w-9 h-9 flex items-center justify-center rounded-full bg-white/10 hover:bg-white/20 transition-colors"
      @click="close"
    >
      X
    </button>
    <div class="flex w-125 flex-col gap-5 rounded-2xl bg-white p-8 shadow-2xl">
      <div class="flex gap-3">
        <label
          class="flex-1 rounded-xl border border-gray-200 bg-gray-50 px-4 py-3 text-center text-sm font-medium text-gray-700 transition hover:border-gray-300 hover:bg-gray-100 cursor-pointer"
        >
          Bilder auswählen
          <input
            type="file"
            multiple
            accept="image/* video/*"
            class="hidden"
            @change="handleFileSelect"
          />
        </label>

        <label
          class="flex-1 rounded-xl border border-gray-200 bg-gray-50 px-4 py-3 text-center text-sm font-medium text-gray-700 transition hover:border-gray-300 hover:bg-gray-100 cursor-pointer"
        >
          Ordner auswählen
          <input
            type="file"
            webkitdirectory
            multiple
            class="hidden"
            @change="handleFileSelect"
          />
        </label>
      </div>

      <div
        v-if="selectedFiles.length"
        class="rounded-xl border border-gray-200 bg-gray-50 p-4"
      >
        <p
          class="mb-3 text-xs font-semibold uppercase tracking-wide text-gray-400"
        >
          {{ selectedFiles.length }} Dateien ausgewählt
        </p>

        <ul class="max-h-48 space-y-1 overflow-y-auto">
          <li
            v-for="file in selectedFiles"
            :key="file.name"
            class="truncate rounded-lg px-3 py-2 text-sm text-gray-600 hover:bg-white"
          >
            {{ file.name }}
          </li>
        </ul>
      </div>

      <button
        @click="uploadFiles"
        :disabled="selectedFiles.length === 0"
        class="w-full rounded-xl bg-gray-900 px-4 py-3 text-sm font-medium text-white transition hover:bg-gray-800 disabled:cursor-not-allowed disabled:bg-gray-200 disabled:text-gray-400"
      >
        Upload {{ selectedFiles.length }} Dateien
      </button>
    </div>
  </div>
</template>
