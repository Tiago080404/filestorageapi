<script setup lang="ts">
import { onMounted, ref } from "vue";

const props = defineProps<{
  selectedFilePath: string;
  fileType: string;
}>();

const emit = defineEmits(["close"]);
let file = ref("");
let changeFile = ref(false);
let newFileName = ref("");

const displayFile = async () => {
  file.value = `${import.meta.env.VITE_API_URL}api/open/${props.selectedFilePath}`;
  console.log(file.value, props.fileType);
};

const close = () => {
  file.value = "";
  emit("close");
};

const downloadFile = async () => {
  console.log();
  const response = await fetch(
    `${import.meta.env.VITE_API_URL}api/download/${props.selectedFilePath}`,
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
    a.setAttribute("download", props.selectedFilePath);
    document.body.appendChild(a);
    a.click();
    URL.revokeObjectURL(url);
    document.body.removeChild(a);
  } else {
    const url = window.URL.createObjectURL(data);
    const a = document.createElement("a");
    a.href = url;
    a.setAttribute("download", props.selectedFilePath);
    document.body.appendChild(a);
    a.click();
    URL.revokeObjectURL(url);
    document.body.removeChild(a);
  }
};

const deleteFile = async () => {
  const response = await fetch(`${import.meta.env.VITE_API_URL}api/remove`, {
    method: "PATCH",
    body: JSON.stringify({
      path: props.selectedFilePath,
    }),
  });
  if (!response.ok) {
    console.log("Could not delete file");
    return;
  } else {
    close();
  }
};
const renameFile = async () => {
  if (newFileName.value.length === 0) {
    return;
  }
  
  console.log(newFileName.value, props.selectedFilePath);
  const response = await fetch(`${import.meta.env.VITE_API_URL}api/rename`, {
    method: "PATCH",
    body: JSON.stringify({
      newname: newFileName.value,
      oldpath: props.selectedFilePath,
    }),
  });
  if (!response.ok) {
    console.log("could not rename");
    return;
  } else {
    close();
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
        class="w-[90vw] h-[75vh] rounded-lg object-contain shadow-2xl ring-1 ring-white/10 border-none"
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
        <span class="text-white/70 text-sm">{{ props.selectedFilePath }}</span>
        <img
          src="../assets/trashcan.svg"
          alt=""
          class="w-6"
          @click="deleteFile"
        />
        <img
          v-if="!changeFile"
          src="../assets/edit-svgrepo-com.svg"
          @click="changeFile = true"
          class="w-6"
        />
        <div v-if="changeFile">
          <input
            type="text"
            v-model="newFileName"
            class="border-gray-200 bg-gray-50 rounded-lg font-medium"
          />
          <img
            src="../assets/edit-svgrepo-com.svg"
            @click="renameFile"
            class="w-6"
          />
        </div>
      </div>
    </div>
  </div>
</template>
