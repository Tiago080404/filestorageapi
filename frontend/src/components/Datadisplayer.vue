<script setup lang="ts">
import { onMounted, ref } from "vue";

const props = defineProps<{
  selectedFileName: string;
  fileType: string;
}>();
const emit = defineEmits(["close"]);

let file = ref("");
const displayFile = async () => {
  file.value = `http://localhost:8080/api/open/${props.selectedFileName}`;
  console.log(file.value, props.fileType);
};

const close = () => {
  file.value = "";
  emit("close");
};

const downloadFile = async () => {
  console.log();
  const response = await fetch(
    `http://localhost:8080/api/download/${props.selectedFileName}`,
    {
      method: "GET",
    },
  );
  if (!response.ok) {
    return;
  }
  const data = await response.blob();
  if (props.fileType === "pdf") {
    var file = new Blob([data], { type: "application/pdf" });
    const url = window.URL.createObjectURL(file);
    const a = document.createElement("a");
    a.href = url;
    a.setAttribute("download", props.selectedFileName);
    document.body.appendChild(a);
    a.click();
    URL.revokeObjectURL(url);
    document.body.removeChild(a);
  } else {
    const url = window.URL.createObjectURL(data);
    const a = document.createElement("a");
    a.href = url;
    a.setAttribute("download", props.selectedFileName);
    document.body.appendChild(a);
    a.click();
    URL.revokeObjectURL(url);
    document.body.removeChild(a);
  }
};

onMounted(async () => {
  await displayFile();
});
</script>
<template>
  <div
    v-if="file"
    class="fixed inset-0 z-50 flex items-center justify-center bg-black/80 backdrop-blur-sm"
  >
    <button
      class="fixed top-3 right-4 text-white text-xl w-9 h-9 flex items-center justify-center rounded-full bg-white/10 hover:bg-white/20 transition-colors"
      @click="close"
    >
      X
    </button>
    <div class="relative flex flex-col items-center">
      <img
        v-if="props.fileType === 'img'"
        :src="file"
        class="max-w-[90vw] max-h-[75vh] object-contain rounded-lg shadow-2xl ring-1 ring-white/10"
      />
      <iframe
        v-if="props.fileType === 'pdf'"
        :src="file"
        class="w-[90vw] h-[75vh] rounded-lg object-contain  shadow-2xl ring-1 ring-white/10 border-none"
      ></iframe>
      <video
        v-if="props.fileType === 'vid'"
        :src="file"
        controls
        class="max-w-[90vw] max-h-[75vh] object-contain rounded-lg shadow-2xl ring-1 ring-white/10"
      ></video>
      <div class="flex items-center gap-3 mt-4">
        <img
          src="../assets/downloadicon.svg"
          alt=""
          class="w-6 h-6 invert cursor-pointer hover:opacity-80 transition-opacity"
          @click="downloadFile"
        />
        <span class="text-white/70 text-sm">{{ props.selectedFileName }}</span>
      </div>
    </div>
  </div>
</template>
