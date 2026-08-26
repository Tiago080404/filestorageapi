<script setup lang="ts">
import { ref } from "vue";

let selectedFiles = ref<File[]>([]);
const handleFileSelect = (e: Event) => {
  const input = e.target as HTMLInputElement;
  const filesAsArray = Array.from(input?.files || []);
  selectedFiles.value = selectedFiles.value.concat(filesAsArray);
};
const uploadFiles = async () => {
  const url = "http://localhost:8080/api/upload";
  const formData = new FormData();
  for (let i = 0; i < selectedFiles.value.length; i++) {
    formData.append("files[]", selectedFiles.value[i]);
    const response = await fetch(url, {
      method: "POST",
      body: formData,
    });
    if (!response.ok) {
      console.log("could not upload");
      return;
    }
    console.log(await response.text());
  }
};
</script>
<template>
  <div class="">
    <input type="file" @change="handleFileSelect" />
    <li v-for="file in selectedFiles">{{ file.name }}:{{ file.size }}</li>
    <button @click="uploadFiles">Upload</button>
  </div>
</template>
